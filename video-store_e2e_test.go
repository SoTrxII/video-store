//go:build end2end
// +build end2end

/*
*
Attention passengers, all aboard this shipwreck !
No mock, no stub, not a single fake allowed, this is the real deal.
Please press F for your video-hosting quota, we're going to go with some real scenario !
*/
package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	videos_controller "video-manager/controller/videos"
	progress_broker "video-manager/internal/progress-broker"
	video_hosting "video-manager/internal/video-hosting"
)

const (
	// Resources directory relative path
	ResDir = "./resources/test/"
	// Key to use when upload to the object storage
	SampleObjectStoreKey = "key"
	// Main application url
	AppBasePath = "http://localhost:8080/v1/"
)

var (
	// Sample metadata for a playlist
	samplePlaylistMeta = video_hosting.ItemMetadata{
		Description: "testDesc",
		Title:       "testTitle",
		Visibility:  "unlisted",
	}
	// Sample metadata for a video
	sampleVideoMeta = videos_controller.CreateVideoBody{
		ItemMetadata: video_hosting.ItemMetadata{
			Description: "this is a video",
			Title:       "look a this",
			Visibility:  "unlisted",
		},
		StorageKey: SampleObjectStoreKey,
	}
	// Contains all event received so far
	eventStack []progress_broker.UploadInfos
)

func Setup(t *testing.T) {
	// Start a standalone dapr client, as we won't use the app code for these tests
	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("Dapr not started !")
	}
	// Upload an asset to the backend object storage, to be later use as the video
	err = copy(&c, filepath.Join(ResDir, "video.mp4"), SampleObjectStoreKey)
	if err != nil {
		// Obj store not configured properly
		t.Fatalf(err.Error())
	}

	// make a listener for upload event in the broker
	s := daprd.NewService(":8081")
	var eventSub = &common.Subscription{
		PubsubName: "pubsub",
		Topic:      DefaultPubSubTopic,
		Route:      "/",
	}
	if err := s.AddTopicEventHandler(eventSub, onBrokerEvent); err != nil {
		t.Fatalf("error adding topic subscription: %v", err)
	}

	go func() {
		err := s.Start()
		if err != nil {
			log.Fatalf("error starting subscription server")
		}
	}()

	// start the app
	go main()
}

/*
Real world scenario. As we are fancy, we're also going to use a custom thumbnail.
So :
  - Create a playlist p
  - Create a video v
  - Set thumbnail of v to an image i
  - Insert v in p
  - Update v
  - Delete v
  - Delete p
*/
func Test_VideoLifecycle(t *testing.T) {
	Setup(t)

	// Let's create a playlist as an external user would
	res := sendRequest(t, http.MethodPost, getAppEndpoint(t, "playlists"), samplePlaylistMeta, false)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var resPlaylist video_hosting.Playlist
	parseBodyInto(t, res.Body, &resPlaylist)
	assert.Equal(t, samplePlaylistMeta.Title, resPlaylist.Title)
	assert.Equal(t, samplePlaylistMeta.Description, resPlaylist.Description)

	// Then upload a video
	res = sendRequest(t, http.MethodPost, getAppEndpoint(t, "videos"), sampleVideoMeta, false)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var resVideo video_hosting.Video
	parseBodyInto(t, res.Body, &resVideo)
	assert.Equal(t, sampleVideoMeta.Title, resVideo.Title)
	assert.Equal(t, sampleVideoMeta.Description, resVideo.Description)
	// The progress broker should have fired at least one event
	assert.GreaterOrEqual(t, 1, len(eventStack))
	// And one of them should be a "done" event
	found := false
	for _, evt := range eventStack {
		if evt.State == progress_broker.Done {
			found = true
		}
	}
	assert.True(t, found)

	// Change the thumbnail of the uploaded video
	thumbEd := getAppEndpoint(t, fmt.Sprintf("/videos/%s/thumbnail", resVideo.Id))
	//thumbEd := getAppEndpoint(t, fmt.Sprintf("/videos/%s/thumbnail", "0yRc7Bb0MPg"))
	res = sendRequest(t, http.MethodPost, thumbEd, filepath.Join(ResDir, "test.jpg"), true)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	// Insert the created video into the created playlist
	addVidToPEd := getAppEndpoint(t, fmt.Sprintf("/playlists/%s/videos/%s", resPlaylist.Id, resVideo.Id))
	//addVidToPEd := getAppEndpoint(t, fmt.Sprintf("/playlists/%s/videos/%s", "PLjm3adukVo4Qrg32BQA1ZWj5n2FDG4NaA", "0yRc7Bb0MPg"))
	res = sendRequest(t, http.MethodPut, addVidToPEd, nil, false)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	// Update v : Change title
	resVideoCpy := resVideo
	resVideoCpy.Title += "mod"
	resVideoCpy.Description += "mod"
	updateVidEd := getAppEndpoint(t, fmt.Sprintf("/videos/%s", resVideo.Id))
	res = sendRequest(t, http.MethodPut, updateVidEd, resVideoCpy, false)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var resVideoCpyRet video_hosting.Video
	parseBodyInto(t, res.Body, &resVideoCpyRet)
	assert.Equal(t, resVideoCpy.Title, resVideoCpyRet.Title)
	assert.Equal(t, resVideoCpy.Description, resVideoCpyRet.Description)

	// Delete v
	deleteVidEd := getAppEndpoint(t, fmt.Sprintf("/videos/%s", resVideo.Id))
	res = sendRequest(t, http.MethodDelete, deleteVidEd, nil, false)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	// Delete p
	deletePlaylistEd := getAppEndpoint(t, fmt.Sprintf("/playlists/%s", resPlaylist.Id))
	res = sendRequest(t, http.MethodDelete, deletePlaylistEd, nil, false)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)

}

