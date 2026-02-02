package coreurl

import (
	"fmt"
	"time"
)

// URL represents a shortened URL
type URL struct {
	// Record ID from database
	Id int `json:"id"`

	// Short code representing the URL
	// example: aZ3kLm
	Code string `json:"code"`

	// Original long URL
	// example: https://example.com/some/very/long/path
	Long string `json:"long"`

	// Time when the URL was created
	// example: 2025-01-28T12:34:56Z
	CreatedAt time.Time `json:"created_at"`
}

// URL in string format for debugging purposes
func (u *URL) String() string {
	return fmt.Sprintf("Code: %v, Long: %v", u.Code, u.Long)
}
