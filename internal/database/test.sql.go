// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: test.sql

package database

import (
	"context"
)

const getCities = `-- name: GetCities :many
SELECT id, name, test FROM test
`

func (q *Queries) GetCities(ctx context.Context) ([]Test, error) {
	rows, err := q.db.QueryContext(ctx, getCities)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Test
	for rows.Next() {
		var i Test
		if err := rows.Scan(&i.ID, &i.Name, &i.Test); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
