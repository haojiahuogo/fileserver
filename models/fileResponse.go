package models

//返回的结构
type FileResponse struct {
	FileId string `json:"fileid,omitempty" description:"图片名"`
	Uri    string `json:"uri,omitempty" description:"地址"`
}
