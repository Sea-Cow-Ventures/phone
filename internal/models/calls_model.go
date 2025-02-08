package models

type Call struct {
	ID          string `db:"id"`
	From        string `db:"fromNumber"`
	To          string `db:"toNumber"`
	Direction   string `db:"direction"`
	CreatedDate string `db:"createdDate"`
	CallerName  string `db:"callerName"`
	HandledBy   string `db:"handledBy"`
}

type CallPages struct {
	Calls       []Call `json:"calls"`
	TotalPages  int    `json:"totalPages"`
	CurrentPage int    `json:"currentPage"`
}
