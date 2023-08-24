package bootstrap

import (
	"log"
	"os"

	"github.com/dgomezlikeyoube/ms_domain/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func DBconnection() (*gorm.DB, error) {
	// Configuración de la base de datos
	dbConfig := DatabaseConfig{
		Host:     os.Getenv("DATABASE_HOST_MSQL"),
		Port:     os.Getenv("DATABASE_PORT_MSQL"),
		User:     os.Getenv("DATABASE_USER_MSQL"),
		Password: os.Getenv("DATABASE_PASS_MSQL"),
		DBName:   os.Getenv("DATABASE_DBNAME_MSQL"),
	}

	// Construye la cadena de conexión
	dsn := dbConfig.User + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + dbConfig.Port + ")/" + dbConfig.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	// Abre la conexión a la base de datos
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}

	if os.Getenv("DATABASE_MIGRATE") == "true" {
		if err := db.AutoMigrate(&domain.Product{}); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

}
