package models

import "errors"

var (
	// Returned when Get method does not find any sql records matching the query
	ErrorNoRecord = errors.New("models: no matching record found")

	// Returned when a user tries to login with the wrong email or password
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// Returned when a user tries to signup with an already used email
	ErrDuplicatedEmail = errors.New("models: duplicate email")
)
