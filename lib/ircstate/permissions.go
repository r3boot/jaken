package ircstate

import (
	"fmt"
	"log"
)

func (s *State) GetPermissionRowId(nickname, role string) int {
	nicknameId := s.GetNickNameRowId(nickname)
	roleId := s.GetRoleRowId(role)

	if nicknameId == -1 || roleId == -1 {
		return -1
	}

	var rowId int
	sql := "SELECT rowid FROM permissions WHERE nickname_id = ? AND role_id = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(nicknameId, roleId).Scan(&rowId)
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

func (s *State) HasPermission(nickname, role string) bool {
	return s.GetPermissionRowId(nickname, role) != -1
}

func (s *State) AddPermission(nickname, role string) error {
	nicknameId := s.GetNickNameRowId(nickname)
	roleId := s.GetRoleRowId(role)

	if nicknameId == -1 || roleId == -1 {
		return fmt.Errorf("nickname or role not found")
	}

	sql := "INSERT INTO permissions VALUES (?, ?)"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(nicknameId, roleId).Scan()
	switch {
	case err.Error() == "sql: no rows in result set":
	case err != nil:
		{
			log.Fatalf("statement.QueryRow: %v", err)
		}
	}

	return nil
}

func (s *State) RemovePermission(nickname, role string) error {
	nicknameId := s.GetNickNameRowId(nickname)
	roleId := s.GetRoleRowId(role)

	sql := "DELETE FROM permissions WHERE nickname_id = ? AND role_id = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(nicknameId, roleId).Scan()
	switch {
	case err.Error() == "sql: no rows in result set":
	case err != nil:
		{
			log.Fatalf("statement.QueryRow: %v", err)
		}
	}

	return nil
}

func (s *State) RemovePermissions(nickname string) error {
	nicknameId := s.GetNickNameRowId(nickname)

	sql := "DELETE FROM permissions WHERE nickname_id = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(nicknameId).Scan()
	switch {
	case err.Error() == "sql: no rows in result set":
	case err != nil:
		{
			log.Fatalf("statement.QueryRow: %v", err)
		}
	}

	return nil
}

func (s *State) ListPermissions(nickname string) []string {
	var (
		permissions []string
	)

	nicknameId := s.GetNickNameRowId(nickname)

	sql := "SELECT role_id FROM permissions WHERE nickname_id = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	rows, err := statement.Query(nicknameId)
	for rows.Next() {
		var roleId int
		innerErr := rows.Scan(&roleId)
		if innerErr != nil {
			log.Fatalf("rows.Scan: %v", err)
		}
		role := s.GetRoleById(roleId)

		permissions = append(permissions, role)
	}
	switch {
	case err == nil:
		return permissions
	case err.Error() == "sql: no rows in result set":
		return nil
	case err != nil:
		log.Fatalf("statement.QueryRow: %v", err)
	}

	return nil
}
