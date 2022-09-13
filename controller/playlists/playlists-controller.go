package playlists_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	object_storage "video-manager/internal/object-storage"
	video_hosting "video-manager/internal/video-hosting"
	video_store_service "video-manager/pkg/video-store-service"
)

type PlaylistController[B object_storage.BindingProxy] struct {
	Service *video_store_service.VideoStoreService[B]
}

// ShowAccount godoc
// @Summary      Creates a new playlist
// @Description  Creates a new playlist on the remote video hosting platform
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param 		 meta body video_hosting.ItemMetadata true "Required data to create a playlist"
// @Success      200  {object}  video_hosting.Playlist
// @Failure      400	"Required metata are wrong in some ways"
// @Failure      500
// @Router       /playlists [post]
func (vc *PlaylistController[S]) Create(c *gin.Context) {
	var target video_hosting.ItemMetadata
	if err := c.BindJSON(&target); err != nil {
		c.String(http.StatusBadRequest, `invalid body provided: %s !`, err.Error())
		return
	}
	vid, err := vc.Service.VidHost.CreatePlaylist(&target)
	if err != nil {
		if re, ok := err.(*video_hosting.RequestError); ok {
			c.String(re.StatusCode, re.Error())
		} else {
			_ = c.Error(err)
		}
		return
	}
	c.SecureJSON(http.StatusOK, vid)
}

// ShowAccount godoc
// @Summary      Get a playlist
// @Description  Retrieve a playlist by ID
// @Tags         playlists
// @Produce      json
// @Param        id   path      int  true  "Playlist ID"
// @Success      200  {object}  video_hosting.Playlist
// @Failure      400
// @Failure      404  {string}  string "No playlist with this ID"
// @Failure      500
// @Router       /playlists/{id} [get]
func (vc *PlaylistController[S]) Retrieve(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, `No id provided !`)
		return
	}
	playlist, err := vc.Service.VidHost.RetrievePlaylist(id)
	if err != nil {
		if re, ok := err.(*video_hosting.RequestError); ok {
			c.String(re.StatusCode, re.Error())
		} else {
			_ = c.Error(err)
		}
		return
	}
	c.SecureJSON(http.StatusOK, playlist)
}

// ShowAccount godoc
// @Summary      Update a playlist
// @Description  Update a playlist by ID if it exists
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Playlist ID"
// @Param 		 playlist body video_hosting.Playlist true "Updated playlist"
// @Success      200 {object}  video_hosting.Playlist
// @Failure      400
// @Failure      404  {string}  string "No playlist with this ID"
// @Failure      500
// @Router       /playlists/{id} [put]
func (vc *PlaylistController[S]) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, `No id provided !`)
		return
	}
	var target video_hosting.Playlist
	if err := c.BindJSON(&target); err != nil {
		c.String(http.StatusBadRequest, `invalid body provided: %s !`, err.Error())
		return
	}
	vid, err := vc.Service.VidHost.UpdatePlaylist(id, &target)
	if err != nil {
		if re, ok := err.(*video_hosting.RequestError); ok {
			c.String(re.StatusCode, re.Error())
		} else {
			_ = c.Error(err)
		}
		return
	}
	c.SecureJSON(http.StatusOK, vid)
}

// ShowAccount godoc
// @Summary      Delete a playlist
// @Description  Delete the playlist by ID if it exists
// @Tags         playlists
// @Produce      json
// @Param        id   path      int  true  "Playlist ID"
// @Success      204
// @Failure      400
// @Failure      404  {string}  string "No playlist with this ID"
// @Failure      500
// @Router       /playlists/{id} [delete]
func (vc *PlaylistController[S]) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, `No id provided !`)
		return
	}
	err := vc.Service.VidHost.DeletePlaylist(id)
	if err != nil {
		if re, ok := err.(*video_hosting.RequestError); ok {
			c.String(re.StatusCode, re.Error())
		} else {
			_ = c.Error(err)
		}
		return
	}
	c.String(http.StatusNoContent, "")
}

// ShowAccount godoc
// @Summary      Add a video to the selected playlist
// @Description  Add an existing playlist to an existing video
// @Tags         playlists
// @Produce      json
// @Param        pid   path      int  true  "Playlist ID"
// @Param        vid   path      int  true  "Video ID"
// @Success      204
// @Failure      400
// @Failure      404  {string}  string "Either the playlist or video don't exists"
// @Failure      500
// @Router       /playlists/{pid}/videos/{vid} [put]
func (vc *PlaylistController[S]) AddVideo(c *gin.Context) {
	pId := c.Param("pid")
	if pId == "" {
		c.String(http.StatusBadRequest, `No playlist id provided !`)
		return
	}
	vId := c.Param("pid")
	if vId == "" {
		c.String(http.StatusBadRequest, `No video id provided !`)
		return
	}
	err := vc.Service.VidHost.AddVideoToPlaylist(vId, pId)
	if err != nil {
		if re, ok := err.(*video_hosting.RequestError); ok {
			c.String(re.StatusCode, re.Error())
		} else {
			_ = c.Error(err)
		}
		return
	}
	c.String(http.StatusNoContent, "")
}
