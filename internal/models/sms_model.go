package models

type Message struct {
	From     string `db:"fromNumber"`
	To       string `db:"toNumber"`
	Body     string `db:"body"`
	SentDate string `db:"sentDate"`
}
