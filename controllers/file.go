package controllers

import (
	"fileserver/models"

	"github.com/astaxie/beego"
)

// 文件相关接口
type FileController struct {
	beego.Controller
}

// @Title 保存文件
// @Summary 保存文件
// @Description 保存文件
// @Param	file		formData 	file	true		"文件"
// @router / [post]
func (this *FileController) Post() {
	file, header, e := this.GetFile("file")
	defer file.Close()
	if e != nil {
		this.Data["json"] = map[string]string{
			"error": e.Error(),
		}
		this.ServeJSON()
		return
	}
	f := models.NewFileInfo(file, header)
	saveOK := f.SaveFile()
	if !saveOK {
		this.Data["json"] = map[string]string{
			"error": "文件保存失败",
		}
		this.ServeJSON()
		return
	}
	result := map[string]models.FileResponse{
		"file": models.FileResponse{
			FileId: f.FileID,
			Uri:    f.Path.UriPath + f.FileID + f.FileExt,
		},
	}

	this.Data["json"] = result
	this.ServeJSON()
}
