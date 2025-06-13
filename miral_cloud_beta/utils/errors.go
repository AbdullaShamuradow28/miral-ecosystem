package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// StorageError - аналог Django StorageError
type StorageError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Reason  string `json:"reason"`
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("StorageError: %s (Code: %d, Reason: %s)", e.Message, e.Code, e.Reason)
}

// ToJSON - для преобразования ошибки в JSON
func (e *StorageError) ToJSON() []byte {
	data, _ := json.Marshal(e)
	return data
}

// NewStorageError - вспомогательная функция для создания StorageError
func NewStorageError(message string, code int, reason string) *StorageError {
	return &StorageError{
		Message: message,
		Code:    code,
		Reason:  reason,
	}
}

// SendErrorResponse - Отправляет стандартизированный JSON-ответ об ошибке
func SendErrorResponse(c *gin.Context, err error) {
	if storageErr, ok := err.(*StorageError); ok {
		c.Data(storageErr.Code, "application/json", storageErr.ToJSON())
	} else {
		// Для любых других ошибок, возвращаем общую ошибку сервера
		genericErr := NewStorageError("Internal Server Error", http.StatusInternalServerError, err.Error())
		c.Data(http.StatusInternalServerError, "application/json", genericErr.ToJSON())
	}
}