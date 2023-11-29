# Transaction Processor

This project includes a Lambda function that processes transaction records from an S3 bucket, interacts with a PostgreSQL database, and sends summary reports via email.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Docker and Docker Compose
- AWS CLI
- AWS SAM CLI
- Go (version specified in `go.mod`) Developed in 1.21.4

### Installing

A step-by-step series of examples that tell you how to get a development environment running.

1. **Clone the repository**

```sh
git clone https://github.com/DiegoSan99/tsx-processor
cd transaction-processor
```

2. **Start the PostgreSQL database**

```sh
docker-compose up -d
```

3.  **Init database**
    Make sure to apply the database schema and any initial data setups such as the init-db.sql script.

4.  **Local Lambda Testing**

    a. Ensure you have a template.yaml for the project. The one provided works too

    b. Prepare a test-event.json file with the structure expected by your Lambda function. You can use the one provided also.

    c. Build the main.go file for a linux image. The GOARCH variable might vary depending on you local environment

    ```sh
    GOOS=linux GOARCH=amd64 go build -o main cmd/tps/main.go
    ```

    d. Invoke the function with the next command

    ```sh
    sam local invoke "TransactionProcessorFunction" -e test-event.json --env-vars env.json
    ```