// Parse a response body into the given interface
// Make the test fail if it cannot be done
func parseBodyInto[T interface{}](t *testing.T, body io.ReadCloser, dest *T) {
	resBody, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = json.Unmarshal(resBody, dest)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

// Send a http request with the given payload.
// Payload can either be any interface that will be converted to json or if "isFile" is true, a filepath
// in which case the file content will be used as the request payload
func sendRequest(t *testing.T, method string, endpoint string, payload interface{}, isFile bool) *http.Response {
	var mPayload []byte
	var err error
	mPayloadBuff := bytes.NewBuffer([]byte{})
	contentType := ""
	if payload != nil {
		// If the payload is a file
		if isFile {
			// Read it
			mPayload, err = os.ReadFile(fmt.Sprintf("%v", payload))
			contentType = "octet/stream"

		} else {
			// else Serialize the payload into a buffer
			mPayload, err = json.Marshal(payload)
			contentType = "application/json"
		}

		if err != nil {
			t.Fatalf(err.Error())
		}
		mPayloadBuff = bytes.NewBuffer(mPayload)
	}

	// Prepare the request
	req, err := http.NewRequest(method, endpoint, mPayloadBuff)
	if contentType != "" {
		req.Header.Set("content-type", contentType)
	}
	if err != nil {
		t.Fatalf(err.Error())
	}
	// And sent it
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf(err.Error())
	}
	return res
}

// Return a properly formed main app url
// Makes the test fail if the endpoint can be formed
func getAppEndpoint(t *testing.T, uri string) string {
	endpoint, err := url.JoinPath(AppBasePath, uri)
	if err != nil {
		t.Fatalf(err.Error())
	}
	return endpoint
}

// Executed each time a new event is received from the progress_broker
func onBrokerEvent(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	fmt.Printf("event - PubsubName: %s, Topic: %s, ID: %s, Data: %s", e.PubsubName, e.Topic, e.ID, e.Data)
	// Parse the progress event
	var evt progress_broker.UploadInfos
	/*type tmp struct {
		RecordId string `json:"recordId"`
		//State progress_broker.State `json:"uploadState"`
	}*/
	data, mErr := json.Marshal(e.Data)
	if mErr != nil {
		fmt.Errorf("could not parse event : %s", mErr.Error())
	}
	//var evtTmp tmp
	marshErr := json.Unmarshal(data, &evt)
	if marshErr != nil {
		fmt.Errorf("could not parse event : %s", marshErr.Error())
	}
	fmt.Printf("%+v\n", evt)
	eventStack = append(eventStack, evt)
	return false, nil
}

// Copy a file into the remove object storage
func copy(daprClient *client.Client, src string, keyName string) error {
	rawContent, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	b64Content := make([]byte, base64.StdEncoding.EncodedLen(len(rawContent)))
	base64.StdEncoding.Encode(b64Content, rawContent)
	_, err = (*daprClient).InvokeBinding(context.Background(), &client.InvokeBindingRequest{
		Name:      "object-store",
		Operation: "create",
		Data:      b64Content,
		Metadata: map[string]string{
			"key": keyName,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
