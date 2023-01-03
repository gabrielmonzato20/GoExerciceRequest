package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

type ContacaoDb struct {
	Cotacao string
}

func NewCotacaoDb(cotacao string) *ContacaoDb {
	return &ContacaoDb{
		Cotacao: cotacao,
	}
}

func main() {
	db, err := sql.Open("sqlite3", "datasorce/datasorce.db")
	if err != nil {
		panic(err)
	}
	err = CreateDatabase()
	if err != nil {
		panic(err)
	}
	err = CreateTable(db)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/cotacao", CotacaoDolar)
	http.ListenAndServe(":8080", nil)

}
func CreateDatabase() error {
	_, err := os.Create("datasorce/datasorce.db")
	if err != nil {
		return err
	}
	return nil
}
func CreateTable(db *sql.DB) error {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS exchange(id INTEGER PRIMARY KEY, dt datetime default current_timestamp, value  VARCHAR(30))")
	if err != nil {
		return err
	}
	statement.Exec()
	return nil
}
func InsertInto(db *sql.DB, cotacao ContacaoDb) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	stmt, err := db.PrepareContext(ctx, "INSERT  INTO  exchange (value) VALUES  (?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(cotacao.Cotacao)
	if err != nil {
		return err
	}
	return nil
}
func CotacaoDolar(w http.ResponseWriter, r *http.Request) {

	cotacao, err := RequestCotacao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	db, err := sql.Open("sqlite3", "datasorce/datasorce.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer db.Close()

	cotacaoDb := NewCotacaoDb(cotacao.USDBRL.Bid)
	err = InsertInto(db, *cotacaoDb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)

}

func RequestCotacao() (*Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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
