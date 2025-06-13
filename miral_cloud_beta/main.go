package main

import (
	"log"
	"miral_cloud_go/database"
	"miral_cloud_go/handlers"
	"miral_cloud_go/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "miral_cloud_go/docs"
)

// @title Miral Cloud Storage API
// @version 1.0
// @description Backend API для облачного хранилища файлов.
// @contact.name Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1

func main() {
	
	db := database.InitDB()

	
	fileService := services.NewFileService(db)

	
	fileHandler := handlers.NewFileHandler(fileService)

	
	router := gin.Default()

	
	v1 := router.Group("/api/v1")
	{
		files := v1.Group("/files")
		{
			files.POST("", fileHandler.CreateFile)        
			files.GET("", fileHandler.ListFiles)            
			files.GET("/:id", fileHandler.GetFileByID)     
			files.GET("/:id/download", fileHandler.DownloadFile) 
			files.PUT("/:id", fileHandler.UpdateFile)       
			files.DELETE("/:id", fileHandler.DeleteFile)    
		}
		
	}

	
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	
	port := ":8080"
	log.Printf("Miral Cloud Go API запущен на порту %s", port)
	log.Fatal(router.Run(port))
}