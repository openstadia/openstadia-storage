package connections

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/minio/madmin-go"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/openstadia/openstadia-storage/configuration"
	"github.com/openstadia/openstadia-storage/models"
)

func InitAllConnections(configStore *configuration.ConfigStore) (*ConnectionsStore, error) {
	endpoint := configStore.GetMinioSettings().MinioDomain
	accessKey := configStore.GetMinioSettings().MinioAccessKey
	secretKey := configStore.GetMinioSettings().MinioSecretKey
	useSsl := configStore.GetMinioSettings().MinioSsl

	connections := models.Connections{}

	minioClient, err := InitMainClient(endpoint, accessKey, secretKey, useSsl)
	if err != nil {
		return nil, err
	}
	connections.MinioClient = minioClient

	madminClient, err := InitMadminClient(endpoint, accessKey, secretKey, useSsl)
	if err != nil {
		return nil, err
	}
	connections.MadminClient = madminClient

	//database, err := StartBoltDb()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//connections.Database = database

	store := CreateConnectionsStore(&connections)

	return &store, nil
}

func InitMainClient(endpoint string, accessKey string, secretKey string, useSsl bool) (*minio.Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSsl,
	})
	if err != nil || !minioClient.IsOnline() {
		return nil, fmt.Errorf("failed to initialize the MinIO client connection: %w", err)
	}

	return minioClient, nil
}

func InitMadminClient(endpoint string, accessKey string, secretKey string, useSsl bool) (*madmin.AdminClient, error) {
	mdmClient, err := madmin.New(endpoint, accessKey, secretKey, useSsl)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the Madmin client connection: %w", err)
	}

	return mdmClient, nil
}

func StartBoltDb() (*bolt.DB, error) {
	db, err := bolt.Open("storageUsers.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	//defer func(db *bolt.DB) {
	//	err = db.Close()
	//	if err != nil {
	//		log.Println(fmt.Errorf("failed to connect to Bolt database: %w", err))
	//	}
	//}(db)

	return db, nil
}
