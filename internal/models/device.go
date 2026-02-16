package models

import "time"

type Device struct {
    ID        string    `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    IPAddress string    `json:"ip_address" db:"ip_address"`
    Location  string    `json:"location" db:"location"`
    Status    string    `json:"status" db:"status"`
    Version   string    `json:"version" db:"version"`
    LastSeen  time.Time `json:"last_seen" db:"last_seen"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateDeviceRequest struct {
    Name      string `json:"name" binding:"required"`
    IPAddress string `json:"ip_address" binding:"required"`
    Location  string `json:"location" binding:"required"`
}

type UpdateDeviceRequest struct {
    Name      string `json:"name"`
    IPAddress string `json:"ip_address"`
    Location  string `json:"location"`
    Status    string `json:"status"`
}
