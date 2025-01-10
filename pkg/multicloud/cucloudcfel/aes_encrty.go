package cucloudcfel

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

// d 函数，假设使用 sha256 并返回十六进制字符串
func md5Str(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// PKCS7Padding 实现 PKCS#7 填充
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

// encryptPwd 使用 AES-CBC 加密
func encryptPwd(username, password, timestamp string) (string, error) {
	
	// 转换密钥和 IV 为二进制数据
	key := []byte(md5Str(username)[:32]) // 使用 SHA-256 的前 32 字符（AES-256 密钥）
	// iv := []byte(o)       // 16 字节 IV

	// 创建 AES 加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 确保明文符合块大小
	plaintext := PKCS7Padding([]byte(password), block.BlockSize())

	// 加密
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, []byte(md5Str(timestamp)[2:18]))
	mode.CryptBlocks(ciphertext, plaintext)

	// 返回 Base64 编码的加密结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
