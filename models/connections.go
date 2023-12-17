package models

import (
	"github.com/boltdb/bolt"
	"github.com/minio/madmin-go"
	"github.com/minio/minio-go/v7"
)

type Connections struct {
	MinioClient  *minio.Client
	MadminClient *madmin.AdminClient
	Database     *bolt.DB
}
