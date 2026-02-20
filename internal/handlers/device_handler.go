package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"portofolionetworkapi/internal/database"
	"portofolionetworkapi/internal/models"
)

func ListDevices(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, name, ip_address, location, status, created_at, updated_at
		FROM devices ORDER BY id DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var d models.Device
		if err := rows.Scan(&d.ID, &d.Name, &d.IPAddress,
			&d.Location, &d.Status, &d.CreatedAt, &d.UpdatedAt); err != nil {
			continue
		}
		devices = append(devices, d)
	}

	if devices == nil {
		devices = []models.Device{}
	}

	c.JSON(http.StatusOK, gin.H{"data": devices, "total": len(devices)})
}

func CreateDevice(c *gin.Context) {
	var req models.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id string
	err := database.DB.QueryRow(`
		INSERT INTO devices (name, ip_address, location, status)
		VALUES ($1, $2, $3, 'unknown')
		RETURNING id
	`, req.Name, req.IPAddress, req.Location).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "device created"})
}

func UpdateDevice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec(`
		UPDATE devices SET name=$1, ip_address=$2, location=$3, status=$4, updated_at=NOW()
		WHERE id=$5
	`, req.Name, req.IPAddress, req.Location, req.Status, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "device updated"})
}

func DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	_, err := database.DB.Exec("DELETE FROM devices WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "device deleted"})
}
