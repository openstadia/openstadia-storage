package storage

import (
	"github.com/gorilla/mux"
	//"github.com/openstadia/openstadia-storage/connection/connections/operations/bolt_operations"
	"github.com/openstadia/openstadia-storage/connection/connections/operations/minio_operations"
	"github.com/openstadia/openstadia-storage/connection/routing/api"
	"github.com/openstadia/openstadia-storage/connection/routing/authentification"
	"github.com/openstadia/openstadia-storage/models"
	"log"
	"net/http"
)

func SetUpStorageRouter() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/files_list", authentification.NewEnsureAuth(getFilesList)).Methods("GET")
	r.Handle("/storage_info", authentification.NewEnsureAuth(getOnlyGeneralStorageInfo)).Methods("GET")
	r.Handle("/file_download", authentification.NewEnsureAuth(getFile)).Methods("POST")
	r.Handle("/file", authentification.NewEnsureAuth(deleteFile)).Methods("DELETE")
	r.Handle("/file", authentification.NewEnsureAuth(uploadFile)).Methods("POST")
	//r.Handle("/test", authentification.NewEnsureAuth(test1)).Methods("GET")
	return r
}

func checkStorageExists(w http.ResponseWriter, u *models.HubUser) bool {
	exists, err := minio_operations.BucketExists(u)
	if err != nil {
		log.Println("Error at checking if a bucket exists:\n" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	if !exists {
		w.WriteHeader(http.StatusNotFound)
	}
	return exists
}

func tryToCreateStorage(w http.ResponseWriter, u *models.HubUser) bool {
	err := minio_operations.CreateBucketIfNeeded(u)
	if err != nil {
		log.Println("Error at creating a bucket:\n" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return true
}

func getFilesList(w http.ResponseWriter, _ *http.Request, u *models.HubUser) {
	if !tryToCreateStorage(w, u) {
		return
	}
	list, err := minio_operations.GetUserFilesList(u)
	if err != nil {
		log.Println("Error at GetUserFilesList:\n" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	quota, err := minio_operations.GetBucketQuota(u)
	if err != nil {
		log.Println("Error at GetBucketQuota:\n" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	list.TotalSpace = quota.Quota

	api.EncodeToResponse(w, list)
}

func getOnlyGeneralStorageInfo(w http.ResponseWriter, _ *http.Request, u *models.HubUser) {
	if !tryToCreateStorage(w, u) {
		return
	}
	info, err := minio_operations.GetGeneralStorageInfo(u)
	if err != nil {
		log.Println("Error at getting storage info:\n" + err.Error())
		return
	}

	api.EncodeToResponse(w, info)
}

func getFile(w http.ResponseWriter, r *http.Request, u *models.HubUser) {
	if !tryToCreateStorage(w, u) {
		return
	}
	var file models.GetFile
	if !api.DecodeBody(w, r, &file) {
		return
	}
	url, err := minio_operations.GetFile(u, file)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	api.EncodeToResponse(w, models.GetFileLink{Url: url.String()})
}

func deleteFile(w http.ResponseWriter, r *http.Request, u *models.HubUser) {
	if !checkStorageExists(w, u) {
		return
	}
	var files models.DeleteFile
	if !api.DecodeBody(w, r, &files) {
		return
	}

	// | working way |
	failed, exactErrors := minio_operations.DeleteFiles(u, files)
	if failed != nil && len(failed) != 0 {
		for _, value := range exactErrors {
			log.Println("Error at deleting file(s):\n" + value.Error())
		}
	}

	api.EncodeToResponse(w, models.UnsuccessfulFileOperationsList{List: failed})

	// | broken way |
	//failed, exactErrors := minio_operations.BROKEN_DeleteMultipleFiles(u, files)
	//if failed != nil && len(failed) != 0 {
	//	for _, value := range exactErrors {
	//		log.Println(value.Error())
	//	}
	//	err := json.NewEncoder(w).Encode(failed)
	//	if err != nil {
	//		log.Println("Error at encoding the deleteFile response:\n" + err.Error())
	//		w.WriteHeader(http.StatusInternalServerError)
	//	}
	//	return
	//}

}

func uploadFile(w http.ResponseWriter, r *http.Request, u *models.HubUser) {
	if !checkStorageExists(w, u) {
		return
	}
	var file models.UploadFile
	if !api.DecodeBody(w, r, &file) {
		return
	}
	url, err := minio_operations.UploadFile(u, file)
	if err != nil {
		return
	}
	api.EncodeToResponse(w, models.UploadFileLink{Url: url.String()})
}

//func test1(w http.ResponseWriter, r *http.Request, u *models.HubUser) {
//	bolt_operations.GetUserQuota(u)
//}
