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
	video_hosting "video-manager/internal/video-hosting"
	video_store_service "video-manager/pkg/video-store-service"
)

type mocked struct {
	videoStore       *mock_video_hosting.MockIVideoHost
	objectStoreProxy *mock_object_storage.MockBindingProxy
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
	vss := video_store_service.VideoStoreService[*mock_object_storage.MockBindingProxy, *mock_progress_broker.MockPubSubProxy]{
		ObjStore: objectStore,
		VidHost:  vidHost,
	}
	controller := VideoController[*mock_object_storage.MockBindingProxy, *mock_progress_broker.MockPubSubProxy]{Service: &vss}
	gin.SetMode(gin.TestMode)
	return &mocked{
		videoStore:       vidHost,
		objectStoreProxy: proxy,
		controller:       &controller,
	}
}

func Test_VideoController_Delete_Error_Id(t *testing.T) {
	deps := Setup(t)
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
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	// Setting the deletion to fail
	deps.videoStore.EXPECT().DeleteVideo(gomock.Any()).Return(fmt.Errorf("test"))

	deps.controller.Delete(c)
	assert.Equal(t, 1, len(c.Errors))
}

func Test_VideoController_Delete_Error_NotFound(t *testing.T) {
	deps := Setup(t)
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
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	// Setting the deletion to pass
	deps.videoStore.EXPECT().DeleteVideo(gomock.Any()).Return(nil)

	deps.controller.Delete(c)
	assert.Equal(t, 0, len(c.Errors))
}

func Test_VideoController_Update_Error_Id(t *testing.T) {
	deps := Setup(t)
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
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Update(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_VideoController_Update_Error_NotFound(t *testing.T) {
	deps := Setup(t)
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
	deps := Setup(t)
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
	deps := Setup(t)
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
	deps := Setup(t)
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
	deps := Setup(t)
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
	deps := Setup(t)
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
	deps := Setup(t)
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

func Test_VideoController_Create_Error_Key(t *testing.T) {
	deps := Setup(t)
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
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Testing with an empty body
	deps.controller.Create(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_VideoController_Create_Error_TitleValidation(t *testing.T) {
	deps := Setup(t)
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
	deps := Setup(t)
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
	}
	setJsonAsBody(t, c, body)
	// Testing with an empty body
	deps.controller.Create(c)
	assert.Equal(t, http.StatusOK, w.Code)
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
