package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {

	_, err := CreateFile()
	if err != nil {
		panic(err)
	}
	cotacao, err := RequestCotacao()
	if err != nil {
		panic(err)
	}
	err = SaveCotacaoFile(*cotacao)

}
func CreateFile() (*os.File, error) {
	file, err := os.Create("datasorce/cotacao.txt")
	if err != nil {
		return nil, err
	}
	return file, nil
}
func RequestCotacao() (*Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data Cotacao
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
func SaveCotacaoFile(cotacao Cotacao) error {
	file, err := os.OpenFile("datasorce/cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = fmt.Fprintln(file, fmt.Sprintf("Dólar: %s \n", cotacao.USDBRL.Bid))

	if err != nil {
		return err
	}
	defer file.Close()
	if err != nil {
		return err
	}
	return nil
}
