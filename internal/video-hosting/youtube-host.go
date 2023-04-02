package video_hosting

import (
	"context"
	"fmt"
	"github.com/senseyeio/duration"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"io"
	"net/http"
	"strings"
	"time"
)

func (ytP YoutubeVideoStore) CreateVideo(meta *ItemMetadata, uploadContent io.Reader, onProgress *ProgressFunc) (*Video, error) {
	call := ytP.Service.Videos.Insert([]string{"id", "snippet", "status", "contentDetails"}, &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Description: meta.Description,
			Title:       meta.Title,
			CategoryId:  ytP.Options.CategoryId,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: string(meta.Visibility),
		},
	})

	// The progress callback is optional
	if onProgress != nil {
		call.ProgressUpdater(googleapi.ProgressUpdater(*onProgress))
	}

	ytVid, err := call.Media(uploadContent).Do()
	if err != nil {
		return nil, handleGoogleApiError(err)
	}

	// We are forced to make a separate API call to get all the files details.
	ytVid2, err := ytP.getYoutubeVideoById(ytVid.Id)
	if err != nil {
		return nil, handleGoogleApiError(err)
	}

	return toGenericVideo(ytVid2)
}

func (ytP YoutubeVideoStore) RetrieveVideo(id string) (*Video, error) {
	ytVid, err := ytP.getYoutubeVideoById(id)
	if err != nil {
		return nil, handleGoogleApiError(err)
	}
	return toGenericVideo(ytVid)
}

func (ytP YoutubeVideoStore) UpdateVideo(id string, replacement *Video) (*Video, error) {
	ytVid, err := ytP.getYoutubeVideoById(id)
	if err != nil {
		return nil, err
	}
	err = patchYoutubeVideo(ytVid, replacement)
	if err != nil {
		return nil, err
	}
	call := ytP.Service.Videos.Update([]string{"snippet", "status", "contentDetails", "id"}, ytVid)
	updated, err := call.Do()
	if err != nil {
		return nil, handleGoogleApiError(err)
	}
	return toGenericVideo(updated)
}

func (ytP YoutubeVideoStore) DeleteVideo(id string) error {
	call := ytP.Service.Videos.Delete(id)
	return call.Do()
}

func (ytP YoutubeVideoStore) GetVideoAccessPrefix() string {
	return getYoutubePrefix()
}

func (ytP YoutubeVideoStore) CreatePlaylist(meta *ItemMetadata) (*Playlist, error) {
	call := ytP.Service.Playlists.Insert([]string{"snippet", "status", "contentDetails"}, &youtube.Playlist{
		Snippet: &youtube.PlaylistSnippet{
			Description: meta.Description,
			Title:       meta.Title,
		},
		Status: &youtube.PlaylistStatus{
			PrivacyStatus: string(meta.Visibility),
		},
	})
	YtPlaylist, err := call.Do()
	if err != nil {
		return nil, handleGoogleApiError(err)
	}
	playlist, err := toGenericPlaylist(YtPlaylist)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

func (ytP YoutubeVideoStore) RetrievePlaylist(id string) (*Playlist, error) {
	ytPlaylist, err := ytP.getYoutubePlaylistById(id)
	if err != nil {
		return nil, handleGoogleApiError(err)
	}
	playlist, err := toGenericPlaylist(ytPlaylist)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

func (ytP YoutubeVideoStore) UpdatePlaylist(id string, replacement *Playlist) (*Playlist, error) {
	currentPlaylist, err := ytP.getYoutubePlaylistById(id)
	if err != nil {
		return nil, err
	}
	err = patchYoutubePlaylist(currentPlaylist, replacement)
	if err != nil {
		return nil, err
	}
	call := ytP.Service.Playlists.Update([]string{"snippet", "status", "contentDetails"}, currentPlaylist)
	updated, err := call.Do()
	if err != nil {
		return nil, handleGoogleApiError(err)
	}
	return toGenericPlaylist(updated)
}

func (ytP YoutubeVideoStore) DeletePlaylist(id string) error {
	call := ytP.Service.Playlists.Delete(id)
	err := call.Do()
	if err != nil {
		return handleGoogleApiError(err)
	}
	return nil
}

func (ytP YoutubeVideoStore) UpdateVideoThumbnail(videoId string, thumbnailContent io.Reader) error {
	call := ytP.Service.Thumbnails.Set(videoId)
	call.Media(thumbnailContent)
	_, err := call.Do()
	if err != nil {
		return handleGoogleApiError(err)
	}
	return nil
}

func (ytP YoutubeVideoStore) AddVideoToPlaylist(videoId string, playlistId string) error {
	call := ytP.Service.PlaylistItems.Insert([]string{"snippet"}, &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: playlistId,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: videoId,
			},
		},
	})
	_, err := call.Do()
	if err != nil {
		return handleGoogleApiError(err)
	}
	return nil
}

