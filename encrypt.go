package yee

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

type AesEncrypt struct {
	Key string
	Iv  string
}

func (crypt *AesEncrypt) getKey() []byte {

	keyLen := len(crypt.Key)
	if keyLen < 16 {
		panic("The aes key length must not be less than 16")
	}
	arrKey := []byte(crypt.Key)
	if keyLen >= 32 {
		return arrKey[:32]
	}
	if keyLen >= 24 {
		return arrKey[:24]
	}
	return arrKey[:16]
}


func (crypt *AesEncrypt) Encrypt(plantText []byte) ([]byte, error) {
	key := crypt.getKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plantText = crypt.PKCS7Padding(plantText, block.BlockSize())

	blockModel := cipher.NewCBCEncrypter(block, []byte(crypt.Iv)[:aes.BlockSize])

	ciphertext := make([]byte, len(plantText))

	blockModel.CryptBlocks(ciphertext, plantText)
	return ciphertext, nil
}


func (crypt *AesEncrypt) Decrypt(src []byte) (strDesc string, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	key := crypt.getKey()
	keyBytes := key
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	blockModel := cipher.NewCBCDecrypter(block, []byte(crypt.Iv)[:aes.BlockSize])
	plantText := make([]byte, len(src))
	blockModel.CryptBlocks(plantText, src)
	plantText = crypt.PKCS7UnPadding(plantText)
	return string(plantText), nil
}


func (crypt *AesEncrypt) PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

func (crypt *AesEncrypt) PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (crypt *AesEncrypt) EnPwdCode(pwdStr string) string {
	pwd := []byte(pwdStr)
	result, err := crypt.Encrypt(pwd)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(result)
}

func (crypt *AesEncrypt) DePwdCode(pwd string) string {
	temp, _ := hex.DecodeString(pwd)
	res, _ := crypt.Decrypt(temp)
	return res
}
