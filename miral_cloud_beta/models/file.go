package models

import "time"

// File - метаданные о файле. Сам файл хранится отдельно (локально или S3).
type File struct {
	gorm.Model
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name"`        // Оригинальное имя файла
	Comment       string    `json:"comment"`     // Комментарий к файлу
	FilePath      string    `json:"file_path"`   // Путь к файлу в хранилище (например, uploaded_files/uuid_originalfilename.enc)
	IsEncrypted   bool      `json:"is_encrypted"` // Флаг, указывающий, зашифрован ли файл в хранилище
	EncryptionKey []byte    `json:"-" gorm:"type:bytea"` // Ключ шифрования (скрыт из JSON)
	IV            []byte    `json:"-" gorm:"type:bytea"` // Вектор инициализации (скрыт из JSON)
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserID        string    `json:"user_id"`    // ID пользователя, которому принадлежит файл
	User          MCUser    `json:"user" gorm:"foreignKey:UserID"`
	Private       bool      `json:"private" gorm:"default:false"` // Приватный ли файл
	GroupID       *uint     `json:"group_id"`   // ID группы (может быть nil)
	Group         *Group    `json:"group" gorm:"foreignKey:GroupID"`
}

// EncryptedFile - специфические метаданные для зашифрованных файлов.
// Поскольку у нас теперь шифрование на стороне Go, мы можем хранить
// ключи и IV вместе с File, или иметь отдельную структуру.
// Я предлагаю добавить поля в File для простоты, т.к. каждый файл может быть зашифрован.
// Если файл зашифрован, то File.IsEncrypted = true, и мы используем эти поля.
/*
type EncryptedFile struct {
	gorm.Model
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"user_id"`
	User      MCUser    `json:"user" gorm:"foreignKey:UserID"`
	Private   bool      `json:"private" gorm:"default:false"`
	EncryptionKey []byte `json:"-" gorm:"type:bytea"` // Ключ шифрования
	IV        []byte    `json:"-" gorm:"type:bytea"` // Вектор инициализации
	FilePath  string    `json:"file_path"` // Путь к зашифрованному файлу в хранилище
}
*/