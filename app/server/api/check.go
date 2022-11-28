package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/ismdeep/chunk-upload-demo/internal/schema"
	"github.com/ismdeep/chunk-upload-demo/pkg/util"
)

// CheckData check data
func CheckData(c *gin.Context) {
	taskID := c.Param("task_id")

	filePath := filepath.Clean(fmt.Sprintf("%v/%v", DataRoot, taskID))
	infoFilePath := filepath.Clean(fmt.Sprintf("%v/%v.json", DataRoot, taskID))

	raw, err := os.ReadFile(infoFilePath)
	if err != nil {
		fail(c, err)
		return
	}

	var info schema.UploadTaskReq
	if err := json.Unmarshal(raw, &info); err != nil {
		fail(c, err)
		return
	}

	_, size, err := util.GetFileNameSize(filePath)
	if err != nil {
		fail(c, err)
		return
	}

	if size != info.Size {
		fail(c, err)
		return
	}

	sha512Value, err := util.SHA512File(filePath)
	if err != nil {
		fail(c, err)
		return
	}

	if sha512Value != info.SHA512 {
		fail(c, errors.New("check failed"))
		return
	}

	_ = os.Remove(filePath)
	_ = os.Remove(infoFilePath)

	success(c, nil)
}
