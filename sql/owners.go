package sql

import "github.com/realOkeani/wolf-dynasty-api/models"

var (
	createOwnersTable = `
		CREATE TABLE IF NOT EXISTS Teams (
			ID VARCHAR(255),
			NAME VARCHAR(255),
			CreatedAt DATETIME,
			UpdatedAt DATETIME,
			EMAIL VARCHAR(255),
			PRIMARY KEY (ID)
		)`

	insertOwner = `
		INSERT INTO Owners (ID, NAME, CreatedAt, UpdatedAt, EMAIL) VALUES (?, ?, ?, ?, ?)
  `
	getAllOwners = `SELECT ID, NAME, CreatedAt, UpdatedAt, EMAIL FROM Owners`

	getOwner = `SELECT ID, NAME, CreatedAt, UpdatedAt, EMAIL FROM Owners WHERE ID = ?`

	updateOwner = `UPDATE Owners SET NAME = ?, EMAIL = ?, UpdatedAt = ? WHERE ID = ?`

	deleteOwner = `DELETE FROM Owners WHERE ID=?`
)

func (c client) AddOwner(owner models.Owner) (models.Owner, error) {
	_, err := c.DB.Exec(insertOwner,
		owner.ID,
		owner.Name,
		owner.CreatedAt,
		owner.UpdatedAt,
		owner.Email)

	return owner, err
}

func (c client) GetOwners() ([]models.Owner, error) {
	var owners []models.Owner

	err := c.DB.Select(&owners, getAllOwners)
	if err != nil {
		return owners, err
	}

	return owners, nil
}

func (c client) GetOwner(ownerID string) (models.Owner, error) {
	var owner models.Owner

	err := c.DB.Get(&owner, getOwner, ownerID)
	if err != nil {
		return owner, err
	}

	return owner, nil
}

func (c client) UpdateOwner(owner models.Owner) (models.Owner, error) {
	_, err := c.DB.Exec(updateOwner,
		owner.Name,
		owner.UpdatedAt,
		owner.ID,
		owner.Email)

	return owner, err
}

func (c client) DeleteOwner(owner models.Owner) error {
	_, err := c.DB.Exec(deleteOwner,
		owner.ID)

	return err
}