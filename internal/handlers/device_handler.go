package handlers

import (
	"time"
	
	"github.com/gin-gonic/gin"
	"portofolionetworkapi/internal/database"
	"portofolionetworkapi/internal/models"
)

func ListDevices(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, name, ip_address, location, status, version, last_seen, created_at, updated_at
		FROM devices ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch devices"})
		return
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device
		err := rows.Scan(&device.ID, &device.Name, &device.IPAddress, &device.Location,
			&device.Status, &device.Version, &device.LastSeen, &device.CreatedAt, &device.UpdatedAt)
		if err != nil {
			continue
		}
		
		if time.Since(device.LastSeen) < 2*time.Minute {
			device.Status = "online"
		} else {
			device.Status = "offline"
		}
		
		devices = append(devices, device)
	}

	// Get device count and limit info
	limitReached, count, _ := database.IsDeviceLimitReached()
	timeUntilReset := database.GetTimeUntilReset()

	c.JSON(200, gin.H{
		"success": true,
		"data":    devices,
		"total":   len(devices),
		"meta": gin.H{
			"limit_reached":    limitReached,
			"max_devices":      database.MaxDevices,
			"current_count":    count,
			"reset_in_minutes": int(timeUntilReset.Minutes()),
		},
	})
}

func CreateDevice(c *gin.Context) {
	// Check device limit
	limitReached, count, err := database.IsDeviceLimitReached()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to check device limit"})
		return
	}
	
	if limitReached {
		resetIn := database.GetTimeUntilReset()
		c.JSON(400, gin.H{
			"error":   "Demo device limit reached",
			"message": "Maximum 25 devices allowed. Database will auto-reset in " + formatDuration(resetIn),
			"limit":   database.MaxDevices,
			"current": count,
			"reset_in_minutes": int(resetIn.Minutes()),
		})
		return
	}

	var req models.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var deviceID string
	err = database.DB.QueryRow(`
		INSERT INTO devices (name, ip_address, location, status, last_seen)
		VALUES ($1, $2, $3, 'online', NOW()) RETURNING id
	`, req.Name, req.IPAddress, req.Location).Scan(&deviceID)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create device"})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"message": "Device created",
		"id":      deviceID,
	})
}

func UpdateDevice(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	_, err := database.DB.Exec(`
		UPDATE devices 
		SET name = COALESCE(NULLIF($1, ''), name),
			ip_address = COALESCE(NULLIF($2, ''), ip_address),
			location = COALESCE(NULLIF($3, ''), location),
			status = COALESCE(NULLIF($4, ''), status),
			updated_at = NOW()
		WHERE id = $5
	`, req.Name, req.IPAddress, req.Location, req.Status, id)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update device"})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Device updated"})
}

func DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	_, err := database.DB.Exec("DELETE FROM devices WHERE id = $1", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete device"})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Device deleted"})
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	if minutes < 1 {
		return "less than 1 minute"
	}
	if minutes == 1 {
		return "1 minute"
	}
	return string(rune(minutes)) + " minutes"
}
