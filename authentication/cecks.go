package fauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CeckToken(token string) (FaceToken, error) {
	if len(token) < 200 || len(token) > 400 {
		return FaceToken{}, errors.New("tocken errato")
	}
	request := fmt.Sprintf(serverFacebookTemplate, token[7:], serverToken) // :7 per rimuovere Bearer
	//fmt.Println(request)
	res, err := http.Get(request)
	if err != nil {
		return FaceToken{}, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return FaceToken{}, errors.New("impossibile raggiungere il facebook")
	}
	fmt.Println(string(data))
	var possibleError FacebookError
	err = json.Unmarshal(data, &possibleError)
	if !possibleError.Data.IsValid {
		return FaceToken{}, possibleError
	}
	var ftk FaceToken
	err = json.Unmarshal(data, &ftk)
	if err != nil {
		return FaceToken{}, errors.New("impossibile parsare tocken")
	}
	return ftk, nil
}
