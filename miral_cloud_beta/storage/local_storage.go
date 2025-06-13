package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"miral_cloud_go/utils" // Используем константы для лимитов
)

const StorageBaseDir = "./miral_cloud_files"

func init() {
	// Убедимся, что директория для хранения файлов существует
	if _, err := os.Stat(StorageBaseDir); os.IsNotExist(err) {
		err := os.Mkdir(StorageBaseDir, 0755)
		if err != nil {
			panic(fmt.Sprintf("Не удалось создать директорию для хранения файлов: %v", err))
		}
		fmt.Printf("Директория для хранения файлов создана: %s\n", StorageBaseDir)
	}
}

// SaveFile сохраняет файл на локальную файловую систему
// filename - уникальное имя, под которым файл будет сохранен (например, UUID)
// fileHeader - заголовок файла из multipart-формы
func SaveFile(filename string, fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size > utils.FILE_UPLOAD_MAX_SIZE {
		return "", utils.NewStorageError("Размер файла превышает лимит", 413, "File too large")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("не удалось открыть загруженный файл: %w", err)
	}
	defer src.Close()

	filePath := filepath.Join(StorageBaseDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось создать файл на диске: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("не удалось скопировать данные файла: %w", err)
	}

	return filePath, nil
}

// GetFile читает файл с локальной файловой системы
func GetFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, utils.NewStorageError("Файл не найден", 404, "File not found on disk")
		}
		return nil, fmt.Errorf("не удалось прочитать файл с диска: %w", err)
	}
	return data, nil
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return utils.NewStorageError("Файл не найден для удаления", 404, "File not found for deletion")
		}
		return fmt.Errorf("не удалось удалить файл с диска: %w", err)
	}
	return nil
}