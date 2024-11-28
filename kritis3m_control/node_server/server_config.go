package node_server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/philslol/kritis3m_scale/kritis3m_control/controller"
	"github.com/rs/zerolog"
)

// ZerologMiddleware returns a gin.HandlerFunc that logs requests using zerolog.
func ZerologMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process the request
		c.Next()

		// Log request and response details based on log level
		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		// Choose log level based on status code
		logEntry := logger.With().
			Str("method", method).
			Str("path", path).
			Str("client_ip", clientIP).
			Int("status", status).
			Dur("duration", duration).
			Logger()

		// Log at different levels depending on the HTTP status
		if status >= 500 {
			// Log as an error for server errors
			logEntry.Error().Msg("Server error occurred")
		} else if status >= 400 {
			// Log as a warning for client errors
			logEntry.Warn().Msg("Client error occurred")
		} else if logger.GetLevel() == zerolog.DebugLevel {
			// Log as debug for all requests if the log level is Debug
			logEntry.Debug().Msg("Request handled")
		} else {
			// Log as info for successful responses
			logEntry.Info().Msg("Request handled")
		}
	}
}

func Init(ctrl_logger controller.LogController,
	ctrl_heartbeat controller.NodeHeartbeatController,
	ctrl_register controller.NodeRegisterController,
	logger zerolog.Logger, mode string) *gin.Engine {

	gin.SetMode(mode)
	router := gin.New()
	router.Use(ZerologMiddleware(logger))

	// router.Use(service.ErrorHandler(log.Logger))
	router.Use(gin.Recovery())
	api := router.Group("/api/node/:serialnumber")
	{
		initial := api.Group("/initial")
		initial.GET("/register", ctrl_register.InitialAssignConfiguration)

		operation := api.Group("/operation/version/:version_number")
		operation.GET("register/instructed", ctrl_register.InstructedAssignConfiguration)
		operation.GET("/heartbeat", ctrl_heartbeat.RespondHeartbeatRequest)

		logging_api := api.Group("/logger")
		logging_api.POST("/active_con", ctrl_logger.PushActiveConnections)
		logging_api.POST("/con_request", ctrl_logger.PushConnectionRequests)
		logging_api.POST("/err", ctrl_logger.PushErr)
		logging_api.POST("/info", ctrl_logger.PushInfo)
		logging_api.POST("/warn", ctrl_logger.PushWarn)
		logging_api.POST("/debug", ctrl_logger.PushDebug)
	}
	return router

}
