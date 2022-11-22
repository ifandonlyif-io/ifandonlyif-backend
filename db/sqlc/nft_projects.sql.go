// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: nft_projects.sql

package db

import (
	"context"
)

const listNftProjects = `-- name: ListNftProjects :many
SELECT id, name, contract_address, collection_name, image_uri FROM nft_projects
`

func (q *Queries) ListNftProjects(ctx context.Context) ([]NftProject, error) {
	rows, err := q.db.QueryContext(ctx, listNftProjects)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []NftProject{}
	for rows.Next() {
		var i NftProject
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ContractAddress,
			&i.CollectionName,
			&i.ImageUri,
		); err != nil {
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