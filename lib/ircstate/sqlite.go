package ircstate

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const (
	initializeTableSql = `
CREATE TABLE IF NOT EXISTS nicknames (
    nickname varchar
);

CREATE TABLE IF NOT EXISTS hostmasks (
    hostmask varchar,
    nickname_id integer,
    FOREIGN KEY (nickname_id) REFERENCES nicknames(rowid)
		ON DELETE CASCADE
		ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS roles (
    name varchar not null
);

CREATE TABLE IF NOT EXISTS permissions (
    nickname_id integer,
    role_id integer,
    FOREIGN KEY (nickname_id) REFERENCES nicknames (rowid)
		ON DELETE CASCADE
		ON UPDATE NO ACTION,
    FOREIGN KEY (role_id) REFERENCES roles (rowid) 
		ON DELETE CASCADE 
		ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS bindings (
    command varchar,
	role_id integer,
	FOREIGN KEY (role_id) REFERENCES roles(rowid)
		ON DELETE CASCADE
    	ON UPDATE NO ACTION
);

DELETE FROM roles WHERE name = 'admin';
INSERT INTO roles VALUES ('admin');

DELETE FROM bindings WHERE role_id = 1;
INSERT INTO bindings VALUES
    ('test', 1),
	('meet', 1),
	('forget', 1),
	('role', 1),
    ('perm', 1),
	('binding', 1);
`
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
