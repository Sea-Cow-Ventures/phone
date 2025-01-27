package main

import (
	"fmt"

	lookupsV2 "github.com/twilio/twilio-go/rest/lookups/v2"

	"go.uber.org/zap"
)

func lookupPhoneNumber(phoneNumber string) (*lookupsV2.LookupsV2PhoneNumber, error) {
	fields := "caller_name"
	params := lookupsV2.FetchPhoneNumberParams{Fields: &fields}

	lookup, err := t.LookupsV2.FetchPhoneNumber(phoneNumber, &params)
	if err != nil {
		fmt.Println("Error fetching phone number:", err)
	}

	logger.Info("Phone number", zap.Any("phoneNumber", lookup))

	err = insertPhoneLookup(lookup)
	if err != nil {
		return nil, fmt.Errorf("error inserting phone number lookup data: %w", err)
	}

	return lookup, nil
}
func insertPhoneLookup(lookup *lookupsV2.LookupsV2PhoneNumber) error {
	query := `
        INSERT INTO phoneNumberLookups (
            callingCountryCode,
            countryCode,
            phoneNumber,
            nationalFormat,
            valid,
            callerName,
            callerType,
            url
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `

	callerName := ""
	callerType := ""

	if callerNameMap, ok := (*lookup.CallerName).(map[string]interface{}); ok {
		if name, exists := callerNameMap["caller_name"]; exists {
			callerName = name.(string)
		}
		if cType, exists := callerNameMap["caller_type"]; exists {
			callerType = cType.(string)
		}
	}

	_, err := db.Exec(query,
		lookup.CallingCountryCode,
		lookup.CountryCode,
		lookup.PhoneNumber,
		lookup.NationalFormat,
		lookup.Valid,
		callerName,
		callerType,
		lookup.Url,
	)
	return err
}
