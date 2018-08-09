package internal

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

// 生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetMd5(s string) []byte {
	h := md5.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

// 对字符串进行SHA1哈希
func GetSha1String(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 加密密码明文
func UserPwdEncrypt(password string, salt string) string {
	if salt == "" {
		salt = "znn_"
	}
	return GetSha1String(GetSha1String(string(GetMd5(password))+salt) + GetMd5String(salt))
}

// 生成用户唯一token
func createToken(username string, salt string) string {
	if salt == "" {
		salt = "znn_"
	}
	return GetSha1String(string(GetMd5(username+salt)) + Rand().Hex())
}

// HmacSha1
func HmacSha1ToString(k, v string) string {
	key := []byte(v)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(k))
	return hex.EncodeToString(mac.Sum(nil))
}
