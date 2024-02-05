package main

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

// Aircraft represents an aircraft with a serial number and a manufacturer
type Aircraft struct {
    SerialNumber string `json:"serial_number" gorm:"primaryKey"`
    Manufacturer string `json:"manufacturer"`
    Flights      []Flight `json:"flights" gorm:"foreignKey:AircraftSerialNumber"`
}

// Flight represents a flight with a departure and arrival airport, date and time, and an aircraft
type Flight struct {
    ID                    uint   `json:"id" gorm:"primaryKey"`
    DepartureAirport      string `json:"departure_airport"`
    ArrivalAirport        string `json:"arrival_airport"`
    DepartureDateTime     string `json:"departure_date_time"`
    ArrivalDateTime       string `json:"arrival_date_time"`
    AircraftSerialNumber  string `json:"aircraft_serial_number"`
    Aircraft              Aircraft `json:"aircraft" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Connect to a SQLite database
db, err := gorm.Open(sqlite.Open("fleet.db"), &gorm.Config{})
if err != nil {
    panic("failed to connect to database")
}

// Migrate the models
db.AutoMigrate(&Aircraft{}, &Flight{})

// Create a Gin router with default middleware
r := gin.Default()

// Use a group to organize the routes for the /albums resource
v1 := r.Group("/api/v1")
{
    v1.GET("/aircrafts", getAircrafts)
    v1.GET("/aircrafts/:serial_number", getAircraft)
    v1.POST("/aircrafts", createAircraft)
    v1.PUT("/aircrafts/:serial_number", updateAircraft)
    v1.DELETE("/aircrafts/:serial_number", deleteAircraft)
    v1.GET("/flights", getFlights)
    v1.GET("/flights/:id", getFlight)
    v1.POST("/flights", createFlight)
    v1.PUT("/flights/:id", updateFlight)
    v1.DELETE("/flights/:id", deleteFlight)
    v1.GET("/reports", getReports)
}

// Run the server on port 8080
r.Run(":8080")

// getAircrafts returns a list of all aircrafts as JSON
func getAircrafts(c *gin.Context) {
    // Get the database from the context
    db := c.MustGet("db").(*gorm.DB)

    // Declare a variable to store the aircrafts
    var aircrafts []Aircraft

    // Find all aircrafts and preload their flights
    db.Preload("Flights").Find(&aircrafts)

    // Return the aircrafts as JSON
    c.JSON(http.StatusOK, gin.H{"data": aircrafts})
}

// getAircraft returns a single aircraft by its serial number as JSON
func getAircraft(c *gin.Context) {
    // Get the database from the context
    db := c.MustGet("db").(*gorm.DB)

    // Get the serial number from the URL parameter
    serialNumber := c.Param("serial_number")

    // Declare a variable to store the aircraft
    var aircraft Aircraft

    // Find the aircraft by its serial number and preload its flights
    if err := db.Preload("Flights").Where("serial_number = ?", serialNumber).First(&aircraft).Error; err != nil {
        // If the aircraft is not found, return an error
        c.JSON(http.StatusNotFound, gin.H{"error": "Aircraft not found"})
        return
    }

    // Return the aircraft as JSON
    c.JSON(http.StatusOK, gin.H{"data": aircraft})
}

// createAircraft adds a new aircraft from JSON received in the request body
func createAircraft(c *gin.Context) {
    // Get the database from the context
    db := c.MustGet("db").(*gorm.DB)

    // Validate the input
    var input CreateAircraftInput
    if err := c.ShouldBindJSON(&input); err != nil {
        // If the input is invalid, return an error
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Create an aircraft with the input data
    aircraft := Aircraft{SerialNumber: input.SerialNumber, Manufacturer: input.Manufacturer}

    // Save the aircraft to the database
    if err := db.Create(&aircraft).Error; err != nil {
        // If the save fails, return an error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return the created aircraft as JSON
    c.JSON(http.StatusCreated, gin.H{"data": aircraft})
}

// updateAircraft updates an existing aircraft by its serial number with JSON received in the request body
func updateAircraft(c *gin.Context) {
    // Get the database from the context
    db := c.MustGet("db").(*gorm.DB)

    // Get the serial number from the URL parameter
    serialNumber := c.Param("serial_number")

    // Declare a variable to store the aircraft
    var aircraft Aircraft

    // Find the aircraft by its serial number
    if err := db.Where("serial_number = ?", serialNumber).First(&aircraft).Error; err != nil {
        // If the aircraft is not found, return an error
        c.JSON(http.StatusNotFound, gin.H{"error": "Aircraft not found"})
        return
    }

    // Validate the input
    var input UpdateAircraftInput
    if err := c.ShouldBindJSON(&input); err != nil {
        // If the input is invalid, return an error
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Update the aircraft with the input data
    if err := db.Model(&aircraft).Updates(input).Error; err != nil {
        // If the update fails, return an error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return the updated aircraft as JSON
    c.JSON(http.StatusOK, gin.H{"data": aircraft})
}

// deleteAircraft deletes an existing aircraft by its serial number
func deleteAircraft(c *gin.Context) {
    // Get the database from the context
    db := c.MustGet("db").(*gorm.DB)

    // Get the serial number from the URL parameter
    serialNumber := c.Param("serial_number")

    // Declare a variable to store the aircraft
    var aircraft Aircraft

    // Find the aircraft by its serial number
    if err := db.Where("serial_number = ?", serialNumber).First(&aircraft).Error; err != nil {
        // If the aircraft is not found, return an error
        c.JSON(http.StatusNotFound, gin.H{"error": "Aircraft not found"})
        return
    }

    // Delete the aircraft from the database
    if err := db.Delete(&aircraft).Error; err != nil {
        // If the delete fails, return an error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return a success message
    c.JSON(http.StatusOK, gin.H{"data": "Aircraft deleted"})
}

// getFlights returns a list of all flights as JSON
func getFlights(c *gin.Context) {
    // Get the database from the context
    db := c.MustGet("db").(*gorm.DB)

    // Declare a variable to store the flights
    var flights []Flight

    // Find all flights and preload their aircrafts
    db.Preload("Aircraft").Find(&flights)

    // Return the flights as JSON
    c.JSON(http.StatusOK, gin.H{"data": flights})
}

// getFlight returns a single flight by its ID as JSON
func getFlight(c *gin.Context) {
    // Gets the database from the context
    db := c.MustGet("db").(*gorm.DB)

    // Gets the ID from the URL parameter
    id := c.Param("id")

    // Declares a variable to store the flight
    var flight Flight

    // Find the flight by its ID and preload its aircraft
    if err := db.Preload("Aircraft").Where("id = ?", id).First(&flight).Error; err != nil {
        // If the flight is not found, returnS an error
        c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
        return
    }
