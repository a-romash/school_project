package database

import (
	"project/pkg/config"
	"project/pkg/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	RegisterNewUser(*model.User) error
	GetUser(login string) (*model.User, error)
	UpdateUser(*model.User) error
	DeleteUser(id int) error

	NewToken(login string) (*model.Token, error)
	ValidateToken(token uuid.UUID) error
	DeleteToken(token uuid.UUID) error
	DeleteTokensByLogin(login string) error

	CreateNewTest(*model.Test) error
	GetSolutionsForTest(id int) (*model.Solution, error)
	GetTest(id string) (*model.Test, error)
	PutSolutionToTest(*model.Test) error
	UpdateTest(*model.Test) error
	UpdateStatistics(*model.Solution) error

	Init(config *config.IConfig) (*pgxpool.Config, error)
}
