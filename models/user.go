package models

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
}

// CREATE TABLE users
// (
// id bigint unsigned NOT NULL AUTO_INCREMENT,
// name longtext,
// email varchar(191) UNIQUE,
// password longtext,
// PRIMARY KEY (id)
// );
