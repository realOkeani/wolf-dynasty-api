package models

import "time"

type Owner struct {
	ID        string    `json:"id" db:"ID"`
	Name      string    `json:"name" db:"NAME"`
	CreatedAt time.Time `json:"createdAt" db:"CreatedAt"`
	UpdatedAt time.Time `json:"updatedAt" db:"UpdatedAt"`
	Email 	  string 	`json:"email" db:"EMAIL"`
}
