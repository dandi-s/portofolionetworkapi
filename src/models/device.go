package models

import "time"

type Device struct {
	ID        int       `json:"id"`
	DeviceID  string    `json:"device_id"`
	Name      string    `json:"name"`
	IPAddress string    `json:"ip_address"`
	Region    string    `json:"region"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateDeviceRequest struct {
	DeviceID  string `json:"device_id" binding:"required"`
	Name      string `json:"name"      binding:"required"`
	IPAddress string `json:"ip_address" binding:"required"`
	Region    string `json:"region"`
}

type UpdateDeviceRequest struct {
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
	Region    string `json:"region"`
	Status    string `json:"status"`
}