//go:build integration
// +build integration

package video_hosting

// As the youtube video store is basically only calls to an external API call, there is no point to do unit testing.
// The only type of testing useful here will be integration testing
import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	ResDir = "../../resources/test"
)

func Setup(t *testing.T) *YoutubeVideoStore {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	store, err := NewYoutubeStore(ctx, &YoutubeStoreCredentials{
		ClientId:     os.Getenv("YT_CLIENT_ID"),
		ClientSecret: os.Getenv("YT_CLIENT_SECRET"),
		RefreshToken: os.Getenv("YT_REFRESH_TOKEN"),
	}, nil)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	return store
}

func Test_YoutubeStore_PlaylistLifecycle(t *testing.T) {
	store := Setup(t)
	// Create a playlist
	p, err := store.CreatePlaylist(&ItemMetadata{
		Description: "test-go-api",
		Title:       "test",
		Visibility:  Unlisted,
	})
	assert.Nil(t, err)

	// Propagation time
	time.Sleep(10 * time.Second)
	// Change its title and check that the changes propagated
	const changedTitle = "test2"
	p.Title = changedTitle
	p2, err := store.UpdatePlaylist(p.Id, p)
	assert.Nil(t, err)
	assert.Equal(t, changedTitle, p2.Title)
	time.Sleep(10 * time.Second)
	p3, err := store.RetrievePlaylist(p.Id)
	assert.Nil(t, err)
	assert.Equal(t, changedTitle, p3.Title)

	// Finally delete the playlist
	err = store.DeletePlaylist(p.Id)
	assert.Nil(t, err)
}

func Test_YoutubeStore_VideoLifecycle(t *testing.T) {
	store := Setup(t)
	// Create a video
	f, err := os.Open(filepath.Join(ResDir, "video.mp4"))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	v, err := store.CreateVideo(&ItemMetadata{
		Description: "test-go-api",
		Title:       "test",
		Visibility:  Unlisted,
	}, f, nil)
	assert.Nil(t, err)
	assert.Equal(t, v.Duration, int64(10))
	// Propagation time
	time.Sleep(10 * time.Second)
	// Change its title and check that the changes propagated
	const changedTitle = "test2"
	v.Title = changedTitle
	v2, err := store.UpdateVideo(v.Id, v)
	assert.Nil(t, err)
	assert.Equal(t, changedTitle, v2.Title)
	time.Sleep(10 * time.Second)
	v3, err := store.RetrieveVideo(v.Id)
	assert.Nil(t, err)
	assert.Equal(t, changedTitle, v3.Title)

	// Finally delete the video
	//err = store.DeleteVideo(v.Id)
	assert.Nil(t, err)
}

func Test_YoutubeStore_SetThumbnail(t *testing.T) {
	store := Setup(t)
	// Create a video
	f, err := os.Open(filepath.Join(ResDir, "video.mp4"))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	v, err := store.CreateVideo(&ItemMetadata{
		Description: "test-go-api",
		Title:       "test",
		Visibility:  Unlisted,
	}, f, nil)
	assert.Nil(t, err)
	f, err = os.Open(filepath.Join(ResDir, "test.jpg"))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	err = store.UpdateVideoThumbnail(v.Id, f)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	err = store.DeleteVideo(v.Id)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func Test_YoutubeStore_AddVideoToPlaylist(t *testing.T) {
	store := Setup(t)
	// Create a playlist
	p, err := store.CreatePlaylist(&ItemMetadata{
		Description: "test-go-api",
		Title:       "test",
		Visibility:  Unlisted,
	})
	assert.Nil(t, err)
	// Create a video
	f, err := os.Open(filepath.Join(ResDir, "video.mp4"))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	v, err := store.CreateVideo(&ItemMetadata{
		Description: "test-go-api",
		Title:       "test",
		Visibility:  Unlisted,
	}, f, nil)
	assert.Nil(t, err)

	err = store.AddVideoToPlaylist(v.Id, p.Id)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	//err = store.DeleteVideo(v.Id)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	// err = store.DeletePlaylist(p.Id)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}
