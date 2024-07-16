package mocks

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type MockRows struct {
	Rows    []pgx.Row
	cusrsor int
}

func (sr MockRows) Close() {

}
func (sr *MockRows) Err() error {
	panic("implement me")
}
func (sr *MockRows) CommandTag() pgconn.CommandTag {
	panic("implement me")
}
func (sr *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	panic("implement me")
}
func (sr *MockRows) Next() bool {
	return sr.cusrsor < len(sr.Rows)
}
func (sr *MockRows) Scan(dest ...any) error {
	err := sr.Rows[sr.cusrsor].Scan(dest...)
	if err != nil {
		return err
	}
	sr.cusrsor++
	return nil
}
func (sr *MockRows) Values() ([]any, error) {
	panic("implement me")
}
func (sr *MockRows) RawValues() [][]byte {
	panic("implement me")
}
func (sr *MockRows) Conn() *pgx.Conn {
	panic("implement me")
}

type ErrorRow struct {
	Error error
}

func (er *ErrorRow) Scan(dest ...any) error {
	return er.Error
}

type InsertRow struct {
	Id int
}

func (ir *InsertRow) Scan(dest ...any) error {
	id := dest[0].(*int)
	*id = ir.Id
	return nil
}
