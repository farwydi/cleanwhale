package tonic

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const healthPath = "/health"

// AddHealthHandler add to gin.Engine HealthHandler.
func AddHealthHandler(r *gin.Engine) {
	r.GET(healthPath, HealthHandler)
}

// HealthHandler make handler initiates status 200.
var HealthHandler = func(c *gin.Context) {
	c.Status(http.StatusOK)
}
