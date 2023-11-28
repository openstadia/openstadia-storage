package minio_operations

import (
	"context"
	"errors"
	"github.com/minio/madmin-go"
	"github.com/minio/minio-go/v7"
	"github.com/openstadia/openstadia-storage/configuration"
	"github.com/openstadia/openstadia-storage/connection/connections"
	"github.com/openstadia/openstadia-storage/connection/crud"
	"github.com/openstadia/openstadia-storage/models"
	"net/url"
	"strconv"
)

func GetUserBucketName(u *models.HubUser) string {
	return "user-" + strconv.Itoa(u.Id)
}

func objectInfoToStorageObject(info minio.ObjectInfo) models.StorageObject {
	obj := models.StorageObject{}
	obj.Name = info.Key
	obj.Size = uint64(info.Size)
	//obj.TimeCreated =
	obj.TimeEdited = info.LastModified
	return obj
}

func BucketExists(u *models.HubUser) (bool, error) {
	exists, err := connections.Client.BucketExists(context.Background(), GetUserBucketName(u))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func CreateBucketIfNeeded(u *models.HubUser) error {
	bucketName := GetUserBucketName(u)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exists, err := BucketExists(u)
	if err != nil {
		return err
	}
	if !exists {
		//ctx2, cancel2 := context.WithCancel(context.Background())
		//defer cancel2()
		err = connections.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		userInfo := crud.GetStorageUserInfo(u)
		err = connections.MadminClient.SetBucketQuota(ctx, bucketName, &madmin.BucketQuota{Quota: userInfo.AllowedTotalSpace, Type: "hard"})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetUserFilesList(u *models.HubUser) (*models.ObjectsList, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	objectCh := connections.Client.ListObjects(ctx, GetUserBucketName(u), minio.ListObjectsOptions{})
	var objects []models.StorageObject
	var usedSpace uint64 = 0
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, objectInfoToStorageObject(object))
		usedSpace += uint64(object.Size)
	}

	result := models.ObjectsList{Objects: objects, UsedSpace: usedSpace}
	if len(objects) == 0 {
		result.Objects = []models.StorageObject{}
	}

	return &result, nil
}

func GetBucketQuota(u *models.HubUser) (*madmin.BucketQuota, error) {
	bucketStats, err := connections.MadminClient.GetBucketQuota(context.Background(), GetUserBucketName(u))
	if err != nil {
		return nil, err
	}
	return &bucketStats, err
}

func GetGeneralStorageInfo(u *models.HubUser) (*models.StorageInfo, error) {
	list, err := GetUserFilesList(u)
	if err != nil {
		return nil, err
	}
	result := &models.StorageInfo{}

	quota, err := GetBucketQuota(u)
	if err != nil {
		return nil, err
	}

	result.TotalSpace = quota.Quota
	result.UsedSpace = 0
	for _, obj := range list.Objects {
		result.UsedSpace += obj.Size
	}
	return result, nil
}

func DeleteFiles(u *models.HubUser, files models.DeleteFile) ([]models.UnsuccessfulFileOperation, []error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var failedOperations = []models.UnsuccessfulFileOperation{} // code analysis tools say to replace this, but then it will be initially null
	var exactErrors []error
	for _, path := range files.Paths {
		err := connections.Client.RemoveObject(ctx, GetUserBucketName(u), path, minio.RemoveObjectOptions{}) // ForceDelete??? It's not documented
		if err != nil {
			failedOperations[len(failedOperations)] = models.UnsuccessfulFileOperation{Path: path}
			exactErrors[len(exactErrors)] = err
			continue
		}
	}
	//if len(failedOperations) == 0 {
	//	return nil, nil
	//}
	return failedOperations, exactErrors
}

func GetFile(u *models.HubUser, file models.GetFile) (*url.URL, error) { // later think about downloading multiple in an archive
	_, err := connections.Client.StatObject(context.Background(), GetUserBucketName(u), file.Path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	reqParams := make(url.Values)
	// will need to be changed to only the file part of the entire path when folders are made and the path becomes actually a path
	reqParams.Set("response-content-disposition", "attachment; filename=\""+file.Path+"\"")
	presignedURL, err := connections.Client.PresignedGetObject(context.Background(), GetUserBucketName(u), file.Path, configuration.OperationSettings.FileDownloadLinksExpiration, reqParams)
	if err != nil {
		return nil, err
	}
	return presignedURL, nil
}

// BROKEN_DeleteMultipleFiles It's BROKEN, leaving just in case someone decides to finish it
func BROKEN_DeleteMultipleFiles(u *models.HubUser, files models.DeleteFile) ([]models.UnsuccessfulFileOperation, []error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var failedOperations = []models.UnsuccessfulFileOperation{} // code analysis tools say to replace this, but then it will be initially null
	var exactErrors []error

	objs := make(chan minio.ObjectInfo)
	for _, path := range files.Paths {
		go func(path string) {
			obj, err := connections.Client.StatObject(ctx, GetUserBucketName(u), path, minio.StatObjectOptions{})
			if err != nil {
				return
			}
			objs <- obj
		}(path)
	}
	close(objs)

	err := connections.Client.RemoveObjects(ctx, GetUserBucketName(u), objs, minio.RemoveObjectsOptions{}) // ForceDelete??? It's not documented
	if err != nil {
		////errorText := "Errors on removing objects:\n"
		for errorPart := range err {
			failedOperations[len(failedOperations)] = models.UnsuccessfulFileOperation{Path: errorPart.ObjectName}
			exactErrors[len(exactErrors)] = errorPart.Err
			////errorText += errorPart.ObjectName + ":" + errorPart.Err.Error() + "; "
		}
		////return errors.New(errorText)
	}
	return failedOperations, exactErrors
}

func UploadFile(u *models.HubUser, file models.UploadFile) (*url.URL, error) { // later think about downloading multiple in an archive
	_, err := connections.Client.StatObject(context.Background(), GetUserBucketName(u), file.Path, minio.GetObjectOptions{})
	if err == nil {
		return nil, errors.New("The file already exists!")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	presignedURL, err := connections.Client.PresignedPutObject(ctx, GetUserBucketName(u), file.Path, configuration.OperationSettings.FileUploadLinksExpiration)
	if err != nil {
		return nil, err
	}

	return presignedURL, nil
}
