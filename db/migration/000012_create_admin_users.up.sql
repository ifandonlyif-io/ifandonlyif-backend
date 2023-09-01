create table "admin_users"
(
    "id"       uuid      DEFAULT gen_random_uuid() PRIMARY KEY,
    name       text      not null,
    email      text      not null,
    password   text      not null,
    is_admin   boolean   not null,
    created_at timestamp DEFAULT  now()
);

insert into "admin_users" (name, email, password, is_admin)
    values ('admin', 'admin@tokimi.space', '$2a$12$RWRv.nvhe2.Ku7Z6k3uVOeNFsj9Vgirc3MWc.ozaFCTuPhHDGisJu', true);
-- password: A5haenXNxGNd7H