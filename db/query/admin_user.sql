-- name: GetAllAdminUsers :many
select * from admin_users;

-- name: GetAdminUserByID :one
select * from admin_users where "id" = $1 limit 1;

-- name: GetAdminUserByEmail :one
select * from admin_users where "email" = $1 limit 1;

-- name: CreateAdminUser :one
insert into admin_users ("name", "email", "password", "is_admin") values ($1, $2, $3, $4) RETURNING *;

-- name: DeleteAdminUser :exec
delete from admin_users where "id" = $1;

-- name: UpdateAdminUser :one
update admin_users set "name" = $2, "email" = $3, "password" = $4 where "id" = $1 RETURNING *;
