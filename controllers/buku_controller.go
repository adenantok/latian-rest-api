package controllers

import (
	"errors"
	"latian-rest-api/config"
	"latian-rest-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BukuService adalah interface yang mendefinisikan operasi CRUD untuk buku
type BukuService interface {
	GetBuku() ([]models.Buku, error)          // Mengambil semua buku
	AddBuku(buku models.Buku) error           // Menambahkan buku baru
	GetBukuById(id int) (*models.Buku, error) // Mengambil buku berdasarkan ID
	UpdateBuku(buku models.Buku) error        // Memperbarui data buku
	DeleteBuku(id int) error                  // Menghapus buku berdasarkan ID
}

// bukuService adalah struct yang mengimplementasikan BukuService
type bukuService struct{}

// BukuController bertanggung jawab untuk menangani permintaan HTTP terkait buku
type BukuController struct {
	service BukuService
}

// NewBukuService menginisialisasi service buku
func NewBukuService() BukuService {
	return &bukuService{}
}

// NewBukuController menginisialisasi controller buku dengan service
func NewBukuController(service BukuService) *BukuController {
	return &BukuController{service: service}
}

// GetBukuHandler menangani permintaan GET untuk mendapatkan semua data buku
func (bc *BukuController) GetBukuHandler(c *gin.Context) {
	buku, err := bc.service.GetBuku()
	if err != nil {
		// Jika terjadi kesalahan saat mengambil data, kirim respons error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data buku"})
		return
	}
	// Berhasil mendapatkan data buku
	c.JSON(http.StatusOK, gin.H{"data": buku})
}

// GetBuku mengambil semua buku dari database
func (bs *bukuService) GetBuku() ([]models.Buku, error) {
	var bukuList []models.Buku
	// Menggunakan GORM untuk mengambil semua data buku
	if err := config.DB.Find(&bukuList).Error; err != nil {
		return nil, err
	}
	return bukuList, nil
}

// AddBukuHandler menangani permintaan POST untuk menambah buku
func (bc *BukuController) AddBukuHandler(c *gin.Context) {
	var buku models.Buku

	// Bind JSON dari request ke struct buku
	if err := c.ShouldBindJSON(&buku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// Validasi: Pastikan 'judul' tidak kosong dan 'harga' ada
	var validationErrors []string
	if buku.Judul == "" {
		validationErrors = append(validationErrors, "Field 'judul' tidak boleh kosong")
	}
	if buku.Harga == nil {
		validationErrors = append(validationErrors, "Field 'harga' tidak boleh kosong")
	}

	// Jika validasi gagal, kirim respons error
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Panggil service untuk menambahkan buku ke database
	if err := bc.service.AddBuku(buku); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambah data buku"})
		return
	}

	// Berhasil menambah buku
	c.JSON(http.StatusOK, gin.H{"message": "Data buku berhasil ditambahkan"})
}

// AddBuku menambah buku ke database menggunakan GORM
func (bs *bukuService) AddBuku(buku models.Buku) error {
	if err := config.DB.Create(&buku).Error; err != nil {
		return err
	}
	return nil
}

// GetBukuByIdHandler menangani permintaan GET berdasarkan ID buku
func (bc *BukuController) GetBukuByIdHandler(c *gin.Context) {
	id := c.Param("id")

	// Konversi ID dari string ke int
	bukuId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Panggil service untuk mendapatkan buku berdasarkan ID
	buku, err := bc.service.GetBukuById(bukuId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "buku tidak ditemukan"})
		return
	}

	// Berhasil mendapatkan buku berdasarkan ID
	c.JSON(http.StatusOK, gin.H{"data": buku})
}

// GetBukuById mengambil buku dari database berdasarkan ID
func (bs *bukuService) GetBukuById(id int) (*models.Buku, error) {
	var buku models.Buku
	if err := config.DB.First(&buku, id).Error; err != nil {
		return nil, err
	}
	return &buku, nil
}

// UpdateBuku menangani permintaan PUT untuk memperbarui buku
func (bc *BukuController) UpdateBuku(c *gin.Context) {
	var buku models.Buku
	if err := c.ShouldBindJSON(&buku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// Validasi input: cek apakah 'id', 'judul', dan 'harga' tidak kosong
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

	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	if err := bc.service.UpdateBuku(buku); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data buku"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data buku berhasil diperbarui"})
}

// UpdateBuku memperbarui data buku di database
func (bs *bukuService) UpdateBuku(buku models.Buku) error {
	if err := config.DB.Save(&buku).Error; err != nil {
		return err
	}
	return nil
}

// DeleteBukuHandler menangani permintaan DELETE berdasarkan ID buku
func (bc *BukuController) DeleteBukuHandler(c *gin.Context) {
	id := c.Param("id")
	bukuId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	err = bc.service.DeleteBuku(bukuId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Buku dengan ID " + id + " berhasil dihapus"})
}

// DeleteBuku menghapus buku dari database berdasarkan ID
func (bs *bukuService) DeleteBuku(id int) error {
	var buku models.Buku

	// Cari buku berdasarkan ID
	if err := config.DB.First(&buku, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("buku tidak ditemukan")
		}
		return err
	}

	// Hapus buku jika ditemukan
	if err := config.DB.Delete(&buku, id).Error; err != nil {
		return err
	}
	return nil
}
