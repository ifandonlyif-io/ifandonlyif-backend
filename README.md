 ## Installation
 `go mod tidy`
 
`go get github.com/stretchr/testify`

[sqlc](https://docs.sqlc.dev/en/stable/overview/install.html)

`go get github.com/labstack/echo/v4`

[swagger](http://localhost:1323/swagger/index.html#)
 ## Commands
 
>   install docker postgres
`make postgres`

`make createdb`

`make dropdb`

`make migrateup`

`make migratedown`

> generates sql query for go 
`make sqlc` 