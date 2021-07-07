package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/th3g00dc0d3r/codebank/domain"
	"github.com/th3g00dc0d3r/codebank/infrastructure/repository"
)

func main() {
	db := setupDb()
	defer db.Close()

	cc := domain.NewCreditCard()
	cc.Number = "1234"
	cc.Name = "Rui"
	cc.ExpirationMonth = 7
	cc.ExpirationYear = 2021
	cc.CVV = 123
	cc.Limit = 1000
	cc.Balance = 0

	repo := repository.NewTransactionRepositoryDb(db)
	err := repo.CreateCreditCard(*cc)
	if err != nil {
		fmt.Println(err)
	} 
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"db",
		"5432",
		"postgres",
		"root",
		"codebank",
	)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("ERROR CONNECTING TO DATABASE")
	}

	return db
}
