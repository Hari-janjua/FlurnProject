package Routes

import (
	"FlurnProject/Controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/seats", Controllers.GetAllSeats)
	router.GET("/seats/:id", Controllers.GetSeatDetailsById)
	router.POST("/booking", Controllers.CreateBooking)
	router.GET("/bookings", Controllers.GetBookingDetails)

	return router

}
