package video_store_service

import (
	"fmt"
	object_storage "video-manager/internal/object-storage"
	"video-manager/internal/video-hosting"
)

// UploadVideoFromStorage Upload a video identified on the object storage by "storageKey" to the video hosting platform
func (vsc *VideoStoreService[B]) UploadVideoFromStorage(storageKey string, meta *video_hosting.ItemMetadata) (*video_hosting.Video, error) {

	if meta == nil {
		return nil, fmt.Errorf("no video metadata provided, aborting")
	}
	// Get the content of the file to upload and buffer it into memory
	reader, err := vsc.ObjStore.Buffer(storageKey)
	if err != nil {
		return nil, fmt.Errorf("error while downloading video from object storage : %w", err)
	}

	// Upload the buffered content to the video storage
	vid, err := vsc.VidHost.CreateVideo(meta, *reader, nil)
	if err != nil {
		return nil, fmt.Errorf("error while uploading video : %w", err)
	}

	return vid, err
}

type VideoStoreService[B object_storage.BindingProxy] struct {
	// Backend object storage
	ObjStore *object_storage.ObjectStorage[B]
	// Video hosting platform
	VidHost video_hosting.IVideoHost
}
