package database

import (
	"fmt"
	"time"

	"aidan/phone/internal/models"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func ReadCalls(page int, limit int) (models.CallPages, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	calls := []models.Call{}

	var totalCalls int
	countQuery := `
		SELECT COUNT(*) 
		FROM calls c
		LEFT JOIN agents a ON c.toNumber = a.number 
		WHERE NOT (c.direction = 'outbound-dial' AND a.number IS NOT NULL)`

	err := db.QueryRow(countQuery).Scan(&totalCalls)
	if err != nil {
		return models.CallPages{}, err
	}

	totalPages := (totalCalls + limit - 1) / limit

	query := `
		SELECT 
			c.id, 
			c.fromNumber, 
			c.toNumber, 
			c.direction, 
			c.createdDate, 
			COALESCE(c.callerName, '') AS callerName,
			COALESCE(ha.name, '') AS handledBy
		FROM calls c
		LEFT JOIN agents a ON c.toNumber = a.number 
		LEFT JOIN handledCalls hc ON c.id = hc.callId
		LEFT JOIN agents ha ON hc.agentId = ha.id
		WHERE NOT (c.direction = 'outbound-dial' AND a.number IS NOT NULL) 
		ORDER BY c.createdDate DESC 
		LIMIT ? OFFSET ?`

	rows, err := db.Queryx(query, limit, offset)
	if err != nil {
		return models.CallPages{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var call models.Call
		err := rows.StructScan(&call)
		if err != nil {
			return models.CallPages{}, err
		}
		calls = append(calls, call)
	}

	return models.CallPages{
		Calls:       calls,
		TotalPages:  totalPages,
		CurrentPage: page,
	}, nil
}

func InsertCall(call twilioApi.ApiV2010Call) error {
	query := `
		INSERT INTO calls (
			fromNumber, toNumber, direction, updatedDate, price, uri,
			accountSid, status, callSid, sentDate, createdDate,
			priceUnit, apiVersion, parentCallSid,
			toFormatted, fromFormatted, phoneNumberSid, answeredBy,
			forwardedFrom, groupSid, callerName, queueTime, trunkSid
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	updatedDate, err := time.Parse(time.RFC1123Z, *call.DateUpdated)
	if err != nil {
		return fmt.Errorf("failed to parse updated date: %w", err)
	}

	sentDate, err := time.Parse(time.RFC1123Z, *call.StartTime)
	if err != nil {
		return fmt.Errorf("failed to parse sent date: %w", err)
	}

	createdDate, err := time.Parse(time.RFC1123Z, *call.DateCreated)
	if err != nil {
		return fmt.Errorf("failed to parse created date: %w", err)
	}

	_, err = db.Exec(
		query,
		call.From,
		call.To,
		call.Direction,
		updatedDate.Format("2006-01-02 15:04:05"),
		call.Price,
		call.Uri,
		call.AccountSid,
		call.Status,
		call.Sid,
		sentDate.Format("2006-01-02 15:04:05"), // Format time for MySQL DATETIME
		createdDate.Format("2006-01-02 15:04:05"), // Format time for MySQL DATETIME
		call.PriceUnit,
		call.ApiVersion,
		call.ParentCallSid,
		call.ToFormatted,
		call.FromFormatted,
		call.PhoneNumberSid,
		call.AnsweredBy,
		call.ForwardedFrom,
		call.GroupSid,
		call.CallerName,
		call.QueueTime,
		call.TrunkSid,
	)

	if err != nil {
		return fmt.Errorf("failed to insert call log: %w", err)
	}

	return nil
}

func DoesCallExist(callSid string) (bool, error) {
	query := "SELECT COUNT(*) FROM calls WHERE callSid = ?"

	var count int
	err := db.Get(&count, query, callSid)
	if err != nil {
		return false, fmt.Errorf("failed to check callSid existence: %w", err)
	}

	return count > 0, nil
}

func MarkCallHandled(callId, agentId int) error {
	query := `
		INSERT INTO handledCalls (callId, agentId) 
		VALUES (?, ?) 
		ON DUPLICATE KEY UPDATE agentId = ?
	`
	_, err := db.Exec(query, callId, agentId, agentId)
	if err != nil {
		return err
	}
	return nil
}
