package image

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"treesystem/common"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// PhotoWatermark 图片添加水印并保存到新位置
// @sPath 图片路径
// @waterPath 水印图片路径
// @savePath 存储图片路径
// @quality 图片质量百分比 （值为1~100）
func PhotoWatermark(sPath, waterPath, savePath string, quality int) (bool, error) {
	//打开原始图片
	imgFile, err := os.Open(sPath)
	defer common.DeferFileCols(imgFile)
	if err != nil {
		fmt.Println("打开图片出错")
		fmt.Println(err)
		os.Exit(-1)
		return false, err
	}
	return readerWatermark(imgFile, waterPath, savePath, quality)
}

// BuffPhotoWatermark 图片buff添加水印后保存
// @data buff图片切片
// @waterPath 水印图片路径
// @savePath 存储图片路径
// @quality 图片质量百分比 （值为1~100）
func BuffPhotoWatermark(data []byte, waterPath, savePath string, quality int) (bool, error) {
	//打开原始图片
	imgFile := bytes.NewReader(data)
	return readerWatermark(imgFile, waterPath, savePath, quality)
}

// readerWatermark 将Reader流图片加上水印
// @r Reader流图片
// @waterPath 水印图片路径
// @savePath 存储图片路径
// @quality 图片质量百分比 （值为1~100）
func readerWatermark(r io.Reader, waterPath, savePath string, quality int) (bool, error) {
	img, err := jpeg.Decode(r)
	if err != nil {
		fmt.Println("图片解码时出错")
		fmt.Println(img)
		os.Exit(-1)
		return false, err
	}
	//水印图片
	wmbFile, err := os.Open(waterPath)
	if err != nil {
		fmt.Println("打开水印图片出错")
		fmt.Println(err)
		os.Exit(-1)
		return false, err
	}
	defer common.DeferFileCols(wmbFile)
	wmbImg, err := png.Decode(wmbFile)
	if err != nil {
		fmt.Println("水印图片解码时出错")
		fmt.Println(err)
		os.Exit(-1)
		return false, err
	}
	//把水印写在右下角，并向0坐标偏移10个像素
	offset := image.Pt(img.Bounds().Dx()-wmbImg.Bounds().Dx()-10, img.Bounds().Dy()-wmbImg.Bounds().Dy()-10)
	b := img.Bounds()
	//根据b画布的大小新建一个新图像
	m := image.NewRGBA(b)

	//image.Pt代表Point结构体，目标的源点，即(0,0)
	//draw.Src源图像透过遮罩后，替换掉目标图像
	//draw.Over源图像透过遮罩后，覆盖在目标图像上（类似图层）
	draw.Draw(m, b, img, image.Pt(0, 0), draw.Src)
	draw.Draw(m, wmbImg.Bounds().Add(offset), wmbImg, image.Pt(0, 0), draw.Over)

	dirs := strings.Split(savePath, "\\")
	//获取文件存储的路径
	vPath := strings.Join(dirs[0:len(dirs)-1], "\\")
	//判断路径是否存在
	bl, err := common.PathExists("", vPath)
	if err != nil {
		return false, err
	}
	//路径不存在时先创建路径
	if !bl {
		//创建文件存储路径
		bl, err = common.CreatePathDir(vPath)
		if err != nil {
			return bl, err
		}
	}
	//生成新图片,并设置图片质量
	imgW, err := os.Create(savePath)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	imgOptions := jpeg.Options{Quality: quality}
	err = jpeg.Encode(imgW, m, &imgOptions)
	if err != nil {
		return false, err
	}
	defer common.DeferFileCols(imgW)

	fmt.Println("添加水印图片结束请查看")
	return true, nil
}

// PhotoWatermarkText 图片添加水印后保存
// @data 图片路径
// @waterText 水印图片文字
// @ttfPath 水印字体路径
// @savePath 存储图片路径
// @fontSize 字体大小
// @fontDpi 字体质量
// @quality 图片质量百分比 （值为1~100）
func PhotoWatermarkText(sPath, waterText, ttfPath, savePath string, fontSize float64, fontDpi float64, quality int) (bool, error) {
	//打开原始图片
	imgFile, err := os.Open(sPath)
	defer common.DeferFileCols(imgFile)
	if err != nil {
		fmt.Println("打开图片出错")
		fmt.Println(err)
		os.Exit(-1)
		return false, err
	}
	return readerWatermarkText(imgFile, waterText, ttfPath, savePath, fontSize, fontDpi, quality)
}

// BuffPhotoWatermarkText 图片buff添加水印后保存
// @data buff图片切片
// @waterText 水印图片文字
// @ttfPath 水印字体路径
// @savePath 存储图片路径
// @fontSize 字体大小
// @fontDpi 字体质量
// @quality 图片质量百分比 （值为1~100）
func BuffPhotoWatermarkText(data []byte, waterText, ttfPath, savePath string, fontSize float64, fontDpi float64, quality int) (bool, error) {
	//打开buff流
	imgFile := bytes.NewReader(data)
	return readerWatermarkText(imgFile, waterText, ttfPath, savePath, fontSize, fontDpi, quality)
}

