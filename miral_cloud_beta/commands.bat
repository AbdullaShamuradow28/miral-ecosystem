go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get golang.org/x/crypto/bcrypt
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/sqlite # или gorm.io/driver/postgres
go get github.com/google/uuid

go install github.com/swaggo/swag/cmd/swag@latest
swag init

go mod init miral_cloud_go
go mod tidy
swag init
go run .
