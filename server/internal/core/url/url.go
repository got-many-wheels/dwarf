package coreurl

import (
	"fmt"
	"time"
)

type URL struct {
	// Code to represent the shortened url as a call name
	Code string `json:"code" bson:"code"`
	// The original url
	Long string `json:"long" bson:"long"`
	// The date of the url was created
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// URL in string format for debugging purposes
func (u *URL) String() string {
	return fmt.Sprintf("Code: %v, Long: %v", u.Code, u.Long)
}
