package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/philslol/kritis3m_scale/kritis3m_control/service"
)

type NodeRegisterController interface {
	InitialAssignConfiguration(c *gin.Context)
	InstructedAssignConfiguration(c *gin.Context)
}

type NodeHardbeatController interface {
	RespondHardbeatRequest(c *gin.Context)
}

type LogController interface {
	PushActiveConnections(c *gin.Context)
	PushConnectionRequests(c *gin.Context)
	PushInfo(c *gin.Context)
	PushWarn(c *gin.Context)
	PushErr(c *gin.Context)
	PushDebug(c *gin.Context)
}

// ------------------------   implementations ------------------------------

// ------------------------  Register

func NewNodeRegisterControllerImpl(svc service.NodeRegisterService) NodeRegisterControllerImpl {
	return NodeRegisterControllerImpl{svc: svc}
}

type NodeRegisterControllerImpl struct {
	svc service.NodeRegisterService
}

func (ctrl NodeRegisterControllerImpl) InitialAssignConfiguration(c *gin.Context) {
	ctrl.svc.InitialAssignConfiguration(c)
}
func (ctrl NodeRegisterControllerImpl) InstructedAssignConfiguration(c *gin.Context) {
	ctrl.svc.InstructedAssignConfiguration(c)
}

// ------------------------ Hardbeat

func NewNodeHardbeatControllerImpl(svc service.NodeHardbeatServiceImpl) NodeHardbeatControllerImpl {
	return NodeHardbeatControllerImpl{svc: svc}
}

type NodeHardbeatControllerImpl struct {
	svc service.NodeHardbeatService
}

func (ctrl NodeHardbeatControllerImpl) RespondHardbeatRequest(c *gin.Context) {
	ctrl.svc.RespondHardbeatRequest(c)
}

// ------------------------LogController

func NewLogControllerImpl(svc service.LogService) LogController { return LogControllerImpl{svc: svc} }

type LogControllerImpl struct {
	svc service.LogService
}

func (ctrl LogControllerImpl) PushActiveConnections(c *gin.Context) {
	ctrl.svc.PushActiveConnections(c)
}
func (ctrl LogControllerImpl) PushConnectionRequests(c *gin.Context) {
	ctrl.svc.PushConnectionRequests(c)
}
func (ctrl LogControllerImpl) PushInfo(c *gin.Context)  { ctrl.svc.PushInfo(c) }
func (ctrl LogControllerImpl) PushWarn(c *gin.Context)  { ctrl.svc.PushWarn(c) }
func (ctrl LogControllerImpl) PushErr(c *gin.Context)   { ctrl.svc.PushErr(c) }
func (ctrl LogControllerImpl) PushDebug(c *gin.Context) { ctrl.svc.PushErr(c) }
