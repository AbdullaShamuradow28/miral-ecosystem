package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

// GenerateAESKey генерирует случайный 32-байтовый (256-битный) ключ AES.
func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32) // 256-bit key
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("ошибка генерации ключа AES: %w", err)
	}
	return key, nil
}

// GenerateIV генерирует случайный 16-байтовый (128-битный) вектор инициализации (IV).
func GenerateIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize) // 16 bytes for AES-CBC
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("ошибка генерации IV: %w", err)
	}
	return iv, nil
}

// Encrypt encrypts plaintext using AES-CBC with the given key and IV.
// It applies PKCS7 padding.
func Encrypt(plaintext []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания шифра AES: %w", err)
	}

	paddedPlaintext := pkcs7Pad(plaintext, aes.BlockSize)

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-CBC with the given key and IV.
// It removes PKCS7 padding.
func Decrypt(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания шифра AES: %w", err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("размер зашифрованных данных не кратен размеру блока AES")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	paddedPlaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(paddedPlaintext, ciphertext)

	plaintext, err := pkcs7Unpad(paddedPlaintext, aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("ошибка удаления PKCS7 отступа: %w", err)
	}

	return plaintext, nil
}

// pkcs7Pad pads the data to a multiple of blockSize.
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7Unpad unpads the data.
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("невозможно распаковать пустые данные")
	}
	if len(data)%blockSize != 0 {
		return nil, fmt.Errorf("размер данных не кратен размеру блока")
	}
	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize || len(data) < padding {
		return nil, fmt.Errorf("недопустимое значение отступа PKCS7")
	}
	// Проверяем, что все символы отступа одинаковы
	for i := 0; i < padding; i++ {
		if data[len(data)-1-i] != byte(padding) {
			return nil, fmt.Errorf("недопустимый отступ PKCS7: несовпадающие значения")
		}
	}
	return data[:len(data)-padding], nil
}