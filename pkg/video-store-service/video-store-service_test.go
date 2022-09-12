package video_store_service

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
	mock_object_storage "video-manager/internal/mock/object-storage"
	mock_video_hosting "video-manager/internal/mock/video-hosting"
	object_storage "video-manager/internal/object-storage"
	video_hosting "video-manager/internal/video-hosting"
)

type mocked struct {
	videoStore       *mock_video_hosting.MockIVideoHost
	objectStoreProxy *mock_object_storage.MockBindingProxy
	service          VideoStoreService[*mock_object_storage.MockBindingProxy]
}

func Setup(t *testing.T) *mocked {
	ctx := context.Background()
	dir, err := os.MkdirTemp("", "assets")
	if err != nil {
		t.Fatal(err)
	}
	objStoreCtrl := gomock.NewController(t)
	proxy := mock_object_storage.NewMockBindingProxy(objStoreCtrl)
	objectStore := object_storage.NewObjectStorage[*mock_object_storage.MockBindingProxy](&ctx, dir, proxy)
	vidCtrl := gomock.NewController(t)
	vidHost := mock_video_hosting.NewMockIVideoHost(vidCtrl)
	vss := VideoStoreService[*mock_object_storage.MockBindingProxy]{
		ObjStore: objectStore,
		VidHost:  vidHost,
	}
	return &mocked{
		videoStore:       vidHost,
		objectStoreProxy: proxy,
		service:          vss,
	}
}

func TestVideoStoreService_UploadFromObjectStore_DownloadError(t *testing.T) {
	deps := Setup(t)
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
	deps := Setup(t)
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
	deps := Setup(t)
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
	deps := Setup(t)
	_, err := deps.service.UploadVideoFromStorage("test", nil)
	assert.NotNil(t, err)
}
