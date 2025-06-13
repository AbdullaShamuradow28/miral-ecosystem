package handlers

import (
	"fmt"
	"miral_cloud_go/services"
	"miral_cloud_go/utils"
	"net/http"
	"strconv" // Для парсинга ID
	"strings" // Для проверки URL-путей

	"github.com/gin-gonic/gin"
)

// FileHandler содержит экземпляр FileService
type FileHandler struct {
	FileService *services.FileService
}

// NewFileHandler создает новый FileHandler
func NewFileHandler(fileService *services.FileService) *FileHandler {
	return &FileHandler{FileService: fileService}
}

// CreateFile godoc
// @Summary Загрузить новый файл
// @Description Загружает файл и сохраняет его метаданные в базе данных. Опционально шифрует файл.
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Название файла"
// @Param comment formData string false "Комментарий к файлу"
// @Param file formData file true "Сам файл"
// @Param user_id formData string true "ID пользователя, загружающего файл"
// @Param private formData boolean false "Приватный ли файл (по умолчанию false)"
// @Param group_id formData integer false "ID группы (опционально)"
// @Param is_encrypted formData boolean false "Шифровать ли файл (по умолчанию false)"
// @Success 201 {object} models.File "Файл успешно загружен"
// @Failure 400 {object} utils.StorageError "Неверные данные запроса или файл отсутствует"
// @Failure 413 {object} utils.StorageError "Размер файла превышает лимит"
// @Failure 500 {object} utils.StorageError "Внутренняя ошибка сервера"
// @Router /files [post]
func (h *FileHandler) CreateFile(c *gin.Context) {
	var req services.CreateFileRequest

	// c.ShouldBind проверяет наличие полей из формы и выполняет валидацию
	if err := c.ShouldBind(&req); err != nil {
		// Ошибка валидации от Gin, обычно это Missing field или Invalid format
		errStr := err.Error()
		if strings.Contains(errStr, "required") && strings.Contains(errStr, "'file'") {
			c.JSON(http.StatusBadRequest, utils.NewStorageError("Поле 'file' обязательно", http.StatusBadRequest, "File field required"))
			return
		}
		c.JSON(http.StatusBadRequest, utils.NewStorageError(fmt.Sprintf("Неверные данные запроса: %v", err), http.StatusBadRequest, "Bad request data"))
		return
	}

	file, err := h.FileService.CreateFile(req)
	if err != nil {
		// Преобразуем ошибку в StorageError, если это наша кастомная ошибка
		if storageErr, ok := err.(*utils.StorageError); ok {
			c.JSON(storageErr.Code, storageErr)
		} else {
			// Если это другая внутренняя ошибка, возвращаем 500
			log.Printf("Ошибка при создании файла: %v", err)
			c.JSON(http.StatusInternalServerError, utils.NewStorageError(fmt.Sprintf("Внутренняя ошибка сервера: %v", err), http.StatusInternalServerError, "Internal server error"))
		}
		return
	}

	c.JSON(http.StatusCreated, file)
}

// GetFileByID godoc
// @Summary Получить метаданные файла по ID
// @Description Возвращает метаданные файла по его уникальному ID.
// @Tags Files
// @Produce json
// @Param id path int true "ID файла"
// @Success 200 {object} models.File "Метаданные файла"
// @Failure 400 {object} utils.StorageError "Неверный ID файла"
// @Failure 404 {object} utils.StorageError "Файл не найден"
// @Failure 500 {object} utils.StorageError "Внутренняя ошибка сервера"
// @Router /files/{id} [get]
func (h *FileHandler) GetFileByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewStorageError("Неверный ID файла", http.StatusBadRequest, "Invalid file ID"))
		return
	}

	file, err := h.FileService.GetFileByID(uint(id))
	if err != nil {
		if storageErr, ok := err.(*utils.StorageError); ok {
			c.JSON(storageErr.Code, storageErr)
		} else {
			log.Printf("Ошибка при получении файла по ID: %v", err)
			c.JSON(http.StatusInternalServerError, utils.NewStorageError(fmt.Sprintf("Внутренняя ошибка сервера: %v", err), http.StatusInternalServerError, "Internal server error"))
		}
		return
	}

	c.JSON(http.StatusOK, file)
}

