package main

import (
	"crypto/sha256"
	"fmt"
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	ClientSecret string `json:"clientSecret"`
	PwdHash      string `json:"pwdHash"`
}

type UserStatus struct {
	ID       string `json:"id"`
	isOnline bool   `json:"isOnline"`
}

func VerifyPassword(pwd string, hash string) bool {
	if hash == GeneratePasswordHash(pwd) {
		return true
	}
	return false
}

func GeneratePasswordHash(pwd string) string {
	h := sha256.New()
	h.Write([]byte(pwd))
	return fmt.Sprintf("%x", h.Sum(nil))
}
