package models

import (
	"crypto/sha256"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Username     string             `json:"username"`
	ClientSecret string             `json:"clientSecret"`
	PwdHash      string             `json:"pwdHash"`
}

type UserStatus struct {
	ID       string `json:"id"`
	IsOnline bool   `json:"isOnline"`
}

func (u *User) VerifyPassword(pwd string) bool {
	return u.PwdHash == generatePasswordHash(pwd)
}

func generatePasswordHash(pwd string) string {
	h := sha256.New()
	h.Write([]byte(pwd))
	return fmt.Sprintf("%x", h.Sum(nil))
}
