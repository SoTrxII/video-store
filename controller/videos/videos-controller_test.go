package videos_controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	mock_object_storage "video-manager/internal/mock/object-storage"
	mock_progress_broker "video-manager/internal/mock/progress-broker"
	mock_video_hosting "video-manager/internal/mock/video-hosting"
	object_storage "video-manager/internal/object-storage"
	progress_broker "video-manager/internal/progress-broker"
	video_hosting "video-manager/internal/video-hosting"
	video_store_service "video-manager/pkg/video-store-service"
)

type mocked struct {
	videoStore       *mock_video_hosting.MockIVideoHost
	objectStoreProxy *mock_object_storage.MockBindingProxy
	brokerProxy      *mock_progress_broker.MockPubSubProxy
	controller       *VideoController[*mock_object_storage.MockBindingProxy, *mock_progress_broker.MockPubSubProxy]
}

var (
	sampleVid = video_hosting.Video{
		Id:           "testId",
		Title:        "testTitle",
		Description:  "testDescription",
		CreatedAt:    time.Unix(1662202180, 0).UTC(),
		Duration:     0,
		Visibility:   "unlisted",
		ThumbnailUrl: "",
	}
	sampleMetadata = video_hosting.ItemMetadata{
		Description: "testDescription",
		Title:       "testTitle",
		Visibility:  "unlisted",
	}
)

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

	vss := video_store_service.VideoStoreService[*mock_object_storage.MockBindingProxy, *mock_progress_broker.MockPubSubProxy]{
		ObjStore: objectStore,
		VidHost:  vidHost,
	}
	// The broker can either be included in the controller or set to nil
	// in which case no progress event will be sent
	if initBroker {
		vss.EvtBroker = broker
	}
	controller := VideoController[*mock_object_storage.MockBindingProxy, *mock_progress_broker.MockPubSubProxy]{Service: &vss}
	gin.SetMode(gin.TestMode)
	return &mocked{
		videoStore:       vidHost,
		objectStoreProxy: objStoreProxy,
		brokerProxy:      psProxy,
		controller:       &controller,
	}

}

