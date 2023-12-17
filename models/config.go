package models

import "time"

type MinioSettings struct {
	MinioDomain    string
	MinioAccessKey string
	MinioSecretKey string
	MinioSsl       bool
}

type HubSettings struct {
	HubDomain string
}

type OperationSettings struct {
	FileDownloadLinksExpiration time.Duration
	FileUploadLinksExpiration   time.Duration
}

type Configuration struct {
	MinioSettings    MinioSettings
	HubSettings      HubSettings
	OperationRelated OperationSettings
}
