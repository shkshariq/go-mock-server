
# Go Mock API

This project is a simple and flexible mock API server built using Go. It allows users to create and update dynamic API endpoints with custom responses, and serves those endpoints with the stored responses.

## Prerequisites

- Go 1.13 or higher
- MySQL database

## Installation

1. Download the project source code:

```bash
$ git clone https://github.com/yourusername/yourrepository.git
$ cd yourrepository
```

2. Install the required Go packages:

```bash
$ go get github.com/go-sql-driver/mysql
```

3. Update the `initDB` function in the source code with your MySQL credentials.

4. Build the project:

```bash
$ go build -o go-mock-api
```

5. Run the built binary:

```bash
$ ./go-mock-api
```

This will start the mock API server on port 8080.

## Usage

### Create or update a mock API

To create or update a mock API, send a POST request to `/create_update_api?api_path=/api/mock-test` with the following JSON payload:
```json
{
	"status_code": 200,
	"headers": {
		"Content-Type": "application/json"
	},
	"body": {
		"message": "Hello Gopher"
	}
}
```

Replace the values in the payload with the desired path, status code, headers, and response body.


### Access the mock API

To access the mock API, send a request to the corresponding path:

```bash
curl "http://localhost:8080/api/mock-test"
```

Response
```json
{
		"message": "Hello Gopher"
}
```

## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT).