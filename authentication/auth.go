package fauth

import (
	"encoding/json"
	"log"
)

type FacebookError struct {
	Data struct {
		Err struct {
			code    string `josn:"code"`
			message string `josn:"message"`
		} `json:"error"`
		IsValid bool `json:"is_valid"`
	} `josn:"data"`
}

func (f FacebookError) Error() string {
	r, err := json.Marshal(f)
	if err != nil {
		log.Fatal("Errore non valido")
		return "errore fatale gesione di questo errore"
	}
	return string(r)
}

type FaceToken struct {
	Data struct {
		AppID               string `json:"app_id"`
		Type                string `json:"type"`
		Application         string `json:"application"`
		DataAccessExpiresAt int    `json:"data_access_expires_at"`
		ExpiresAt           int    `json:"expires_at"`
		IsValid             bool   `json:"is_valid"`
		IssuedAt            int    `json:"issued_at"`
		Metadata            struct {
			AuthType string `json:"auth_type"`
			Sso      string `json:"sso"`
		} `json:"metadata"`
		Scopes []string `json:"scopes"`
		UserID string   `json:"user_id"`
	} `json:"data"`
}
