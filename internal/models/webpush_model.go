package models

type SubscriptionRequest struct {
	Endpoint  string `json:"endpoint" db:"endpoint"`
	UserAgent string `json:"userAgent" db:"userAgent"`
	Keys      struct {
		Auth   string `json:"auth"`
		P256dh string `json:"p256dh"`
	} `json:"keys"`
}

type DBSubscription struct {
	Endpoint  string `db:"endpoint"`
	AuthKey   string `db:"authKey"`
	P256dhKey string `db:"p256dhKey"`
}