// DownloadFile godoc
// @Summary Скачать файл
// @Description Скачивает файл по его ID. Если файл зашифрован, он будет дешифрован перед отправкой.
// @Tags Files
// @Produce octet-stream
// @Param id path int true "ID файла"
// @Success 200 {file} byte "Содержимое файла"
// @Failure 400 {object} utils.StorageError "Неверный ID файла"
// @Failure 404 {object} utils.StorageError "Файл не найден"
// @Failure 500 {object} utils.StorageError "Внутренняя ошибка сервера"
// @Router /files/{id}/download [get]
func (h *FileHandler) DownloadFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewStorageError("Неверный ID файла", http.StatusBadRequest, "Invalid file ID"))
		return
	}

	fileContent, err := h.FileService.DownloadFile(uint(id))
	if err != nil {
		if storageErr, ok := err.(*utils.StorageError); ok {
			c.JSON(storageErr.Code, storageErr)
		} else {
			log.Printf("Ошибка при скачивании файла: %v", err)
			c.JSON(http.StatusInternalServerError, utils.NewStorageError(fmt.Sprintf("Внутренняя ошибка сервера: %v", err), http.StatusInternalServerError, "Internal server error"))
		}
		return
	}

	// Получаем метаданные файла, чтобы определить имя файла для Content-Disposition
	fileMeta, err := h.FileService.GetFileByID(uint(id))
	if err != nil {
		// Если не удалось получить метаданные (хотя файл успешно скачан),
		// можно вернуть дефолтное имя или просто байты.
		// Здесь для простоты вернем 500, хотя файл есть. В реальной системе нужно продумать.
		log.Printf("Ошибка при получении метаданных файла для скачивания: %v", err)
		c.JSON(http.StatusInternalServerError, utils.NewStorageError("Не удалось получить метаданные файла для заголовка", http.StatusInternalServerError, "Internal server error"))
		return
	}

	// Устанавливаем заголовок Content-Disposition для правильного скачивания
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileMeta.Name))
	c.Data(http.StatusOK, "application/octet-stream", fileContent)
}

// UpdateFile godoc
// @Summary Обновить метаданные файла
// @Description Обновляет название, комментарий, приватность или группу файла.
// @Tags Files
// @Accept json
// @Produce json
// @Param id path int true "ID файла для обновления"
// @Param file body services.UpdateFileRequest true "Данные для обновления файла"
// @Success 200 {object} models.File "Метаданные файла успешно обновлены"
// @Failure 400 {object} utils.StorageError "Неверные данные запроса или ID"
// @Failure 404 {object} utils.StorageError "Файл не найден"
// @Failure 500 {object} utils.StorageError "Внутренняя ошибка сервера"
// @Router /files/{id} [put]
func (h *FileHandler) UpdateFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewStorageError("Неверный ID файла", http.StatusBadRequest, "Invalid file ID"))
		return
	}

	var req services.UpdateFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewStorageError(fmt.Sprintf("Неверные данные запроса: %v", err), http.StatusBadRequest, "Bad request data"))
		return
	}

	file, err := h.FileService.UpdateFile(uint(id), req)
	if err != nil {
		if storageErr, ok := err.(*utils.StorageError); ok {
			c.JSON(storageErr.Code, storageErr)
		} else {
			log.Printf("Ошибка при обновлении файла: %v", err)
			c.JSON(http.StatusInternalServerError, utils.NewStorageError(fmt.Sprintf("Внутренняя ошибка сервера: %v", err), http.StatusInternalServerError, "Internal server error"))
		}
		return
	}

	c.JSON(http.StatusOK, file)
}

// DeleteFile godoc
// @Summary Удалить файл
// @Description Удаляет файл из хранилища и его метаданные из базы данных.
// @Tags Files
// @Param id path int true "ID файла для удаления"
// @Success 204 "Файл успешно удален"
// @Failure 400 {object} utils.StorageError "Неверный ID файла"
// @Failure 404 {object} utils.StorageError "Файл не найден"
// @Failure 500 {object} utils.StorageError "Внутренняя ошибка сервера"
// @Router /files/{id} [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewStorageError("Неверный ID файла", http.StatusBadRequest, "Invalid file ID"))
		return
	}

	err = h.FileService.DeleteFile(uint(id))
	if err != nil {
		if storageErr, ok := err.(*utils.StorageError); ok {
			c.JSON(storageErr.Code, storageErr)
		} else {
			log.Printf("Ошибка при удалении файла: %v", err)
			c.JSON(http.StatusInternalServerError, utils.NewStorageError(fmt.Sprintf("Внутренняя ошибка сервера: %v", err), http.StatusInternalServerError, "Internal server error"))
		}
		return
	}

	c.Status(http.StatusNoContent) // 204 No Content
}

// ListFiles godoc
// @Summary Получить список файлов
// @Description Возвращает список файлов с возможностью фильтрации по пользователю, приватности или группе.
// @Tags Files
// @Produce json
// @Param user_id query string false "Фильтр по ID пользователя"
// @Param private query boolean false "Фильтр по приватности (true/false)"
// @Param group_id query int false "Фильтр по ID группы"
// @Success 200 {array} models.File "Список файлов"
// @Failure 500 {object} utils.StorageError "Внутренняя ошибка сервера"
// @Router /files [get]
func (h *FileHandler) ListFiles(c *gin.Context) {
	var filter services.ListFilesFilter
	// c.ShouldBindQuery привязывает параметры запроса (?user_id=...) к структуре
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewStorageError(fmt.Sprintf("Неверные параметры запроса: %v", err), http.StatusBadRequest, "Invalid query parameters"))
		return
	}

	files, err := h.FileService.ListFiles(filter)
	if err != nil {
		if storageErr, ok := err.(*utils.StorageError); ok {
			c.JSON(storageErr.Code, storageErr)
		} else {
			log.Printf("Ошибка при получении списка файлов: %v", err)
			c.JSON(http.StatusInternalServerError, utils.NewStorageError(fmt.Sprintf("Внутренняя ошибка сервера: %v", err), http.StatusInternalServerError, "Internal server error"))
		}
		return
	}

	c.JSON(http.StatusOK, files)
}