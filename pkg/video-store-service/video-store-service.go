package video_store_service

import (
	"fmt"
	"io"
	"math"
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

// Fired while uploading a video
type uploadDone struct {
	// Video Id
	Id string `json:"id"`
	// URL prefix to watch videos on the url
	WatchPrefix string `json:"watchPrefix"`
	// Video duration
	Duration int64 `json:"duration"`
}

// Fired while uploading a video
type uploadResult struct {
	Result *video_hosting.Video
	Error  error
}

type VideoStoreOptions struct {
	// Number of time to retry calls made to the object store.
	// Each call will be followed by a wait time of (2^attempt)s
	objStoreMaxRetry int8
}

// UploadVideoFromStorage Upload a video identified on the object storage by "storageKey" to the video hosting platform
func (vsc *VideoStoreService[B, P]) UploadVideoFromStorage(jobId string, storageKey string, meta *video_hosting.ItemMetadata) (*video_hosting.Video, error) {

	if meta == nil {
		return nil, fmt.Errorf("no video metadata provided, aborting")
	}
	// Get the content of the file to upload and buffer it into memory
	// So there may be a race condition here.
	// As far as I understand, object uploaded on a storage aren't available immediately after upload, there is a slight
	// delay that might be caused by the configured B64 decoding. Still, as the file gets bigger, this delay gets longer.
	// So we actually can't trust the Buffer to work the first time around.
	var reader *io.Reader
	var err error
	// Using "<=", we make sure the loop in entered at least once, event if max retry is 0
	for attempts := int8(0); attempts <= vsc.opt.objStoreMaxRetry; attempts++ {
		reader, err = vsc.ObjStore.Buffer(storageKey)
		if err == nil {
			break
		}
		log.Warnf("error in attempt %d at dowloading the video from the object storage %s", attempts, err.Error())
		// We will use a linear backoff strategy, we don't have any collision whatsoever, we just want to wait until the video is available
		// The sum 2^n from 0 to 10 = 2047 ~= 30min  of total wait, this is way more than enough, as more will be over an
		// http session time. Plus, if the waiting time is really because of the b64 decoding, it's a 0(n) time complexity algorithm
		delaySecs := int64(math.Pow(2, float64(attempts)))
		time.Sleep(time.Duration(delaySecs) * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("error while downloading video from object storage : %w", err)
	}

	// Progress routine, post upload progress on the event broker if it has defined
	var onProgress video_hosting.ProgressFunc
	quit := make(chan uploadResult, 1)
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
		go vsc.startProgressRoutine(jobId, time.Second, pgChannel, quit)
	}

	// Upload the buffered content to the video storage
	vid, err := vsc.VidHost.CreateVideo(meta, *reader, &onProgress)

	// Wait for the event broker goroutine
	if vsc.EvtBroker != nil {
		// Send the error/nil to the buffered error channel
		quit <- uploadResult{
			Result: vid,
			Error:  err,
		}
		// And wait for the progress channel to be closed by the goroutine
		<-pgChannel
	} else {
		close(quit)
		close(pgChannel)
	}

	if err != nil {
		return nil, fmt.Errorf("error while uploading video : %w", err)
	}

	return vid, err
}

// Periodically send progress to the event broker
// If an error is passed in errorCh, send Error, if nil is passed, send Done instead
func (vsc *VideoStoreService[B, P]) startProgressRoutine(jobId string, every time.Duration, pgChannel chan uploadProgress, resCh chan uploadResult) {
	ticker := time.NewTicker(every)
	for {
		select {
		// Each time the tickers tick
		case <-ticker.C:
			// Try to publish a progress event if any is available
			select {
			case pg := <-pgChannel:
				sErr := vsc.EvtBroker.SendProgress(progress_broker.UploadInfos{
					JobId: jobId,
					State: progress_broker.InProgress,
					Data:  pg,
				})
				if sErr != nil {
					log.Errorf("Could not send event to progress broker : %s", sErr.Error())
				}
			default:
				// No progress to post
			}
		// Else, at each iteration, check if the upload finished in any way
		case res := <-resCh:
			ticker.Stop()
			// If an error is detected, send error, else send done
			state := progress_broker.Done
			var data interface{}
			if res.Error != nil {
				state = progress_broker.Error
				data = uploadError{Message: res.Error.Error()}
			} else {
				data = uploadDone{
					Id:          res.Result.Id,
					WatchPrefix: res.Result.WatchPrefix,
					Duration:    res.Result.Duration,
				}
			}
			sErr := vsc.EvtBroker.SendProgress(progress_broker.UploadInfos{
				JobId: jobId,
				State: state,
				Data:  data,
			})
			if sErr != nil {
				log.Errorf("Could not send event to progress broker : %s", sErr.Error())
			}
			close(pgChannel)
			return
		}
	}
}

func (vsc *VideoStoreService[B, P]) SetVideoThumbnailFromStorage(vidId, thumbStorageKey string) error {
	reader, err := vsc.ObjStore.Buffer(thumbStorageKey)
	if err != nil {
		return fmt.Errorf("error while downloading thumbnail from object storage : %w", err)
	}

	return vsc.VidHost.UpdateVideoThumbnail(vidId, *reader)
}

type VideoStoreService[B object_storage.BindingProxy, P progress_broker.PubSubProxy] struct {
	// Backend object storage
	ObjStore *object_storage.ObjectStorage[B]
	// Event broker to send notification into
	EvtBroker *progress_broker.ProgressBroker[P]
	// Video hosting platform
	VidHost video_hosting.IVideoHost
	// Customize behaviour of the service
	// Not using a pointer will initialize a struct will default values
	opt VideoStoreOptions
}
