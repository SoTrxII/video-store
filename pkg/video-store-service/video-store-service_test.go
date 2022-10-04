package video_store_service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
	mock_object_storage "video-manager/internal/mock/object-storage"
	mock_progress_broker "video-manager/internal/mock/progress-broker"
	mock_video_hosting "video-manager/internal/mock/video-hosting"
	object_storage "video-manager/internal/object-storage"
	progress_broker "video-manager/internal/progress-broker"
	video_hosting "video-manager/internal/video-hosting"
)

type mocked struct {
	videoStore       *mock_video_hosting.MockIVideoHost
	objectStoreProxy *mock_object_storage.MockBindingProxy
	brokerProxy      *mock_progress_broker.MockPubSubProxy
	service          VideoStoreService[*mock_object_storage.MockBindingProxy, *mock_progress_broker.MockPubSubProxy]
}

func Setup(t *testing.T, initBroker bool) *mocked {
	ctx := context.Background()
	dir, err := os.MkdirTemp("", "assets")
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	// Initialize object storage
	objStoreProxy := mock_object_storage.NewMockBindingProxy(ctrl)
	objectStore := object_storage.NewObjectStorage[*mock_object_storage.MockBindingProxy](&ctx, dir, objStoreProxy)

	// Initialize video host
	vidCtrl := gomock.NewController(t)
	vidHost := mock_video_hosting.NewMockIVideoHost(vidCtrl)

	//  Initialize event broker
	psProxy := mock_progress_broker.NewMockPubSubProxy(ctrl)
	broker, err := progress_broker.NewProgressBroker[*mock_progress_broker.MockPubSubProxy](&ctx, &psProxy, progress_broker.NewBrokerOptions{
		Component: "",
		Topic:     "",
	})
	if err != nil {
		t.Fatal(err)
	}

	vss := VideoStoreService[*mock_object_storage.MockBindingProxy, *mock_progress_broker.MockPubSubProxy]{
		ObjStore: objectStore,
		VidHost:  vidHost,
	}
	// The broker can either be included in the controller or set to nil
	// in which case no progress event will be sent
	if initBroker {
		vss.EvtBroker = broker
	}
	return &mocked{
		videoStore:       vidHost,
		objectStoreProxy: objStoreProxy,
		brokerProxy:      psProxy,
		service:          vss,
	}
}

func TestVideoStoreService_UploadFromObjectStore_DownloadError(t *testing.T) {
	deps := Setup(t, false)
	// Setup the proxy to fail to simulate a download error
	deps.objectStoreProxy.EXPECT().InvokeBinding(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("test"))

	_, err := deps.service.UploadVideoFromStorage("test", &video_hosting.ItemMetadata{
		Description: "desc",
		Title:       "title",
		Visibility:  "unlisted",
	})
	assert.NotNil(t, err)
}

func TestVideoStoreService_UploadFromObjectStore_CreateVideoError(t *testing.T) {
	deps := Setup(t, false)
	// Setup the proxy to fail to simulate a download error
	deps.objectStoreProxy.EXPECT().InvokeBinding(gomock.Any(), gomock.Any()).Return(&client.BindingEvent{Data: []byte("a")}, nil)
	deps.videoStore.EXPECT().CreateVideo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("test"))
	_, err := deps.service.UploadVideoFromStorage("test", &video_hosting.ItemMetadata{
		Description: "desc",
		Title:       "title",
		Visibility:  "unlisted",
	})
	assert.NotNil(t, err)
	fmt.Println(err)
}

func TestVideoStoreService_UploadFromObjectStore_Ok(t *testing.T) {
	deps := Setup(t, false)
	deps.objectStoreProxy.EXPECT().InvokeBinding(gomock.Any(), gomock.Any()).Return(&client.BindingEvent{Data: []byte("a")}, nil)
	deps.videoStore.EXPECT().CreateVideo(gomock.Any(), gomock.Any(), gomock.Any()).Return(&video_hosting.Video{
		Id:           "test",
		Title:        "test",
		Description:  "test",
		CreatedAt:    time.Time{},
		Duration:     0,
		Visibility:   "unlisted",
		ThumbnailUrl: "",
	}, nil)
	_, err := deps.service.UploadVideoFromStorage("test", &video_hosting.ItemMetadata{
		Description: "desc",
		Title:       "title",
		Visibility:  "unlisted",
	})
	assert.Nil(t, err)
}

func TestVideoStoreService_UploadFromObjectStore_InvalidMetadata(t *testing.T) {
	deps := Setup(t, false)
	_, err := deps.service.UploadVideoFromStorage("test", nil)
	assert.NotNil(t, err)
}

