/**
 * Security Audit
 * Security auditing and vulnerability assessment tools
 */

package security

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type SecurityAudit struct {
	db *sql.DB
}

type VulnerabilityReport struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
	Endpoint    string    `json:"endpoint"`
	CreatedAt   time.Time `json:"createdAt"`
}

type SecurityEvent struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	UserID    string    `json:"userId"`
	IPAddress string    `json:"ipAddress"`
	UserAgent string    `json:"userAgent"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewSecurityAudit(db *sql.DB) *SecurityAudit {
	return &SecurityAudit{db: db}
}

func (sa *SecurityAudit) RunSecurityScan() ([]VulnerabilityReport, error) {
	var vulnerabilities []VulnerabilityReport
	
	// Check for SQL injection vulnerabilities
	sqlVulns := sa.checkSQLInjection()
	vulnerabilities = append(vulnerabilities, sqlVulns...)
	
	// Check for weak authentication
	authVulns := sa.checkAuthentication()
	vulnerabilities = append(vulnerabilities, authVulns...)
	
	// Check for data exposure
	dataVulns := sa.checkDataExposure()
	vulnerabilities = append(vulnerabilities, dataVulns...)
	
	// Store vulnerabilities
	for _, vuln := range vulnerabilities {
		sa.storeVulnerability(vuln)
	}
	
	return vulnerabilities, nil
}

func (sa *SecurityAudit) checkSQLInjection() []VulnerabilityReport {
	var vulnerabilities []VulnerabilityReport
	
	// Check for parameterized queries usage
	query := `
		SELECT query_text FROM query_performance_logs 
		WHERE query_text LIKE '%' || ? || '%' 
		OR query_text LIKE '%' || ? || '%'
	`
	
	rows, err := sa.db.Query(query, "'+", "';")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var queryText string
			if rows.Scan(&queryText) == nil {
				vulnerabilities = append(vulnerabilities, VulnerabilityReport{
					ID:          generateID(),
					Type:        "sql_injection",
					Severity:    "high",
					Description: "Potential SQL injection vulnerability detected",
					Endpoint:    "database",
					CreatedAt:   time.Now(),
				})
			}
		}
	}
	
	return vulnerabilities
}

func (sa *SecurityAudit) checkAuthentication() []VulnerabilityReport {
	var vulnerabilities []VulnerabilityReport
	
	// Check for weak passwords
	query := `SELECT COUNT(*) FROM users WHERE LENGTH(password_hash) < 60`
	var weakPasswords int
	sa.db.QueryRow(query).Scan(&weakPasswords)
	
	if weakPasswords > 0 {
		vulnerabilities = append(vulnerabilities, VulnerabilityReport{
			ID:          generateID(),
			Type:        "weak_authentication",
			Severity:    "medium",
			Description: fmt.Sprintf("%d users with potentially weak password hashing", weakPasswords),
			Endpoint:    "authentication",
			CreatedAt:   time.Now(),
		})
	}
	
	return vulnerabilities
}

func (sa *SecurityAudit) checkDataExposure() []VulnerabilityReport {
	var vulnerabilities []VulnerabilityReport
	
	// Check for unencrypted sensitive data
	query := `SELECT COUNT(*) FROM users WHERE email NOT LIKE '%@%' OR email IS NULL`
	var invalidEmails int
	sa.db.QueryRow(query).Scan(&invalidEmails)
	
	if invalidEmails > 0 {
		vulnerabilities = append(vulnerabilities, VulnerabilityReport{
			ID:          generateID(),
			Type:        "data_validation",
			Severity:    "low",
			Description: "Invalid email formats detected in user data",
			Endpoint:    "user_data",
			CreatedAt:   time.Now(),
		})
	}
	
	return vulnerabilities
}

func (sa *SecurityAudit) LogSecurityEvent(eventType, userID, ipAddress, userAgent, details string) error {
	eventID := generateID()
	
	query := `
		INSERT INTO security_events (id, type, user_id, ip_address, user_agent, details, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := sa.db.Exec(query, eventID, eventType, userID, ipAddress, userAgent, details, time.Now())
	return err
}

func (sa *SecurityAudit) GetSecurityEvents(limit int) ([]SecurityEvent, error) {
	query := `
		SELECT id, type, user_id, ip_address, user_agent, details, created_at
		FROM security_events
		ORDER BY created_at DESC
		LIMIT ?
	`
	
	rows, err := sa.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var events []SecurityEvent
	for rows.Next() {
		var event SecurityEvent
		err := rows.Scan(&event.ID, &event.Type, &event.UserID, &event.IPAddress, &event.UserAgent, &event.Details, &event.CreatedAt)
		if err != nil {
			continue
		}
		events = append(events, event)
	}
	
	return events, nil
}

func (sa *SecurityAudit) storeVulnerability(vuln VulnerabilityReport) error {
	query := `
		INSERT INTO vulnerability_reports (id, type, severity, description, endpoint, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	_, err := sa.db.Exec(query, vuln.ID, vuln.Type, vuln.Severity, vuln.Description, vuln.Endpoint, vuln.CreatedAt)
	return err
}

func (sa *SecurityAudit) ValidateInput(input string) bool {
	// Check for common injection patterns
	patterns := []string{
		`<script`,
		`javascript:`,
		`onload=`,
		`onerror=`,
		`'.*OR.*'`,
		`".*OR.*"`,
		`UNION.*SELECT`,
		`DROP.*TABLE`,
	}
	
	input = strings.ToLower(input)
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, input)
		if matched {
			return false
		}
	}
	
	return true
}

func (sa *SecurityAudit) CheckRateLimit(userID, endpoint string, limit int, window time.Duration) bool {
	query := `
		SELECT COUNT(*) FROM security_events 
		WHERE user_id = ? AND details LIKE ? AND created_at > ?
	`
	
	var count int
	sa.db.QueryRow(query, userID, "%"+endpoint+"%", time.Now().Add(-window)).Scan(&count)
	
	return count < limit
}

func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (sa *SecurityAudit) SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}