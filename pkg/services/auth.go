package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nicolas-ggd/ch-mod/internal/db/models"
	"github.com/Nicolas-ggd/ch-mod/internal/db/models/request"
	"github.com/Nicolas-ggd/ch-mod/pkg/common"
	"github.com/Nicolas-ggd/ch-mod/pkg/repository"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"os"
	"strings"
)

var secret = os.Getenv("PRIVATE_SECRET")

type IAuthService interface {
	Register(model *request.UserRegisterRequest) (*models.User, error)
	Login(model *request.LoginRequest) (*string, error)
	Logout(userId uint) error
	CheckJWT(token string) (*models.TokenClaim, error)
	SetPassword(credentials request.SetPasswordRequest, hash string) error
	VerifyCredentials(email string) error
}

type AuthService struct {
	userRepository  repository.UserRepository
	tokenRepository repository.TokenRepository
}

func NewAuthService(authRepo repository.UserRepository, tokenRepository repository.TokenRepository) *AuthService {
	return &AuthService{userRepository: authRepo, tokenRepository: tokenRepository}
}

func (as *AuthService) Register(model *request.UserRegisterRequest) (*models.User, error) {
	err := common.ComparePasswords(model.Password, model.ConfirmPassword)
	if err != nil {
		return nil, err
	}

	user, err := as.userRepository.Register(model.ToModel())
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (as *AuthService) Login(model *request.LoginRequest) (*string, error) {
	user, err := as.userRepository.GetByEmail(model.Email)
	if err != nil {
		return nil, err
	}

	err = common.CompareHashAndPasswordBcrypt(user.Password, model.Password)
	if err != nil {
		return nil, fmt.Errorf("user email or password is incorrect")
	}

	token, err := as.tokenRepository.CreateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	userToken, err := as.userRepository.GetAuthToken(user.ID)
	if err != nil {
		userToken = &models.UserToken{
			UserID: user.ID,
			Type:   models.Auth,
		}
	}

	userToken.Hash = []byte(token)

	err = as.tokenRepository.CreateToken(userToken)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (as *AuthService) Logout(userId uint) error {
	err := as.tokenRepository.DeleteToken(userId, models.Auth)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) SetPassword(credentials request.SetPasswordRequest, hash string) error {
	token, err := as.userRepository.GetUserTokenByHash(hash)
	if err != nil {
		return err
	}

	if credentials.Password != credentials.PasswordConfirmation {
		return fmt.Errorf("password and confirmation password doesn't match")
	}

	err = as.userRepository.ChangePassword(&credentials, token.UserID)
	if err != nil {
		return err
	}

	err = as.tokenRepository.DeleteToken(token.UserID, models.Validation)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) CheckJWT(token string) (*models.TokenClaim, error) {

	extracted, err := ExtractBearerToken(token)
	if err != nil {
		return nil, err
	}

	//if user token exists
	user, err := as.userRepository.GetUserTokenByHash(extracted)
	if err != nil {
		return nil, err
	}

	if user.UserID == 0 {
		return nil, fmt.Errorf("unauthorized")
	}

	//validate that user token
	userObj, err := ParseJWTClaims(token)
	if err != nil {
		return nil, err
	}

	return userObj, nil
}

func (as *AuthService) VerifyCredentials(email string) error {
	user, err := as.userRepository.GetByEmail(email)
	if err != nil {
		return err
	}

	randomStr := common.GenerateRandomString(20)

	usrToken := &models.UserToken{
		UserID: user.ID,
		Type:   models.Validation,
		Hash:   []byte(randomStr),
	}

	err = as.tokenRepository.UpdateToken(usrToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = as.tokenRepository.CreateToken(usrToken)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	return nil
}

func ValidateJWT(jwtToken string) (*jwt.Token, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, errors.New("bad jwt token")
	}

	return token, nil
}

func ParseJWTClaims(header string) (*models.TokenClaim, error) {
	token, err := ExtractBearerToken(header)
	if err != nil {
		return nil, fmt.Errorf("failed to extract bearer token: %v", err)
	}

	parsedToken, err := ValidateJWT(token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate JWT token: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unexpected type for JWT claims")
	}

	var tokenClaim models.TokenClaim

	marshalled, err := json.Marshal(claims["user"])
	if err != nil {
		return nil, fmt.Errorf("marshalling error")
	}

	err = json.Unmarshal(marshalled, &tokenClaim)
	if err != nil {
		return nil, fmt.Errorf("error converting 'user' claim to float64")
	}

	return &tokenClaim, nil
}

func ExtractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("invalid credentials, missing header, or bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 || jwtToken[0] != "Bearer" {
		return "", errors.New("incorrectly formatted or missing 'Bearer' keyword in authorization header")
	}

	return jwtToken[1], nil
}
