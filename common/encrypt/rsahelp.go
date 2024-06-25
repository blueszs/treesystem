package encrypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"treesystem/common"
)

// GenerateRSAKey 生成RSA私钥和公钥，保存到文件中
// @bits 密钥位数（1024、2048等）
// @savePath 密钥保存路径
// @fileName 密钥前缀名称
func GenerateRSAKey(bits int, savePath, fileName string) error {
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//使用pem格式对x509输出的内容进行编码
	//创建文件保存私钥
	privateKeyPath := fmt.Sprintf("%s\\%s_private.pem", savePath, fileName)
	privateFile, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}
	defer common.DeferFileCols(privateFile)
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	//将数据保存到文件
	err = pem.Encode(privateFile, &privateBlock)
	if err != nil {
		return err
	}
	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}
	//pem格式编码
	//创建用于保存公钥的文件
	publicKeyPath := fmt.Sprintf("%s\\%s_public.pem", savePath, fileName)
	publicFile, err := os.Create(publicKeyPath)
	if err != nil {
		return err
	}
	defer common.DeferFileCols(publicFile)
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	//保存到文件
	err = pem.Encode(publicFile, &publicBlock)
	return err
}

// RSAEncrypt RSA加密
// @plainText 待加密数据
// @publicKeyPath 公钥路径
func RSAEncrypt(plainText []byte, publicKeyPath string) ([]byte, error) {
	//打开文件
	file, err := os.Open(publicKeyPath)
	if err != nil {
		return nil, err
	}
	defer common.DeferFileCols(file)
	//读取文件的内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}
	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		return nil, err
	}
	//返回密文
	return cipherText, nil
}

// RSADecrypt RSA解密
// @cipherText 待解密数据
// @privateKeyPath 私钥路径
func RSADecrypt(cipherText []byte, privateKeyPath string) ([]byte, error) {
	//打开文件
	file, err := os.Open(privateKeyPath)
	if err != nil {
		return nil, err
	}
	defer common.DeferFileCols(file)
	//获取文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}
	//pem解码
	block, _ := pem.Decode(buf)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//对密文进行解密
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		return nil, err
	}
	//返回明文
	return plainText, nil
}

// RsaEncryptWithSha1Base64 加密：采用sha1算法加密后转base64格式
// @originalData 待加密数据
// @publicKey 公钥Base64字符串
func RsaEncryptWithSha1Base64(originalData, publicKey string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", err
	}
	pubKey, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), []byte(originalData))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// RsaDecryptWithSha1Base64 解密：对采用sha1算法加密后转base64格式的数据进行解密（私钥PKCS1格式）
// @encryptedData 待解密数据
// @privateKey 私钥Base64字符串
func RsaDecryptWithSha1Base64(encryptedData, privateKey string) (string, error) {
	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	key, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	prvKey, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		return "", err
	}
	originalData, err := rsa.DecryptPKCS1v15(rand.Reader, prvKey, encryptedDecodeBytes)
	if err != nil {
		return "", err
	}
	return string(originalData), err
}

// RsaSignWithSha1Hex 签名：采用sha1算法进行签名并输出为hex格式（私钥PKCS8格式）
// @data 待签名数据
// @prvKey 私钥
// @shaV 签名方式（crypto.SHA1、crypto.SHA256）
func RsaSignWithSha1Hex(data string, prvKey string, shaV crypto.Hash) (string, error) {
	keyBytes, err := hex.DecodeString(prvKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(keyBytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		return "", err
	}
	h := sha1.New()
	_, err = h.Write([]byte(data))
	if err != nil {
		fmt.Println("err", err)
		return "", err
	}
	hash := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), shaV, hash[:])
	if err != nil {
		fmt.Printf("Error from signing: %v\n", err)
		return "", err
	}
	out := hex.EncodeToString(signature)
	return out, nil
}

// RsaVerySignWithSha1Base64 验签：对采用sha1算法进行签名后转base64格式的数据进行验签
// @originalData 待验证签名数据
// @signData 签名
// @pubKey 公钥
// @shaV 签名方式（crypto.SHA1、crypto.SHA256）
func RsaVerySignWithSha1Base64(originalData, signData, pubKey string, shaV crypto.Hash) (bool, error) {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return false, err
	}
	public, _ := base64.StdEncoding.DecodeString(pubKey)
	pub, err := x509.ParsePKIXPublicKey(public)
	if err != nil {
		return false, err
	}
	hash := sha1.New()
	_, err = hash.Write([]byte(originalData))
	if err != nil {
		return false, err
	}
	err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), shaV, hash.Sum(nil), sign)
	if err != nil {
		return false, err
	}
	return true, nil
}
