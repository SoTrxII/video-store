package video_store_service

import (
	"fmt"
	"time"
	"video-manager/internal/logger"
	object_storage "video-manager/internal/object-storage"
	progress_broker "video-manager/internal/progress-broker"
	"video-manager/internal/video-hosting"
)

var (
	log = logger.Build()
)

// Fired when an error occured while uploading a video
type uploadError struct {
	// Error message
	Message string `json:message`
}

// Fired while uploading a video
type uploadProgress struct {
	// Nb bytes uploaded
	Current int64 `json:"current"`
	// Total bytes to upload
	Total int64 `json:"total"`
}

// UploadVideoFromStorage Upload a video identified on the object storage by "storageKey" to the video hosting platform
func (vsc *VideoStoreService[B, P]) UploadVideoFromStorage(storageKey string, meta *video_hosting.ItemMetadata) (*video_hosting.Video, error) {

	if meta == nil {
		return nil, fmt.Errorf("no video metadata provided, aborting")
	}
	// Get the content of the file to upload and buffer it into memory
	reader, err := vsc.ObjStore.Buffer(storageKey)
	if err != nil {
		return nil, fmt.Errorf("error while downloading video from object storage : %w", err)
	}

	// Progress routine, post upload progress on the event broker if it has defined
	var onProgress video_hosting.ProgressFunc
	quit := make(chan error, 1)
	pgChannel := make(chan uploadProgress)
	if vsc.EvtBroker != nil {
		onProgress = func(current int64, total int64) {
			select {
			case pgChannel <- uploadProgress{
				Current: current,
				Total:   total,
			}:
			default:
				// pgChannel is full
			}
		}
		go vsc.startProgressRoutine(storageKey, time.Second, pgChannel, quit)
	}

	// Upload the buffered content to the video storage
	vid, err := vsc.VidHost.CreateVideo(meta, *reader, &onProgress)
	quit <- err
	if err != nil {
		return nil, fmt.Errorf("error while uploading video : %w", err)
	}

	return vid, err
}

// Periodically send progress to the event broker
// If an error is passed in errorCh, send Error, if nil is passed, send Done instead
func (vsc *VideoStoreService[B, P]) startProgressRoutine(storageKey string, every time.Duration, pgChannel chan uploadProgress, errorCh chan error) {
	ticker := time.NewTicker(every)
	for {
		select {
		// Each time the tickers tick
		case <-ticker.C:
			// Try to publish a progress event if any is available
			select {
			case pg := <-pgChannel:
				sErr := vsc.EvtBroker.SendProgress(progress_broker.UploadInfos{
					RecordId:    storageKey,
					UploadState: progress_broker.InProgress,
					Data:        pg,
				})
				if sErr != nil {
					log.Errorf("Could not send event to progress broker : %s", sErr.Error())
				}
			default:
				// No progress to post
			}
		// Else, at each iteration, check if the upload finished in any way
		case err := <-errorCh:
			ticker.Stop()
			// If an error is detected, send error, else send done
			state := progress_broker.Done
			var data interface{}
			if err != nil {
				state = progress_broker.Error
				data = uploadError{Message: err.Error()}
			}
			sErr := vsc.EvtBroker.SendProgress(progress_broker.UploadInfos{
				RecordId:    storageKey,
				UploadState: state,
				Data:        data,
			})
			if sErr != nil {
				log.Errorf("Could not send event to progress broker : %s", sErr.Error())
			}
			close(pgChannel)
			return
		}
	}
}

type VideoStoreService[B object_storage.BindingProxy, P progress_broker.PubSubProxy] struct {
	// Backend object storage
	ObjStore *object_storage.ObjectStorage[B]
	// Event broker to send notification into
	EvtBroker *progress_broker.ProgressBroker[P]
	// Video hosting platform
	VidHost video_hosting.IVideoHost
}
