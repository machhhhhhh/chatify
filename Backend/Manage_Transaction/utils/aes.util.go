package utils

import (
	"bytes"
	"chatify/configs"
	global_types "chatify/types"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func AESEncrypted(payload *global_types.IObjectAES) (string, error) {

	var key []byte = []byte(configs.ENV.AESSetting.AES_KEY)
	var iv []byte = []byte(configs.ENV.AESSetting.AES_IV)

	if len(key) != 32 || len(iv) != aes.BlockSize {
		return "", fmt.Errorf("invalid AES key/IV size: key=%d iv=%d", len(key), len(iv))
	}

	plain, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	var padded []byte = PKCS7Pad(plain) // always appends 1–16 bytes

	blk, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipher_text := make([]byte, len(padded))
	cipher.NewCBCEncrypter(blk, iv).CryptBlocks(cipher_text, padded)

	return base64.StdEncoding.EncodeToString(cipher_text), nil
}

func AESDecrypted(encrypted string) (*global_types.IObjectAES, error) {
	var enc string = strings.ReplaceAll(encrypted, " ", "+")

	cipher_text, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return nil, err
	}

	var key []byte = []byte(configs.ENV.AESSetting.AES_KEY)
	var iv []byte = []byte(configs.ENV.AESSetting.AES_IV)

	if len(key) != 32 || len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("invalid AES key/IV size: key=%d iv=%d", len(key), len(iv))
	}

	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cipher_text)%aes.BlockSize != 0 {
		return nil, errors.New("Cipher text is not multiple of block size")
	}

	cipher.NewCBCDecrypter(blk, iv).CryptBlocks(cipher_text, cipher_text)

	data, err := PKCS7Unpad(cipher_text)
	if err != nil {
		return nil, err
	}

	var payload global_types.IObjectAES
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func PKCS7Pad(src []byte) []byte {
	var block_size int = aes.BlockSize                        // 16
	var pad_length int = block_size - (len(src) % block_size) // yields 1..16
	// if len(src)%bs == 0, padLen == 16
	var pad []byte = bytes.Repeat([]byte{byte(pad_length)}, pad_length)
	return append(src, pad...)
}

// PKCS7Unpad removes padding from decrypted data, verifying correctness.
func PKCS7Unpad(data []byte) ([]byte, error) {
	var num int = len(data)
	if num == 0 {
		return nil, errors.New("pkcs7: data is empty")
	}

	// AES block size (16) — for true PKCS#5 you'd use 8, but AES uses 16.
	var block_size int = aes.BlockSize
	if num%block_size != 0 {
		return nil, fmt.Errorf("pkcs7: data is not a multiple of block size (%d)", block_size)
	}

	// Value of the last byte is the pad length
	var pad_length int = int(data[num-1])
	if pad_length == 0 || pad_length > block_size {
		return nil, fmt.Errorf("pkcs7: invalid padding length %d", pad_length)
	}

	// Verify that each of the final padLen bytes equals padLen
	for i := num - pad_length; i < num; i++ {
		if int(data[i]) != pad_length {
			return nil, fmt.Errorf("pkcs7: invalid padding byte at position %d: got %d, want %d", i, data[i], pad_length)
		}
	}

	// Return data without the padding
	return data[:num-pad_length], nil
}
