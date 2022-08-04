package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

//===================Contract-User====================//

type Service interface {
	RegisterUser(RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
}

//===================Struct-Call=====================//
type service struct {
	repository Repository
}

//===============Pointer-To-Service==================//

func NewService(repository Repository) *service {
	return &service{repository}
}

//==================Service-Register=================//

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {

	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)
	user.Role = "user"
	newUser, err := s.repository.Save(user)

	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

//====================Func-Login======================//
func (s *service) Login(input LoginInput) (User, error) {

	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return user, err
	}

	return user, nil

}

//===============Func-Validation-or-IsEmailAvailable=================//

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {

	email := input.Email
	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil

}

//================Func-SaveAvatar===================//

func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {

	user, err := s.repository.FindById(ID)

	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation

	updatedUser, err := s.repository.Update(user)

	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil

}
