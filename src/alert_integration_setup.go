package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	// ✅ Sama dengan import di main.go — dari internal/
	"portofolionetworkapi/internal/handlers"
	"portofolionetworkapi/internal/services"
)

func setupAlertIntegration(router *gin.Engine) {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	odooURL := os.Getenv("ODOO_URL")
	odooDB := os.Getenv("ODOO_DB")
	odooUser := os.Getenv("ODOO_USER")
	odooPass := os.Getenv("ODOO_PASSWORD")

	if botToken == "" || chatID == "" {
		log.Println("[WARN] TELEGRAM_BOT_TOKEN / TELEGRAM_CHAT_ID not set")
	}
	if odooURL == "" {
		log.Println("[WARN] ODOO_URL not set — Odoo ticketing disabled")
	}

	telegram := services.NewTelegramService(botToken, chatID)
	odoo := services.NewOdooService(odooURL, odooDB, odooUser, odooPass)

	if odooURL != "" {
		if err := odoo.Login(); err != nil {
			log.Printf("[WARN] Odoo login failed: %v", err)
		} else {
			log.Println("[OK] Odoo connected")
		}
	}

	teamID := envInt("ODOO_TEAM_ID", 1)
	defaultUserID := envInt("ODOO_DEFAULT_USER_ID", 1)

	orchestrator := services.NewAlertOrchestrator(telegram, odoo, teamID)
	alertHandler := handlers.NewAlertHandler(orchestrator, defaultUserID)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/webhooks/zabbix", alertHandler.HandleZabbixWebhook)
		v1.POST("/alerts/test", alertHandler.HandleZabbixWebhook)
	}

	log.Println("[OK] Alert routes registered")
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}