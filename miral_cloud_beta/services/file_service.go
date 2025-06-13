package services

import (
	"bytes" // Добавлен
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"miral_cloud_go/database" // Возможно, тебе потребуется переименовать это, если оно используется для инициализации DB, а не пакета "database"
	"miral_cloud_go/encryption"
	"miral_cloud_go/models"
	"miral_cloud_go/storage" // Добавлен, предполагается, что это твой модуль для работы с файловой системой
	"miral_cloud_go/utils"
	"os"     // Добавлен
	"path/filepath" // Добавлен
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileService - сервис для работы с файлами
type FileService struct {
	DB *gorm.DB
}

// NewFileService создает новый экземпляр FileService
func NewFileService(db *gorm.DB) *FileService {
	return &FileService{DB: db}
}

// CreateFileRequest - DTO для создания файла
type CreateFileRequest struct {
	Name        string                `form:"name" binding:"required"`
	Comment     string                `form:"comment"`
	File        *multipart.FileHeader `form:"file" binding:"required"`
	UserID      string                `form:"user_id" binding:"required"` // ID пользователя, загружающего файл
	Private     bool                  `form:"private"`
	GroupID     *uint                 `form:"group_id"`       // Необязательное поле
	IsEncrypted bool                  `form:"is_encrypted"` // Флаг: шифровать ли файл
}

// CreateFile загружает, (опционально) шифрует и сохраняет файл, а также его метаданные.
func (s *FileService) CreateFile(req CreateFileRequest) (*models.File, error) {
	// Проверка квоты пользователя (это будет более сложная логика, пока заглушка)
	// Для реального использования: получить пользователя, подсчитать его текущее использование storageUsed
	// и проверить, не превышает ли он STORAGE_MAX_SIZE после добавления нового файла.
	// Сейчас просто проверяем размер загружаемого файла.

	if req.File.Size > utils.FILE_UPLOAD_MAX_SIZE { // utils.FILE_UPLOAD_MAX_SIZE должен быть определен
		return nil, utils.NewStorageError("Размер файла превышает лимит", 413, "File too large")
	}

	// Генерация уникального имени файла для хранения
	fileUUID := uuid.New().String()
	// Используем OriginalFilename, а не Filename из заголовка. Filename может включать путь.
	storedFilename := fmt.Sprintf("%s_%s", fileUUID, filepath.Base(req.File.Filename)) 

	// Читаем содержимое файла
	fileContent, err := req.File.Open()
	if err != nil {
		return nil, utils.NewStorageError(fmt.Sprintf("Не удалось открыть загруженный файл: %v", err), 500, "File open error")
	}
	defer fileContent.Close()

	fileBytes, err := io.ReadAll(fileContent)
	if err != nil {
		return nil, utils.NewStorageError(fmt.Sprintf("Не удалось прочитать содержимое файла: %v", err), 500, "File read error")
	}

	var encryptionKey, iv []byte
	var finalFileBytes []byte = fileBytes
	var isEncrypted = req.IsEncrypted

	if isEncrypted {
		// Генерируем ключ и IV
		encryptionKey, err = encryption.GenerateAESKey() // encryption.GenerateAESKey() должен быть определен
		if err != nil {
			return nil, utils.NewStorageError(fmt.Sprintf("Не удалось сгенерировать ключ шифрования: %v", err), 500, "Encryption key generation error")
		}
		iv, err = encryption.GenerateIV() // encryption.GenerateIV() должен быть определен
		if err != nil {
			return nil, utils.NewStorageError(fmt.Sprintf("Не удалось сгенерировать IV: %v", err), 500, "IV generation error")
		}

		// Шифруем данные
		finalFileBytes, err = encryption.Encrypt(fileBytes, encryptionKey, iv) // encryption.Encrypt() должен быть определен
		if err != nil {
			return nil, utils.NewStorageError(fmt.Sprintf("Не удалось зашифровать файл: %v", err), 500, "File encryption error")
		}
		storedFilename += ".enc" // Добавляем расширение для зашифрованных файлов
	}

	// Сохраняем файл на диске/хранилище через storage пакет
	// storage.SaveFile должен быть определен и принимать имя файла и байты
	filePathInStorage, err := storage.SaveFile(storedFilename, bytes.NewReader(finalFileBytes)) 
	if err != nil {
		return nil, utils.NewStorageError(fmt.Sprintf("Не удалось сохранить файл на сервере: %v", err), 500, "File save error")
	}

	// Создаем запись о файле в БД
	file := models.File{
		Name:          req.Name,
		Comment:       req.Comment,
		FilePath:      filePathInStorage, // Сохраняем путь, который вернул storage.SaveFile
		IsEncrypted:   isEncrypted,
		EncryptionKey: encryptionKey,
		IV:            iv,
		UserID:        req.UserID, // Предполагаем, что UserID - это строка UUID
		Private:       req.Private,
		GroupID:       req.GroupID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.DB.Create(&file).Error; err != nil {
		// Если не удалось сохранить в БД, попробуем удалить файл с диска
		_ = storage.DeleteFile(filePathInStorage) // storage.DeleteFile должен быть определен
		return nil, utils.NewStorageError(fmt.Sprintf("Не удалось сохранить метаданные файла: %v", err), 500, "File metadata save error")
	}

	return &file, nil
}

// GetFileByID получает метаданные файла по ID
func (s *FileService) GetFileByID(id uint) (*models.File, error) {
	var file models.File
	// Preload "User" и "Group" для автоматической подгрузки связанных данных
	// Предполагается, что models.File имеет связи с models.MCUser (если она есть) и models.Group
	if err := s.DB.Preload("User").Preload("Group").First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewStorageError("Файл не найден", 404, "File not found")
		}
		return nil, utils.NewStorageError(fmt.Sprintf("Ошибка при получении файла: %v", err), 500, "Database error")
	}
	return &file, nil
}

