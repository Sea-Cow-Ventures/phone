package main

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

func readCalls(page int, limit int) (CallResponse, error) {
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
