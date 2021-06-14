#!/bin/sh

set -o pipefail

printf "======================= Initialising database ======================="

printf "\n+++++++++++ Creating Schema myreads +++++++++++\n"

curl --location --request POST $1 \
--header 'Content-Type: application/json' \
--header "Authorization: Basic $2" \
--data-raw '{
	"operation":"create_schema",
	"schema": "myreads"
}'

printf "\n+++++++++++ Creating table users +++++++++++\n"

curl --location --request POST $1 \
--header 'Content-Type: application/json' \
--header "Authorization: Basic $2" \
--data-raw '{
	"operation":"create_table",
	"schema":"myreads",
	"table":"users",
	"hash_attribute": "id"
}'

printf "\n+++++++++++ Creating table books +++++++++++\n"

curl --location --request POST $1 \
--header 'Content-Type: application/json' \
--header "Authorization: Basic $2" \
--data-raw '{
	"operation":"create_table",
	"schema":"myreads",
	"table":"books",
	"hash_attribute": "id"
}'

printf "\n+++++++++++ Inserting dummy data in users table to initialise +++++++++++\n"

curl --location --request POST $1 \
--header 'Content-Type: application/json' \
--header "Authorization: Basic $2" \
--data-raw '{
  "operation":"sql",
  "sql": "INSERT INTO myreads.users (name, email, password) VALUES('\''dummy'\'', '\''summy'\'', '\''dymmy'\'')"
}'

printf "\n+++++++++++ Inserting dummy data in books table to initialise +++++++++++\n"

curl --location --request POST $1 \
--header 'Content-Type: application/json' \
--header "Authorization: Basic $2" \
--data-raw '{
  "operation":"sql",
  "sql": "INSERT INTO myreads.books (name, userid, status, image, author) VALUES('\''dummy'\'', '\''dummy'\'', '\''dummy'\'', '\''dummy'\'', '\''dummy'\'')"
}'

printf "\n======================= Starting Server =======================\n"

HARPERDB_HOST=$1 HARPERDB_UNAME=$3 HARPERDB_PSWD=$4 go run main.go