//go:build integration
// +build integration

// Integration testing for object-storage. Dapr must be booted up for this to run
package object_storage

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"path"
	"testing"
)

const (
	BucketPath    = "../../resources/bucket/"
	DaprComponent = "object-store"
	// Name of the file used for testing
	TestFileName = "test.txt"
	// Name of the asset used to test download
	TestDownloadAssetKey = "testdl"
	// Name of the asset used to test deletion
	TestDeleteAssetKey = "testdelete"
)

type e2eTestSuite struct {
	suite.Suite
	client   client.Client
	objStore *ObjectStorage[client.Client]
}

func (s *e2eTestSuite) SetupSuite() {
	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		s.Fail(err.Error())
	}
	// Check if the sidecar is up
	daprClient, err := client.NewClient()
	if err != nil {
		fmt.Println("DAPR IS NOT RUNNING")
		s.Fail(err.Error())
	}
	s.client = daprClient
	ctx := context.Background()
	// Copy a file in the bucket directory to download it later
	err = copy(&s.client, path.Join(ResPath, TestFileName), TestDownloadAssetKey)
	err = copy(&s.client, path.Join(ResPath, TestFileName), TestDeleteAssetKey)
	if err != nil {
		s.Fail(err.Error())
	}
	s.objStore = &ObjectStorage[client.Client]{
		assetsPath:    dir,
		componentName: DaprComponent,
		client:        &daprClient,
		ctx:           &ctx,
	}

}

func (s *e2eTestSuite) TestDownload_Int() {
	expectedPath, err := s.objStore.Download(TestDownloadAssetKey)
	if err != nil {
		s.T().Fatal(err)
	}
	actual, err := os.ReadFile(*expectedPath)
	if err != nil {
		s.T().Fatal(err)
	}
	expected, err := os.ReadFile(path.Join(ResPath, TestDownloadAssetKey))

	assert.Equal(s.T(), string(expected), string(actual))
}

func (s *e2eTestSuite) TestDownload_Int_NotExists() {
	_, err := s.objStore.Download("notexists")
	if err == nil {
		s.T().Fatal(err)
	}
	fmt.Println(err)
	assert.Contains(s.T(), err.Error(), "does not exist")
}

func (s *e2eTestSuite) TestUpload_Int() {
	file := "audio.m4a"
	err := s.objStore.Upload(path.Join(ResPath, file), file)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *e2eTestSuite) TestDelete_Int() {
	// Check that the file can be downloaded
	_, err := s.objStore.Download(TestDeleteAssetKey)
	if err != nil {
		s.T().Fatal(err)
	}
	// Then delete it
	err = s.objStore.Delete(TestDeleteAssetKey)
	if err != nil {
		s.T().Fatal(err)
	}
	// And check that it cannot be downloaded anymore
	res, err := s.objStore.Download(TestDeleteAssetKey)
	if err == nil {
		s.T().Fatal(err)
	}
	_ = res
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func copy(daprClient *client.Client, src string, keyName string) error {
	rawContent, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	b64Content := make([]byte, base64.StdEncoding.EncodedLen(len(rawContent)))
	base64.StdEncoding.Encode(b64Content, rawContent)
	_, err = (*daprClient).InvokeBinding(context.Background(), &client.InvokeBindingRequest{
		Name:      DaprComponent,
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