func TestVideoStoreService_UploadFromObjectStore_WithProgress(t *testing.T) {
	deps := Setup(t, true)
	deps.brokerProxy.
		EXPECT().
		PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	deps.objectStoreProxy.EXPECT().InvokeBinding(gomock.Any(), gomock.Any()).Return(&client.BindingEvent{Data: []byte("a")}, nil)
	deps.videoStore.EXPECT().CreateVideo(gomock.Any(), gomock.Any(), gomock.Any()).Return(&video_hosting.Video{
		Id:           "test",
		Title:        "test",
		Description:  "test",
		CreatedAt:    time.Time{},
		Duration:     0,
		Visibility:   "unlisted",
		ThumbnailUrl: "",
	}, nil)
	_, err := deps.service.UploadVideoFromStorage("test", &video_hosting.ItemMetadata{
		Description: "desc",
		Title:       "title",
		Visibility:  "unlisted",
	})
	assert.Nil(t, err)
}

func TestVideoStoreService_ProgressRoutine_Ok(t *testing.T) {
	deps := Setup(t, true)

	// Prepare two progress event to send to the routine

	// The first one will simulate a mid-upload progress event
	progressEvent := uploadProgress{
		Current: 1,
		Total:   1,
	}
	pgEvt, err := json.Marshal(progress_broker.UploadInfos{
		RecordId: "test",
		State:    progress_broker.InProgress,
		Data:     progressEvent,
	})
	if err != nil {
		t.Fatal(err)
	}
	deps.brokerProxy.
		EXPECT().
		PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), string(pgEvt), gomock.Any()).
		Times(2)

	// And the second one a "Done" event, returned when the upload succeeded
	doneEvt, err := json.Marshal(progress_broker.UploadInfos{
		RecordId: "test",
		State:    progress_broker.Done,
		Data: uploadDone{
			Id:          "test",
			WatchPrefix: "",
			Duration:    0,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	deps.brokerProxy.
		EXPECT().
		PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), string(doneEvt), gomock.Any()).
		Times(1)

	// Make the channels and start the routine
	pgChannel := make(chan uploadProgress)
	resCh := make(chan uploadResult)
	go deps.service.startProgressRoutine("test", time.Second, pgChannel, resCh)

	// Send the two progress events first
	pgChannel <- progressEvent
	pgChannel <- progressEvent
	// And signal the subprocess to end with no error
	resCh <- uploadResult{
		Result: &video_hosting.Video{
			Id:           "test",
			Title:        "test",
			Description:  "test",
			CreatedAt:    time.Time{},
			Duration:     0,
			Visibility:   "unlisted",
			ThumbnailUrl: "test",
		},
		Error: nil,
	}

	// In the meantime, make the main process wait for 10s max
	// after which the test will fail. This prevents infinite loops
	ticker := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-pgChannel:
			fmt.Println("SubProcess stopped")
			return
		case <-ticker.C:
			t.Fail()
			fmt.Println("Timeout")
			return
		}
	}
}

func TestVideoStoreService_ProgressRoutine_Error(t *testing.T) {
	deps := Setup(t, true)

	// Prepare two progress event to send to the routine

	// The first one will simulate a mid-upload progress event
	progressEvent := uploadProgress{
		Current: 1,
		Total:   1,
	}
	pgEvt, err := json.Marshal(progress_broker.UploadInfos{
		RecordId: "test",
		State:    progress_broker.InProgress,
		Data:     progressEvent,
	})
	if err != nil {
		t.Fatal(err)
	}
	deps.brokerProxy.
		EXPECT().
		PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), string(pgEvt), gomock.Any()).
		Times(2)

	// And the second one a "Done" event, returned when the upload succeeded
	errorEvt, err := json.Marshal(progress_broker.UploadInfos{
		RecordId: "test",
		State:    progress_broker.Error,
		Data:     uploadError{Message: "test"},
	})
	if err != nil {
		t.Fatal(err)
	}
	deps.brokerProxy.
		EXPECT().
		PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), string(errorEvt), gomock.Any()).
		Times(1)

	// This should never be called for an error
	deps.videoStore.
		EXPECT().
		GetVideoAccessPrefix().
		Times(0)
	// Make the channels and start the routine
	pgChannel := make(chan uploadProgress)
	resCh := make(chan uploadResult)
	go deps.service.startProgressRoutine("test", time.Second, pgChannel, resCh)

	// Send the two progress events first
	pgChannel <- progressEvent
	pgChannel <- progressEvent
	// And signal the subprocess to end with an error
	resCh <- uploadResult{
		Result: nil,
		Error:  fmt.Errorf("test"),
	}
	// In the meantime, make the main process wait for 10s max
	// after which the test will fail. This prevents infinite loops
	ticker := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-pgChannel:
			fmt.Println("SubProcess stopped")
			return
		case <-ticker.C:
			t.Fail()
			fmt.Println("Timeout")
			return
		}
	}
}
