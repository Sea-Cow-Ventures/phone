package main

type Call struct {
	ID          string `db:"id"`
	From        string `db:"fromNumber"`
	To          string `db:"toNumber"`
	Direction   string `db:"direction"`
	CreatedDate string `db:"createdDate"`
	CallerName  string `db:"callerName"`
}

func readCalls(page int, limit int) ([]Call, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	calls := []Call{}

	query := `
		SELECT c.id, c.fromNumber, c.toNumber, c.direction, c.createdDate, COALESCE(c.callerName, '') AS callerName
		FROM calls c
		LEFT JOIN agents a ON c.toNumber = a.number 
		WHERE NOT (c.direction = 'outbound-dial' AND a.number IS NOT NULL) 
		ORDER BY c.createdDate DESC 
		LIMIT ? OFFSET ?`

	rows, err := db.Queryx(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var call Call
		err := rows.StructScan(&call)
		if err != nil {
			return nil, err
		}
		calls = append(calls, call)
	}

	return calls, nil
}
