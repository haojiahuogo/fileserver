package models

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/disintegration/imaging"
)

//文件信息
type ImageInfo struct {
	FileID     string
	File       multipart.File
	FileHeader *multipart.FileHeader
	Path       *PathInfo
	FullName   string
	FileExt    string
	FileSize   int64
	zipArgs    string
	ZipImg     []*ZipImage
}

const (
	_IMAGETYPE = "image"
)

//创建文件信息对象
func NewImageInfo(file multipart.File, fileHeader *multipart.FileHeader, zipArgs string) *ImageInfo {
	ii := &ImageInfo{
		FileID:     NewObjectId().Hex(),
		File:       file,
		FileHeader: fileHeader,
		Path:       NewPathInfo(_IMAGETYPE),
		FileSize:   file.(Sizer).Size(),
	}
	ii.FileExt = filepath.Ext(ii.FileHeader.Filename)
	ii.FullName = filepath.Join(ii.Path.FullPath, ii.FileID) + ii.FileExt
	ii.zipArgs = zipArgs
	return ii
}

//保存文件到本地
func (this *ImageInfo) SaveFile() bool {
	beego.Debug("开始写文件")
	//创建文件保存目录
	this.Path.CreateDirectory()

	//保存原始图片
	beego.Debug("创建原始图片:", this.FullName)
	f, err := os.OpenFile(this.FullName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return false
	}
	defer f.Close()
	io.Copy(f, this.File)
	beego.Debug("保存完毕.")

	//保存压缩图
	beego.Debug(this.zipArgs)
	if this.zipArgs != "" {
		this.ZipImg = NewZipImageArray(this.zipArgs, this.Path, this.FileExt)
		beego.Debug("this.ZipImg length:", len(this.ZipImg))
		rawImage, e := imaging.Open(this.FullName)
		if e != nil {
			beego.Debug("原图打开失败:", e.Error())
			beego.Debug("rawImageFullName:", this.FullName)
			return false
		}
		if len(this.ZipImg) > 0 {
			for i := 0; i < len(this.ZipImg); i++ {
				iofile, _ := os.Open(this.FullName)
				defer iofile.Close()
				result := this.ZipImg[i].SaveFile(rawImage, iofile)
				if !result {
					this.ZipImg[i] = nil
				}
			}
		}
		beego.Debug(this.FullName)
	}
	return true

}
