package mhash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash/crc32"
	"hash/fnv"

	"github.com/cespare/xxhash/v2"
	"github.com/spaolacci/murmur3"
	"golang.org/x/crypto/sha3"
)

type this struct{}

func (t *this) New() *this {
	return &this{}
}

// MD5 计算字符串的 MD5 哈希值并返回十六进制编码的字符串。
func (t *this) MD5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// SHA1计算字符串的 SHA-1 哈希值并返回十六进制编码的字符串。
func (t *this) SHA1(s string) string {
	hash := sha1.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// SHA3 计算字符串的 SHA3-256 哈希值并返回十六进制编码的字符串。
func (t *this) SHA3(s string) string {
	hash := sha3.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// SHA256 计算字符串的 SHA-256 哈希值并返回十六进制编码的字符串。
func (t *this) SHA256(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// SHA512 计算字符串的 SHA-512 哈希值并返回十六进制编码的字符串。
func (t *this) SHA512(s string) string {
	hash := sha512.Sum512([]byte(s))
	return hex.EncodeToString(hash[:])
}

// FNV 计算字符串的 FNV-1a 哈希值并返回 uint64 类型的哈希值。
func (t *this) FNV(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

// CRC32 计算字符串的 CRC32 哈希值并返回 uint32 类型的哈希值。
func (t *this) CRC32(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}

// Murmur 计算字符串的 MurmurHash3 哈希值并返回 uint32 类型的哈希值。
func (t *this) Murmur(s string) uint32 {
	return murmur3.Sum32([]byte(s))
}

// Jenkins 计算字符串的 Jenkins Hash 哈希值并返回 uint32 类型的哈希值。
// 由于 Go 中没有官方的 Jenkins Hash 实现，这里提供一个简单的实现。
func (t *this) Jenkins(s string) uint32 {
	var hash uint32 = 0
	for i := 0; i < len(s); i++ {
		hash += uint32(s[i])
		hash += (hash << 10)
		hash ^= (hash >> 6)
	}
	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)
	return hash
}

// XXHash64 计算字符串的 XXHash64 哈希值并返回 uint64 类型的哈希值。
func (t *this) XXHash64(s string) uint64 {
	return xxhash.Sum64([]byte(s))
}
