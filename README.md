# ConcurrentTaskService
Concurrent TaskService is a Go application designed to handle incoming HTTP requests concurrently using Goroutines. The project implements a worker pool with concurrently

## Project Setup and Usage

### Requirements

##### Go 1.21.5
##### Docker and Docker Compose installed on your system.

#### Make Commands:
#### all: Builds and starts both the application and database services.
#### down: Stops and removes the running containers.
#### re: Cleans up all containers and volumes, then rebuilds and starts everything from scratch.
#### clean: Removes all containers, volumes, and unused Docker images for a thorough cleanup. " Warning: This will remove all Docker images on your system."
#### db: Starts only the MySQL database service.
#### app: Starts only the application service.

## Accessing Swagger UI:

Once the application is running access the Swagger UI documentation at: http://localhost:8080/swagger
## Using Postman Collection:

A Postman collection has been included for convenient API testing. Import the collection to explore and interact with the API endpoints.

#### Additional Information:
The MySQL database data is persisted in the /home/mysql directory on your host machine.
To rebuild specific services, use docker-compose up --build <service-name> (e.g., docker-compose up --build app).
