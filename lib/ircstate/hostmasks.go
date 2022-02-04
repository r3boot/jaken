package ircstate

import "log"

func (s *State) GetHostmaskRowId(hostmask string) int {
	var rowId int
	sql := "SELECT rowid FROM hostmasks WHERE hostmask = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(hostmask).Scan(&rowId)
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

func (s *State) HasHostmask(hostmask string) bool {
	return s.GetNickNameRowId(hostmask) != -1
}

func (s *State) GetHostmaskNicknameId(hostmask string) int {
	var rowId int
	sql := "SELECT rowid FROM hostmasks WHERE hostmask = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(hostmask).Scan(&rowId)
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

func (s *State) GetNicknameForHostmask(hostmask string) string {
	var nicknameId int
	sql := "SELECT nickname_id FROM hostmasks WHERE hostmask = ?"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(hostmask).Scan(&nicknameId)
	switch {
	case err == nil:
		return s.GetNickNameById(nicknameId)
	case err.Error() == "sql: no rows in result set":
		return ""
	case err != nil:
		log.Fatalf("statement.QueryRow: %v", err)
	}

	return ""
}

func (s *State) AddHostmaskToNickname(hostmask, nickname string) {
	nicknameId := s.GetNickNameRowId(nickname)

	sql := "INSERT INTO hostmasks VALUES (?, ?)"

	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatalf("db.Prepare: %v", err)
	}
	defer statement.Close()

	err = statement.QueryRow(hostmask, nicknameId).Scan()
	switch {
	case err.Error() == "sql: no rows in result set":
	case err != nil:
		{
			log.Fatalf("statement.QueryRow: %v", err)
		}
	}
}
