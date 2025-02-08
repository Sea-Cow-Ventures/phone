package models

type VoiceWebhook struct {
	AccountSid        string `json:"AccountSid" form:"AccountSid" validate:"required"`
	ApiVersion        string `json:"ApiVersion" form:"ApiVersion" validate:"required"`
	CallSid           string `json:"CallSid" form:"CallSid" validate:"required"`
	CallStatus        string `json:"CallStatus" form:"CallStatus" validate:"required"`
	Called            string `json:"Called" form:"Called" validate:"required"`
	CalledCity        string `json:"CalledCity" form:"CalledCity" validate:"required"`
	CalledCountry     string `json:"CalledCountry" form:"CalledCountry" validate:"required"`
	CalledState       string `json:"CalledState" form:"CalledState" validate:"required"`
	CalledZip         string `json:"CalledZip" form:"CalledZip" validate:"required"`
	Caller            string `json:"Caller" form:"Caller" validate:"required"`
	CallerCity        string `json:"CallerCity" form:"CallerCity" validate:"required"`
	CallerCountry     string `json:"CallerCountry" form:"CallerCountry" validate:"required"`
	CallerName        string `json:"CallerName" form:"CallerName" validate:"required"`
	CallerState       string `json:"CallerState" form:"CallerState" validate:"required"`
	CallerZip         string `json:"CallerZip" form:"CallerZip" validate:"required"`
	Digits            string `json:"Digits" form:"Digits" validate:"required"`
	Direction         string `json:"Direction" form:"Direction" validate:"required"`
	FinishedOnKey     string `json:"FinishedOnKey" form:"FinishedOnKey" validate:"required"`
	From              string `json:"From" form:"From" validate:"required"`
	FromCity          string `json:"FromCity" form:"FromCity" validate:"required"`
	FromCountry       string `json:"FromCountry" form:"FromCountry" validate:"required"`
	FromState         string `json:"FromState" form:"FromState" validate:"required"`
	FromZip           string `json:"FromZip" form:"FromZip" validate:"required"`
	To                string `json:"To" form:"To" validate:"required"`
	ToCity            string `json:"ToCity" form:"ToCity" validate:"required"`
	ToCountry         string `json:"ToCountry" form:"ToCountry" validate:"required"`
	ToState           string `json:"ToState" form:"ToState" validate:"required"`
	ToZip             string `json:"ToZip" form:"ToZip" validate:"required"`
	Msg               string `json:"msg" form:"msg" validate:"required"`
	RecordingUrl      string `json:"RecordingUrl" form:"RecordingUrl"`
	RecordingDuration string `json:"RecordingDuration" form:"RecordingDuration"`
}
