package seed

import (
	"log"
	"users-grpc/internal/models"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	users := []models.User{
		{Name: "Alice", Age: 28, Address: "Jakarta"},
		{Name: "Bob", Age: 35, Address: "Bandung"},
		{Name: "Charlie", Age: 22, Address: "Surabaya"},
	}
	for _, u := range users {
		var existing models.User
		if err := db.Where("name = ? AND address = ?", u.Name, u.Address).First(&existing).Error; err == nil {
			continue
		}
		if err := db.Create(&u).Error; err != nil {
			log.Printf("seed error for %s: %v", u.Name, err)
		}
	}
}
