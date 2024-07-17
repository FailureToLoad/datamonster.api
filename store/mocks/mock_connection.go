package mocks

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type MockConnection struct {
	Rows pgx.Rows
	Row  pgx.Row
	err  error
}

func (c MockConnection) Close() {
	fmt.Println("Close called")
}
func (c MockConnection) Begin(ctx context.Context) (pgx.Tx, error) {
	panic("not implemented")
}
func (c MockConnection) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	tag := pgconn.NewCommandTag("tag")
	if c.err != nil {
		return tag, c.err
	}
	return tag, nil
}
func (c MockConnection) Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error) {
	if c.err != nil {
		return nil, c.err
	}
	if c.Rows == nil {
		panic("rows field not set")
	}
	return c.Rows, nil
}
func (c MockConnection) QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row {
	if c.Row == nil {
		panic("row field not set")
	}

	return c.Row
}

func (c *MockConnection) SetRows(rows pgx.Rows) {
	c.Rows = rows
}

func (c *MockConnection) SetRow(row pgx.Row) {
	c.Row = row
}

func (c *MockConnection) SetError(err error) {
	c.err = err
}
