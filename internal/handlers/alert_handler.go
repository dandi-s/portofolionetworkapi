package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"portofolionetworkapi/internal/services"
)

type AlertHandler struct {
	orchestrator        *services.AlertOrchestrator
	defaultAssignUserID int
}

func NewAlertHandler(orchestrator *services.AlertOrchestrator, defaultUserID int) *AlertHandler {
	return &AlertHandler{
		orchestrator:        orchestrator,
		defaultAssignUserID: defaultUserID,
	}
}

type ZabbixWebhookRequest struct {
	EventID   string `json:"event_id" binding:"required"`
	Device    string `json:"device"   binding:"required"`
	IP        string `json:"ip"        binding:"required"`
	Severity  string `json:"severity"  binding:"required"`
	Problem   string `json:"problem"   binding:"required"`
	Status    string `json:"status"    binding:"required"`
	Customers int    `json:"customers_affected"`
	SLA       string `json:"sla_status"`
	Timestamp string `json:"timestamp"`
}

func (h *AlertHandler) HandleZabbixWebhook(c *gin.Context) {
	var req ZabbixWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}

	ts := time.Now()
	if req.Timestamp != "" {
		if parsed, err := time.Parse(time.RFC3339, req.Timestamp); err == nil {
			ts = parsed
		}
	}
	if req.SLA == "" {
		req.SLA = deriveSLA(req.Severity)
	}

	alert := services.AlertPayload{
		EventID:   req.EventID,
		Device:    req.Device,
		IP:        req.IP,
		Severity:  req.Severity,
		Problem:   req.Problem,
		Status:    req.Status,
		Timestamp: ts,
		Customers: req.Customers,
		SLA:       req.SLA,
	}

	assignUser := h.defaultAssignUserID
	if raw := c.Query("assign_user"); raw != "" {
		if uid, err := strconv.Atoi(raw); err == nil {
			assignUser = uid
		}
	}

	log.Printf("[WEBHOOK] %s [%s] %s → %s", req.EventID, req.Severity, req.Device, req.Status)
	result := h.orchestrator.HandleAlert(alert, assignUser)
	c.JSON(http.StatusOK, result)
}

func deriveSLA(severity string) string {
	switch severity {
	case "DISASTER", "HIGH":
		return "BREACHED"
	case "AVERAGE":
		return "AT_RISK"
	default:
		return "OK"
	}
}