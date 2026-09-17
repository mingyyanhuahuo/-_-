package service

import "time"

const (
	defaultPageSize = 10
	maxPageSize     = 50
)

const (
	maxUploadFileSize = 5 << 20 // 5MB
)
const fileCleanGrace = 24 * time.Hour
