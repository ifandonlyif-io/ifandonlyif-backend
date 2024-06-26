// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: gas_price.sql

package db

import (
	"context"
	"database/sql"
)

const createGasPrice = `-- name: CreateGasPrice :one
INSERT INTO gas_prices (
    average
) VALUES (
  $1
) RETURNING id, average, created_at
`

func (q *Queries) CreateGasPrice(ctx context.Context, average sql.NullString) (GasPrice, error) {
	row := q.db.QueryRowContext(ctx, createGasPrice, average)
	var i GasPrice
	err := row.Scan(&i.ID, &i.Average, &i.CreatedAt)
	return i, err
}

const getAveragePriceByLastDay = `-- name: GetAveragePriceByLastDay :many
 SELECT COALESCE(average) AS average,
  COALESCE(created_at) AS created_at
  FROM gas_prices
  ORDER BY created_at DESC
  LIMIT 24
`

type GetAveragePriceByLastDayRow struct {
	Average   sql.NullString `json:"average"`
	CreatedAt sql.NullTime   `json:"createdAt"`
}

func (q *Queries) GetAveragePriceByLastDay(ctx context.Context) ([]GetAveragePriceByLastDayRow, error) {
	rows, err := q.db.QueryContext(ctx, getAveragePriceByLastDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAveragePriceByLastDayRow{}
	for rows.Next() {
		var i GetAveragePriceByLastDayRow
		if err := rows.Scan(&i.Average, &i.CreatedAt); err != nil {
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
