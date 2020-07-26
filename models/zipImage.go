package models

import (
	"encoding/json"
	"image"

	"os"
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/disintegration/imaging"
)

type ZipImage struct {
	Name     string `json:"name"`     //图片名称,如果有返回时使用该名称返回
	Mode     string `json:"mode"`     //处理模式Resize|Fit|Fill
	FillType string `json:"fillType"` //Fill类型
	W        int    `json:"w"`        //宽 1*1为原始尺寸
	H        int    `json:"h"`        //高 1*1为原始尺寸
	Delete   int    `json:"delete"`   //是否删除原图 1删除 0 不删
	FileID   string
	Path     *PathInfo
	FileExt  string
	FullName string
}

func NewZipImageArray(jsonContent string, Path *PathInfo, fileExt string) []*ZipImage {
	zia := []*ZipImage{}
	e := json.Unmarshal([]byte(jsonContent), &zia)
	if e != nil {
		beego.Debug("jsonContent:", jsonContent)
		beego.Debug("zipImage对象创建失败,请检查zipArgs.", e)
		return nil
	}

	for i := 0; i < len(zia); i++ {
		zia[i].FileID = NewObjectId().Hex()
		zia[i].Path = Path
		zia[i].FileExt = fileExt
		//拼接压缩图完成文件名
		zia[i].FullName = filepath.Join(Path.FullPath, zia[i].FileID) + zia[i].FileExt
	}
	return zia
}

//压缩图片处理
func (this *ZipImage) SaveFile(rawImage image.Image, osfile *os.File) bool {
	c, _, err := image.DecodeConfig(osfile)
	if err == nil && this.W == 1 && this.H == 1 {
		this.W = c.Width
		this.H = c.Height
	}
	beego.Debug("压缩图处理开始.")
	var dst *image.NRGBA
	switch this.Mode {
	case "Resize":
		//按固定大小缩放会造成图片变形
		dst = imaging.Resize(rawImage, this.W, this.H, imaging.Lanczos)
	case "Fit":
		//等比例缩放
		dst = imaging.Fit(rawImage, this.W, this.H, imaging.Lanczos)
	case "Fill":
		//按照固定模式裁剪缩放
		var anchor imaging.Anchor
		switch this.FillType {
		case "Center":
			//裁剪中间部分
			anchor = imaging.Center
		case "TopLeft":
			//裁剪左上部分
			anchor = imaging.TopLeft
		case "Top":
			//裁上部分
			anchor = imaging.Top
		case "TopRight":
			//裁剪右上部分
			anchor = imaging.TopRight
		case "Left":
			anchor = imaging.Left
		case "Right":
			anchor = imaging.Right
		case "BottomLeft":
			anchor = imaging.BottomLeft
		case "Bottom":
			anchor = imaging.Bottom
		case "BottomRight":
			anchor = imaging.BottomRight
		default:
			anchor = imaging.Center
		}
		dst = imaging.Fill(rawImage, this.W, this.H, anchor, imaging.Lanczos)
	}
	beego.Debug("缩放完成\n Stride:", dst.Stride, "\n Rect:", dst.Rect, "\n Pix:", len(dst.Pix))
	beego.Debug(this.FullName)
	e := imaging.Save(dst, this.FullName)
	if e != nil {
		beego.Debug("压缩图保存失败:", e.Error())
		beego.Debug("压缩图文件名:", this.FullName)
		return false
	}
	return true
}
