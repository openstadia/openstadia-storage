package models

type HubUser struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
	IsActive bool   `json:"is_active"`
	Email    string `json:"email"`
}

type User struct {
	Id       int             `json:"id"`
	UserInfo StorageUserInfo `json:"user_info"`
}

type StorageUserInfo struct {
	TotalSpaceAvailable   uint64 `json:"total_space_available" binding:"required"`
	StorageFeatureAllowed bool   `json:"storage_feature_allowed"`
}
