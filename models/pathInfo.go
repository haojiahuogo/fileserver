package models

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

//路径信息
type PathInfo struct {
	_typeName string
	Now       time.Time
	FullPath  string
	UriPath   string
}

//文件基础路径
var basePath = filepath.Join(beego.AppPath, beego.AppConfig.String("filePath"))

//创建目录信息
func NewPathInfo(typeName string) *PathInfo {
	pi := &PathInfo{
		_typeName: typeName,
		Now:       time.Now(),
	}
	pi._CreatePath()
	return pi
}

//创建物理目录
func (this *PathInfo) CreateDirectory() string {
	_CreateDirectory(this.FullPath)
	return this.FullPath

}

func _CreateDirectory(path string) bool {
	exists, _ := _PathExists(path)
	if !exists {
		os.MkdirAll(path, 0777)
	}
	return true
}

//判断目录是否存在
func _PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//创建路径变量
func (this *PathInfo) _CreatePath() {
	this.FullPath = filepath.Join(basePath, this._typeName,
		strconv.Itoa(this.Now.Year()), strconv.Itoa(int(this.Now.Month())), strconv.Itoa(this.Now.Day()),
		strconv.Itoa(this.Now.Hour()), strconv.Itoa(this.Now.Minute()))
	this.UriPath = "/" + this._typeName + "/" +
		strconv.Itoa(this.Now.Year()) + "/" +
		strconv.Itoa(int(this.Now.Month())) + "/" +
		strconv.Itoa(this.Now.Day()) + "/" +
		strconv.Itoa(this.Now.Hour()) + "/" +
		strconv.Itoa(this.Now.Minute()) + "/"
}
