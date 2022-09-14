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
	video_store_service "video-manager/pkg/video-store-service"
)

const (
	Port                = 8080
	DefaultDaprGrpcPort = 500001
	// env
	DAPR_GRPC_PORT    = "DAPR_GRPC_PORT"
	OBJECT_STORE_NAME = "OBJECT_STORE_NAME"
	GIN_MODE          = "GIN_MODE"
	//PUBSUB_NAME           = "PUBSUB_NAME"
	//PUBSUB_TOPIC_PROGRESS = "PUBSUB_TOPIC_PROGRESS"
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
			videos.POST(":id/thumbnail", vidCtrl.SetThumbnail)
		}
		playlists := v1.Group("/playlists")
		{
			playlists.POST("", playlistCtrl.Create)
			playlists.GET(":id", playlistCtrl.Retrieve)
			playlists.PUT(":id", playlistCtrl.Update)
			playlists.DELETE(":id", playlistCtrl.Delete)
			playlists.PUT(":id/video/:vid", playlistCtrl.AddVideo)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	err = router.Run(fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatalf(err.Error())
	}
}

// Resolve the pseudo DI-container
func resolveDI(ctx *context.Context) (videos_controller.VideoController[client.Client], playlists_controller.PlaylistController[client.Client]) {
	// From bottom to top:
	// - Instanciate Dapr...
	// TODO This should be configurable
	proxy, err := makeDaprClient(2000)
	if err != nil {
		log.Fatalf("Error during init : %s", err.Error())
	}
	storeName := ""
	if storeName = os.Getenv(OBJECT_STORE_NAME); storeName == "" {
		log.Fatalf("Error during init : No dapr store defined !")
	}
	objStore, err := object_storage.NewDaprObjectStorage(ctx, proxy, storeName)
	if err != nil {
		log.Fatalf("Error during init : %s", err.Error())
	}
	// With Dapr and the storage client, we can then resolve the video store service...
	storeService, err := video_store_service.MakeVideoStoreService[client.Client](*ctx, video_store_service.Youtube, *objStore)
	// With in turn give us the controllers
	return videos_controller.VideoController[client.Client]{Service: storeService}, playlists_controller.PlaylistController[client.Client]{Service: storeService}
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
