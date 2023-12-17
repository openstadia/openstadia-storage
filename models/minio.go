package models

import "time"

type StorageObject struct {
	Name       string    `json:"name"`
	Size       uint64    `json:"size_total"`
	TimeEdited time.Time `json:"time_edited"`
}

type ObjectsList struct {
	Objects []StorageObject `json:"objects"`
}

type StorageInfo struct {
	UsedSpace  uint64 `json:"used_space"`
	TotalSpace uint64 `json:"total_space"`
}
