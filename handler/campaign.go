package handler

import (
	"net/http"
	"startup-backend-api/campaign"
	"startup-backend-api/images/helper"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

//==================Campaign-Handler====================//

func (h *campaignHandler) GetCampaigns(c *gin.Context) {

	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)

	if err != nil {

		response := helper.APIResponse("Error to get campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return

	}

	response := helper.APIResponse("List of campaign", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))

	c.JSON(http.StatusOK, response)

}

//===============CampaignByID-Handler===================//

func (h *campaignHandler) GetCampaign(c *gin.Context) {

	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Detail campaign", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))

	c.JSON(http.StatusOK, response)

}
