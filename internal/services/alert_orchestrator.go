package services

import (
	"fmt"
	"log"
	"time"
)

type AlertPayload struct {
	EventID   string    `json:"event_id"`
	Device    string    `json:"device"`
	IP        string    `json:"ip"`
	Severity  string    `json:"severity"`
	Problem   string    `json:"problem"`
	Status    string    `json:"status"`
	Customers int       `json:"customers_affected"`
	SLA       string    `json:"sla_status"`
	Timestamp time.Time `json:"timestamp"`
}

type AlertOrchestrator struct {
	telegram *TelegramService
	odoo     *OdooService
	teamID   int
}

func NewAlertOrchestrator(telegram *TelegramService, odoo *OdooService, teamID int) *AlertOrchestrator {
	return &AlertOrchestrator{
		telegram: telegram,
		odoo:     odoo,
		teamID:   teamID,
	}
}

type HandleAlertResult struct {
	TelegramSent bool   `json:"telegram_sent"`
	TicketID     int    `json:"ticket_id,omitempty"`
	Message      string `json:"message"`
	Error        string `json:"error,omitempty"`
}

func (o *AlertOrchestrator) HandleAlert(alert AlertPayload, assignUserID int) HandleAlertResult {
	result := HandleAlertResult{
		Message: "Alert processed successfully",
	}

	// 1. Format and send Telegram Message
	msg := fmt.Sprintf(
		"🚨 <b>Network Alert: %s</b>\n"+
			"<b>Device:</b> %s (%s)\n"+
			"<b>Severity:</b> %s\n"+
			"<b>Problem:</b> %s\n"+
			"<b>Impact:</b> %d Customers\n"+
			"<b>SLA Status:</b> %s\n"+
			"<b>Time:</b> %s",
		alert.Status,
		alert.Device,
		alert.IP,
		alert.Severity,
		alert.Problem,
		alert.Customers,
		alert.SLA,
		alert.Timestamp.Format(time.RFC1123),
	)

	if err := o.telegram.SendMessage(msg); err != nil {
		log.Printf("[ERROR] Failed to send Telegram message: %v", err)
		result.Error += fmt.Sprintf("Telegram Error: %v; ", err)
	} else {
		result.TelegramSent = true
	}

	// 2. Only create Odoo ticket if severity is high enough and status is PROBLEM
	//    Many setups only want tickets on new problems, not on resolution initially.
	//    We can adjust this logic, but let's default to creating tickets on PROBLEM.
	if alert.Status == "PROBLEM" && (alert.Severity == "DISASTER" || alert.Severity == "HIGH" || alert.Severity == "AVERAGE") {
		// Attempt login if not already done or handle session logic inside login
		if o.odoo.url != "" {
			title := fmt.Sprintf("[%s] %s - %s", alert.Severity, alert.Device, alert.Problem)
			desc := fmt.Sprintf(
				"Event ID: %s\nDevice: %s\nIP: %s\nSeverity: %s\nProblem: %s\nCustomers Affected: %d\nSLA: %s",
				alert.EventID, alert.Device, alert.IP, alert.Severity, alert.Problem, alert.Customers, alert.SLA,
			)

			ticketID, err := o.odoo.CreateTicket(title, desc, o.teamID, assignUserID)
			if err != nil {
				log.Printf("[ERROR] Failed to create Odoo ticket: %v", err)
				result.Error += fmt.Sprintf("Odoo Error: %v", err)
			} else {
				result.TicketID = ticketID

				// Send a followup message to Telegram with the ticket ID
				ticketMsg := fmt.Sprintf("✅ Ticket #%d created for issue on %s", ticketID, alert.Device)
				o.telegram.SendMessage(ticketMsg)
			}
		}
	}

	if result.Error != "" {
		result.Message = "Processed with errors"
	}

	return result
}
