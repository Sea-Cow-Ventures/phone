package main

import (
	"fmt"
	"time"
)

type BlockedPhone struct {
	Id          int       `db:"id"`
	PhoneNumber string    `db:"phoneNumber"`
	BlockDate   time.Time `db:"blockDate"`
}

func getBlockedNumbers() ([]BlockedPhone, error) {
	query :=
		`SELECT
			id,
			phoneNumber,
			blockDate
		FROM blockedPhoneNumbers`

	var blocks []BlockedPhone
	err := db.Select(&blocks, query)
	if err != nil {
		return []BlockedPhone{}, fmt.Errorf("failed to query blockedPhoneNumbers %w", err)
	}

	return blocks, nil
}