// readerWatermark 将Reader流图片加上水印
// @r Reader流图片
// @waterText 水印图片文字
// @ttfPath 水印字体路径
// @savePath 存储图片路径
// @fontSize 字体大小
// @fontDpi 字体质量
// @quality 图片质量百分比 （值为1~100）
func readerWatermarkText(r io.Reader, waterText, ttfPath, savePath string, fontSize float64, fontDpi float64, quality int) (bool, error) {
	img, err := jpeg.Decode(r)
	if err != nil {
		fmt.Println("图片解码时出错")
		fmt.Println(img)
		os.Exit(-1)
		return false, err
	}
	//水印图片的宽度
	srcWidth := 100
	//水印图片的高度
	srcHeight := 40
	imgText := image.NewRGBA(image.Rect(0, 0, srcWidth, srcHeight))
	//读取字体数据
	fontBytes, err := os.ReadFile(ttfPath)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	//载入字体数据
	fs, err := opentype.Parse(fontBytes)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	mFace, err := opentype.NewFace(fs, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     fontDpi,
		Hinting: font.HintingNone,
	})
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	//把水印写在右下角，并向0坐标偏移10个像素
	offset := image.Pt(img.Bounds().Dx()-imgText.Bounds().Dx()-10, img.Bounds().Dy()-imgText.Bounds().Dy()-10)
	b := img.Bounds()
	//根据b画布的大小新建一个新图像
	m := image.NewRGBA(b)
	//image.Pt代表Point结构体，目标的源点，即(0,0)
	//draw.Src源图像透过遮罩后，替换掉目标图像
	//draw.Over源图像透过遮罩后，覆盖在目标图像上（类似图层）
	draw.Draw(m, b, img, image.Pt(0, 0), draw.Src)
	//draw.Draw(m, wmbImg.Bounds().Add(offset), wmbImg, image.Pt(0, 0), draw.Over)
	photoAddText(imgText, 20, 20, mFace, waterText)
	draw.Draw(m, imgText.Bounds().Add(offset), imgText, image.Pt(0, 0), draw.Over)
	dirs := strings.Split(savePath, "\\")
	//获取文件存储的路径
	vPath := strings.Join(dirs[0:len(dirs)-1], "\\")
	//判断路径是否存在
	bl, err := common.PathExists("", vPath)
	if err != nil {
		return false, err
	}
	//路径不存在时先创建路径
	if !bl {
		//创建文件存储路径
		_, err = common.CreatePathDir(vPath)
		if err != nil {
			return false, err
		}
	}
	//生成新图片,并设置图片质量
	imgW, err := os.Create(savePath)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	imgOptions := jpeg.Options{Quality: quality}
	err = jpeg.Encode(imgW, m, &imgOptions)
	if err != nil {
		return false, err
	}
	defer common.DeferFileCols(imgW)
	return true, nil
}

// photoAddText 图片添加文字
// @img 图像结构
// @x x横轴
// @y y垂直轴
// @mFace 字体
// @label 文字
func photoAddText(img *image.RGBA, x, y int, mFace font.Face, label string) {
	col := color.RGBA{R: 255, A: 255}
	point := fixed.Point26_6{X: fixed.Int26_6(x * 70), Y: fixed.Int26_6(y * 70)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: mFace,
		Dot:  point,
	}
	d.DrawString(label)
}

// Trimming 裁剪图片
func Trimming(srcFilename, trimPath string, rowSize, columnSize int) (bool, error) {
	src, width, height, err := loadImage(srcFilename)
	if err != nil {
		log.Println("load image fail..")
		return false, err
	}
	pWidth := width / rowSize
	pHeight := height / columnSize
	var imgSum [][]image.Image
	for i := 0; i < rowSize; i++ {
		var imgRow []image.Image
		for j := 0; j < columnSize; j++ {
			img, err := photoCopy(src, i*pWidth, j*pHeight, pWidth, pHeight)
			if err != nil {
				log.Println("image copy fail...")
				return false, err
			}
			imgRow = append(imgRow, img)
			//saveErr := saveImage(trimPath+"//"+GetRandomStr(16)+".jpg", img)
			saveErr := saveImage(trimPath+"//"+strconv.Itoa(i)+"_"+strconv.Itoa(j)+".jpg", img)
			if saveErr != nil {
				log.Println("save image fail..")
				return false, err
			}
		}
		imgSum = append(imgSum, imgRow)
	}
	CombinePhoto(imgSum, "C:\\Users\\Blue\\dwhelper\\"+common.GetRandomStr(16)+".jpg", 10)
	return true, nil
}

