package services

import (
	"auth/backend/internal/config"
	"auth/backend/internal/database"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Service struct {
	tokenRepo *database.TokenRepo
	cfg       *config.AppConfig
}

type TokenResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func NewService(tokenRepo *database.TokenRepo, cfg *config.AppConfig) *Service {
	return &Service{tokenRepo: tokenRepo, cfg: cfg}
}

func (s *Service) HashRefresh(refresh string) ([]byte, error) {
	hashedRefresh, err := bcrypt.GenerateFromPassword([]byte(refresh), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}
	return hashedRefresh, nil
}

func (s *Service) NewJWT(key string, ipAddr string, lifeTime float64) (string, error) {
	payload := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(lifeTime)).Unix(),
		"ip":  ipAddr,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)
	tokenSigned, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenSigned, nil
}

func (s *Service) NewRefresh() (string, error) {
	b := make([]byte, 32)

	now := time.Now().Unix()
	source := rand.NewSource(now)
	r := rand.New(source)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func (s *Service) ValidateToken(requestRefresh string, currRefresh []byte) bool {
	err := bcrypt.CompareHashAndPassword(currRefresh, []byte(requestRefresh))
	if err != nil {
		return false
	}
	return true
}

func (s *Service) Login(userId string, ipAddr string) (TokenResponse, error) {
	res := TokenResponse{}
	accessLifeTime, err := strconv.ParseFloat(s.cfg.AccessLifeTime, 64)
	if err != nil {
		return res, err
	}

	res.Access, err = s.NewJWT(s.cfg.Secret, ipAddr, accessLifeTime)
	if err != nil {
		return res, err
	}

	res.Refresh, err = s.NewRefresh()
	if err != nil {
		return res, err
	}

	hashedRefresh, err := s.HashRefresh(res.Refresh)
	if err != nil {
		return res, err
	}

	err = s.tokenRepo.DeleteRefreshByUserID(userId)
	if err != nil {
		return res, err
	}

	err = s.tokenRepo.InsertRefresh(userId, string(hashedRefresh), ipAddr)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *Service) Refresh(refresh string, userID string, ipAddr string) (TokenResponse, error) {
	res := TokenResponse{}

	currRefresh, err := s.tokenRepo.GetRefreshByUserID(userID)
	if err != nil {
		return res, err
	}

	if currRefresh.IpAddr != ipAddr {
		// todo: отправка варнинга на почту
		log.Println("EMAIL WARNING")
		return res, fmt.Errorf("EMAIL WARNING")
	}

	refreshOk := s.ValidateToken(refresh, currRefresh.Refresh)
	if refreshOk == false {
		return res, fmt.Errorf("Invalid refresh")
	}

	return res, err
}
