postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine
createdb:
	docker exec postgres14 createdb --username=root --owner=root iffdev
dropdb:
	docker exec -it postgres14 dropdb iffdev
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/iffdev?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/iffdev?sslmode=disable" -verbose down
sqlc: 
	sqlc generate

proto:
	rm -f doc/swagger/*.swagger.json \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=ifandonlyif
# proto/*.proto
# statik -src=./doc/swagger -dest=./doc

.PHONY: postgres createdb dropdb migrateup migratedown proto
