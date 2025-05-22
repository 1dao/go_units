package share

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

type Crypto struct {
	key    []byte
	method uint8
}

const (
	EncryptNone = 0x00
	EncryptXOR  = 0x01
	EncryptAES  = 0x02
)

func NewCrypto(key []byte, method uint8) *Crypto {
	return &Crypto{key: key, method: method}
}

func (c *Crypto) Encrypt(data []byte) ([]byte, error) {
	switch c.method {
	case EncryptNone:
		return data, nil
	case EncryptXOR:
		return c.xorEncrypt(data), nil
	case EncryptAES:
		return c.aesEncrypt(data)
	default:
		return nil, errors.New("unknown encryption method")
	}
}

func (c *Crypto) Decrypt(data []byte) ([]byte, error) {
	switch c.method {
	case EncryptNone:
		return data, nil
	case EncryptXOR:
		return c.xorDecrypt(data), nil
	case EncryptAES:
		return c.aesDecrypt(data)
	default:
		return nil, errors.New("unknown decryption method")
	}
}

func (c *Crypto) xorEncrypt(data []byte) []byte {
	encrypted := make([]byte, len(data))
	keyLen := len(c.key)
	for i := range data {
		encrypted[i] = data[i] ^ c.key[i%keyLen]
	}
	return encrypted
}

func (c *Crypto) xorDecrypt(data []byte) []byte {
	// XOR decryption is identical to XOR encryption
	return c.xorEncrypt(data)
}

func (c *Crypto) aesEncrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)

	// 打包格式：nonce(12字节) + ciphertext
	result := make([]byte, 12+len(ciphertext))
	copy(result[:12], nonce)
	copy(result[12:], ciphertext)
	return result, nil
}

func (c *Crypto) aesDecrypt(data []byte) ([]byte, error) {
	if len(data) < 12 {
		return nil, errors.New("invalid cipher text")
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := data[:12]
	ciphertext := data[12:]
	return aesgcm.Open(nil, nonce, ciphertext, nil)
}
