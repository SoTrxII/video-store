package videos_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	object_storage "video-manager/internal/object-storage"
	progress_broker "video-manager/internal/progress-broker"
	video_hosting "video-manager/internal/video-hosting"
	video_store_service "video-manager/pkg/video-store-service"
)

type VideoController[B object_storage.BindingProxy, P progress_broker.PubSubProxy] struct {
	Service *video_store_service.VideoStoreService[B, P]
}

// POST body required to create a new video on the hosting platform
// from the backend object storage
type CreateVideoBody struct {
	// Required metadata to upload a video
	video_hosting.ItemMetadata
	// Key to retrieve the video from the object storage
	StorageKey string `json:"storageKey" binding:"required"`
}

// ShowAccount godoc
// @Summary      Upload a video
// @Description  Upload a video from the object storage to the video hosting platform
// @Tags         videos
// @Accept       json
// @Produce      json
// @Param 		 videometa body CreateVideoBody true "Required data to upload a video"
// @Success      200  {object}  video_hosting.Video
// @Failure      400
// @Failure      404  {string}  string "No video with this ID"
// @Failure      500
// @Router       /videos [post]
func (vc *VideoController[S, P]) Create(c *gin.Context) {
	var target CreateVideoBody
	if err := c.BindJSON(&target); err != nil {
		c.String(http.StatusBadRequest, `invalid body provided: %s !`, err.Error())
		return
	}
	if target.StorageKey == "" {
		c.String(http.StatusBadRequest, `No storage key provided, aborting !`)
		return
	}
	// TODO : Progress
	vid, err := vc.Service.UploadVideoFromStorage(target.StorageKey, &video_hosting.ItemMetadata{
		Description: target.Description,
		Title:       target.Title,
		Visibility:  target.Visibility,
	})
	if err != nil {
		if re, ok := err.(*video_hosting.RequestError); ok {
			c.String(re.StatusCode, re.Error())
		} else {
			c.Status(http.StatusInternalServerError)
			_ = c.Error(err)
		}
		return
	}
	c.SecureJSON(http.StatusOK, vid)
}

// ShowAccount godoc
// @Summary      Get a video
// @Description  Retrieve a video by ID
// @Tags         videos
// @Produce      json
// @Param        id   path      int  true  "Video ID"
// @Success      200  {object}  video_hosting.Video
// @Failure      400
// @Failure      404  {string}  string "No video with this ID"
// @Failure      500
// @Router       /videos/{id} [get]
func (vc *VideoController[S, P]) Retrieve(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, `No id provided !`)
		return
	}
	vid, err := vc.Service.VidHost.RetrieveVideo(id)
	if err != nil {
		if re, ok := err.(*video_hosting.RequestError); ok {
			c.String(re.StatusCode, re.Error())
		} else {
			c.Status(http.StatusInternalServerError)
			_ = c.Error(err)
		}
		return
	}
	c.SecureJSON(http.StatusOK, vid)
}

// ShowAccount godoc
// @Summary      Update a video
// @Description  Update the video by ID if it exists
// @Tags         videos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Video ID"
// @Param 		 video body video_hosting.Video true "Updated video"
// @Success      200 {object}  video_hosting.Video
// @Failure      400
// @Failure      404  {string}  string "No video with this ID"
// @Failure      500
// @Router       /videos/{id} [put]
func (vc *VideoController[S, P]) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, `No id provided !`)
		return
	}
	var target video_hosting.Video
	if err := c.BindJSON(&target); err != nil {
		c.String(http.StatusBadRequest, `invalid body provided: %s !`, err.Error())
		return
	}
	vid, err := vc.Service.VidHost.UpdateVideo(id, &target)
	if err != nil {
		if re, ok := err.(*video_hosting.RequestError); ok {
			c.String(re.StatusCode, re.Error())
		} else {
			c.Status(http.StatusInternalServerError)
			_ = c.Error(err)
		}
		return
	}
	c.SecureJSON(http.StatusOK, vid)
}

// ShowAccount godoc
// @Summary      Delete a video
// @Description  Delete the video by ID if it exists
// @Tags         videos
// @Produce      json
// @Param        id   path      int  true  "Video ID"
// @Success      204
// @Failure      400
// @Failure      404  {string}  string "No video with this ID"
// @Failure      500
// @Router       /videos/{id} [delete]
func (vc *VideoController[S, P]) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, `No id provided !`)
		return
	}
	err := vc.Service.VidHost.DeleteVideo(id)
	if err != nil {
		if re, ok := err.(*video_hosting.RequestError); ok {
			c.String(re.StatusCode, re.Error())
		} else {
			c.Status(http.StatusInternalServerError)
			_ = c.Error(err)
		}
		return
	}
	c.String(http.StatusNoContent, "")
}

// ShowAccount godoc
// @Summary      Set the thumbnail of a video
// @Description  Set the thumbnail of an existing video on the remote video hosting platform
// @Tags         videos
// @Accept       octet-stream
// @Param        key   path      int  true  "Video ID"
// @Param 		 thumb body octet-stream true "Thumbnail content"
// @Success      204
// @Failure      500
// @Router       /videos/{id}/thumbnail [post]
func (vc *VideoController[S, P]) SetThumbnail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, `No id provided !`)
		return
	}
	err := vc.Service.VidHost.UpdateVideoThumbnail(id, c.Request.Body)
	if err != nil {
		if re, ok := err.(*video_hosting.RequestError); ok {
			c.String(re.StatusCode, re.Error())
		} else {
			c.Status(http.StatusInternalServerError)
			_ = c.Error(err)
		}
		return
	}
	c.String(http.StatusNoContent, "")
}
