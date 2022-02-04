package ircstate

import "log"

func (s *State) GetRoleRowId(role string) int {
	var rowId int
	sql := "SELECT rowid FROM roles WHERE name = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(role).Scan(&rowId)
	switch {
	case err == nil:
		return rowId
	case err.Error() == "sql: no rows in result set":
		return -1
	case err != nil:
		log.Fatalf("statement.QueryRow: %v", err)
	}

	return -1
}

func (s *State) HasRole(role string) bool {
	return s.GetRoleRowId(role) != -1
}

func (s *State) GetRoleById(rowId int) string {
	var name string
	sql := "SELECT name FROM roles WHERE rowid = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(rowId).Scan(&name)
	switch {
	case err == nil:
		return name
	case err.Error() == "sql: no rows in result set":
		return ""
	case err != nil:
		log.Fatalf("statement.QueryRow: %v", err)
	}

	return ""
}

func (s *State) AddRole(role string) {
	sql := "INSERT INTO roles VALUES (?)"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(role).Scan()
	switch {
	case err.Error() == "sql: no rows in result set":
	case err != nil:
		{
			log.Fatalf("statement.QueryRow: %v", err)
		}
	}
}

func (s *State) RemoveRole(role string) {
	sql := "DELETE FROM roles WHERE name = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(role).Scan()
	switch {
	case err.Error() == "sql: no rows in result set":
	case err != nil:
		{
			log.Fatalf("statement.QueryRow: %v", err)
		}
	}
}

func (s *State) GetRoles() []string {
	var (
		roles []string
	)

	sql := "SELECT name FROM roles"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	rows, err := statement.Query()
	for rows.Next() {
		var name string
		innerErr := rows.Scan(&name)
		if innerErr != nil {
			log.Fatalf("rows.Scan: %v", err)
		}
		roles = append(roles, name)
	}
	switch {
	case err == nil:
		return roles
	case err.Error() == "sql: no rows in result set":
		return nil
	case err != nil:
		log.Fatalf("statement.QueryRow: %v", err)
	}

	return nil
}
