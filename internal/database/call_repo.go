package database

import (
	"fmt"
	"time"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Call struct {
	ID          string `db:"id"`
	From        string `db:"fromNumber"`
	To          string `db:"toNumber"`
	Direction   string `db:"direction"`
	CreatedDate string `db:"createdDate"`
	CallerName  string `db:"callerName"`
}

type CallResponse struct {
	Calls       []Call `json:"calls"`
	TotalPages  int    `json:"totalPages"`
	CurrentPage int    `json:"currentPage"`
}

func ReadCalls(page int, limit int) (CallResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	calls := []Call{}

	var totalCalls int
	countQuery := `
		SELECT COUNT(*) 
		FROM calls c
		LEFT JOIN agents a ON c.toNumber = a.number 
		WHERE NOT (c.direction = 'outbound-dial' AND a.number IS NOT NULL)`

	err := db.QueryRow(countQuery).Scan(&totalCalls)
	if err != nil {
		return CallResponse{}, err
	}

	totalPages := (totalCalls + limit - 1) / limit

	query := `
		SELECT c.id, c.fromNumber, c.toNumber, c.direction, c.createdDate, COALESCE(c.callerName, '') AS callerName
		FROM calls c
		LEFT JOIN agents a ON c.toNumber = a.number 
		WHERE NOT (c.direction = 'outbound-dial' AND a.number IS NOT NULL) 
		ORDER BY c.createdDate DESC 
		LIMIT ? OFFSET ?`

	rows, err := db.Queryx(query, limit, offset)
	if err != nil {
		return CallResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var call Call
		err := rows.StructScan(&call)
		if err != nil {
			return CallResponse{}, err
		}
		calls = append(calls, call)
	}

	return CallResponse{
		Calls:       calls,
		TotalPages:  totalPages,
		CurrentPage: page,
	}, nil
}

/*func DialNumber(agentNumber string, toNumber string) error {
	database.Logger.Info("message")
	params := &twilioApi.CreateCallParams{}
	params.SetTo(agentNumber)
	params.SetFrom(database.Cnf.PhoneNumber)
	params.SetUrl("https://" + database.Cnf.UrlBasePath + "/connectAgent?toNumber=" + toNumber)

	resp, err := server.T.Api.CreateCall(params)
	if err != nil {
		database.Logger.Infof("Failed to initiate call: %v", err)
		return fmt.Errorf("failed to initiate call: %w", err)
	}

	database.Logger.Infof("Call initiated with SID: %s", *resp.Sid)
	return nil
}*/

/*func ReadAccountCallHistory() error {
	database.Logger.Info("Reading account call history")

	calls, err := server.T.Api.ListCall(&twilioApi.ListCallParams{})
	if err != nil {
		return err
	}

	sort.Slice(calls, func(i, j int) bool {
		createdDateI, errI := time.Parse(time.RFC1123Z, *calls[i].DateCreated)
		if errI != nil {
			fmt.Println("Error parsing created date for index", i, ":", errI)
			return false
		}

		createdDateJ, errJ := time.Parse(time.RFC1123Z, *calls[j].DateCreated)
		if errJ != nil {
			fmt.Println("Error parsing created date for index", j, ":", errJ)
			return false
		}

		return createdDateI.Before(createdDateJ)
	})

	var insertedCalls int

	for _, call := range calls {
		exists, err := doesCallExist(*call.Sid)
		if err != nil {
			database.Logger.Errorf("Error reading messages: %w", zap.Error(err))
		}
		if !exists {
			err = insertCall(call)
			insertedCalls++
		}
		if err != nil {
			database.Logger.Errorf("Error inserting message: %w", zap.Error(err))
		}
	}

	database.Logger.Info("Read account call history", zap.Int("NewCalls", insertedCalls))

	return nil
}*/

func insertCall(call twilioApi.ApiV2010Call) error {
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

func doesCallExist(callSid string) (bool, error) {
	query := "SELECT COUNT(*) FROM calls WHERE callSid = ?"

	var count int
	err := db.Get(&count, query, callSid)
	if err != nil {
		return false, fmt.Errorf("failed to check callSid existence: %w", err)
	}

	return count > 0, nil
}
