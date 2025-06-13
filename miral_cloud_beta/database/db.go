package database

import (
	"log"
	"miral_cloud_go/models" // Подключаем наши модели
	"os"

	"gorm.io/driver/sqlite" // Используем SQLite для простоты, можно заменить на postgres
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	// Указываем путь к файлу SQLite базы данных
	dbPath := "./miral.db"

	// Удаляем файл базы данных, если он существует (для чистого запуска при прототипировании)
	// В продакшене этого делать не нужно!
	if _, err := os.Stat(dbPath); err == nil {
		os.Remove(dbPath)
		log.Printf("Существующая база данных %s удалена для чистого запуска.", dbPath)
	}


	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Логируем запросы GORM
	})

	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	log.Println("Подключение к базе данных успешно!")

	// Автоматическая миграция схем
	err = DB.AutoMigrate(
		&models.MCUser{},
		&models.Group{},
		&models.File{},
		&models.Project{},
		&models.Collection{},
		&models.Document{},
	)
	if err != nil {
		log.Fatalf("Ошибка при миграции базы данных: %v", err)
	}
	log.Println("Миграция базы данных успешно завершена.")
}