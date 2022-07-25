package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"vote/vote/app/config"
)

//AES加解密

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize//需要padding的数目
	//只要少于256就能放到一个byte中，默认的blockSize=16(即采用16*8=128, AES-128长的密钥)
	//最少填充1个byte，如果原文刚好是blocksize的整数倍，则再填充一个blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)//生成填充的文本
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//加密
func AesEncrypt(data string) (string, error) {
	origData := []byte(data)
	key := []byte(config.AES.AesKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	pass64 := base64.StdEncoding.EncodeToString(crypted)
	return pass64, nil
}

//解密
func AesDecrypt(crypted string) (string, error) {
	bytesPass, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	key := []byte(config.AES.AesKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(bytesPass))
	blockMode.CryptBlocks(origData, bytesPass)
	origData = PKCS5UnPadding(origData)
	return string(origData), nil
}
