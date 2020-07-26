package models

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/astaxie/beego"
)

//文件信息
type FileInfo struct {
	FileID     string
	File       multipart.File
	FileHeader *multipart.FileHeader
	Path       *PathInfo
	FullName   string
	FileExt    string
	FileSize   int64 //文件大小(字节)
}

const (
	_FILETYPE = "file"
)

// 获取大小
type Sizer interface {
	Size() int64
}

//创建文件对象
func NewFileInfo(file multipart.File, fileHeader *multipart.FileHeader) *FileInfo {
	f := &FileInfo{
		FileID:     NewObjectId().Hex(),
		File:       file,
		FileHeader: fileHeader,
		Path:       NewPathInfo(_FILETYPE),
		FileSize:   file.(Sizer).Size(),
	}
	beego.Debug("上传文件的大小为", file.(Sizer).Size())
	f.FileExt = filepath.Ext(f.FileHeader.Filename)
	f.FullName = filepath.Join(f.Path.FullPath, f.FileID) + f.FileExt
	return f
}

//保存文件
func (this *FileInfo) SaveFile() bool {
	beego.Debug("开始写文件")
	//创建文件保存目录
	this.Path.CreateDirectory()

	//保存原始图片
	beego.Debug("文件路径:", this.FullName)
	f, err := os.OpenFile(this.FullName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return false
	}
	defer f.Close()
	io.Copy(f, this.File)
	beego.Debug("保存完毕.")
	return true

}
