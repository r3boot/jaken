package ircstate

import (
	"log"
)

func (s *State) GetBindingRowId(command, role string) int {
	roleId := s.GetRoleRowId(role)
	if roleId == -1 {
		return -1
	}

	var rowId int
	sql := "SELECT rowid FROM bindings WHERE command = ? AND role_id = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(command, roleId).Scan(&rowId)
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

func (s *State) HasBinding(command, role string) bool {
	return s.GetBindingRowId(command, role) != -1
}

func (s *State) GetRoleId(command string) int {
	var roleId int

	sql := "SELECT role_id FROM bindings WHERE command = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(command).Scan(&roleId)
	switch {
	case err == nil:
		return roleId
	case err.Error() == "sql: no rows in result set":
		return -1
	case err != nil:
		log.Fatalf("statement.QueryRow: %v", err)
	}

	return -1
}

func (s *State) GetBindingRole(command string) string {
	roleId := s.GetRoleId(command)
	if roleId == -1 {
		return ""
	}

	return s.GetRoleById(roleId)
}

func (s *State) AddBinding(command, role string) {
	roleId := s.GetRoleRowId(role)
	if roleId == -1 {
		return
	}

	sql := "INSERT INTO bindings VALUES (?, ?)"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(command, roleId).Scan()
	switch {
	case err.Error() == "sql: no rows in result set":
	case err != nil:
		{
			log.Fatalf("statement.QueryRow: %v", err)
		}
	}

	return
}

func (s *State) RemoveBinding(command, role string) {
	roleId := s.GetRoleRowId(role)
	if roleId == -1 {
		return
	}

	sql := "DELETE FROM bindings WHERE command = ? AND role_id = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(command, roleId).Scan()
	switch {
	case err.Error() == "sql: no rows in result set":
	case err != nil:
		{
			log.Fatalf("statement.QueryRow: %v", err)
		}
	}
}

func (s *State) GetBindingsForRole(role string) []string {
	var (
		bindings []string
	)

	roleId := s.GetRoleRowId(role)
	if roleId == -1 {
		return nil
	}

	sql := "SELECT command FROM bindings WHERE role_id = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	rows, err := statement.Query(&roleId)
	for rows.Next() {
		var command string
		innerErr := rows.Scan(&command)
		if innerErr != nil {
			log.Fatalf("rows.Scan: %v", err)
		}
		bindings = append(bindings, command)
	}
	switch {
	case err == nil:
		return bindings
	case err.Error() == "sql: no rows in result set":
		return nil
	case err != nil:
		log.Fatalf("statement.QueryRow: %v", err)
	}

	return nil
}
