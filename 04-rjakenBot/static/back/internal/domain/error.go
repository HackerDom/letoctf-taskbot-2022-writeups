package domain

type UserError interface {
	error
	UserError() string
}
