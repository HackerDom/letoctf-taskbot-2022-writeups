package storage

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"

	"github.com/HackerDom/letoctf-taskbot-2022-tasks/filestore/pkg/basestore"
	"github.com/HackerDom/letoctf-taskbot-2022-tasks/filestore/pkg/basestore/adapter"
)

var usersTable = goqu.Dialect("postgres").
	From("users").
	Prepared(true)

type DbStorage interface {
	CreateUser(ctx context.Context, username string, pass []byte, id uuid.UUID) (uuid.UUID, error)
	GetUserId(ctx context.Context, username string, pass []byte) (uuid.UUID, error)
}

//go:embed migrations/migration.sql
var migration []byte

func NewDbStorage(connPool *pgx.ConnPool) (DbStorage, error) {
	if _, err := connPool.Exec(string(migration)); err != nil {
		return nil, fmt.Errorf("migrate database failed: %v", err)
	}

	createUserQuery, _, err := usersTable.Insert().
		Rows(goqu.Record{"id": "", "username": "", "pass": ""}).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("prebuild query to create user failed: %v", err)
	}

	getUserIdQuery, _, err := usersTable.Select(
		goqu.C("id"),
	).Where(
		goqu.C("username").Eq("username"),
		goqu.C("pass").Eq("pass"),
	).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("prebuild query to get user id failed: %v", err)
	}

	return &dbStorage{
		BaseStore: basestore.New(
			adapter.NewPGXAdapter(connPool),
		),
		createUserQuery: createUserQuery,
		getUserIdQuery:  getUserIdQuery,
	}, nil
}

type dbStorage struct {
	*basestore.BaseStore

	createUserQuery string
	getUserIdQuery  string
}

func (s *dbStorage) CreateUser(
	ctx context.Context,
	username string,
	pass []byte,
	id uuid.UUID,
) (uuid.UUID, error) {
	if id == uuid.Nil {
		var err error
		id, err = uuid.NewUUID()
		if err != nil {
			return uuid.Nil, fmt.Errorf("generate id failed: %v", err)
		}
	}

	if _, err := s.Q(ctx).Exec(ctx, s.createUserQuery, id.String(), pass, username); err != nil {
		return uuid.Nil, fmt.Errorf("exec query to create user failed: %v", err)
	}

	return id, nil
}

func (s *dbStorage) GetUserId(ctx context.Context, username string, pass []byte) (uuid.UUID, error) {
	var rawId string
	row := s.Q(ctx).QueryRow(ctx, s.getUserIdQuery, username, pass)

	if err := row.Scan(&rawId); err != nil {
		return uuid.Nil, errors.Wrap(err, "scan user id failed")
	}

	return uuid.MustParse(rawId), nil
}
