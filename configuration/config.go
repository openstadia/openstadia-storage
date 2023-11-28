package configuration

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type TypeStorageSettings struct {
	HubDomain      string
	MinioDomain    string
	MinioAccessKey string
	MinioSecretKey string
	MinioSsl       bool
}

type TypeOperationSettings struct {
	FileDownloadLinksExpiration time.Duration
	FileUploadLinksExpiration   time.Duration
}

var Settings TypeStorageSettings
var OperationSettings TypeOperationSettings

func tryToLoadValue(key string, target *string, printWarnings bool) bool {
	value := os.Getenv(key)
	if value == "" {
		if printWarnings {
			print("Config parameter '" + key + "' is not set!")
		}
		return false
	}
	*target = value
	return true
}

func tryToLoadWithDefault(key string, target *string, def string) bool {
	value := os.Getenv(key)
	if value == "" {
		*target = def
	}
	*target = value
	return true
}

func tryToLoadIntWithDefault(key string, target *int, def int) {
	value := os.Getenv(key)
	if value == "" {
		*target = def
	}
	val, err := strconv.Atoi(value)
	if err != nil {
		return
	}
	*target = val
}

func tryToLoadDurationWithDefault(key string, target *time.Duration, def time.Duration) {
	value := os.Getenv(key)
	if value == "" {
		*target = def
	}
	val, err := strconv.Atoi(value)
	if err != nil {
		return
	}
	*target = time.Second * time.Duration(val)
}

func tryToLoadBoolValue(key string, target *bool, printWarnings bool) bool {
	var boolStr string
	if !tryToLoadValue(key, &boolStr, printWarnings) {
		return false
	}
	if strings.ToLower(boolStr) == "true" {
		*target = true
	} else if strings.ToLower(boolStr) == "false" {
		*target = false
	} else {
		return false
	}
	return true
}

func Init() bool {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, check if it exists!")
	}

	if !tryToLoadValue("HUB_DOMAIN", &Settings.HubDomain, true) {
		return false
	}
	if !tryToLoadValue("MINIO_DOMAIN", &Settings.MinioDomain, true) {
		return false
	}
	if !tryToLoadValue("MINIO_ACCESS_KEY", &Settings.MinioAccessKey, true) {
		return false
	}
	if !tryToLoadValue("MINIO_SECRET_KEY", &Settings.MinioSecretKey, true) {
		return false
	}
	if !tryToLoadBoolValue("MINIO_SSL", &Settings.MinioSsl, true) {
		return false
	}
	tryToLoadDurationWithDefault("STORAGE_FILE_DOWNLOAD_EXPIRATION_TIME", &OperationSettings.FileDownloadLinksExpiration, 24*60*60) // seconds
	tryToLoadDurationWithDefault("STORAGE_FILE_UPLOAD_EXPIRATION_TIME", &OperationSettings.FileUploadLinksExpiration, 24*60*60)     // seconds
	return true
}
