package sql

import (
	"github.com/jmoiron/sqlx"
	"github.com/realOkeani/wolf-dynasty-api/models"
)

//go:generate counterfeiter . Client
type Client interface {
	GetOwners() ([]models.Owners, error)
}

type client struct {
	*sqlx.DB
}

func NewTeamsClient(db *sqlx.DB) Client {
	db.MustExec(createOwnersTable)

	return &client{
		DB: db,
	}
}

