package image

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// CreateQR 创建二维码保持在指定位置返回是否成功，失败时返回错误对象
// @qrData 二维码原始数据
// @savePath 二维码存储地址
func CreateQR(qrData, savePath string) (bool, error) {
	// Create the barcode
	qrCode, err := qr.Encode(qrData, qr.M, qr.Auto)
	if err != nil {
		return false, err
	}
	// Scale the barcode to 200x200 pixels
	qrCode, err = barcode.Scale(qrCode, 200, 200)
	if err != nil {
		return false, err
	}
	// create the output file
	file, err := os.Create(savePath)
	if err != nil {
		return false, err
	}
	disclose := func() {
		er := file.Close()
		println(er.Error())
	}
	defer disclose()
	// encode the barcode as png
	err = png.Encode(file, qrCode)
	if err != nil {
		return false, err
	}
	return true, err
}

// CreateQRBytes 创建二维码返回buff数据，失败时返回错误对象
// @qrData 二维码原始数据
func CreateQRBytes(qrData string) ([]byte, error) {
	// Create the barcode
	qrCode, err := qr.Encode(qrData, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}
	// Scale the barcode to 200x200 pixels
	qrCode, err = barcode.Scale(qrCode, 200, 200)
	if err != nil {
		return nil, err
	}
	//file, err := os.Create(savePath)
	data, err := ImageToBuff(qrCode)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateQRBase64 创建二维码并返回Base64字符串，失败时返回错误对象
// @qrData 二维码原始数据
func CreateQRBase64(qrData string) (string, error) {
	// Create the barcode
	qrCode, err := qr.Encode(qrData, qr.M, qr.Auto)
	if err != nil {
		return "", err
	}
	// Scale the barcode to 200x200 pixels
	qrCode, err = barcode.Scale(qrCode, 200, 200)
	if err != nil {
		return "", err
	}
	//file, err := os.Create(savePath)
	data, err := ImageToBuff(qrCode)
	if err != nil {
		return "", err
	}
	//开辟存储空间
	dist := make([]byte, 50000)
	base64.StdEncoding.Encode(dist, data) //buff转成base64
	//这里要注意，因为申请的固定长度数组，所以没有被填充完的部分需要去掉，负责输出可能出错
	index := bytes.IndexByte(dist, 0)
	baseImage := dist[0:index]
	return string(baseImage), nil
}

// ImageToBuff 将image对象生成buff数据
func ImageToBuff(img image.Image) ([]byte, error) {
	//开辟一个新的空buff
	emptyBuff := bytes.NewBuffer(nil)
	////img写入到buff
	err := jpeg.Encode(emptyBuff, img, nil)
	if err != nil {
		return nil, err
	}
	baseImage := emptyBuff.Bytes()
	return baseImage, nil
}
