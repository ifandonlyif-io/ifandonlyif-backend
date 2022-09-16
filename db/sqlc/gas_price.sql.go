// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
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

func (q *Queries) CreateGasPrice(ctx context.Context, average sql.NullInt32) (GasPrice, error) {
	row := q.db.QueryRowContext(ctx, createGasPrice, average)
	var i GasPrice
	err := row.Scan(&i.ID, &i.Average, &i.CreatedAt)
	return i, err
}

const getAveragePriceByLastDay = `-- name: GetAveragePriceByLastDay :many
  SELECT created_at,
    average
  FROM gas_prices
  ORDER BY created_at DESC
  LIMIT 24
`

type GetAveragePriceByLastDayRow struct {
	CreatedAt sql.NullTime  `json:"createdAt"`
	Average   sql.NullInt32 `json:"average"`
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
		if err := rows.Scan(&i.CreatedAt, &i.Average); err != nil {
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
