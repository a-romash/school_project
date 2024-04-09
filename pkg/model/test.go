package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Test struct {
	Author    string        `db:"author" json:"author"`
	AuthorId  int           `db:"author_id" json:"author_id"`
	Title     string        `db:"title" json:"title"`
	Questions []interface{} `db:"questions" json:"questions"`
	Id        uuid.UUID     `db:"id" json:"id"`
	Answers   []interface{} `db:"answers" json:"answers"`
	Max_score int           `db:"max_score" json:"max_score"`
	Created   time.Time     `db:"created" json:"created"`
	Updated   time.Time     `db:"updated" json:"updated"`
}

func CreateTest(author, title string, authorId int, questions []interface{}, answers []interface{}) (test *Test) {
	test = &Test{
		Author:    author,
		AuthorId:  authorId,
		Title:     title,
		Questions: questions,
		Max_score: len(questions),
		Answers:   answers,
		Id:        uuid.New(),
	}
	return test
}

func (t *Test) GetJson() (data []byte, err error) {
	data, err = json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return data, nil
}
