package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"fmt"
)

const (
	DES = iota
	AES
)

// Encrypt
func Encrypt(src []byte, algo int, key []byte) ([]byte, error) {
	var block cipher.Block
	var err error
	switch algo {
	case AES:
		block, err = aes.NewCipher(key)
	case DES:
		block, err = des.NewCipher(key)
	default:
		return nil, fmt.Errorf("unsupported encrypt algo %d", algo)
	}
	// Handle errors uniformly
	if err != nil {
		fmt.Printf("Encryption failed, err: %v\n", err)
		return nil, err
	}
	// Padding
	pkcs, err := NewPkcs7(block.BlockSize())
	if err != nil {
		fmt.Printf("Encryption failed, err: %v\n", err)
		return nil, err
	}
	src, err = pkcs.Padding(src)
	if err != nil {
		fmt.Printf("Encryption failed, err: %v\n", err)
		return nil, err
	}
	encrypting := cipher.NewCBCEncrypter(block, key) // Use CBC mode for encryption
	dest := make([]byte, len(src))                   // Cipher text has the same length as plain text
	encrypting.CryptBlocks(dest, src)
	return dest, nil
}

// Decrypt
func Decrypt(in []byte, algo int, key []byte) ([]byte, error) {
	var block cipher.Block
	var err error
	switch algo {
	case AES:
		block, err = aes.NewCipher(key)
	case DES:
		block, err = des.NewCipher(key)
	default:
		return nil, fmt.Errorf("unsupported encrypt algo %d", algo)
	}
	// Handle errors uniformly
	if err != nil {
		return nil, err
	}
	// Decryption
	decrypting := cipher.NewCBCDecrypter(block, key) // Use CBC mode for decryption
	transition := make([]byte, len(in))
	decrypting.CryptBlocks(transition, in)
	// Remove padding
	pkcs, err := NewPkcs7(block.BlockSize())
	if err != nil {
		return nil, err
	}
	out, err := pkcs.UnPadding(transition)
	if err != nil {
		return nil, err
	}
	return out, nil
}
