package model

import "github.com/google/uuid"

type Solution struct {
	Author  string    `db:"author" json:"author"`
	Class   string    `db:"class" json:"class"`
	Answers []string  `db:"answers" json:"answers"`
	Result  int       `db:"result" json:"result"`
	TestId  uuid.UUID `db:"test_id" json:"test_id"`
	Id      int       `db:"id" json:"id"`
}
