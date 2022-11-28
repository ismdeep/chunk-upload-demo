package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ismdeep/chunk-upload-demo/internal/schema"
)

// NewUploadTask new upload task
func NewUploadTask(c *gin.Context) {
	var req schema.UploadTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, err)
		return
	}

	taskID := uuid.NewString()
	f, err := os.Create(fmt.Sprintf("%v/%v", DataRoot, taskID))
	if err != nil {
		fail(c, err)
		return
	}

	if err := f.Truncate(req.Size); err != nil {
		fail(c, err)
		return
	}

	raw, err := json.Marshal(req)
	if err != nil {
		fail(c, err)
		return
	}

	if err := os.WriteFile(fmt.Sprintf("%v/%v.json", DataRoot, taskID), raw, 0600); err != nil {
		fail(c, err)
		return
	}

	success(c, schema.UploadTaskResp{TaskID: taskID})
}
