package sqlite

import (
	"currency-service/internal/storage"

	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" //init sqlite3 driver
)

type Storage struct {
	db * sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fnPath = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fnPath, err)
	}

	const create string = `
		CREATE TABLE IF NOT EXISTS currency (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			based TEXT(3),
			date TEXT(10),
			eur NUMERIC,
			jpy NUMERIC,
			uah NUMERIC
		);`

	if _, err := db.Exec(create); err != nil {
		return nil, fmt.Errorf("%s: %w", fnPath, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveCurrency(lr *storage.LatestRates) error {
	const fnPath = "storage.sqlite.SaveCurrency"

	stmt, err := s.db.Prepare("INSERT INTO currency(based, date, eur, jpy, uah) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", fnPath, err)
	}

	_, err = stmt.Exec(lr.Base, lr.Date, lr.EUR, lr.JPY, lr.UAH)
	if err != nil {
		return fmt.Errorf("%s: %w", fnPath, err)
	}

	return err
}

func (s *Storage) GetCurrencyByDate(date string) (*storage.LatestRates, error) {
	const fnPath = "storage.sqlite.GetCurrencyByDate"

	stmt, err := s.db.Prepare("SELECT based, date, eur, jpy, uah FROM currency WHERE date=?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fnPath, err)
	}

	var rates storage.LatestRates

	err = stmt.QueryRow(date).Scan(&rates.Base, &rates.Date, &rates.EUR, &rates.JPY, &rates.UAH)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fnPath, err)
	}

	return &rates, nil
}