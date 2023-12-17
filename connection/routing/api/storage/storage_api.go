package storage

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/openstadia/openstadia-storage/configuration"
	"github.com/openstadia/openstadia-storage/connection/connections"
	"github.com/openstadia/openstadia-storage/connection/connections/operations"
	"github.com/openstadia/openstadia-storage/connection/routing/authentification"
	"github.com/openstadia/openstadia-storage/models"
	"github.com/openstadia/openstadia-storage/utils"
	"log"
	"net/http"
)

func SetUpStorageRouter(configStore *configuration.ConfigStore, connectionsStore *connections.ConnectionsStore) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/files_list", authentification.NewEnsureAuth(getFilesList, configStore, connectionsStore)).Methods("GET")
	r.Handle("/storage_info", authentification.NewEnsureAuth(getBucketStorageInfo, configStore, connectionsStore)).Methods("GET")
	r.Handle("/file_download", authentification.NewEnsureAuth(getFile, configStore, connectionsStore)).Methods("POST")
	r.Handle("/file", authentification.NewEnsureAuth(deleteFile, configStore, connectionsStore)).Methods("DELETE")
	r.Handle("/file", authentification.NewEnsureAuth(uploadFile, configStore, connectionsStore)).Methods("POST")
	return r
}

func checkStorageExists(w http.ResponseWriter, u *models.User, connectionsStore *connections.ConnectionsStore) bool {
	exists, err := operations.BucketExists(u, connectionsStore)
	if err != nil {
		log.Println(fmt.Errorf("failed to checking if a bucket exists: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	if !exists {
		w.WriteHeader(http.StatusNotFound)
	}
	return exists
}

func tryToCreateStorage(w http.ResponseWriter, u *models.User, connectionsStore *connections.ConnectionsStore) bool {
	err := operations.CreateBucketAndQuotaIfNeeded(u, connectionsStore)
	if err != nil {
		log.Println(fmt.Errorf("failed to create a bucket: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return true
}

func getFilesList(w http.ResponseWriter, _ *http.Request, u *models.User, _ *configuration.ConfigStore, connectionsStore *connections.ConnectionsStore) {
	if !tryToCreateStorage(w, u, connectionsStore) {
		return
	}
	list, err := operations.GetUserFilesList(u, connectionsStore)
	if err != nil {
		log.Println(fmt.Errorf("failed to get files list: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.EncodeToResponse(w, list)
}

func getBucketStorageInfo(w http.ResponseWriter, _ *http.Request, u *models.User, _ *configuration.ConfigStore, connectionsStore *connections.ConnectionsStore) {
	if !tryToCreateStorage(w, u, connectionsStore) {
		return
	}
	info, err := operations.GetBucketStorageInfo(u, connectionsStore)
	if err != nil {
		log.Println(fmt.Errorf("failed to get storage info: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.EncodeToResponse(w, info)
}

func getFile(w http.ResponseWriter, r *http.Request, u *models.User, configStore *configuration.ConfigStore, connectionsStore *connections.ConnectionsStore) {
	if !tryToCreateStorage(w, u, connectionsStore) {
		return
	}
	var file models.GetFile
	if !utils.DecodeBody(w, r, &file) {
		return
	}
	url, err := operations.GetFile(u, file, configStore, connectionsStore)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	utils.EncodeToResponse(w, models.GetFileLink{Url: url.String()})
}

func deleteFile(w http.ResponseWriter, r *http.Request, u *models.User, _ *configuration.ConfigStore, connectionsStore *connections.ConnectionsStore) {
	if !checkStorageExists(w, u, connectionsStore) {
		return
	}
	var files models.DeleteFiles
	if !utils.DecodeBody(w, r, &files) {
		return
	}

	failed := operations.DeleteFiles(u, connectionsStore, files)
	if failed.List != nil && len(failed.List) != 0 {
		for _, failedOperation := range failed.List {
			log.Println(fmt.Errorf("failed to delete file \"%s\": %w", failedOperation.Path, failedOperation.Error))
		}
	}

	utils.EncodeToResponse(w, failed)
}

func uploadFile(w http.ResponseWriter, r *http.Request, u *models.User, configStore *configuration.ConfigStore, connectionsStore *connections.ConnectionsStore) {
	if !checkStorageExists(w, u, connectionsStore) {
		return
	}
	var file models.UploadFile
	if !utils.DecodeBody(w, r, &file) {
		return
	}
	url, err := operations.UploadFile(u, file, configStore, connectionsStore)
	if err != nil {
		return
	}
	utils.EncodeToResponse(w, models.UploadFileLink{Url: url.String()})
}
