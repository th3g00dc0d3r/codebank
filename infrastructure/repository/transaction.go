package repository

import (
	"database/sql"
	"errors"

	"github.com/th3g00dc0d3r/codebank/domain"
)

type TransactionRepositoryDb struct {
	db *sql.DB
}

func NewTransactionRepositoryDb(db *sql.DB) *TransactionRepositoryDb {
	return &TransactionRepositoryDb{db: db}
}

func (t *TransactionRepositoryDb) SaveTransaction(transaction domain.Transaction, creditCard domain.CreditCard) error {

	stmt, err := t.db.Prepare(
		`INSERT INTO transactions 
		(
			id, credit_card_id, amount, status, description, store, created_at
		) 
		VALUES 
		(
			$1, $2, $3, $4, $5, $6, $7
		)`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		transaction.ID,
		transaction.CreditCardId,
		transaction.Amount,
		transaction.Status,
		transaction.Description,
		transaction.Store,
		transaction.CreatedAt,
	)

	if err != nil {
		return err
	}

	if transaction.Status == "approved" {
		err = t.updateBalance(creditCard)
		if err != nil {
			return err
		}
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionRepositoryDb) updateBalance(creditCard domain.CreditCard) error {
	_, err := t.db.Exec(
		"UPDATE credit_cards SET balance = $1 WHERE id = $2",
		creditCard.Balance, creditCard.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionRepositoryDb) CreateCreditCard(creditCard domain.CreditCard) error {
	smtp, err := t.db.Prepare(
		`INSERT INTO credit_cards 
		(
			id, name, number, expiration_month, expiration_year, CVV, balance, balance_limit
		)
		VALUES
		(
			$1, $2, $3, $4, $5, $6, $7, $8
		)`)

	if err != nil {
		return err
	}

	_, err = smtp.Exec(
		creditCard.ID,
		creditCard.Name,
		creditCard.Number,
		creditCard.ExpirationMonth,
		creditCard.ExpirationYear,
		creditCard.CVV,
		creditCard.Balance,
		creditCard.Limit,
	)

	if err != nil {
		return err
	}

	err = smtp.Close()
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionRepositoryDb) GetCreditCard(creditCard domain.CreditCard) (domain.CreditCard, error) {
	var c domain.CreditCard
	smtp, err := t.db.Prepare("SELECT id, balance, balance_limit FROM credit_cards WHERE number = $1")
	if err != nil {
		return c, err
	}
	if err = smtp.QueryRow(creditCard.Number).Scan(&c.ID, &c.Balance, &c.Limit); err != nil {
		return c, errors.New("CREDIT CARD DOES NOT EXIST")
	}

	return c, nil

}