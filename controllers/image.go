package controllers

import (
	"fileserver/models"

	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

// 图片相关接口
type ImageController struct {
	beego.Controller
}

const (
	_ImageAllowType = ".jpg|.jpeg|.gif|.png"
)

// @Title 保存图片
// @Description 保存图片同时可生成缩略图和大图 [{"mode":"Fill","w":160,"h":160,"name":"simg"},{"mode":"Fit","w":800,"h":10000,"name":"bimg"}]
// @Summary 保存图片
// @Param	file		formData 	file	true		"文件"
// @Param	zipstr	formData 	string	true		"str"
// @router / [post]
func (this *ImageController) Post() {
	beego.Info("接收到请求:", time.Now().String())
	zipArgs := this.GetString("zipstr", "")
	beego.Debug("zipArgs:", zipArgs)

	file, header, e := this.GetFile("file")
	if e != nil {
		this.Data["json"] = map[string]string{
			"error": e.Error(),
		}
		this.ServeJSON()
		return
	}
	defer file.Close()

	imageInfo := models.NewImageInfo(file, header, zipArgs)
	allow := strings.Index(_ImageAllowType, strings.ToLower(imageInfo.FileExt))
	if allow < 0 {
		this.Data["json"] = map[string]string{
			"error": imageInfo.FileExt + "后缀文件不允许上传!",
		}
		this.ServeJSON()
		return
	}
	beego.Debug(imageInfo.FileSize)
	saveOK := imageInfo.SaveFile()
	if !saveOK {
		this.Data["json"] = map[string]string{
			"error": "图片保存失败",
		}
		this.ServeJSON()
		return
	}
	result := map[string]models.FileResponse{}
	//原始文件
	result["raw"] = models.FileResponse{
		FileId: imageInfo.FileID,
		Uri:    imageInfo.Path.UriPath + imageInfo.FileID + imageInfo.FileExt,
	}
	//返回缩略图和压缩图
	if len(imageInfo.ZipImg) > 0 {
		for i := 0; i < len(imageInfo.ZipImg); i++ {
			num := strconv.Itoa(i)
			if imageInfo.ZipImg[i] == nil {
				result[num] = models.FileResponse{}
				continue
			}
			if imageInfo.ZipImg[i].Name != "" {
				num = imageInfo.ZipImg[i].Name
			}
			result[num] = models.FileResponse{
				FileId: imageInfo.ZipImg[i].FileID,
				Uri:    imageInfo.Path.UriPath + imageInfo.ZipImg[i].FileID + imageInfo.ZipImg[i].FileExt,
			}
			beego.Debug(imageInfo.Path.FullPath + "/" + imageInfo.ZipImg[i].FileID + imageInfo.ZipImg[i].FileExt)
		}
	}

	this.Data["json"] = result
	this.ServeJSON()
}
