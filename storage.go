package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5"
)

type PostgresStorage struct {
	conn *pgx.Conn
}

func NewPostgresStorage(ctx context.Context, dbUrl string) (*PostgresStorage, error) {
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		conn: conn,
	}, nil
}

func (p *PostgresStorage) NextSeqValue(ctx context.Context) (*big.Int, error) {
	row := p.conn.QueryRow(ctx, "SELECT nextval('serial')")

	var valStr string
	if err := row.Scan(&valStr); err != nil {
		return nil, err
	}

	val := &big.Int{}
	val, ok := val.SetString(valStr, 10)
	if !ok {
		return nil, fmt.Errorf("invalid sequence value %s", valStr)
	}

	return val, nil
}

func (p *PostgresStorage) Save(ctx context.Context, url string, slug string) error {
	_, err := p.conn.Exec(ctx,
		"INSERT INTO urls (url, slug) VALUES (@url, @slug)",
		pgx.NamedArgs{
			"url":  url,
			"slug": slug,
		},
	)

	return err
}

func (p *PostgresStorage) Get(ctx context.Context, slug string) (string, error) {
	row := p.conn.QueryRow(ctx, "SELECT url FROM urls WHERE slug = @slug", pgx.NamedArgs{"slug": slug})

	var url string
	if err := row.Scan(&url); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return url, nil
}
