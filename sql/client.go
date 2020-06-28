package sql

import (
	"github.com/jmoiron/sqlx"
	"github.com/realOkeani/wolf-dynasty-api/models"
)

//go:generate counterfeiter . Client
type Client interface {
	GetOwners() ([]models.Owner, error)
	AddOwner(models.Owner) (models.Owner, error)
	GetOwner(string) (models.Owner, error)
	UpdateOwner(models.Owner) (models.Owner, error)
	DeleteOwner(models.Owner) error
}

type client struct {
	*sqlx.DB
}

// func NewOwnersClient(db *sqlx.DB) Client {
// 	db.MustExec(createOwnersTable)

// 	return &client{
// 		DB: db,
// 	}
// }
