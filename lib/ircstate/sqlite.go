package ircstate

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const (
	initializeTableSql = `CREATE TABLE IF NOT EXISTS users (hostmask varchar[64] not null primary key, nickname varchar[64]);`
)

type State struct {
	storage string
	db      *sql.DB
}

func New(fname string) (*State, error) {
	var err error

	state := &State{
		storage: fname,
	}

	state.db, err = sql.Open("sqlite3", state.storage)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	_, err = state.db.Exec(initializeTableSql)
	if err != nil {
		log.Printf("%q: %s\n", err, initializeTableSql)
		return nil, fmt.Errorf("db.Exec: %v", err)
	}

	return state, nil
}

func (s *State) Close() {
	s.db.Close()
}

func (s *State) HasNickname(requestedNickname string) bool {
	var nickname string
	sql := "SELECT nickname FROM users WHERE nickname = ? LIMIT 1"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(requestedNickname).Scan(&nickname)
	switch {
	case nickname == requestedNickname:
		return true
	case err.Error() == "sql: no rows in result set":
		return false
	case err != nil:
		log.Fatalf("statement.QueryRow: %v", err)
	}

	return false
}

func (s *State) HasHostmask(requestedHostmask string) bool {
	var hostmask string
	sql := "SELECT hostmask FROM users WHERE hostmask = ? LIMIT 1"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(requestedHostmask).Scan(&hostmask)
	switch {
	case hostmask == requestedHostmask:
		return true
	case err.Error() == "sql: no rows in result set":
		return false
	case err != nil:
		log.Fatalf("statement.QueryRow: %v", err)
	}

	return false
}

func (s *State) AddUser(nickname, hostmask string) {
	sql := "INSERT INTO users VALUES (?, ?)"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(hostmask, nickname).Scan()
	switch {
	case err.Error() == "sql: no rows in result set":
	case err != nil:
		{
			log.Fatalf("statement.QueryRow: %v", err)
		}
	}
}

func (s *State) RemoveUser(nickname string) {
	sql := "DELETE FROM users WHERE nickname = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(nickname).Scan()
	switch {
	case err.Error() == "sql: no rows in result set":
	case err != nil:
		{
			log.Fatalf("statement.QueryRow: %v", err)
		}
	}
}
