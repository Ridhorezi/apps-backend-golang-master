package transaction

import "startup-backend-api/user"

type GetCampaignTransactionsInput struct {
	ID   int       `uri:"id" binding:"required"`
	User user.User // relationship
}

type CreateTransactionInput struct {
	Amount     int       `json:"amount" binding:"required"`
	CampaignID int       `json:"campaign_id" binding:"required"`
	User       user.User // relationship
}
