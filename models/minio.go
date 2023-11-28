package models

import "time"

type StorageObject struct {
	Name       string    `json:"name"`
	Size       uint64    `json:"sizeTotal"`
	TimeEdited time.Time `json:"time_edited"`
}

type ObjectsList struct {
	// this would always be called with the general info request nearby, so it would retrieve and iterate files twice
	UsedSpace  uint64          `json:"used_space"`
	TotalSpace uint64          `json:"total_space"`
	Objects    []StorageObject `json:"objects"`
}

type StorageInfo struct {
	UsedSpace  uint64 `json:"used_space"`
	TotalSpace uint64 `json:"total_space"`
}
