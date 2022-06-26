.PHONY: fill-data
fill-data:
	docker exec -i insider_project mysql -u root -proot < ./services/database/data.sql

.PHONY: create-db
create-db:
	docker exec -i insider_project mysql -u root -proot < ./services/database/create.sql

.PHONY: setup-db
setup-db:\
	create-db\
	fill-data

.PHONY: docker-up
docker-up:
	docker-compose up -d

.PHONY: wait
wait:
	sleep 10

.PHONY: api-up
api-up:\
	docker-up\
	wait\
	setup-db

