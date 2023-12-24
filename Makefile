DC := docker-compose -f ./docker-compose.yml

all:
	@mkdir -p /home/mysql
	@$(DC) up

down:
	@$(DC) down

re: clean all

clean:
	@$(DC) down
	@rm -rf /home/mysql
	docker system prune -a -f

db:
	@mkdir -p /home/mysql
	@$(DC) up -d --build mysql

app:
	@mkdir -p /home/mysql
	@$(DC) up -d --build app

.PHONY: all down re clean db app