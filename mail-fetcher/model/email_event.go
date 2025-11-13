package model

type EmailUploadedEvent struct {
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
	FilePath string `json:"file_path"`
	S3Key    string `json:"s3_key"`
}
