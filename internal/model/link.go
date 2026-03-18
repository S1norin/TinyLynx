package model

import "time"

type Link struct {
	ID            int
	OriginalLink string
	ShortCode    string
	CreatedAt     time.Time
}

type Analytics struct {
	ID        int
	LinkID    int
	CreatedAt time.Time
	IPAddress string
	UserAgent string
	Referrer  string
	Country   string
	Device    string
	Browser   string
	Platform  string
}
