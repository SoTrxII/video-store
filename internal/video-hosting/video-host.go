package video_hosting

import (
	"io"
	"time"
)

// ProgressFunc Function called to handle long operation progress
type ProgressFunc func(current int64, total int64)

// A videos-store is an interface for any videos hosting platform
type IVideoHost interface {

	/* Videos CRUD */

	// CreateVideo Upload a new video on the hosting platform
	// onProgress is an optional callback, set it to null to ignore it
	CreateVideo(meta *ItemMetadata, uploadContent io.Reader, onProgress *ProgressFunc) (*Video, error)
	// RetrieveVideo Search an existing video given its ID.
	// Return nil if the video with this specific ID doesn't exists
	RetrieveVideo(id string) (*Video, error)
	// UpdateVideo Update the info of the video identified by id with the infos of replacement
	// /!\ Some attributes ("id", "creationDate", "duration") are READ-ONLY
	UpdateVideo(id string, replacement *Video) (*Video, error)
	// DeleteVideo Delete an existing video from the remote video hosting platform
	DeleteVideo(id string) error

	/* Playlist CRUD */

	// CreatePlaylist Creates a playlist on the remote video hosting platform
	CreatePlaylist(meta *ItemMetadata) (*Playlist, error)
	// RetrievePlaylist Search an existing playlist with its ID.
	// Return nil if the playlist with this specific ID doesn't exists
	RetrievePlaylist(id string) (*Playlist, error)
	// UpdatePlaylist Update the info of the playlist identified by id with the infos of replacement
	// /!\ Some attributes ("id", "creationDate") are READ-ONLY
	UpdatePlaylist(id string, replacement *Playlist) (*Playlist, error)
	// DeletePlaylist Delete an existing playlist from the remote video hosting platform
	DeletePlaylist(id string) error

	/* Utilities */

	// AddVideoToPlaylist Add an existing video to an existing playlist on the hosting platform
	// TODO
	AddVideoToPlaylist(videoId string, playlistId string) error
	// UpdateVideoThumbnail Set the thumbnail for a video
	UpdateVideoThumbnail(videoId string, thumbnailContent io.Reader) error
}

// Video A video hosted on a video storage website
type Video struct {
	Id string `json:"id"`
	// Video display name
	Title string `json:"title" validate:"updatable"`
	// Video description
	Description string `json:"description" validate:"updatable"`
	// Creation date
	CreatedAt time.Time `json:"createdAt"`
	// Video duration in seconds
	Duration int64 `json:"duration"`
	// public/private/unlisted
	Visibility Visibility `json:"visibility" validate:"updatable"`
	// Playlist thumbnail
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
}

// Playlist A collection of videos hosted on a video storage
// website
type Playlist struct {
	Id string `json:"id"`
	// Number of video in this playlist
	ItemCount int64 `json:"itemCount"`
	// Playlist display name
	Title string `json:"title" validate:"updatable"`
	// Short description about what's in the playlist
	Description string `json:"description" validate:"updatable"`
	// Creation date
	CreatedAt time.Time `json:"createdAt"`
	// public/private/unlisted
	Visibility Visibility `json:"visibility" validate:"updatable"`
	// Playlist thumbnail
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
}

// Represent the visibility of an object on the video storage
type Visibility string

const (
	Public   Visibility = "public"
	Private  Visibility = "private"
	Unlisted Visibility = "unlisted"
)

// All metadata about the item to update
type ItemMetadata struct {
	// Short text describing the content of the item
	// Youtube actually limits to 5000 bytes, which *isn't* 5000 characters
	// https://developers.google.com/youtube/v3/docs/videos#properties
	Description string `json:"description" binding:"max=1000"`
	// Title of the item
	// The max character limitation is currently taken from the Yt docs
	// https://developers.google.com/youtube/v3/docs/videos#properties
	// This may change if another provider is requiring less than 100 characters
	Title string `json:"title" binding:"required,max=100"`
	// Visibility of the item
	Visibility Visibility `json:"visibility" binding:"required"`
}

// This error is only thrown when an error
type RequestError struct {
	StatusCode int
	Err        error
}

func (r *RequestError) Error() string {
	return r.Err.Error()
}
