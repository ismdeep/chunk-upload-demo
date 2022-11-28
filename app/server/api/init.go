package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// DataRoot data root
var DataRoot string

func init() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	DataRoot = fmt.Sprintf("%v/.cache", workDir)
	if err := os.MkdirAll(DataRoot, 0750); err != nil {
		panic(err)
	}
}

// Eng inst
var Eng *gin.Engine

func init() {
	Eng = gin.Default()

	Eng.POST("/api/upload-tasks", NewUploadTask)
	Eng.PUT("/api/upload-tasks/:task_id/chunks/:index/:size", PutChunkData)
	Eng.PUT("/api/upload-tasks/:task_id/sha512-check", CheckData)
}