// Retrieve a youtube video with the provided ID
// Errors if not found
func (ytP YoutubeVideoStore) getYoutubeVideoById(id string) (*youtube.Video, error) {
	call := ytP.Service.Videos.List([]string{"contentDetails", "id", "snippet", "status", "fileDetails"})
	call.Id(id)
	res, err := call.Do()
	if err != nil {
		return nil, err
	}
	if len(res.Items) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return res.Items[0], nil
}

// Retrieve a youtube playlist with the provided ID
// Errors if not found
func (ytP YoutubeVideoStore) getYoutubePlaylistById(id string) (*youtube.Playlist, error) {
	call := ytP.Service.Playlists.List([]string{"snippet", "status", "contentDetails"})
	call.Id(id)
	res, err := call.Do()
	if err != nil {
		return nil, err
	}
	if len(res.Items) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return res.Items[0], nil
}

// Return the prefix in which we can plug an ID to watch a video
func getYoutubePrefix() string {
	return "https://www.youtube.com/watch?v="
}

// Converts a Youtube-specific video in a generic video
// /!\ The youtube video input must contain the parts "fileDetails", "id", "snippet" and "status"
func toGenericVideo(in *youtube.Video) (*Video, error) {
	if in.Snippet == nil || in.Status == nil {
		return nil, fmt.Errorf(`Missing some required parts (snippet or status`)
	}
	creationDate, err := time.Parse(time.RFC3339, in.Snippet.PublishedAt)
	if err != nil {
		return nil, err
	}
	thumbUrl := ""
	if in.Snippet.Thumbnails != nil && in.Snippet.Thumbnails.Default != nil {
		thumbUrl = in.Snippet.Thumbnails.Default.Url
	}

	// Duration has to be parsed for the fileDetails instead of contentdetails
	// as the contentDetails duration field is set after the video has been processed
	duration := int64(0)
	if in.FileDetails != nil {
		duration = int64(in.FileDetails.DurationMs / 1000)
	}
	return &Video{
		Id:           in.Id,
		Title:        in.Snippet.Title,
		Description:  in.Snippet.Description,
		CreatedAt:    creationDate,
		Duration:     duration,
		Visibility:   Visibility(in.Status.PrivacyStatus),
		ThumbnailUrl: thumbUrl,
		WatchPrefix:  getYoutubePrefix(),
	}, nil
}

func iSO8601DurationToSeconds(in string) (*int64, error) {
	const (
		HourInSeconds   = 3600
		MinuteInSeconds = 60
	)
	d, err := duration.ParseISO8601(in)
	if err != nil {
		return nil, err
	}
	vidDuration := int64(0)
	vidDuration += int64(d.TH * HourInSeconds)
	vidDuration += int64(d.TM * MinuteInSeconds)
	vidDuration += int64(d.TS)
	return &vidDuration, nil
}

// Converts a Youtube-specific playlist in a generic playlist
// /!\ The youtube playlist input must contain the parts "snippet", "contentDetails" and "status"
func toGenericPlaylist(in *youtube.Playlist) (*Playlist, error) {
	if in.Snippet == nil || in.Status == nil || in.ContentDetails == nil {
		return nil, fmt.Errorf(`Missing some required parts (snippet, contentDetails or status`)
	}
	creationDate, err := time.Parse(time.RFC3339, in.Snippet.PublishedAt)
	if err != nil {
		return nil, err
	}
	thumbUrl := ""
	if in.Snippet.Thumbnails != nil && in.Snippet.Thumbnails.Default != nil {
		thumbUrl = in.Snippet.Thumbnails.Default.Url
	}
	return &Playlist{
		Id:           in.Id,
		ItemCount:    in.ContentDetails.ItemCount,
		Title:        in.Snippet.Title,
		Description:  in.Snippet.Description,
		CreatedAt:    creationDate,
		Visibility:   Visibility(in.Status.PrivacyStatus),
		ThumbnailUrl: thumbUrl,
	}, nil
}

