package crud

import (
	"github.com/openstadia/openstadia-storage/models"
)

func GetStorageUserInfo(_ *models.User) models.StorageUserInfo {
	// placeholder
	return models.StorageUserInfo{TotalSpaceAvailable: 0, StorageFeatureAllowed: false}
}
