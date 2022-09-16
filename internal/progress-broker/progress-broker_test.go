package progress_broker

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"testing"
	mock_client "video-manager/internal/mock/dapr"
)

func TestProgressBroker_SendProgress(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	daprClient := mock_client.NewMockClient(ctrl)
	daprClient.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	pg, err := NewProgressBroker[*mock_client.MockClient](&ctx, &daprClient, NewBrokerOptions{
		Component: "",
		Topic:     "",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = pg.SendProgress(UploadInfos{
		RecordId:    "1",
		UploadState: InProgress,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestProgressBroker_SendError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	daprClient := mock_client.NewMockClient(ctrl)
	daprClient.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	pg, err := NewProgressBroker[*mock_client.MockClient](&ctx, &daprClient, NewBrokerOptions{
		Component: "",
		Topic:     "",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = pg.SendProgress(UploadInfos{
		RecordId:    "1",
		UploadState: Error,
		Data:        fmt.Errorf("Test"),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestProgressBroker_SendDone(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	daprClient := mock_client.NewMockClient(ctrl)
	daprClient.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	pg, err := NewProgressBroker[*mock_client.MockClient](&ctx, &daprClient, NewBrokerOptions{
		Component: "",
		Topic:     "",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = pg.SendProgress(UploadInfos{
		RecordId:    "1",
		UploadState: Done,
		Data:        nil,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestProgressBroker_CouldNotSend(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	daprClient := mock_client.NewMockClient(ctrl)
	daprClient.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("test"))
	pg, err := NewProgressBroker[*mock_client.MockClient](&ctx, &daprClient, NewBrokerOptions{
		Component: "",
		Topic:     "",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = pg.SendProgress(UploadInfos{
		RecordId:    "1",
		UploadState: Done,
		Data:        nil,
	})
	if err == nil {
		t.Fatal(err)
	}
}
