package schema

// UploadTaskReq upload task request model
type UploadTaskReq struct {
	Filename string `binding:"required" json:"filename"`
	Size     int64  `binding:"gte=1" json:"size"`
	SHA512   string `binding:"required" json:"sha512"`
}

// UploadTaskResp upload task response model
type UploadTaskResp struct {
	TaskID string `json:"taskId"`
}
