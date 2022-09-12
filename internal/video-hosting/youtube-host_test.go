package video_hosting

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/youtube/v3"
	"testing"
	"time"
)

func TestToGenericPlaylist(t *testing.T) {

	const (
		vidCount     = 10
		cId          = "testChannelId"
		cTitle       = "testChannelTitle"
		id           = "testId"
		desc         = "testdesc"
		creationDate = "2018-08-25T11:12:35Z"
		title        = "testTitle"
		visibility   = "unlisted"
		thumbUrl     = "testThumb"
	)

	// First, without a thumbnail
	ytPlaylist := youtube.Playlist{
		ContentDetails: &youtube.PlaylistContentDetails{
			ItemCount: vidCount,
		},
		Etag: "",
		Id:   id,
		Kind: "",
		Snippet: &youtube.PlaylistSnippet{
			ChannelId:    cId,
			ChannelTitle: cTitle,
			Description:  desc,
			PublishedAt:  creationDate,
			Thumbnails:   nil,
			Title:        title,
		},
		Status: &youtube.PlaylistStatus{
			PrivacyStatus: visibility,
		},
	}
	p, err := toGenericPlaylist(&ytPlaylist)
	assert.Nil(t, err)
	assert.Equal(t, id, p.Id)
	assert.Equal(t, int64(vidCount), p.ItemCount)
	assert.Equal(t, desc, p.Description)
	assert.Equal(t, creationDate, p.CreatedAt.Format(time.RFC3339))
	assert.Equal(t, title, p.Title)
	assert.Equal(t, visibility, string(p.Visibility))
	assert.Equal(t, "", p.ThumbnailUrl)

	// Then, add a thumbnail
	ytPlaylist.Snippet.Thumbnails = &youtube.ThumbnailDetails{
		Default: &youtube.Thumbnail{
			Url: thumbUrl,
		},
	}
	p, err = toGenericPlaylist(&ytPlaylist)
	assert.Nil(t, err)
	assert.Equal(t, id, p.Id)
	assert.Equal(t, int64(vidCount), p.ItemCount)
	assert.Equal(t, desc, p.Description)
	assert.Equal(t, creationDate, p.CreatedAt.Format(time.RFC3339))
	assert.Equal(t, title, p.Title)
	assert.Equal(t, visibility, string(p.Visibility))
	assert.Equal(t, thumbUrl, p.ThumbnailUrl)

	// Error case
	ytPlaylist.Snippet = nil
	_, err = toGenericPlaylist(&ytPlaylist)
	assert.NotNil(t, err)
}

func TestPatchPlaylist(t *testing.T) {
	const (
		vidCount     = 10
		cId          = "testChannelId"
		cTitle       = "testChannelTitle"
		id           = "testId"
		modId        = "testIdModified"
		desc         = "testdesc"
		creationDate = "2018-08-25T11:12:35Z"
		title        = "testTitle"
		modTitle     = "testTitleModified"
		visibility   = "unlisted"
		thumbUrl     = "testThumb"
	)
	ytPlaylist := youtube.Playlist{
		ContentDetails: &youtube.PlaylistContentDetails{
			ItemCount: vidCount,
		},
		Etag: "",
		Id:   id,
		Kind: "",
		Snippet: &youtube.PlaylistSnippet{
			ChannelId:    cId,
			ChannelTitle: cTitle,
			Description:  desc,
			PublishedAt:  creationDate,
			Thumbnails:   nil,
			Title:        title,
		},
		Status: &youtube.PlaylistStatus{
			PrivacyStatus: visibility,
		},
	}
	// Authorized mod, title
	parsed, _ := time.Parse(time.RFC3339, creationDate)
	patch := Playlist{
		Id:        id,
		ItemCount: vidCount,
		// Title is changed, this is authorized
		Title:        modTitle,
		Description:  desc,
		CreatedAt:    parsed,
		Visibility:   visibility,
		ThumbnailUrl: "",
	}
	err := patchYoutubePlaylist(&ytPlaylist, &patch)
	assert.Nil(t, err)
	assert.Equal(t, ytPlaylist.Snippet.Title, patch.Title)

	// Forbidden mod, id
	patch.Id = modId
	err = patchYoutubePlaylist(&ytPlaylist, &patch)
	assert.NotNil(t, err)
	assert.Equal(t, ytPlaylist.Id, id)

}

func TestToGenericVideo(t *testing.T) {
	const (
		duration     = "PT2H"
		id           = "testId"
		desc         = "testdesc"
		creationDate = "2018-08-25T11:12:35Z"
		title        = "testTitle"
		visibility   = "unlisted"
		thumbUrl     = "testThumb"
	)
	// No thumbnail
	ytVid := youtube.Video{
		ContentDetails: &youtube.VideoContentDetails{
			Duration: duration,
		},
		Id: id,
		Snippet: &youtube.VideoSnippet{
			CategoryId:  "",
			Description: desc,
			PublishedAt: creationDate,
			Thumbnails:  nil,
			Title:       title,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: visibility,
		},
	}
	vid, err := toGenericVideo(&ytVid)
	assert.Nil(t, err)
	// 2h -> 7200 secs
	assert.Equal(t, int64(7200), vid.Duration)
	assert.Equal(t, id, vid.Id)
	assert.Equal(t, desc, vid.Description)
	assert.Equal(t, creationDate, vid.CreatedAt.Format(time.RFC3339))
	assert.Equal(t, title, vid.Title)
	assert.Equal(t, visibility, string(vid.Visibility))

	// adding a thumbnail
	ytVid.Snippet.Thumbnails = &youtube.ThumbnailDetails{
		Default: &youtube.Thumbnail{
			Url: thumbUrl,
		},
	}
	vid, err = toGenericVideo(&ytVid)
	assert.Nil(t, err)
	// 2h -> 7200 secs
	assert.Equal(t, int64(7200), vid.Duration)
	assert.Equal(t, id, vid.Id)
	assert.Equal(t, desc, vid.Description)
	assert.Equal(t, creationDate, vid.CreatedAt.Format(time.RFC3339))
	assert.Equal(t, title, vid.Title)
	assert.Equal(t, visibility, string(vid.Visibility))
	assert.Equal(t, thumbUrl, vid.ThumbnailUrl)
}

func TestISO8601DurationToSeconds(t *testing.T) {
	// Err
	parsed, err := iSO8601DurationToSeconds("test")
	assert.Nil(t, parsed)
	assert.NotNil(t, err)

	// Ok
	parsed, err = iSO8601DurationToSeconds("PT12H30M5S")
	assert.Nil(t, err)
	assert.NotNil(t, parsed)
	// "12H30M5S" -> 12 hours, 30 minutes, 5 seconds
	assert.Equal(t, int64(12*3600+30*60+5), *parsed)
}

func TestAssignDefault(t *testing.T) {
	opt := YoutubeStoreOptions{}
	assignDefault(&opt)
	assert.Equal(t, "24", opt.CategoryId)
}
