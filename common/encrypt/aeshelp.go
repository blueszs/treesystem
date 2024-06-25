package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
	"runtime"
)

// AesECBDecrypt ECB模式解密
// @encrypted 已加密切片
// @key 解密的key切片
func AesECBDecrypt(encrypted, key []byte) ([]byte, error) {
	if !validKey(key) {
		return nil, fmt.Errorf("秘钥长度错误,当前传入长度为 %d", len(key))
	}
	if len(encrypted) < 1 {
		return nil, fmt.Errorf("源数据长度不能为0")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(encrypted)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("源数据长度必须是 %d 的整数倍，当前长度为：%d", block.BlockSize(), len(encrypted))
	}
	var dst []byte
	tmpData := make([]byte, block.BlockSize())

	for index := 0; index < len(encrypted); index += block.BlockSize() {
		block.Decrypt(tmpData, encrypted[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}
	dst, err = PKCS5UnPadding(dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// AesECBEncrypt ECB模式加密
// @src 待加密切片
// @key 解密的key切片
func AesECBEncrypt(src, key []byte) ([]byte, error) {
	if !validKey(key) {
		return nil, fmt.Errorf("秘钥长度错误, 当前传入长度为 %d", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(src) < 1 {
		return nil, fmt.Errorf("源数据长度不能为0")
	}
	src = PKCS5Padding(src, block.BlockSize())
	if len(src)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("源数据长度必须是 %d 的整数倍，当前长度为：%d", block.BlockSize(), len(src))
	}
	var dst []byte
	tmpData := make([]byte, block.BlockSize())
	for index := 0; index < len(src); index += block.BlockSize() {
		block.Encrypt(tmpData, src[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}
	return dst, nil
}

// AesCbcEncrypt CBC模式加密
// @plainText 待加密字符串
// @key 密钥
// @ivAes 密钥偏移量
func AesCbcEncrypt(plainText, key, ivAes []byte) ([]byte, error) {
	if !validKey(key) {
		return nil, fmt.Errorf("秘钥长度错误,当前传入长度为 %d", len(key))
	}
	if len(plainText) < 1 {
		return nil, fmt.Errorf("源数据长度不能为0")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	paddingText := PKCS5Padding(plainText, block.BlockSize())

	var iv []byte
	if len(ivAes) != 0 {
		if len(ivAes) != 16 {
			return nil, fmt.Errorf("IV长度错误,当前传入长度为 %d", len(key))
		} else {
			iv = ivAes
		}
	} else {
		iv = ivAes
	} // To initialize the vector, it needs to be the same length as block.blockSize
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(paddingText))
	blockMode.CryptBlocks(cipherText, paddingText)
	return cipherText, nil
}

// AesCbcDecrypt CBC模式解密
// @cipherText 待解密数据
// @key 密钥
// @ivAes 密钥偏移量
func AesCbcDecrypt(cipherText, key, ivAes []byte) ([]byte, error) {
	if !validKey(key) {
		return nil, fmt.Errorf("秘钥长度错误,当前传入长度为 %d", len(key))
	}
	if len(cipherText) < 1 {
		return nil, fmt.Errorf("源数据长度不能为0")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Println("runtime err:", err, "Check that the key or text is correct")
			default:
				log.Println("error:", err)
			}
		}
	}()
	var iv []byte
	if len(ivAes) != 0 {
		if len(ivAes) != 16 {
			return nil, fmt.Errorf("IV长度错误,当前传入长度为 %d", len(key))
		} else {
			iv = ivAes
		}
	} else {
		iv = ivAes
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	paddingText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(paddingText, cipherText)

	plainText, err := PKCS5UnPadding(paddingText)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

// PKCS5Padding PKCS5填充
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// PKCS5UnPadding 去除PKCS5填充
func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unPadding := int(origData[length-1])

	if length < unPadding {
		return nil, fmt.Errorf("invalid unpadding length")
	}
	return origData[:(length - unPadding)], nil
}

// 秘钥长度验证
func validKey(key []byte) bool {
	k := len(key)
	switch k {
	default:
		return false
	case 16, 24, 32:
		return true
	}
}
