### Fleet API
This is a REST API for managing a fleet of aircrafts and flights. It is built with Go, Gin, and Gorm.

Features
CRUD operations for aircrafts and flights
Search flights by departure and arrival airport and date and time range
Generate reports for departure airports and in-flight aircrafts for a given date and time range
Installation
[To run this project, you need to have Go installed on your system. You can download it from here.]

You also need to install the dependencies: Gin and Gorm. You can do that by running:
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/sqlite
Usage
To start the server, run:
go run main.go
The server will listen on port 8080 by default. You can change that by editing the r.Run(":8080") line in the main.go file.

To interact with the API, you can use any HTTP client, such as curl, Postman, or Insomnia. The base URL for the API is http://localhost:8080/api/v1.

The API supports the following endpoints:

GET /aircrafts: Get a list of all aircrafts
GET /aircrafts/:serial_number: Get a single aircraft by its serial number
POST /aircrafts: Create a new aircraft
PUT /aircrafts/:serial_number: Update an existing aircraft by its serial number
DELETE /aircrafts/:serial_number: Delete an existing aircraft by its serial number
GET /flights: Get a list of all flights
GET /flights/:id: Get a single flight by its ID
POST /flights: Create a new flight
PUT /flights/:id: Update an existing flight by its ID
DELETE /flights/:id: Delete an existing flight by its ID
GET /reports: Get a report of departure airports and in-flight aircrafts for a given date and time range
The request and response bodies are in JSON format. Here are some examples of the data structures for each resource:

### Aircraft
{
    "serial_number": "NX429",
    "manufacturer": "Boeing",
    "flights": [
        {
            "id": 1,
            "departure_airport": "LEBL",
            "arrival_airport": "LFPG",
            "departure_date_time": "2024-02-04T10:00:00Z",
            "arrival_date_time": "2024-02-04T12:00:00Z",
            "aircraft_serial_number": "NX429"
        },
        {
            "id": 2,
            "departure_airport": "LFPG",
            "arrival_airport": "EGLL",
            "departure_date_time": "2024-02-04T14:00:00Z",
            "arrival_date_time": "2024-02-04T15:00:00Z",
            "aircraft_serial_number": "NX429"
        }
    ]
}
#### Flight
{
    "id": 1,
    "departure_airport": "LEBL",
    "arrival_airport": "LFPG",
    "departure_date_time": "2024-02-04T10:00:00Z",
    "arrival_date_time": "2024-02-04T12:00:00Z",
    "aircraft_serial_number": "NX429",
    "aircraft": {
        "serial_number": "NX429",
        "manufacturer": "Boeing"
    }
}
### Json
{
    "start_date_time": "2024-02-04T10:00:00Z",
    "end_date_time": "2024-02-04T16:00:00Z",
    "departure_airports": [
        {
            "airport": "LEBL",
            "in_flight_aircraft": 2,
            "aircrafts": [
                {
                    "serial_number": "NX429",
                    "in_flight_time": 120
                },
                {
                    "serial_number": "NX533",
                    "in_flight_time": 90
                }
            ]
        },
        {
            "airport": "LFPG",
            "in_flight_aircraft": 1,
            "aircrafts": [
                {
                    "serial_number": "NX429",
                    "in_flight_time": 60
                }
            ]
        }
    ]
}
### Testings
To run the unit test, run
go test ./...

The tests are located in the files with the suffix _test.go. For example: aircraft_test.go, flight_test.go, report_test.go.
Documentation
To generate the documentation for the API, run:
go doc ./...
