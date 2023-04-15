# Make migrations with goose and CONNECTION_STRING env variable
goose -dir ./db/migrations postgres "$CONNECTION_STRING" up
/project/go-docker/build/myapp