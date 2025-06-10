package custom

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// 必须是16、24或32字节
var Key = []byte("UneLqzFTLwjoch8lYJybbMU4urWbhzm6")

func Encode(key, data []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 创建随机数作为nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// 3. 将加密结果进行base64编码
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return []byte(encoded), nil
}

func Decode(data, key string) (string, error) {
	// 1. 对base64编码的数据进行解码
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.New("invalid base64 format")
	}

	// 2. 使用AES-256-GCM进行解密
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 从密文中提取nonce
	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