// In place update src to match changes in patch
// /!\ The youtube playlist input must contain the parts "snippet", "contentDetails" and "status"
func patchYoutubePlaylist(src *youtube.Playlist, patch *Playlist) error {

	srcCreationDate, err := time.Parse(time.RFC3339, src.Snippet.PublishedAt)
	if err != nil {
		return err
	}
	// Fail on forbidden changes
	//We're not trying to list all possible forbidden attributes, this is just a best-effort
	// to prevent an unnecessary call to the Youtube Data API
	if patch.Id != src.Id || patch.CreatedAt != srcCreationDate {
		return fmt.Errorf(`Attempted to change a read-only attribute (either "id", or "createdAt")`)
	}
	// Update all attributes that can be modified
	src.Snippet.Title = patch.Title
	src.Snippet.Description = patch.Description
	src.Status.PrivacyStatus = string(patch.Visibility)

	return nil
}

// In place update src to match changes in patch
// /!\ The youtube video input must contain the parts "snippet", "contentDetails" and "status"
func patchYoutubeVideo(src *youtube.Video, patch *Video) error {

	srcCreationDate, err := time.Parse(time.RFC3339, src.Snippet.PublishedAt)
	if err != nil {
		return err
	}
	// Fail on forbidden changes.
	// We're not trying to list all possible forbidden attributes, this is just a best-effort
	// to prevent an unnecessary call to the Youtube Data API
	if patch.Id != src.Id || patch.CreatedAt != srcCreationDate {
		return fmt.Errorf(`Attempted to change a read-only attribute (either "id", or "createdAt")`)
	}
	// Update all attributes that can be modified
	src.Snippet.Title = patch.Title
	src.Snippet.Description = patch.Description
	src.Status.PrivacyStatus = string(patch.Visibility)

	return nil
}

func NewYoutubeStore(ctx context.Context, creds *YoutubeStoreCredentials, opt *YoutubeStoreOptions) (*YoutubeVideoStore, error) {
	// Generalist Google oauth config
	config := oauth2.Config{
		// Project ID
		ClientID: creds.ClientId,
		// Project's secret
		ClientSecret: creds.ClientSecret,
		// URL to Google auth services
		Endpoint: google.Endpoint,
		// We want to get access to Youtube, impersonating the user and to be able to upload videos
		Scopes: []string{youtube.YoutubeScope, youtube.YoutubeUploadScope},
	}
	// Initialize an "empty" token, only using the refresh token.
	// We'll let the token source initialize the access token and expiry
	token := &oauth2.Token{RefreshToken: creds.RefreshToken}
	// Using token source, the access token will get auto refreshed
	ytService, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}

	//Assign default values to options and go on
	if opt == nil {
		opt = &YoutubeStoreOptions{}
	}
	assignDefault(opt)
	return &YoutubeVideoStore{Service: ytService, Options: opt}, nil
}

// Assign all default options to the youtube store
func assignDefault(opt *YoutubeStoreOptions) {
	const (
		Entertainment = "24"
	)
	if opt.CategoryId == "" {
		opt.CategoryId = Entertainment
	}
}

// Handle a google api error, extracting the status code
func handleGoogleApiError(err error) error {
	if ge, ok := err.(*googleapi.Error); ok {
		return &RequestError{ge.Code, ge}
	} else if strings.Contains(strings.ToLower(err.Error()), "not found") {
		// Special case : A not found is thrown when a call to a Google APIs succeeded
		// but no result were found
		return &RequestError{http.StatusNotFound, err}
	}
	return err
}

// YoutubeStoreCredentials all info required to authenticate to Youtube Data API v3
type YoutubeStoreCredentials struct {
	// Client ID for a google project (obtained from https://console.cloud.google.com/apis)
	ClientId string
	// Client secret for a google project (obtained from https://console.cloud.google.com/apis)
	ClientSecret string
	// Refresh token obtained via auth (see readme, obtained from https://developers.google.com/oauthplayground/)
	RefreshToken string
}

// YoutubeStoreOptions all options to initialize a Youtube store
type YoutubeStoreOptions struct {
	CategoryId string
}

type YoutubeVideoStore struct {
	Service *youtube.Service
	Options *YoutubeStoreOptions
}
