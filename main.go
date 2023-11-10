package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"strconv"
	"time"
	playlists_controller "video-manager/controller/playlists"
	videos_controller "video-manager/controller/videos"
	_ "video-manager/docs"
	"video-manager/internal/logger"
	object_storage "video-manager/internal/object-storage"
	progress_broker "video-manager/internal/progress-broker"
	video_store_service "video-manager/pkg/video-store-service"
)

const (
	DefaultAppPort              = 8080
	DefaultDaprGrpcPort         = 50001
	DefaultDaprMaxRequestSizeMb = 2000
	// env
	APP_PORT                 = "APP_PORT"
	DAPR_GRPC_PORT           = "DAPR_GRPC_PORT"
	DAPR_MAX_REQUEST_SIZE_MB = "DAPR_MAX_REQUEST_SIZE_MB"
	OBJECT_STORE_NAME        = "OBJECT_STORE_NAME"
	GIN_MODE                 = "GIN_MODE"
	PUBSUB_NAME              = "PUBSUB_NAME"
	PUBSUB_TOPIC_PROGRESS    = "PUBSUB_TOPIC_PROGRESS"

	// Topic to send progress event into
	DefaultPubSubTopic = "upload-state"
)

var (
	log = logger.Build()
)

// @title           Video store
// @version         1.0
// @description     An API to store and retrieve videos from a remote hosting service
// @host      localhost:8080
// @BasePath  /

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Info("No .env file loaded !")
	}
	// Env is loaded after gin is initialized, we must set it manually
	if os.Getenv(GIN_MODE) == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	ctx := context.Background()
	vidCtrl, playlistCtrl := resolveDI(&ctx)
	router := gin.Default()

	router.Use(func() gin.HandlerFunc {
		// Change default logger to an ecs compliant one
		// We use a closure to only build one logger instance
		bufferLogger := logger.Build()
		buf := new(bytes.Buffer)
		bufferLogger.Out = buf
		return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			bufferLogger.Infof(`%s - [%s] "%s %s %s %d %s "%s" %s"`,
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
			return buf.String()
		})
	}())

	// Define all routes
	v1 := router.Group("/v1")
	{
		videos := v1.Group("/videos")
		{
			videos.POST("", vidCtrl.Create)
			videos.GET(":id", vidCtrl.Retrieve)
			videos.PUT(":id", vidCtrl.Update)
			videos.DELETE(":id", vidCtrl.Delete)
			videos.POST(":id/thumbnail/:tId", vidCtrl.SetThumbnail)
		}
		playlists := v1.Group("/playlists")
		{
			playlists.POST("", playlistCtrl.Create)
			playlists.GET(":id", playlistCtrl.Retrieve)
			playlists.PUT(":id", playlistCtrl.Update)
			playlists.DELETE(":id", playlistCtrl.Delete)
			playlists.PUT(":id/videos/:vid", playlistCtrl.AddVideo)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	appPort := DefaultAppPort
	if envPort, err := strconv.ParseInt(os.Getenv("APP_PORT"), 10, 32); err == nil && envPort != 0 {
		appPort = int(envPort)
	}

	log.Infof("Server listening to 0.0.0.0:%d", appPort)
	err = router.Run(fmt.Sprintf(":%d", appPort))
	if err != nil {
		log.Fatalf(err.Error())
	}
}

// Resolve the pseudo DI-container
func resolveDI(ctx *context.Context) (videos_controller.VideoController[client.Client, client.Client], playlists_controller.PlaylistController[client.Client, client.Client]) {
	// From bottom to top:
	// Make a new Dapr instance
	daprMaxRqSize := DefaultDaprMaxRequestSizeMb
	if sizeStr, ok := os.LookupEnv(DAPR_MAX_REQUEST_SIZE_MB); ok {
		if size, err := strconv.ParseUint(sizeStr, 10, 64); err != nil {
			daprMaxRqSize = int(size)
		} else {
			log.Warnf("Invalid dapr max request size supplied %s, defaulting to %d", sizeStr, DefaultAppPort)
		}
	}
	proxy, err := makeDaprClient(daprMaxRqSize)
	if err != nil {
		log.Fatalf("Error during init : %s", err.Error())
	}

	// Resolve the required object storage backend
	storeName := ""
	if storeName = os.Getenv(OBJECT_STORE_NAME); storeName == "" {
		log.Fatalf("Error during init : No dapr store defined !")
	}
	objStore, err := object_storage.NewDaprObjectStorage(ctx, proxy, storeName)
	if err != nil {
		log.Fatalf("Error during init : %s", err.Error())
	}

	// Resolve the optional event broker to send upload progress
	var progressBroker *progress_broker.ProgressBroker[client.Client]
	pubsubName := ""
	if pubsubName = os.Getenv(PUBSUB_NAME); pubsubName != "" {
		topic, exists := os.LookupEnv(PUBSUB_TOPIC_PROGRESS)
		if !exists {
			topic = DefaultPubSubTopic
		}
		log.Infof(`Initializing pubsub with name "%s" and topic "%s"`, pubsubName, topic)
		progressBroker, err = progress_broker.NewProgressBroker[client.Client](ctx, proxy, progress_broker.NewBrokerOptions{
			Component: pubsubName,
			Topic:     topic,
		})
		if err != nil {
			log.Fatalf("Couldn't init pubsub : %s", err.Error())
		}
		log.Infof("Pubsub initialized")
	} else {
		log.Infof("No pubsub name provided. Skipping pubsub initialization")
	}

	// We can then resolve the video store service...
	storeService, err := video_store_service.MakeVideoStoreService[client.Client](*ctx, video_store_service.Youtube, *objStore, progressBroker)
	// With in turn give us the controllers
	vCtrl := videos_controller.VideoController[client.Client, client.Client]{Service: storeService}
	pCtrl := playlists_controller.PlaylistController[client.Client, client.Client]{Service: storeService}
	return vCtrl, pCtrl
}

// Make a custom dapr client with a large max request size, to handle large uploads
func makeDaprClient(maxRequestSizeMB int) (*client.Client, error) {
	var opts []grpc.CallOption

	// Getting dapr grpc port. By default, its 500001
	port := DefaultDaprGrpcPort
	// But the sidecar published a env variable with the real value
	// So we can override the value if it's defined
	if envPort, err := strconv.ParseInt(os.Getenv(DAPR_GRPC_PORT), 10, 32); err != nil && envPort != 0 {
		port = int(envPort)
	}
	opts = append(opts, grpc.MaxCallRecvMsgSize(maxRequestSizeMB*1024*1024))
	conn, err := grpc.Dial(net.JoinHostPort("127.0.0.1", fmt.Sprintf("%d", port)),
		grpc.WithDefaultCallOptions(opts...), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	daprClient := client.NewClientWithConnection(conn)
	return &daprClient, nil
}
