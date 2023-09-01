create table "admin_user_tokens"
(
    "user_id" uuid not null,
    "token" varchar(255) not null
)