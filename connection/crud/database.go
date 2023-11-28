package crud

import "github.com/openstadia/openstadia-storage/models"

func GetStorageUserInfo(u *models.HubUser) models.StorageUserInfo {
	// placeholder
	return models.StorageUserInfo{AllowedTotalSpace: 1024 * 1024 * 1024, StorageFeatureAllowed: true} // 1GB
}
