package cipherutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
)

//使用PKCS7进行填充，IOS也是7
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)] //FIXME：PANIC runtime error: slice bounds out of range [:-137]
}

//
//func Chacha20Decrypt(encryptData, key, nonce []byte) (decrypted []byte, err error) {
//	if len(encryptData) == 0 || len(key) == 0 {
//		return nil, errors.New("invalid parameters")
//	}
//	newKey := key
//	for len(newKey) < 32 {
//		newKey = append(newKey, key...)
//	}
//	encrypter, err := chacha20.NewUnauthenticatedCipher(newKey[0:32], nonce[0:12])
//	if err != nil {
//		return nil, err
//	}
//
//	decrypted = make([]byte, len(encryptData))
//	encrypter.XORKeyStream(decrypted, encryptData)
//
//	decrypted, err = PKCS7UnPaddingV2(decrypted, aes.BlockSize)
//	return
//}

//aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func AesCBCEncrypt(rawData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//填充原文
	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)
	//加密数据
	cipherText := make([]byte, len(rawData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, rawData)

	return cipherText, nil
}

// level3解密
func AesCBCDecrypt(encryptData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(encryptData) < blockSize {
		return nil, errors.New(fmt.Sprintf("encryptData is too short, len:%d", len(encryptData)))
	}

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		return nil, errors.New(fmt.Sprintf("ciphertext is not a multiple of the block size,len=%d", len(encryptData)))
	}

	// 解密数据
	mode := cipher.NewCBCDecrypter(block, iv)
	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData = PKCS7UnPadding(encryptData)
	return encryptData, nil
}

//aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func AesGCMEncrypt(rawData, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//填充原文
	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)
	//加密数据
	aesgcm, e := cipher.NewGCM(block)
	if e != nil {
		return nil, e
	}
	ct := aesgcm.Seal(nil, nonce, rawData, nil)

	return ct, nil
}

// level4 解密
func AesGCMDecrypt(encryptData, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()

	if len(encryptData) < blockSize {
		return nil, errors.New(fmt.Sprintf("encryptData is too short, len:%d", len(encryptData)))
	}
	if len(encryptData)%blockSize != 0 {
		return nil, errors.New(fmt.Sprintf("ciphertext is not a multiple of the block size,len=%d", len(encryptData)))
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plantText, err := aesgcm.Open(nil, nonce, encryptData, nil)

	//解填充
	plantText = PKCS7UnPadding(plantText)
	return plantText, nil
}

// level4 分级解密
func AesGCMDecryptNew(encryptData, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()

	if len(encryptData) < blockSize {
		return nil, errors.New(fmt.Sprintf("encryptData is too short, len:%d", len(encryptData)))
	}
	if len(encryptData)%blockSize != 0 {
		return nil, errors.New(fmt.Sprintf("ciphertext is not a multiple of the block size,len=%d", len(encryptData)))
	}

	aesgcm, err := cipher.NewGCMWithNonceSize(block, len(nonce))
	if err != nil {
		return nil, err
	}
	plantText, err := aesgcm.Open(nil, nonce, encryptData, nil)

	//解填充
	plantText = PKCS7UnPadding(plantText)
	return plantText, nil
}

// 流模式 加密/解密
const (
	ENCRYPT = iota
	DECRYPT
)

type AesStream struct {
	Key     []byte //加密的 key
	Iv      []byte //加密的向量
	Model   int    //ENCRYPT 加密    DECRYPT 解密
	context []byte //处理之后的 明文/密文
	cur     int    //已经读到的数据长度
}

// 写数据
func (a *AesStream) Write(src []byte) (n int, err error) {
	if a.Model == ENCRYPT {
		a.context, _ = AesCBCEncrypt(src, a.Key, a.Iv)
	} else if a.Model == DECRYPT {
		a.context, _ = AesCBCDecrypt(src, a.Key, a.Iv)
	} else {
		return 0, errors.New("模式未知")
	}

	return len(src), nil
}

// 读数据
func (a *AesStream) Read(dst []byte) (n int, err error) {
	allLen := len(a.context)
	wantedLen := len(dst)
	// 已经读的长度,大于等于总长度
	if a.cur >= allLen {
		return 0, io.EOF
	}
	// 还可以读出来的长度
	canRead := allLen - a.cur
	if wantedLen >= canRead {
		// 要读的大于等于可以读的，返回所有可以读的数据
		copy(dst, a.context[a.cur:])
		return canRead, io.EOF
	} else {
		// 要读的数据小于可以读的
		copy(dst, a.context[a.cur:a.cur+wantedLen])
		// 当前读便宜到新的位置
		a.cur += wantedLen
		return wantedLen, nil
	}
}

func PKCS7UnPaddingV2(plantText []byte, blockSize int) ([]byte, error) {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	if unpadding > 16 || unpadding < 0 || unpadding >= length {
		return nil, errors.New("Invalid unpadding")
	}
	return plantText[:(length - unpadding)], nil
}
