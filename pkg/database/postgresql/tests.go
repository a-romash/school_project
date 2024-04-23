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

func (db *Postgresql) UpdateTest(data map[string]interface{}) (err error) {
	const sql_one = `
	DELETE FROM solutions
	WHERE test_id = $1;
	`

	const sql_two = `
	UPDATE tests
	SET title=$1,
		max_score=$2,
		questions=$3,
		answers=$4,
		updated=CURRENT_TIMESTAMP
	WHERE id=$5;
	`

	_, err = db.pool.Exec(context.Background(), sql_one, data["test_id"].(string))
	if err != nil {
		return err
	}

	_, err = db.pool.Exec(context.Background(), sql_two, data["title"], len(data["answers"].([]interface{})), data["questions"], data["answers"], data["test_id"])
	if err != nil {
		return err
	}

	return nil
}

func (db *Postgresql) GetIdTests(userId int) (id []string, err error) {
	const sql = `
	SELECT id
	FROM tests
	WHERE author_id = $1;
	`

	rows, _ := db.pool.Query(context.Background(), sql, userId)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var uuid string
		err := rows.Scan(&uuid)
		if err != nil {
			return []string{}, err
		}
		id = append(id, uuid)
	}

	return id, nil
}

func (db *Postgresql) DeleteTest(test_id string) (err error) {
	const sql_one = `
	DELETE FROM solutions
	WHERE test_id = $1;
	`

	const sql_two = `
	DELETE FROM tests
	WHERE id = $1;
	`

	_, err = db.pool.Exec(context.Background(), sql_one, test_id)
	if err != nil {
		return err
	}

	_, err = db.pool.Exec(context.Background(), sql_two, test_id)
	if err != nil {
		return err
	}
	return nil
}

func (db *Postgresql) CreateNewTest(test *model.Test) (err error) {
	const sql = `
	INSERT INTO tests (id, title, author, author_id, max_score, questions, answers, conversion)
  	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`

	_, err = db.pool.Exec(context.Background(), sql, test.Id, test.Title, test.Author, test.AuthorId, test.Max_score, test.Questions, test.Answers, test.Conversion)
	if err != nil {
		return err
	}
	return nil
}

func (db *Postgresql) CreateNewSolution(solution *model.Solution) (err error) {
	const sql = `
	INSERT INTO solutions (author, class, answers, result, test_id, grade)
  	VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err = db.pool.Exec(context.Background(), sql, solution.Author, solution.Class, solution.Answers, solution.Result, solution.TestId, solution.Grade)
	if err != nil {
		return err
	}
	return nil
}

func (db *Postgresql) GetSolutionsForTest(test_id string) (solutions []model.Solution, err error) {
	const sql = `
	SELECT * FROM solutions
	WHERE test_id=$1
	`

	rows, _ := db.pool.Query(context.Background(), sql, test_id)
	solutions, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.Solution])
	if err != nil {
		return nil, err
	}
	if len(solutions) == 0 {
		return nil, pgx.ErrNoRows
	}
	return solutions, nil
}
