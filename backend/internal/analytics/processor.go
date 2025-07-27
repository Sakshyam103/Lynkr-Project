/**
 * Analytics Processor
 * Event-driven real-time analytics processing
 */

package analytics

import (
	"database/sql"
	"encoding/json"

	// "fmt"
	"time"
)

type Event1 struct {
	Type      string                 `json:"type"`
	UserID    string                 `json:"userId"`
	EventID   string                 `json:"eventId"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

type Processor struct {
	db        *sql.DB
	eventChan chan Event1
	running   bool
}

func NewProcessor(db *sql.DB) *Processor {
	return &Processor{
		db:        db,
		eventChan: make(chan Event1, 1000),
		running:   false,
	}
}

func (p *Processor) Start() {
	p.running = true
	go p.processEvents()
}

func (p *Processor) Stop() {
	p.running = false
	close(p.eventChan)
}

func (p *Processor) PublishEvent(event Event1) {
	if p.running {
		select {
		case p.eventChan <- event:
		default:
			// Channel full, drop event
		}
	}
}

func (p *Processor) processEvents() {
	for event := range p.eventChan {
		switch event.Type {
		case "attendance":
			p.processAttendance(event)
		case "content_view":
			p.processContentView(event)
		case "engagement":
			p.processEngagement(event)
		}
	}
}

func (p *Processor) processAttendance(event Event1) {
	query := `
		INSERT INTO analytics_events (type, user_id, event_id, data, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	data, _ := json.Marshal(event.Data)
	p.db.Exec(query, event.Type, event.UserID, event.EventID, string(data), event.Timestamp)
}

func (p *Processor) processContentView(event Event1) {
	query := `
		UPDATE content SET view_count = view_count + 1 WHERE id = ?
	`
	if contentID, ok := event.Data["contentId"].(string); ok {
		p.db.Exec(query, contentID)
	}
}

func (p *Processor) processEngagement(event Event1) {
	query := `
		INSERT INTO engagement_metrics (user_id, event_id, action, value, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	action := event.Data["action"].(string)
	value := event.Data["value"].(float64)
	p.db.Exec(query, event.UserID, event.EventID, action, value, event.Timestamp)
}
