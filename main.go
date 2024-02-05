package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Aircraft struct {
	SerialNumber string   `json:"serial_number" gorm:"primaryKey"`
	Manufacturer string   `json:"manufacturer"`
	Flights      []Flight `json:"flights" gorm:"foreignKey:AircraftSerialNumber"`
}

type Flight struct {
	ID                   uint     `json:"id" gorm:"primaryKey"`
	DepartureAirport     string   `json:"departure_airport"`
	ArrivalAirport       string   `json:"arrival_airport"`
	DepartureDateTime    string   `json:"departure_date_time"`
	ArrivalDateTime      string   `json:"arrival_date_time"`
	AircraftSerialNumber string   `json:"aircraft_serial_number"`
	Aircraft             Aircraft `json:"aircraft" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func main() {
	// Connect to a SQLite database
	db, err := gorm.Open(sqlite.Open("fleet.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the models
	db.AutoMigrate(&Aircraft{}, &Flight{})

	r := gin.Default()

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
}

func getAircrafts(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var aircrafts []Aircraft

	db.Preload("Flights").Find(&aircrafts)

	c.JSON(http.StatusOK, gin.H{"data": aircrafts})
}

func getAircraft(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	serialNumber := c.Param("serial_number")

	var aircraft Aircraft

	if err := db.Preload("Flights").Where("serial_number = ?", serialNumber).First(&aircraft).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aircraft not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": aircraft})
}

func createAircraft(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	CreateAircraftInput := 
	var input CreateAircraftInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aircraft := Aircraft{SerialNumber: input.SerialNumber, Manufacturer: input.Manufacturer}

	if err := db.Create(&aircraft).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": aircraft})
}

func updateAircraft(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	serialNumber := c.Param("serial_number")

	var aircraft Aircraft

	if err := db.Where("serial_number = ?", serialNumber).First(&aircraft).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aircraft not found"})
		return
	}

	var input UpdateAircraftInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&aircraft).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": aircraft})
}

func deleteAircraft(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	serialNumber := c.Param("serial_number")

	var aircraft Aircraft

	if err := db.Where("serial_number = ?", serialNumber).First(&aircraft).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aircraft not found"})
		return
	}

	if err := db.Delete(&aircraft).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Aircraft deleted"})
}

func getFlights(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var flights []Flight

	db.Preload("Aircraft").Find(&flights)

	c.JSON(http.StatusOK, gin.H{"data": flights})
}

func getFlight(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")

	var flight Flight

	if err := db.Preload("Aircraft").Where("id = ?", id).First(&flight).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}
}

type CreateFlightInput struct {
	DepartureAirport     string `json:"departure_airport"`
	ArrivalAirport       string `json:"arrival_airport"`
	DepartureDateTime    string `json:"departure_date_time"`
	ArrivalDateTime      string `json:"arrival_date_time"`
	AircraftSerialNumber string `json:"aircraft_serial_number"`
}

func createFlight(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var input CreateFlightInput
	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a flight with the input data
	flight := Flight{
		DepartureAirport:     input.DepartureAirport,
		ArrivalAirport:       input.ArrivalAirport,
		DepartureDateTime:    input.DepartureDateTime,
		ArrivalDateTime:      input.ArrivalDateTime,
		AircraftSerialNumber: input.AircraftSerialNumber,
	}

	// Save the flight to the database
	if err := db.Create(&flight).Error; err != nil {
		// If the save fails, return an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the created flight as JSON
	c.JSON(http.StatusOK, gin.H{"data": flight})
}
