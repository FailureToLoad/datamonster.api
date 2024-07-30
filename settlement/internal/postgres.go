package internal

import (
	"context"
	"fmt"
	"github.com/failuretoload/datamonster/store"
)

type PostgresRepo struct {
	pool store.Connection
}

type Settlement struct {
	Id                  int
	Owner               string
	Name                string
	SurvivalLimit       int
	DepartingSurvival   int
	CollectiveCognition int
	CurrentYear         int
}

func New(d store.Connection) *PostgresRepo {
	return &PostgresRepo{pool: d}
}

func (r PostgresRepo) Select(ctx context.Context, userID string) ([]Settlement, error) {
	query := fmt.Sprintf("SELECT * FROM campaign.settlement WHERE owner='%s'", userID)
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return []Settlement{}, err
	}
	defer rows.Close()
	settlements := []Settlement{}
	for rows.Next() {
		var s Settlement
		err := rows.Scan(&s.Id, &s.Owner, &s.Name, &s.SurvivalLimit, &s.DepartingSurvival, &s.CollectiveCognition, &s.CurrentYear)
		if err != nil {
			return settlements, err
		}
		settlements = append(settlements, s)
	}
	return settlements, nil
}

func (r PostgresRepo) Get(ctx context.Context, id, userId string) (Settlement, error) {
	query := `SELECT * FROM campaign.settlement WHERE id = $1 && owner = $2 LIMIT 1`
	var s Settlement
	err := r.pool.QueryRow(ctx, query, id, userId).Scan(&s.Id, &s.Owner, &s.Name, &s.SurvivalLimit, &s.DepartingSurvival, &s.CollectiveCognition, &s.CurrentYear)
	return s, err
}

func (r PostgresRepo) Insert(ctx context.Context, s Settlement) (int, error) {
	insert := "INSERT INTO campaign.settlement (owner, name, survival_limit, departing_survival, collective_cognition, year) "
	values := fmt.Sprintf("VALUES ('%s', '%s', %d, %d, %d, %d) RETURNING id", s.Owner, s.Name, s.SurvivalLimit, s.DepartingSurvival, s.CollectiveCognition, s.CurrentYear)
	query := insert + values
	id := 0
	err := r.pool.QueryRow(ctx, query).Scan(&id)
	return id, err
}
