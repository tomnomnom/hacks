package main

import (
	"database/sql"
	"errors"
	"strings"
)

type domains struct {
	db *sql.DB
}

func (d *domains) names() []string {
	return []string{"domain", "domains"}
}

func (d *domains) setDB(db *sql.DB) {
	d.db = db
}

func (d *domains) initModule() error {
	_, err := d.db.Exec(`
		create table if not exists domains (
			id integer primary key,
			domain text unique not null
		)
	`)

	return err
}

func (d *domains) Add(domain string) error {

	if domain == "" {
		return errors.New("no domain provided")
	}
	domain = normalise(domain)

	_, err := d.db.Exec("insert or ignore into domains (domain) values(?)", domain)
	return err
}

func (d *domains) Delete(domain string) error {
	if domain == "" {
		return errors.New("no domain provided")
	}
	domain = normalise(domain)

	_, err := d.db.Exec("delete from domains where domain = ?", domain)
	return err
}

func (d *domains) All() ([]string, error) {
	rows, err := d.db.Query("select domain from domains")
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	out := []string{}
	for rows.Next() {
		var domain string
		err = rows.Scan(&domain)
		if err != nil {
			return out, err
		}

		out = append(out, domain)
	}

	return out, rows.Err()
}

func normalise(domain string) string {
	return strings.ToLower(domain)
}
