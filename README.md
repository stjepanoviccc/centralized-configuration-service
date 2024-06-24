# Centralized Configuration Service
![Golang](https://img.shields.io/badge/Go-blue?logo=go&logoColor=white)
![Gorilla](https://img.shields.io/badge/Gorilla-yellow?logo=go&logoColor=white)
![Consul](https://img.shields.io/badge/Consul-pink?logo=consul&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-orange?logo=prometheus&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-blue?logo=docker&logoColor=white)
![Jaeger](https://img.shields.io/badge/Jaeger-black?logo=jaeger&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-green?logo=swagger&logoColor=white)
![Postman](https://img.shields.io/badge/Postman-orange?logo=postman&logoColor=black)

## Overview  
This service has been made as an introduction to the domain of distributed systems and cloud computing. It is based on real life problems which come with cloud computing.  

This particular service has been made to solve the problem of automatizing configuration management within our distributed application.
The service enables various microservices/users to access their needed configurations as well as create new configurations, update existing ones and delete them.  

In addition to that we have added various tools used for analytics and DevOps, such as the **Prometheus** library for metrics, **Jaeger** for tracing while managing our workflow with **Github Actions**  
The data is persisted on a K/V NoSQL DB (in our case Consul)

## Technologies Used
The backbone of our service is coded in **Go**  
Data is persisted in **Consul** as Key-Value pairs  
For tracing and metrics we used **Jaeger tracing** paired with **Prometheus**  
Documenting and testing was done using **OpenAPI Swagger**  
Integration and deployment is done using **Github Actions** and **Docker**

## Getting Started:  
Requirements for starting this application are **Docker** and all of it's codependencies, as well as **Docker compose**  
To run the application first navigate to the root folder of the repository and open a terminal with Administrator Privileges within that folder
```shell
docker compose up
```

If you're using linux try running it with elevated privileges

## Postman Collection
[**Postman Docs/Collection**](https://www.postman.com/aleksannder-z/workspace/ars/documentation/30371859-c57d1009-c8cf-4c29-890a-88e412d26ff3) 

## Swagger
**Swagger** is a powerful tool for designing and documenting RESTful web services. Swagger provides a user-friendly interface for developers to visualize and interact with the API's resources without needing to access the backend logic  
This API was made according to the **OpenAPI 2.0** specification.
You can access the documentation of this API via this [link](http://localhost:8000/docs) once the application is running.  

## Idempotency  
**What is Idempotency middleware** ? The idempotency middleware ensures that repeated requests with the same parameters produce the same result, regardless of how many times they are sent. It helps prevent unintended side effects caused by duplicate requests, such as duplicate charges in a payment system or duplicate updates in a database. By generating and storing a unique identifier for each request and its corresponding response, the middleware can check incoming requests against this identifier. If a request with the same identifier is received again, the middleware can retrieve the previous response associated with that identifier and return it without executing the request handler again. This middleware adds an extra layer of reliability and safety to your application, especially in distributed systems where duplicate requests are more likely to occur.  
We are storing Idempotency-Key in our **Consul** DB.  

## Database:  
**Consul** is a NoSQL database designed for storing key-value pairs. We chose Consul for its simplicity and suitability for our project specifications. To access the **Consul UI**, use the port **8500**.  
This will allow you to manage and interact with your persisted data effortlessly.

## Testing:  
We have implemented unit tests for all services in this project.  
These unit tests are designed to ensure the functionality of individual components in isolation.  

## Metrics:
**Prometheus** is an open-source systems monitoring and alerting toolkit designed for reliability and scalability. It collects and stores its metrics as time-series data, providing a powerful query language called **PromQL** to query and visualize the data. You can access **Prometheus UI** on port **9090**.  

**http_total_requests** -> Total number of HTTP requests in last 24h.  
**http_successful_requests** -> Number of successful HTTP requests in last 24h (2xx, 3xx).  
**http_unsuccessful_requests** -> Number of unsuccessful HTTP requests in last 24h (4xx, 5xx).  
**average_request_duration_seconds** -> Average request duration for each endpoint.   
**requests_per_time_unit** -> Number of requests per time unit (e.g., per minute or per second) for each endpoint.  



## Tracing:  
In our application we integrated [Jaeger](https://www.jaegertracing.io/) for distributed tracing to monitor and troubleshoot the performance of our services.

Jaeger gives us insight into the execution flow and helps us with identifying latency problems after our service has been deployed  
The implementation was done using the **Jaeger client for Go** together with the transitive [**OpenTelemetry**](https://opentelemetry.io/) dependencies required
for tracing.


You can access the [Jaeger UI](http://localhost:16686) through this link when the application is running.


## Deploy:  

### Dockerfile  
**We used Multi-Stage build for lighter final image**  
#### BUILD ENVIROMENT  
**FROM golang:1.22-alpine AS build:** Use Go 1.22 on Alpine Linux as the build environment.  
**WORKDIR /app:** Set the working directory inside the container to /app.  
**COPY go.mod go.sum ./:** Copy dependency files into the container.  
**RUN go mod download:** Download Go module dependencies.  
**COPY . .:** Copy all project files into the container.  
**RUN go build -o app .:** Compile the Go application and output the binary as app.    
#### RUNTIME ENVIROMENT    
**FROM alpine:** Use the latest Alpine Linux as the runtime environment.  
**WORKDIR /app:** Set the working directory to /app.  
**COPY --from=build /app/app .:** Copy the built application binary from the build stage.  
**COPY swagger.yaml /app/swagger.yaml:** Copy the swagger.yaml file into the container.  
**EXPOSE 8000:** Expose port 8000 for the application.  
**CMD ["./app"]:** Set the command to run the application.  


### docker-compose.yml  
This docker-compose.yml file is used to define and manage the services required for app.  

**App Service**:  
-image: Specifies the Docker image to use.  
-container_name: Sets the container name.  
-hostname: Sets the hostname for the container.  
-ports: Maps a host port to a container port using an environment variable.  
-depends_on: Ensures that consul and jaeger services are started before this service.  
-networks: Connects the service to a user-defined network.  
-environment: Sets environment variables for the container.    

**Consul Service**:  
-image: Uses the Consul image.  
-ports: Maps the Consul UI port.  
-command: Runs Consul as a server with specific options.  
-volumes: Mounts a volume for persistent data storage.  
-networks: Connects the service to a user-defined network.  

**Prometheus Service**:  
-image: Uses the Prometheus image.  
-ports: Maps the Prometheus UI port.  
-volumes: Mounts directories for Prometheus configuration and data.  
-networks: Connects the service to a user-defined network.  

**Jaeger Service**:  
-image: Uses the Jaeger all-in-one image.  
-ports: Maps Jaeger ports for tracing and the Jaeger UI.  
-networks: Connects the service to a user-defined network.    

## CI pipeline  
Among all of the other tools we used, we also used **Github Actions** to create a CI pipeline.  
In our particular case we have one action which tests, builds and uploads our application to a [**Dockerhub**](https://hub.docker.com/r/aleksannderz57/ars_projekat) remote.

## Authors  

Andrej Stjepanović,  
Software Developer  

Aleksandar Zajac,  
Software Developer  

Dragan Bijelic,  
Software Developer
