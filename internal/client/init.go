package client

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/ismdeep/chunk-upload-demo/internal/schema"
)

// Client model
type Client struct {
	endpoint string
}

// New a client
func New(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

func (receiver *Client) do(request *http.Request, v interface{}) error {
	resp, err := (&http.Client{}).Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Resp struct
	type Resp struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	var r Resp
	if err := json.Unmarshal(body, &r); err != nil {
		return err
	}

	if r.Code != 0 {
		return errors.New(r.Msg)
	}

	// extract v
	raw, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}

	return nil
}

// NewUploadTask new a upload task
func (receiver *Client) NewUploadTask(postData schema.UploadTaskReq) (*schema.UploadTaskResp, error) {
	data, err := json.Marshal(postData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/api/upload-tasks", receiver.endpoint), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var respData schema.UploadTaskResp
	if err := receiver.do(req, &respData); err != nil {
		return nil, err
	}

	return &respData, nil
}

// PutChunk put chunk data
func (receiver *Client) PutChunk(taskID string, index int64, chunkData []byte) error {
	pipeOut, pipeIn := io.Pipe()
	writer := multipart.NewWriter(pipeIn)

	hashCli := sha512.New()
	hashCli.Write(chunkData)
	hash := fmt.Sprintf("%x", hashCli.Sum(nil))

	errChan := make(chan error, 1)
	go func() {
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/api/upload-tasks/%v/chunks/%v/%v", receiver.endpoint, taskID, index, len(chunkData)), pipeOut)
		if err != nil {
			errChan <- err
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("X-Hash", hash)
		if err := receiver.do(req, nil); err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	part, err := writer.CreateFormFile("file", "a")
	if err != nil {
		return err
	}

	if _, err := io.Copy(part, bytes.NewReader(chunkData)); err != nil {
		return err
	}

	_ = writer.Close()
	_ = pipeIn.Close()

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

// Check upload
func (receiver *Client) Check(taskID string) error {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/api/upload-tasks/%v/sha512-check", receiver.endpoint, taskID), nil)
	if err != nil {
		return err
	}

	if err := receiver.do(req, nil); err != nil {
		return err
	}

	return nil
}
