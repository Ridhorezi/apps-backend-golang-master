package transaction

import (
	"startup-backend-api/user"
	"time"
)

//==================From-Table-Transaction-Campaign-User===================//
type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	User       user.User // relationship
}
