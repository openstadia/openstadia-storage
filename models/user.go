package models

type HubUser struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
	IsActive bool   `json:"is_active"`
	Email    string `json:"email"`
}

type StorageUserInfo struct {
	AllowedTotalSpace     uint64 `json:"total_space_available" binding:"required"`
	StorageFeatureAllowed bool   `json:"storage_feature_allowed"`
}
