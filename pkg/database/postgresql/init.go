package postgresql

import (
	"context"
	"log"
	"log/slog"
	"project/pkg/config"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Postgresql struct {
	pool      *pgxpool.Pool
	closeOnce sync.Once
}

func Connect() (db *Postgresql, err error) {
	config := Config()
	return ConnectWithConfig(config)
}

func ConnectWithConfig(config *pgxpool.Config) (db *Postgresql, err error) {
	for i := 0; i < 5; i++ {
		p, err := pgxpool.NewWithConfig(context.Background(), config)
		if err != nil || p == nil {
			time.Sleep(3 * time.Second)
			continue
		}
		log.Printf("pool returned from connect: idk from where so i am really lazy for normal logs tho")
		db = &Postgresql{
			pool: p,
		}
		err = Init(db.pool)
		if err != nil {
			slog.Error("error initing database")
			return nil, err
		}
		slog.Info("database was successfully init")
		return db, nil
	}
	err = errors.Wrap(err, "timed out waiting to connect postgres")
	slog.Error("timed out waiting to connect postgres")
	return nil, err
}

func (db *Postgresql) Close() {
	db.closeOnce.Do(func() {
		db.pool.Close()
	})
}

func Config() *pgxpool.Config {
	const defaultMaxConns = int32(10)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	// Your own Database URL
	DATABASE_URL := config.Config.Postgres_conn

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		slog.Debug("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		slog.Debug("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		slog.Debug("Closed the connection pool to the database!!")
	}

	return dbConfig
}

func Init(p *pgxpool.Pool) (err error) {
	const sql string = `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		lastname VARCHAR(255) NOT NULL,
		school VARCHAR(255) NOT NULL,
		hashedPassword BYTEA NOT NULL
	);

	CREATE TABLE IF NOT EXISTS tokens (
		login VARCHAR(255) REFERENCES users(login),
		token UUID PRIMARY KEY,
		expires_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS tests (
		id UUID PRIMARY KEY,
		author VARCHAR(255) NOT NULL,
		author_id INT REFERENCES users(id),
		solutions_id INT[],
		questions JSONB[] NOT NULL,
		answers VARCHAR(255)[] NOT NULL,
		created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
		updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS solutions (
		author VARCHAR(255) NOT NULL,
		class VARCHAR(255) NOT NULL,
		answers VARCHAR(255)[] NOT NULL, 
		result INT NOT NULL,
		test_id UUID REFERENCES tests(id),
		id SERIAL PRIMARY KEY
	);
	`

	_, err = p.Exec(context.Background(), sql)
	return err
}
