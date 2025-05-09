package infrastructure

import (
	"fmt"
	"func/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDB() *gorm.DB {
	dsn := "host=localhost user=root password=root dbname=devdisitillery port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	db = db.Debug()

	if err := autoMigrate(db); err != nil {
		log.Fatalf("Error al migrar los modelos: %v", err)
	}

	fmt.Println("Conexión a la base de datos establecida correctamente")
	return db
}

func autoMigrate(db *gorm.DB) error {
	models := []interface{}{
		&domain.User{},
		&domain.Project{},
		&domain.Task{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("error al migrar el modelo %T: %v", model, err)
		}
	}

	return nil
}