// DownloadFileRequest - DTO для скачивания файла
type DownloadFileRequest struct {
	FileID uint `json:"file_id" binding:"required"`
}

// DownloadFile получает файл по ID, дешифрует его (если нужно) и возвращает содержимое.
func (s *FileService) DownloadFile(fileID uint) ([]byte, error) {
	file, err := s.GetFileByID(fileID) // Используем уже существующий метод
	if err != nil {
		return nil, err // Передаем ошибку как есть (StorageError или другая)
	}

	// Читаем файл из хранилища
	fileBytes, err := storage.ReadFile(file.FilePath) // storage.ReadFile должен быть определен
	if err != nil {
		return nil, utils.NewStorageError(fmt.Sprintf("Не удалось прочитать файл из хранилища: %v", err), 500, "Storage read error")
	}

	if file.IsEncrypted {
		// Дешифруем, если файл был зашифрован
		if file.EncryptionKey == nil || file.IV == nil {
			return nil, utils.NewStorageError("Отсутствуют ключ или IV для дешифрования", 500, "Missing encryption key/IV")
		}
		decryptedBytes, err := encryption.Decrypt(fileBytes, file.EncryptionKey, file.IV) // encryption.Decrypt() должен быть определен
		if err != nil {
			return nil, utils.NewStorageError(fmt.Sprintf("Не удалось дешифровать файл: %v", err), 500, "File decryption error")
		}
		return decryptedBytes, nil
	}

	return fileBytes, nil // Возвращаем исходные байты, если не зашифрован
}

// UpdateFileRequest - DTO для обновления метаданных файла
type UpdateFileRequest struct {
	Name    *string `json:"name"`    // Используем указатели для опциональных полей
	Comment *string `json:"comment"`
	Private *bool   `json:"private"`
	GroupID *uint   `json:"group_id"`
}

// UpdateFile обновляет метаданные файла.
func (s *FileService) UpdateFile(id uint, req UpdateFileRequest) (*models.File, error) {
	file, err := s.GetFileByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		file.Name = *req.Name
	}
	if req.Comment != nil {
		file.Comment = *req.Comment
	}
	if req.Private != nil {
		file.Private = *req.Private
	}
	if req.GroupID != nil {
		file.GroupID = req.GroupID
	}
	file.UpdatedAt = time.Now()

	if err := s.DB.Save(&file).Error; err != nil {
		return nil, utils.NewStorageError(fmt.Sprintf("Не удалось обновить метаданные файла: %v", err), 500, "File metadata update error")
	}

	return file, nil
}

// DeleteFile удаляет файл из хранилища и его метаданные из БД.
func (s *FileService) DeleteFile(id uint) error {
	file, err := s.GetFileByID(id)
	if err != nil {
		return err
	}

	// Удаляем файл из физического хранилища
	if err := storage.DeleteFile(file.FilePath); err != nil {
		return utils.NewStorageError(fmt.Sprintf("Не удалось удалить файл из хранилища: %v", err), 500, "Storage delete error")
	}

	// Удаляем запись из БД
	if err := s.DB.Delete(&models.File{}, id).Error; err != nil {
		// Ошибка при удалении из БД, возможно, стоит попытаться восстановить файл в хранилище?
		// Это зависит от бизнес-логики и критичности данных.
		return utils.NewStorageError(fmt.Sprintf("Не удалось удалить метаданные файла: %v", err), 500, "File metadata delete error")
	}

	return nil
}

// ListFilesFilter - DTO для фильтрации списка файлов
type ListFilesFilter struct {
	UserID  *string `query:"user_id"`
	Private *bool   `query:"private"`
	GroupID *uint   `query:"group_id"`
}

// ListFiles получает список файлов с возможностью фильтрации.
func (s *FileService) ListFiles(filter ListFilesFilter) ([]models.File, error) {
	var files []models.File
	query := s.DB.Preload("User").Preload("Group") // Загружаем связанные данные

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Private != nil {
		query = query.Where("private = ?", *filter.Private)
	}
	if filter.GroupID != nil {
		query = query.Where("group_id = ?", *filter.GroupID)
	}

	if err := query.Find(&files).Error; err != nil {
		return nil, utils.NewStorageError(fmt.Sprintf("Не удалось получить список файлов: %v", err), 500, "Database error")
	}

	return files, nil
}