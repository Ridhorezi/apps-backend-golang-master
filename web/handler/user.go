package handler

import (
	"net/http"
	"startup-backend-api/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) Index(c *gin.Context) {

	users, err := h.userService.GetAllUsers()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "user_index.html", gin.H{"users": users})

}

func (h *userHandler) Add(c *gin.Context) {

	c.HTML(http.StatusOK, "user_add.html", nil)

}

func (h *userHandler) Create(c *gin.Context) {

	var input user.FormCreateUsersInput

	err := c.ShouldBind(&input)

	if err != nil {
		input.Error = err
		c.HTML(http.StatusOK, "user_add.html", input)
		return
	}

	registerInput := user.RegisterUserInput{}
	registerInput.Name = input.Name
	registerInput.Email = input.Email
	registerInput.Occupation = input.Occupation
	registerInput.Password = input.Password

	_, err = h.userService.RegisterUser(registerInput)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/users")

}
