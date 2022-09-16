package video_store_service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	mock_object_storage "video-manager/internal/mock/object-storage"
	mock_progress_broker "video-manager/internal/mock/progress-broker"
	object_storage "video-manager/internal/object-storage"
	progress_broker "video-manager/internal/progress-broker"
)

func SetupFactory(t *testing.T) (*object_storage.ObjectStorage[*mock_object_storage.MockBindingProxy], *progress_broker.ProgressBroker[*mock_progress_broker.MockPubSubProxy]) {
	ctrl := gomock.NewController(t)
	proxy := mock_object_storage.NewMockBindingProxy(ctrl)
	dir, err := os.MkdirTemp("", "assets")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	psProxy := mock_progress_broker.NewMockPubSubProxy(ctrl)
	broker, err := progress_broker.NewProgressBroker[*mock_progress_broker.MockPubSubProxy](&ctx, &psProxy, progress_broker.NewBrokerOptions{})
	if err != nil {
		t.Fatal(err)
	}
	return object_storage.NewObjectStorage[*mock_object_storage.MockBindingProxy](&ctx, dir, proxy), broker
}
func Test_VideoServiceFactory_MakeYoutubeVideoStoreService_Youtube(t *testing.T) {
	objStore, _ := SetupFactory(t)
	_, err := MakeVideoStoreService[*mock_object_storage.MockBindingProxy, *mock_progress_broker.MockPubSubProxy](context.TODO(), Youtube, *objStore, nil)
	assert.Nil(t, err)
}

func Test_VideoServiceFactory_MakeYoutubeVideoStoreService_Youtube_WithBroker(t *testing.T) {
	objStore, broker := SetupFactory(t)
	_, err := MakeVideoStoreService[*mock_object_storage.MockBindingProxy, *mock_progress_broker.MockPubSubProxy](context.TODO(), Youtube, *objStore, broker)
	assert.Nil(t, err)
}

func Test_VideoServiceFactory_MakeVideoStoreService_Error(t *testing.T) {
	objStore, _ := SetupFactory(t)
	_, err := MakeVideoStoreService[*mock_object_storage.MockBindingProxy, *mock_progress_broker.MockPubSubProxy](context.TODO(), Host(4), *objStore, nil)
	assert.NotNil(t, err)
}
