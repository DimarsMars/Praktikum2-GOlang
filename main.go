package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"column:id; primaryKey"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Age       string    `gorm:"column:age"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func main() {
	dsn := "root:@tcp(localhost:3306)/openapi?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
	}

	router := gin.Default()

	router.GET("/user", func(c *gin.Context) {
		var user []User
		db.Find(&user)
		c.JSON(http.StatusOK, gin.H{"data": user})
	})

	router.POST("/user", func(c *gin.Context) {
		var user User
		c.BindJSON(&user)
		db.Create(&user)
		c.JSON(http.StatusOK, gin.H{"data": user})
	})

	router.GET("/users/:id", func(c *gin.Context) {
		var users User
		id := c.Param("id")

		db.First(&users, id)
		c.JSON(http.StatusOK, gin.H{"data": users})
	})

	router.PUT("/users/:id", func(c *gin.Context) {
		var users User
		id := c.Param("id")

		if err := db.First(&users, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User Not Found"})
			return
		}

		if err := c.ShouldBindJSON(&users); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
			return
		}

		db.Save(&users)
		c.JSON(http.StatusOK, gin.H{"data": users})
	})

	router.DELETE("/users/:id", func(c *gin.Context) {
		var users User
		id := c.Param("id")

		db.Delete(&users, id)
		c.JSON(http.StatusOK, gin.H{"message": "User berhasil dihapus"})
	})

	router.Run(":3000")
}
