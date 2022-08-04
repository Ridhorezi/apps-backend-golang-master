package auth

import "github.com/dgrijalva/jwt-go"

//==================Contract-JWT=====================//

type Service interface {
	GenerateToken(userId int) (string, error)
}

//===================Struct-Call=====================//

type jwtService struct {
}

//====================Secret-Key=====================//

var SECRET_KEY = []byte("STARTUP_s3cr3t_k3y")

//===============Pointer-To-jwtService===============//

func NewService() *jwtService {
	return &jwtService{}
}

//==================Func-jwtService==================//

func (s *jwtService) GenerateToken(userId int) (string, error) {

	claim := jwt.MapClaims{}
	claim["user_id"] = userId
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString(SECRET_KEY)

	if err != nil {
		return signedToken, err
	}

	return signedToken, nil

}
