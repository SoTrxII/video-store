package video_store_service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	mock_object_storage "video-manager/internal/mock/object-storage"
	object_storage "video-manager/internal/object-storage"
)

func SetupFactory(t *testing.T) *object_storage.ObjectStorage[*mock_object_storage.MockBindingProxy] {
	objStoreCtrl := gomock.NewController(t)
	proxy := mock_object_storage.NewMockBindingProxy(objStoreCtrl)
	dir, err := os.MkdirTemp("", "assets")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	return object_storage.NewObjectStorage[*mock_object_storage.MockBindingProxy](&ctx, dir, proxy)
}
func Test_VideoServiceFactory_MakeYoutubeVideoStoreService_Youtube(t *testing.T) {
	objStore := SetupFactory(t)
	_, err := MakeVideoStoreService[*mock_object_storage.MockBindingProxy](context.TODO(), Youtube, *objStore)
	assert.Nil(t, err)
}

func Test_VideoServiceFactory_MakeVideoStoreService_Error(t *testing.T) {
	objStore := SetupFactory(t)
	_, err := MakeVideoStoreService[*mock_object_storage.MockBindingProxy](context.TODO(), Host(4), *objStore)
	assert.NotNil(t, err)
}
