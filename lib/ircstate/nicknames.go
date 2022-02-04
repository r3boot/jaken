package ircstate

import "log"

func (s *State) GetNickNameRowId(nickname string) int {
	var rowId int
	sql := "SELECT rowid FROM nicknames WHERE nickname = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(nickname).Scan(&rowId)
	switch {
	case err == nil:
		return rowId
	case err.Error() == "sql: no rows in result set":
		return -1
	case err != nil:
		log.Fatalf("statement.QueryRow: %v", err)
	default:
		return -1
	}

	return -1
}

func (s *State) HasNickname(nickname string) bool {
	return s.GetNickNameRowId(nickname) != -1
}

func (s *State) GetNickNameById(rowId int) string {
	var nickname string
	sql := "SELECT nickname FROM nicknames WHERE rowid = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(rowId).Scan(&nickname)
	switch {
	case err == nil:
		return nickname
	case err.Error() == "sql: no rows in result set":
		return ""
	case err != nil:
		log.Fatalf("statement.QueryRow: %v", err)
	default:
		return ""
	}

	return ""
}

func (s *State) AddNickname(nickname string) {
	sql := "INSERT INTO nicknames VALUES (?)"

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

func (s *State) RemoveNickname(nickname string) {
	sql := "DELETE FROM nicknames WHERE nickname = ?"

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
