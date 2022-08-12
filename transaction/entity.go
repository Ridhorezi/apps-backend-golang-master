package transaction

import (
	"startup-backend-api/campaign"
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
	PaymentURL string
	User       user.User         // relationship
	Campaign   campaign.Campaign //relationship
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
