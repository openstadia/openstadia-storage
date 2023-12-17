package operations

import (
	"context"
	"fmt"
	"github.com/minio/madmin-go"
	"github.com/minio/minio-go/v7"
	"github.com/openstadia/openstadia-storage/configuration"
	"github.com/openstadia/openstadia-storage/connection/connections"
	"github.com/openstadia/openstadia-storage/connection/crud"
	"github.com/openstadia/openstadia-storage/models"
	"net/url"
)

func GetUserBucketName(u *models.User) string {
	return fmt.Sprintf("user-%d", u.Id)
}

func objectInfoToStorageObject(info minio.ObjectInfo) models.StorageObject {
	obj := models.StorageObject{}
	obj.Name = info.Key
	obj.Size = uint64(info.Size)
	//obj.TimeCreated =
	obj.TimeEdited = info.LastModified

	return obj
}

func BucketExists(u *models.User, connectionsStore *connections.ConnectionsStore) (bool, error) {
	exists, err := connectionsStore.GetMinioClient().BucketExists(context.Background(), GetUserBucketName(u))
	if err != nil {
		return false, err
	}

	return exists, nil
}

func CreateBucketAndQuotaIfNeeded(u *models.User, connectionsStore *connections.ConnectionsStore) error {
	bucketName := GetUserBucketName(u)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exists, err := BucketExists(u, connectionsStore)
	if err != nil {
		return err
	}
	if !exists {
		err = connectionsStore.GetMinioClient().MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		userInfo := crud.GetStorageUserInfo(u)
		err = connectionsStore.GetMadminClient().SetBucketQuota(ctx, bucketName, &madmin.BucketQuota{Quota: userInfo.TotalSpaceAvailable, Type: "hard"})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetUserFilesList(u *models.User, connectionsStore *connections.ConnectionsStore) (*models.ObjectsList, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	objectCh := connectionsStore.GetMinioClient().ListObjects(ctx, GetUserBucketName(u), minio.ListObjectsOptions{})
	var objects = make([]models.StorageObject, 0)
	var usedSpace uint64 = 0
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, objectInfoToStorageObject(object))
		usedSpace += uint64(object.Size)
	}

	result := models.ObjectsList{Objects: objects}
	return &result, nil
}

func GetBucketStorageInfo(u *models.User, connectionsStore *connections.ConnectionsStore) (*models.StorageInfo, error) {
	accountInfo, err := connectionsStore.GetMadminClient().AccountInfo(context.Background(), madmin.AccountOpts{})
	if err != nil {
		return nil, err
	}

	bucketName := GetUserBucketName(u)
	for _, b := range accountInfo.Buckets {
		if b.Name == bucketName {
			bucketInfo := models.StorageInfo{UsedSpace: b.Size, TotalSpace: b.Details.Quota.Quota}
			return &bucketInfo, nil
		}
	}

	return nil, nil
}

func removeBucketObject(u *models.User, connectionsStore *connections.ConnectionsStore, path string) *models.UnsuccessfulFileOperation {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := connectionsStore.GetMinioClient().RemoveObject(ctx, GetUserBucketName(u), path, minio.RemoveObjectOptions{}) // ForceDelete??? It's not documented
	if err != nil {
		return &models.UnsuccessfulFileOperation{Path: path, Error: err}
	}
	return nil
}

func DeleteFiles(u *models.User, connectionsStore *connections.ConnectionsStore, files models.DeleteFiles) models.UnsuccessfulFileOperationsList {
	var failedOperations = make([]models.UnsuccessfulFileOperation, 0)
	for _, path := range files.Paths {
		failedOperation := removeBucketObject(u, connectionsStore, path)
		if failedOperation != nil {
			failedOperations[len(failedOperations)] = *failedOperation
		}
	}

	list := models.UnsuccessfulFileOperationsList{List: failedOperations}
	return list
}

func GetFile(u *models.User, file models.GetFile, configStore *configuration.ConfigStore, connectionsStore *connections.ConnectionsStore) (*url.URL, error) {
	_, err := connectionsStore.GetMinioClient().StatObject(context.Background(), GetUserBucketName(u), file.Path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	reqParams := make(url.Values)
	// will need to be changed to only the file part of the entire path when folders are made and the path becomes actually a path
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Path))

	presignedURL, err := connectionsStore.GetMinioClient().PresignedGetObject(context.Background(), GetUserBucketName(u), file.Path, configStore.GetOperationSettings().FileDownloadLinksExpiration, reqParams)
	if err != nil {
		return nil, err
	}

	return presignedURL, nil
}

func UploadFile(u *models.User, file models.UploadFile, configStore *configuration.ConfigStore, connectionsStore *connections.ConnectionsStore) (*url.URL, error) {
	_, err := connectionsStore.GetMinioClient().StatObject(context.Background(), GetUserBucketName(u), file.Path, minio.GetObjectOptions{})
	if err == nil {
		return nil, fmt.Errorf("failed to upload a file: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	presignedURL, err := connectionsStore.GetMinioClient().PresignedPutObject(ctx, GetUserBucketName(u), file.Path, configStore.GetOperationSettings().FileUploadLinksExpiration)
	if err != nil {
		return nil, err
	}

	return presignedURL, nil
}
