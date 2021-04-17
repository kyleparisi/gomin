#!/bin/bash
set -e

export $(cat .env.default | grep -v ^# | xargs);

echo Starting services:
docker start $APP_NAME-mysql || docker run --name $APP_NAME-mysql -e MYSQL_ROOT_PASSWORD=$DB_PASSWORD -e MYSQL_DATABASE=$DB_DATABASE -p $DB_PORT:$DB_PORT -d mysql
until mysql -h $DB_HOST -u $DB_USERNAME -p$DB_PASSWORD --silent -e "STATUS;"
do
  echo "Waiting for database connection..."
  sleep 1
done
mysql -h $DB_HOST -u $DB_USERNAME -p$DB_PASSWORD --silent -e "DROP DATABASE IF EXISTS $DB_DATABASE;"
mysql -h $DB_HOST -u $DB_USERNAME -p$DB_PASSWORD --silent -e "CREATE DATABASE IF NOT EXISTS $DB_DATABASE;"

echo Running migrations:
./schema/manage.py makemigrations
./schema/manage.py migrate
python generator.py > src/app/types.go

echo Running format:
go fmt

echo Running tests:
APP_DIR=$(pwd) go test -v -p 1 ./...
