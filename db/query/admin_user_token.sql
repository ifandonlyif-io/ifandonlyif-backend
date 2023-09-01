-- name: GetTokenByUserId :one
select * from admin_user_tokens where "user_id" = $1;

-- name: GetUserIdByToken :one
select * from admin_user_tokens where "token" = $1;

-- name: CreateUserToken :one
insert into admin_user_tokens ("user_id", "token") values ($1, $2) RETURNING *;

-- name: DeleteUserToken :exec
delete from admin_user_tokens where "user_id" = $1;

-- name: UpdateUserToken :one
update admin_user_tokens set "token" = $1 where "user_id" = $2 RETURNING *;