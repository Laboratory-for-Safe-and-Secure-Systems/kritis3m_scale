package node_server

import (
	"github.com/gin-gonic/gin"
	"github.com/philslol/kritis3m_scale/kritis3m_control/controller"
)

func Init(ctrl_logger controller.LogController,
	ctrl_heartbeat controller.NodeHeartbeatController,
	ctrl_register controller.NodeRegisterController) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())

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
