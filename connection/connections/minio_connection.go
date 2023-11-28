package connections

import (
	"errors"
	"github.com/minio/madmin-go"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/openstadia/openstadia-storage/configuration"
)

var Client *minio.Client
var MadminClient *madmin.AdminClient

func InitMinioClients() error {
	endpoint := configuration.Settings.MinioDomain
	accessKey := configuration.Settings.MinioAccessKey
	secretKey := configuration.Settings.MinioSecretKey
	useSsl := configuration.Settings.MinioSsl

	err := InitMainClient(endpoint, accessKey, secretKey, useSsl)
	if err != nil {
		return err
	}
	err = InitMadminClient(endpoint, accessKey, secretKey, useSsl)
	if err != nil {
		return err
	}
	return nil
}

func InitMainClient(endpoint string, accessKey string, secretKey string, useSsl bool) error {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSsl,
	})
	if err != nil || !minioClient.IsOnline() {
		return errors.New("failed to initialize the MinIO client connection:\n" + err.Error())
	}

	Client = minioClient
	return nil
}

func InitMadminClient(endpoint string, accessKey string, secretKey string, useSsl bool) error {
	mdmClient, err := madmin.New(endpoint, accessKey, secretKey, useSsl)
	if err != nil {
		return errors.New("failed to initialize the Madmin client connection:\n" + err.Error())
	}
	MadminClient = mdmClient
	return nil
}