func Test_VideoController_Delete_Error_Id(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: ""}}
	deps.controller.Delete(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Testing with no ID
	c.Params = []gin.Param{}
	deps.controller.Delete(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_VideoController_Delete_Error_Deletion(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	// Setting the deletion to fail
	deps.videoStore.EXPECT().DeleteVideo(gomock.Any()).Return(fmt.Errorf("test"))

	deps.controller.Delete(c)
	assert.Equal(t, 1, len(c.Errors))
}

func Test_VideoController_Delete_Error_NotFound(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	// Setting the deletion to fail
	deps.videoStore.EXPECT().DeleteVideo(gomock.Any()).Return(&video_hosting.RequestError{
		StatusCode: 404,
		Err:        fmt.Errorf("not found"),
	})

	deps.controller.Delete(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, 0, len(c.Errors))
}

func Test_VideoController_Delete_Ok(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	// Setting the deletion to pass
	deps.videoStore.EXPECT().DeleteVideo(gomock.Any()).Return(nil)

	deps.controller.Delete(c)
	assert.Equal(t, 0, len(c.Errors))
}

func Test_VideoController_Update_Error_Id(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: ""}}
	deps.controller.Update(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Testing with no ID
	c.Params = []gin.Param{}
	deps.controller.Update(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_VideoController_Update_Error_Payload_NoBody(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Update(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_VideoController_Update_Error_NotFound(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	setJsonAsBody(t, c, sampleVid)
	deps.videoStore.EXPECT().UpdateVideo(gomock.Any(), gomock.Any()).Return(nil, &video_hosting.RequestError{
		StatusCode: 404,
		Err:        fmt.Errorf("not found"),
	})

	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Update(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func Test_VideoController_Update_Error_Other(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	setJsonAsBody(t, c, sampleVid)
	deps.videoStore.EXPECT().UpdateVideo(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("test"))
	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Update(c)
	assert.Equal(t, 1, len(c.Errors))
}

func Test_VideoController_Update_Ok(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	setJsonAsBody(t, c, sampleVid)
	deps.videoStore.EXPECT().UpdateVideo(gomock.Any(), gomock.Any()).Return(&sampleVid, nil)
	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Update(c)
	assert.Equal(t, http.StatusOK, w.Code)
	retBody, err := io.ReadAll(w.Body)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	var updatedVid video_hosting.Video
	err = json.Unmarshal(retBody, &updatedVid)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	assert.Equal(t, sampleVid, updatedVid)
}

func Test_VideoController_Retrieve_Error_Id(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: ""}}
	deps.controller.Retrieve(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Testing with no ID
	c.Params = []gin.Param{}
	deps.controller.Retrieve(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_VideoController_Retrieve_Error_NotFound(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	setJsonAsBody(t, c, sampleVid)
	deps.videoStore.EXPECT().RetrieveVideo(gomock.Any()).Return(nil, &video_hosting.RequestError{
		StatusCode: 404,
		Err:        fmt.Errorf("not found"),
	})
	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Retrieve(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func Test_VideoController_Retrieve_Error_Other(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	setJsonAsBody(t, c, sampleVid)
	deps.videoStore.EXPECT().RetrieveVideo(gomock.Any()).Return(nil, fmt.Errorf("test"))
	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Retrieve(c)
	assert.Equal(t, 1, len(c.Errors))
}

func Test_VideoController_Retrieve_Ok(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	setJsonAsBody(t, c, sampleVid)
	deps.videoStore.EXPECT().RetrieveVideo(gomock.Any()).Return(&sampleVid, nil)
	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Retrieve(c)
	assert.Equal(t, http.StatusOK, w.Code)
	retBody, err := io.ReadAll(w.Body)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	var updatedVid video_hosting.Video
	err = json.Unmarshal(retBody, &updatedVid)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	assert.Equal(t, sampleVid, updatedVid)
}

func jobId(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// No storage key
	body := CreateVideoBody{
		ItemMetadata: video_hosting.ItemMetadata{
			Description: "",
			Title:       "",
			Visibility:  "",
		},
	}
	setJsonAsBody(t, c, body)
	deps.controller.Create(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Empty storage key
	body = CreateVideoBody{
		ItemMetadata: video_hosting.ItemMetadata{
			Description: "",
			Title:       "",
			Visibility:  "",
		},
		StorageKey: "",
	}
	setJsonAsBody(t, c, body)
	deps.controller.Create(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_VideoController_Create_Error_Payload_NoBody(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Testing with an empty body
	deps.controller.Create(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_VideoController_Create_Error_TitleValidation(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Empty storage key
	body := CreateVideoBody{
		ItemMetadata: video_hosting.ItemMetadata{
			Description: "blah",
			Title:       strings.Repeat("test", 30),
			Visibility:  "hh",
		},
		StorageKey: "ddd",
	}
	setJsonAsBody(t, c, body)
	deps.controller.Create(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_VideoController_Create_Ok(t *testing.T) {
	deps := Setup(t, false)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Mocking call to the object storage
	deps.
		objectStoreProxy.
		EXPECT().
		InvokeBinding(gomock.Any(), gomock.Any()).Return(&client.BindingEvent{Data: []byte("aa")}, nil)
	// Mocking call to video host
	deps.
		videoStore.
		EXPECT().
		CreateVideo(gomock.Any(), gomock.Any(), gomock.Any()).Return(&sampleVid, nil)

	body := CreateVideoBody{
		ItemMetadata: video_hosting.ItemMetadata{
			Description: "test",
			Title:       "test",
			Visibility:  "test",
		},
		StorageKey: "test",
		JobId:      "test",
	}
	setJsonAsBody(t, c, body)
	// Testing with an empty body
	deps.controller.Create(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_VideoController_Create_Ok_WithProgress(t *testing.T) {
	deps := Setup(t, true)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Mocking any call to the broker
	deps.brokerProxy.
		EXPECT().
		PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	// Mocking call to the object storage
	deps.
		objectStoreProxy.
		EXPECT().
		InvokeBinding(gomock.Any(), gomock.Any()).Return(&client.BindingEvent{Data: []byte("aa")}, nil)
	// Mocking call to video host
	deps.
		videoStore.
		EXPECT().
		CreateVideo(gomock.Any(), gomock.Any(), gomock.Any()).Return(&sampleVid, nil)

	body := CreateVideoBody{
		ItemMetadata: video_hosting.ItemMetadata{
			Description: "test",
			Title:       "test",
			Visibility:  "test",
		},
		StorageKey: "test",
		JobId:      "test",
	}
	setJsonAsBody(t, c, body)
	// Testing with an empty body
	deps.controller.Create(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVideoController_SetThumbnail_FromStorageKey_Ok(t *testing.T) {
	deps := Setup(t, true)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	deps.
		objectStoreProxy.
		EXPECT().
		InvokeBinding(gomock.Any(), gomock.Any()).Return(&client.BindingEvent{Data: []byte("aa")}, nil)

	deps.videoStore.EXPECT().UpdateVideoThumbnail(gomock.Any(), gomock.Any()).Return(nil)

	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}, {Key: "tId", Value: "meh"}}
	deps.controller.SetThumbnail(c)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestVideoController_SetThumbnail_FromStorageKey_ThumbNotExists(t *testing.T) {
	deps := Setup(t, true)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	deps.
		objectStoreProxy.
		EXPECT().
		InvokeBinding(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("test"))

	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}, {Key: "tId", Value: "meh"}}
	deps.controller.SetThumbnail(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVideoController_SetThumbnail_FromStorageKey_YoutubeError(t *testing.T) {
	deps := Setup(t, true)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	deps.
		objectStoreProxy.
		EXPECT().
		InvokeBinding(gomock.Any(), gomock.Any()).Return(&client.BindingEvent{Data: []byte("aa")}, nil)

	deps.videoStore.EXPECT().UpdateVideoThumbnail(gomock.Any(), gomock.Any()).Return(fmt.Errorf("test"))
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}, {Key: "tId", Value: "meh"}}
	deps.controller.SetThumbnail(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVideoController_SetThumbnail_MissingId(t *testing.T) {
	deps := Setup(t, true)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	deps.controller.SetThumbnail(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVideoController_SetThumbnail_FromBuffer_Ok(t *testing.T) {
	deps := Setup(t, true)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer([]byte("test")))
	deps.videoStore.EXPECT().UpdateVideoThumbnail(gomock.Any(), gomock.Any()).Return(nil)
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}, {Key: "tId", Value: ""}}
	deps.controller.SetThumbnail(c)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

// Set the payload as the JSON body of c
func setJsonAsBody(t *testing.T, c *gin.Context, payload any) {
	buf, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(buf))
	c.Request.Header.Set("Content-Type", "application/json")
}
