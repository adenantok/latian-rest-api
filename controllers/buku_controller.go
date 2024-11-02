package controllers

import (
	"latian-rest-api/config"
	"latian-rest-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BukuService is an interface for managing books
type BukuService interface {
	GetBuku() ([]models.Buku, error)
}

// bukuService is a struct implementing BukuService
type bukuService struct{}

// NewBukuService initializes bukuService
func NewBukuService() BukuService {
	return &bukuService{}
}

// GetBuku fetches all books from the database
func (bs *bukuService) GetBuku() ([]models.Buku, error) {
	var bukuList []models.Buku
	if err := config.DB.Find(&bukuList).Error; err != nil {
		return nil, err
	}
	return bukuList, nil
}

// BukuController manages the book endpoints
type BukuController struct {
	service BukuService
}

// NewBukuController initializes BukuController
func NewBukuController(service BukuService) *BukuController {
	return &BukuController{service: service}
}

// GetBukuHandler handles the GET request for books
func (bc *BukuController) GetBukuHandler(c *gin.Context) {
	buku, err := bc.service.GetBuku()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data buku"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": buku})
}
