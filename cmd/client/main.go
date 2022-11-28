package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ismdeep/chunk-upload-demo/internal/client"
	"github.com/ismdeep/chunk-upload-demo/internal/schema"
	"github.com/ismdeep/chunk-upload-demo/pkg/util"
)

var dataRoot string

// FileChunkShadow file chunk shadow
type FileChunkShadow struct {
	Index int64
	Size  int64
}

func init() {
	// 根据当前工作目录创建程序运行时使用的数据目录
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dataRoot = fmt.Sprintf("%v/.cache", workDir)

	// 如果数据文件已存在，则直接退出（无需创建）
	filePath := filepath.Clean(fmt.Sprintf("%v/000000-0000-0000-0000-000000000000", dataRoot))
	if util.FileIsExists(filePath) {
		return
	}

	// 创建数据目录
	if err := os.MkdirAll(dataRoot, 0750); err != nil {
		panic(err)
	}

	// 创建数据文件
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer func() { // defer 操作，无论后续代码怎么执行，在defer所在的函数退出之前执行，  **作用：** 用于清理环境，比如这里就是对文件进行关闭。
		if err := f.Close(); err != nil {
			fmt.Println("Warn:", err)
		}
	}()

	// 生成随机数据并写入数据文件
	data := make([]byte, 2*1024*1024*1024+5)
	_, _ = rand.Read(data)
	if _, _ = f.Write(data); err != nil {
		panic(err)
	}
}

// UploadFile upload a file
func UploadFile(path string) error {
	fileName, size, err := util.GetFileNameSize(path)
	if err != nil {
		return err
	}

	sha512Value, err := util.SHA512File(path)
	if err != nil {
		return err
	}

	apiClient := client.New("http://127.0.0.1:9000")

	// create upload task
	resp, err := apiClient.NewUploadTask(schema.UploadTaskReq{
		Filename: fileName,
		Size:     size,
		SHA512:   sha512Value,
	})
	if err != nil {
		panic(err)
	}

	// upload
	chunks := make(chan FileChunkShadow, 1024)
	go func() {
		tmpSize := size
		index := int64(0)
		for tmpSize > 0 {
			chunkSize := int64(1 * 1024 * 1024) // chunk size: 1MB
			if chunkSize > tmpSize {
				chunkSize = tmpSize
			}
			chunks <- FileChunkShadow{
				Index: index,
				Size:  chunkSize,
			}
			index += chunkSize
			tmpSize -= chunkSize
		}
		close(chunks)
	}()

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			for chunk := range chunks {
				data, err := util.ReadChunk(path, chunk.Index, chunk.Size)
				if err != nil {
					panic(err)
				}
				if err := apiClient.PutChunk(resp.TaskID, chunk.Index, data); err != nil {
					panic(err)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	// check
	if err := apiClient.Check(resp.TaskID); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := UploadFile(fmt.Sprintf("%v/000000-0000-0000-0000-000000000000", dataRoot)); err != nil {
		panic(err)
	}
}
