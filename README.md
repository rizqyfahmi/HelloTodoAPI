An implementation of Clean Architecture in Go projects based on Echo Framework.

## Requirements
- go 1.17.5
- golang-migrate
- docker & docker-compose

## Create Migration
```
migrate create -ext sql -dir db/migration -seq init_schema
```
Where ```sql``` is the extension of the file, and the directory to store it is ```db/migration```. We use the ```-seq``` flag to generate a sequential version number for the migration file and put ```init_scheme``` as the name of migration

## Run Project
- Use ```make run``` to build and run ```docker-compose up``` and migration of the application
- Use ```make close``` to run ```docker-compose down``` and revert migration of the application
