package model

// Model user represents data about a record user.
type User struct {
	Id        int    `db:"ID" json:"id"`
	Firstname string `db:"FIRSTNAME" json:"firstname"`
	Lastname  string `db:"LASTNAME" json:"lastname"`
	Email     string `db:"EMAIL" json:"email"`
	Tel       string `db:"TEL" json:"tel"`
}