// CombinePhoto 合并图片（类似九宫格合并）
func CombinePhoto(imgSum [][]image.Image, path string, spiltSize int) (bool, error) {
	rowSize := len(imgSum)
	if rowSize == 0 {
		return false, errors.New("图片数组为空")
	}
	columnSize := len(imgSum[0])
	width := (rowSize + 1) * spiltSize
	height := (columnSize + 1) * spiltSize
	for i := 0; i < rowSize; i++ {
		for j := 0; j < columnSize; j++ {
			if i == 0 {
				reg := imgSum[i][j].Bounds()
				height += reg.Max.Y - reg.Min.Y
			}
			if j == 0 {
				reg := imgSum[i][j].Bounds()
				width += reg.Max.X - reg.Min.X
			}
		}
	}
	//填充白色背景
	dstImg := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < dstImg.Bounds().Dx(); x++ {
		for y := 0; y < dstImg.Bounds().Dy(); y++ {
			dstImg.Set(x, y, color.White)
		}
	}
	pointX := 0
	pointY := 0
	for i := 0; i < rowSize; i++ {
		if i == 0 {
			pointX += spiltSize
		} else {
			imrBound := imgSum[i][0].Bounds()
			pointX += imrBound.Max.X - imrBound.Min.X
		}
		for j := 0; j < columnSize; j++ {
			if j == 0 {
				pointY = spiltSize
			}
			imbBound := imgSum[i][j].Bounds()
			draw.Draw(dstImg, image.Rect(pointX, pointY, pointX+(imbBound.Max.X-imbBound.Min.X), pointY+(imbBound.Max.Y-imbBound.Min.Y)), imgSum[i][j], imbBound.Min, draw.Src)
			pointY += imbBound.Max.Y - imbBound.Min.Y + spiltSize
		}
		pointX += spiltSize
	}
	err := saveImage(path, dstImg)
	if err != nil {
		return false, err
	}
	return true, nil
}

// TrimmingPhoto 裁剪图片
func TrimmingPhoto(srcFilename, trimFilename string, pointX, pointY, width, height int) {
	src, _, _, err := loadImage(srcFilename)
	if err != nil {
		log.Println("load image fail..")
	}
	img, err := photoCopy(src, pointX, pointY, width, height)
	if err != nil {
		log.Println("image copy fail...")
	}
	saveErr := saveImage(trimFilename, img)
	if saveErr != nil {
		log.Println("save image fail..")
	}
}

// loadImage 加载图片
// @path 图片物理路径
// @img 图片结构
// @width 图片宽
// @height 图片高
// @err 错误
func loadImage(path string) (img image.Image, width, height int, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	if err != nil {
		return
	}
	reg := img.Bounds()
	width = reg.Max.X - reg.Min.X
	height = reg.Max.Y - reg.Min.Y
	return
}

func saveImage(p string, src image.Image) error {
	f, err := os.OpenFile(p, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return err
	}
	defer f.Close()
	ext := filepath.Ext(p)

	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {

		err = jpeg.Encode(f, src, &jpeg.Options{Quality: 80})

	} else if strings.EqualFold(ext, ".png") {
		err = png.Encode(f, src)
	} else if strings.EqualFold(ext, ".gif") {
		err = gif.Encode(f, src, &gif.Options{NumColors: 256})
	}
	return err
}

// photoCopy复制图片的指定区域
func photoCopy(src image.Image, pointX, pointY, width, height int) (image.Image, error) {

	var subImg image.Image
	if rgbImg, ok := src.(*image.YCbCr); ok {
		subImg = rgbImg.SubImage(image.Rect(pointX, pointY, pointX+width, pointY+height)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.RGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(pointX, pointY, pointX+width, pointY+height)).(*image.RGBA) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.NRGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(pointX, pointY, pointX+width, pointY+height)).(*image.NRGBA) //图片裁剪x0 y0 x1 y1
	} else {
		return subImg, errors.New("图片解码失败")
	}
	return subImg, nil
}

// ImageCompress 图片压缩，按原始比率压缩时只需指定，width或height之一即可
// @sorucePath 源图片路径
// @savePath 压缩图像存储路径
// @width 图像宽度（像素）
// @height 图像高度（像素）
// @quality 图像质量
func ImageCompress(sorucePath, savePath string, width, height, quality int) (bool, error) {
	//需要压缩
	imgfile, err := os.Open(sorucePath)
	if err != nil {
		return false, err
	}
	defer imgfile.Close()
	jpgimg, _, err := image.Decode(imgfile)
	if err != nil {
		return false, err
	}
	var w, h int
	if width > 0 {
		w = width
	} else if height > 0 {
		w = int(jpgimg.Bounds().Max.X * height / jpgimg.Bounds().Max.Y)
	} else {
		w = 150
	}
	if height > 0 {
		h = height
	} else {
		h = int(jpgimg.Bounds().Max.Y * w / jpgimg.Bounds().Max.X)
	}
	img := imaging.Thumbnail(jpgimg, w, h, imaging.Lanczos)
	//保存到新文件中
	newfile, err := os.Create(savePath)
	if err != nil {
		return false, err
	}
	defer newfile.Close()
	// &jpeg.Options{Quality: 10} 图片压缩质量
	err = jpeg.Encode(newfile, img, &jpeg.Options{Quality: quality})
	if err != nil {
		log.Println("Encode::", err)
		return false, err
	}
	return true, nil
}
