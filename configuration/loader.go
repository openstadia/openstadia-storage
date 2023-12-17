package configuration

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/openstadia/openstadia-storage/models"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func tryToLoadValue(key string, target *string) error {
	value := os.Getenv(key)
	if value == "" {
		return fmt.Errorf("config parameter not set: " + key)
	}
	*target = value
	return nil
}

func tryToLoadBoolValue(key string, target *bool) error {
	var boolStr string
	err := tryToLoadValue(key, &boolStr)
	if err != nil {
		return err
	}
	if strings.ToLower(boolStr) == "true" {
		*target = true
	} else if strings.ToLower(boolStr) == "false" {
		*target = false
	} else {
		return fmt.Errorf("illegal bool value at the config key: %w", err)
	}
	return nil
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

func Load() (*ConfigStore, error) {
	config := models.Configuration{}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("couldn't load the .env file, check if it exists!")
	}

	err = tryToLoadValue("HUB_DOMAIN", &config.HubSettings.HubDomain)
	if err != nil {
		return nil, err
	}
	err = tryToLoadValue("MINIO_DOMAIN", &config.MinioSettings.MinioDomain)
	if err != nil {
		return nil, err
	}
	err = tryToLoadValue("MINIO_ACCESS_KEY", &config.MinioSettings.MinioAccessKey)
	if err != nil {
		return nil, err
	}
	err = tryToLoadValue("MINIO_SECRET_KEY", &config.MinioSettings.MinioSecretKey)
	if err != nil {
		return nil, err
	}
	err = tryToLoadBoolValue("MINIO_SSL", &config.MinioSettings.MinioSsl)
	if err != nil {
		return nil, err
	}
	tryToLoadDurationWithDefault("STORAGE_FILE_DOWNLOAD_EXPIRATION_TIME", &config.OperationRelated.FileDownloadLinksExpiration, 24*60*60) // seconds
	tryToLoadDurationWithDefault("STORAGE_FILE_UPLOAD_EXPIRATION_TIME", &config.OperationRelated.FileUploadLinksExpiration, 24*60*60)     // seconds

	store := CreateConfigStore(&config)
	return &store, nil
}
