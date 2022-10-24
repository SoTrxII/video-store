package video_store_service

import (
	"context"
	"fmt"
	"os"
	object_storage "video-manager/internal/object-storage"
	progress_broker "video-manager/internal/progress-broker"
	video_hosting "video-manager/internal/video-hosting"
)

// List of all available video hosting platforms
type Host uint

const (
	Youtube Host = iota
)

// Return an instance of a video storage servcie configured with the provided video host as the backend
func MakeVideoStoreService[T object_storage.BindingProxy, P progress_broker.PubSubProxy](ctx context.Context, host Host, proxy object_storage.ObjectStorage[T], progressBroker *progress_broker.ProgressBroker[P]) (*VideoStoreService[T, P], error) {
	var store video_hosting.IVideoHost
	var err error
	switch host {
	case Youtube:
		store, err = makeYoutubeStoreService(ctx)
	default:
		// This can't actually happen
		err = fmt.Errorf(`the provided host "%v" has no available implementation`, host)
	}
	if err != nil {
		return nil, err
	}

	return &VideoStoreService[T, P]{
		EvtBroker: progressBroker,
		ObjStore:  &proxy,
		VidHost:   store,
		opt:       VideoStoreOptions{objStoreMaxRetry: 10},
	}, nil

}

// Returns an instance of a youtube store
func makeYoutubeStoreService(ctx context.Context) (video_hosting.IVideoHost, error) {
	return video_hosting.NewYoutubeStore(ctx, &video_hosting.YoutubeStoreCredentials{
		ClientId:     os.Getenv("YT_CLIENT_ID"),
		ClientSecret: os.Getenv("YT_CLIENT_SECRET"),
		RefreshToken: os.Getenv("YT_REFRESH_TOKEN"),
	}, nil)
}
