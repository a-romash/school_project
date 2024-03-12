package postgresql

import (
	"context"
	"project/pkg/model"

	"github.com/jackc/pgx/v5"
)

func (db *Postgresql) GetTest(id string) (test *model.Test, err error) {
	const sql = `
	SELECT * FROM tests WHERE id = $1
	`

	rows, _ := db.pool.Query(context.Background(), sql, id)
	tests, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Test])
	if err != nil {
		return nil, err
	}
	if len(tests) == 0 {
		return nil, pgx.ErrNoRows
	}
	test = &tests[0]
	return test, nil
}

func (db *Postgresql) CreateNewTest(test *model.Test) (err error) {
	const sql = `
	INSERT INTO tests (id, author, author_id, questions, answers)
  	VALUES ($1, $2, $3, $4, $5);
	`

	_, err = db.pool.Exec(context.Background(), sql, test.Id, test.Author, test.AuthorId, test.Questions, test.Answers)
	if err != nil {
		return err
	}
	return nil
}
