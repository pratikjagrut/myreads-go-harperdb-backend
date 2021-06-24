#!/bin/sh

# Set $1=true if DB is not initilized.

if [ $1 == true ]
then
	printf "======================= Initialising database ======================="

	printf "\n+++++++++++ Creating Schema myreads +++++++++++\n"

	curl --location --request POST ${HARPERDB_HOST} \
	--header 'Content-Type: application/json' \
	--header "Authorization: Basic ${BASIC_AUTH_TOKEN}" \
	--data-raw '{
		"operation":"create_schema",
		"schema": "myreads"
	}'

	printf "\n+++++++++++ Creating table users +++++++++++\n"

	curl --location --request POST ${HARPERDB_HOST} \
	--header 'Content-Type: application/json' \
	--header "Authorization: Basic ${BASIC_AUTH_TOKEN}" \
	--data-raw '{
		"operation":"create_table",
		"schema":"myreads",
		"table":"users",
		"hash_attribute": "id"
	}'

	printf "\n+++++++++++ Creating table books +++++++++++\n"

	curl --location --request POST ${HARPERDB_HOST} \
	--header 'Content-Type: application/json' \
	--header "Authorization: Basic ${BASIC_AUTH_TOKEN}" \
	--data-raw '{
		"operation":"create_table",
		"schema":"myreads",
		"table":"books",
		"hash_attribute": "id"
	}'

	printf "\n+++++++++++ Inserting dummy data in users table to initialise +++++++++++\n"

	curl --location --request POST ${HARPERDB_HOST} \
	--header 'Content-Type: application/json' \
	--header "Authorization: Basic ${BASIC_AUTH_TOKEN}" \
	--data-raw '{
	"operation":"sql",
	"sql": "INSERT INTO myreads.users (name, email, password) VALUES('\''dummy'\'', '\''dummy'\'', '\''dymmy'\'')"
	}'

	printf "\n+++++++++++ Inserting dummy data in books table to initialise +++++++++++\n"

	curl --location --request POST ${HARPERDB_HOST} \
	--header 'Content-Type: application/json' \
	--header "Authorization: Basic ${BASIC_AUTH_TOKEN}" \
	--data-raw '{
	"operation":"sql",
	"sql": "INSERT INTO myreads.books (name, userid, status, image, author, description) VALUES('\''dummy'\'', '\''dummy'\'', '\''dummy'\'', '\''dummy'\'', '\''dummy'\'', '\''dummy'\'')"
	}'
fi

printf "\n======================= Starting Server =======================\n"

mkdir images
# HARPERDB_HOST=$1 HARPERDB_UNAME=$3 HARPERDB_PSWD=$4 go run main.go

./myreads