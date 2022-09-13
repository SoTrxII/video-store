package object_storage

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"io"
	"os"
	"path/filepath"
)

// ObjectStorage any S3-like storage solution
type ObjectStorage[T BindingProxy] struct {
	// Destination path for all downloads
	assetsPath string
	// Name of the Dapr component to use
	componentName string
	// Client to query the backend storage
	client *T
	// Current running context
	ctx *context.Context
}

// NewDaprObjectStorage Prod ready constructor for an object-storage using Dapr
func NewDaprObjectStorage(ctx *context.Context, daprClient *client.Client, component string) (*ObjectStorage[client.Client], error) {
	dir, err := os.MkdirTemp("", "downloader-")
	if err != nil {
		return nil, err
	}
	return &ObjectStorage[client.Client]{
		assetsPath:    dir,
		componentName: component,
		client:        daprClient,
		ctx:           ctx,
	}, nil
}

// NewObjectStorage General purpose object storage
func NewObjectStorage[T BindingProxy](ctx *context.Context, assetsPath string, client T) *ObjectStorage[T] {
	return &ObjectStorage[T]{
		assetsPath: assetsPath,
		// TODO :: component configurable
		componentName: "",
		client:        &client,
		ctx:           ctx,
	}
}

// Proxy to query the backend storage
type BindingProxy interface {
	// Invoke
	InvokeBinding(ctx context.Context, in *client.InvokeBindingRequest) (out *client.BindingEvent, err error)
}

// Download a file from the backend storage
func (od ObjectStorage[T]) Download(key string) (path *string, err error) {
	reader, err := od.Buffer(key)
	if err != nil {
		return nil, err
	}
	writePath := filepath.Join(od.assetsPath, key)
	output, err := os.Create(writePath)
	defer output.Close()
	io.Copy(output, *reader)
	if err != nil {
		return nil, err
	}
	return &writePath, nil
}

// Buffer the content of a file in memory
func (od ObjectStorage[T]) Buffer(key string) (data *io.Reader, err error) {
	res, err := (*od.client).InvokeBinding(*od.ctx, &client.InvokeBindingRequest{
		Name:      od.componentName,
		Operation: "get",
		Data:      nil,
		Metadata:  map[string]string{"key": key},
	})
	if err != nil {
		return nil, err
	}
	input := bytes.NewBuffer(res.Data)
	decoder := base64.NewDecoder(base64.StdEncoding, input)
	return &decoder, nil
}

// Upload Uploads a file on the backend storage
func (od ObjectStorage[T]) Upload(path string, key string) error {
	b64bytes, err := readFileToB64(path)
	if err != nil {
		return err
	}
	_, err = (*od.client).InvokeBinding(*od.ctx, &client.InvokeBindingRequest{
		Name:      od.componentName,
		Operation: "create",
		Data:      b64bytes,
		Metadata: map[string]string{
			"key": key,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// Delete a file in the remote object storage
func (od ObjectStorage[T]) Delete(key string) error {
	_, err := (*od.client).InvokeBinding(*od.ctx, &client.InvokeBindingRequest{
		Name:      od.componentName,
		Operation: "delete",
		Data:      nil,
		Metadata: map[string]string{
			"key": key,
		},
	})
	return err
}

// Read a file into a base64 bytes-array
func readFileToB64(path string) ([]byte, error) {
	var buf bytes.Buffer
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	b64enc := base64.NewEncoder(base64.StdEncoding, &buf)
	// Make a 10 MB buffer
	p := make([]byte, 10*1024*1024)
	for {
		n, err := reader.Read(p)
		if err != nil {
			if err == io.EOF {
				b64enc.Write(p[:n])
				fmt.Println(string(p[:n]))
				break
			}
			fmt.Println(err)
		}
		b64enc.Write(p[:n])
	}
	err = b64enc.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
