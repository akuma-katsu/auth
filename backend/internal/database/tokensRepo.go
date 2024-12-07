package database

import (
	"database/sql"
	"log"
)

type TokenRepo struct {
	db *sql.DB
}

func NewTokenRepo(db *sql.DB) *TokenRepo {
	return &TokenRepo{db: db}
}

type Token struct {
	UserID  string `db:"user_id"`
	Refresh []byte `db:"Refresh"`
	IpAddr  string `db:"ip"`
}

func (tokenRepo *TokenRepo) InsertRefresh(userId string, hashedRefresh string, ipAddr string) error {
	_, err := tokenRepo.db.Exec("insert into tokens (user_id, Refresh, ip) values ($1, $2, $3)", userId, hashedRefresh, ipAddr)
	if err != nil {
		log.Fatal("Error: cant delete: ", err)
		return err
	}
	return nil
}

func (tokenRepo *TokenRepo) GetRefreshByUserID(id string) (Token, error) {
	token := Token{}
	err := tokenRepo.db.QueryRow("select user_id, Refresh, ip from tokens where user_id = $1", id).Scan(&token.UserID, &token.Refresh, &token.IpAddr)
	return token, err
}

func (tokenRepo *TokenRepo) DeleteRefreshByUserID(id string) error {
	_, err := tokenRepo.db.Exec("delete from tokens where user_id = $1", id)
	if err != nil {
		log.Fatal("Error: cant insert: ", err)
		return err
	}
	return nil
}
