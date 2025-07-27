/**
 * Event Processor
 * Real-time event processing with Go routines
 */

package analytics

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	UserID    string                 `json:"userId"`
	EventID   string                 `json:"eventId"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

type EventProcessor struct {
	eventChan   chan Event
	subscribers map[string][]chan Event
	mu          sync.RWMutex
	running     bool
}

func NewEventProcessor() *EventProcessor {
	return &EventProcessor{
		eventChan:   make(chan Event, 1000),
		subscribers: make(map[string][]chan Event),
		running:     false,
	}
}

// Start begins event processing
func (ep *EventProcessor) Start() {
	ep.running = true
	go ep.processEvents()
}

// Stop halts event processing
func (ep *EventProcessor) Stop() {
	ep.running = false
	close(ep.eventChan)
}

// PublishEvent sends an event to the processing pipeline
func (ep *EventProcessor) PublishEvent(event Event) {
	if ep.running {
		select {
		case ep.eventChan <- event:
		default:
			log.Printf("Event channel full, dropping event: %s", event.ID)
		}
	}
}

// Subscribe registers a channel to receive events of specific types
func (ep *EventProcessor) Subscribe(eventType string, ch chan Event) {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	
	if ep.subscribers[eventType] == nil {
		ep.subscribers[eventType] = make([]chan Event, 0)
	}
	ep.subscribers[eventType] = append(ep.subscribers[eventType], ch)
}

// processEvents handles incoming events
func (ep *EventProcessor) processEvents() {
	for event := range ep.eventChan {
		ep.mu.RLock()
		subscribers := ep.subscribers[event.Type]
		ep.mu.RUnlock()
		
		// Send to all subscribers
		for _, ch := range subscribers {
			select {
			case ch <- event:
			default:
				log.Printf("Subscriber channel full for event type: %s", event.Type)
			}
		}
		
		// Log event for debugging
		eventJSON, _ := json.Marshal(event)
		log.Printf("Processed event: %s", string(eventJSON))
	}
}