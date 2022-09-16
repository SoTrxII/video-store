package playlists_controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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
	controller       *PlaylistController[*mock_object_storage.MockBindingProxy, *mock_progress_broker.MockPubSubProxy]
}

var (
	samplePlaylist = video_hosting.Playlist{
		Id:           "testId",
		Title:        "testTitle",
		Description:  "testDescription",
		CreatedAt:    time.Unix(1662202180, 0).UTC(),
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
	controller := PlaylistController[*mock_object_storage.MockBindingProxy, *mock_progress_broker.MockPubSubProxy]{Service: &vss}
	gin.SetMode(gin.TestMode)
	return &mocked{
		videoStore:       vidHost,
		objectStoreProxy: proxy,
		controller:       &controller,
	}
}

func Test_PlaylistController_Delete_Error_Id(t *testing.T) {
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

func Test_PlaylistController_Delete_Error_Deletion(t *testing.T) {
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	// Setting the deletion to fail
	deps.videoStore.EXPECT().DeletePlaylist(gomock.Any()).Return(fmt.Errorf("test"))

	deps.controller.Delete(c)
	assert.Equal(t, 1, len(c.Errors))
}

func Test_PlaylistController_Delete_Error_NotFound(t *testing.T) {
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	// Setting the deletion to fail
	deps.videoStore.EXPECT().DeletePlaylist(gomock.Any()).Return(&video_hosting.RequestError{
		StatusCode: 404,
		Err:        fmt.Errorf("not found"),
	})

	deps.controller.Delete(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, 0, len(c.Errors))
}

func Test_PlaylistController_Delete_Ok(t *testing.T) {
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	// Setting the deletion to pass
	deps.videoStore.EXPECT().DeletePlaylist(gomock.Any()).Return(nil)

	deps.controller.Delete(c)
	assert.Equal(t, 0, len(c.Errors))
}

func Test_PlaylistController_Update_Error_Id(t *testing.T) {
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

func Test_PlaylistController_Update_Error_Payload_NoBody(t *testing.T) {
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Update(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_PlaylistController_Update_Error_NotFound(t *testing.T) {
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	setJsonAsBody(t, c, samplePlaylist)
	deps.videoStore.EXPECT().UpdatePlaylist(gomock.Any(), gomock.Any()).Return(nil, &video_hosting.RequestError{
		StatusCode: 404,
		Err:        fmt.Errorf("not found"),
	})

	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Update(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func Test_PlaylistController_Update_Error_Other(t *testing.T) {
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	setJsonAsBody(t, c, samplePlaylist)
	deps.videoStore.EXPECT().UpdatePlaylist(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("test"))
	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Update(c)
	assert.Equal(t, 1, len(c.Errors))
}

func Test_PlaylistController_Update_Ok(t *testing.T) {
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	setJsonAsBody(t, c, samplePlaylist)
	deps.videoStore.EXPECT().UpdatePlaylist(gomock.Any(), gomock.Any()).Return(&samplePlaylist, nil)
	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Update(c)
	assert.Equal(t, http.StatusOK, w.Code)
	retBody, err := io.ReadAll(w.Body)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	var updatedPlaylist video_hosting.Playlist
	err = json.Unmarshal(retBody, &updatedPlaylist)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	assert.Equal(t, samplePlaylist, updatedPlaylist)
}

func Test_PlaylistController_Retrieve_Error_Id(t *testing.T) {
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

func Test_PlaylistController_Retrieve_Error_NotFound(t *testing.T) {
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	setJsonAsBody(t, c, samplePlaylist)
	deps.videoStore.EXPECT().RetrievePlaylist(gomock.Any()).Return(nil, &video_hosting.RequestError{
		StatusCode: 404,
		Err:        fmt.Errorf("not found"),
	})
	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Retrieve(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func Test_PlaylistController_Retrieve_Error_Other(t *testing.T) {
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	setJsonAsBody(t, c, samplePlaylist)
	deps.videoStore.EXPECT().RetrievePlaylist(gomock.Any()).Return(nil, fmt.Errorf("test"))
	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Retrieve(c)
	assert.Equal(t, 1, len(c.Errors))
}

func Test_PlaylistController_Retrieve_Ok(t *testing.T) {
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	setJsonAsBody(t, c, samplePlaylist)
	deps.videoStore.EXPECT().RetrievePlaylist(gomock.Any()).Return(&samplePlaylist, nil)
	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}
	deps.controller.Retrieve(c)
	assert.Equal(t, http.StatusOK, w.Code)
	retBody, err := io.ReadAll(w.Body)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	var updatedPlaylist video_hosting.Playlist
	err = json.Unmarshal(retBody, &updatedPlaylist)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	assert.Equal(t, samplePlaylist, updatedPlaylist)
}

func Test_PlaylistController_Create_Error_Key(t *testing.T) {
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Testing with an empty ID
	c.Params = []gin.Param{gin.Param{Key: "key", Value: ""}}
	deps.controller.Create(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Testing with no ID
	c.Params = []gin.Param{}
	deps.controller.Create(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_PlaylistController_Create_Error_Payload_NoBody(t *testing.T) {
	deps := Setup(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Testing with an empty body
	c.Params = []gin.Param{gin.Param{Key: "key", Value: "1"}}
	deps.controller.Create(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
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
