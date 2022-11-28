package api

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ismdeep/parser"
)

// PutChunkData put chunk data
func PutChunkData(c *gin.Context) {
	taskID := c.Param("task_id")

	index, err := parser.ToInt64(c.Param("index"))
	if err != nil {
		fail(c, err)
		return
	}

	size, err := parser.ToInt64(c.Param("size"))
	if err != nil {
		fail(c, err)
		return
	}

	hash := c.Request.Header.Get("X-Hash")
	if hash == "" {
		fail(c, errors.New("X-Hash not found"))
		return
	}

	mr, err := c.Request.MultipartReader()
	if err != nil {
		fail(c, err)
		return
	}

	part, err := mr.NextPart()
	if err != nil {
		fail(c, err)
		return
	}

	if part.FormName() != "file" {
		fail(c, errors.New("wrong form field"))
		return
	}

	data := make([]byte, size)
	if _, err := io.ReadFull(part, data); err != nil {
		fail(c, err)
		return
	}

	hashCli := sha512.New()
	hashCli.Write(data)
	hashNew := fmt.Sprintf("%x", hashCli.Sum(nil))

	f, err := os.OpenFile(fmt.Sprintf("%v/%v", DataRoot, taskID), os.O_RDWR, 0600)
	if err != nil {
		fail(c, err)
		return
	}

	if _, err := f.WriteAt(data, index); err != nil {
		fail(c, err)
		return
	}

	if hashNew != hash {
		fail(c, errors.New("data broken"))
		return
	}

	success(c, nil)
}
