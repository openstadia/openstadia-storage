package models

type DeleteFiles struct {
	Paths []string `json:"path"`
}

type GetFile struct {
	Path string `json:"path"`
}

type UnsuccessfulFileOperation struct {
	Path  string `json:"path"`
	Error error  `json:"error"`
}

type UnsuccessfulFileOperationsList struct {
	List []UnsuccessfulFileOperation `json:"unsuccessful"`
}

type GetFileLink struct {
	Url string `json:"url"`
}

type UploadFile struct {
	Path string `json:"path"`
}

type UploadFileLink struct {
	Url string `json:"url"`
}
