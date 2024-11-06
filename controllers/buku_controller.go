package controllers

import (
	"latian-rest-api/config"
	"latian-rest-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BukuService is an interface for managing books
type BukuService interface {
	GetBuku() ([]models.Buku, error)
	AddBuku(buku models.Buku) error
	GetBukuById(id int) (*models.Buku, error)
	UpdateBuku(buku models.Buku) error
}

// bukuService is a struct implementing BukuService
type bukuService struct{}

// BukuController manages the book endpoints
type BukuController struct {
	service BukuService
}

// NewBukuService initializes bukuService
func NewBukuService() BukuService {
	return &bukuService{}
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

// GetBuku fetches all books from the database
func (bs *bukuService) GetBuku() ([]models.Buku, error) {
	var bukuList []models.Buku
	if err := config.DB.Find(&bukuList).Error; err != nil {
		return nil, err
	}
	return bukuList, nil
}

// AddBukuHandler menangani request HTTP dan response JSON
func (bc *BukuController) AddBukuHandler(c *gin.Context) {
	var buku models.Buku

	// Bind JSON ke struct buku
	if err := c.ShouldBindJSON(&buku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// Validasi: Cek apakah field Judul dan Harga tidak kosong
	var validationErrors []string

	if buku.Judul == "" {
		validationErrors = append(validationErrors, "Field 'judul' tidak boleh kosong")
	}
	if buku.Harga == nil {
		validationErrors = append(validationErrors, "Field 'harga' tidak boleh kosong")
	}

	// Jika ada kesalahan validasi, kirim respons error
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Panggil service untuk menambah buku ke database
	if err := bc.service.AddBuku(buku); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambah data buku"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data buku berhasil ditambahkan"})
}

func (bs *bukuService) AddBuku(buku models.Buku) error {
	if err := config.DB.Create(&buku).Error; err != nil {
		return err
	}
	return nil
}

// GetBukuHandler handles the GET request for books by id
func (bc *BukuController) GetBukuByIdHandler(c *gin.Context) {
	id := c.Param("id")

	// Konversi ID dari string ke int
	bukuId, err := strconv.Atoi(id) // Ubah string menjadi integer.

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	buku, err := bc.service.GetBukuById(bukuId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "buku tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": buku})
}

func (bs *bukuService) GetBukuById(id int) (*models.Buku, error) {
	var buku models.Buku
	if err := config.DB.First(&buku, id).Error; err != nil {
		return nil, err
	}
	return &buku, nil
}

func (bc *BukuController) UpdateBuku(c *gin.Context) {
	var buku models.Buku
	// Bind JSON ke struct buku
	if err := c.ShouldBindJSON(&buku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// Validasi: Cek apakah field Judul dan Harga tidak kosong
	var validationErrors []string

	if buku.Id == nil {
		validationErrors = append(validationErrors, "Field 'id' tidak boleh kosong")
	}
	if buku.Judul == "" {
		validationErrors = append(validationErrors, "Field 'judul' tidak boleh kosong")
	}
	if buku.Harga == nil {
		validationErrors = append(validationErrors, "Field 'harga' tidak boleh kosong")
	}

	// Jika ada kesalahan validasi, kirim respons error
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Panggil service untuk menambah buku ke database
	if err := bc.service.UpdateBuku(buku); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal merebuah data buku"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data buku berhasil dirubah"})
}

func (bs *bukuService) UpdateBuku(buku models.Buku) error {
	if err := config.DB.Save(&buku).Error; err != nil {
		return err
	}
	return nil
}
