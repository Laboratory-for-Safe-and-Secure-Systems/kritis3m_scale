package service

import (
	"github.com/gin-gonic/gin"
	"github.com/philslol/kritis3m_scale/kritis3m_control/db"
	"github.com/philslol/kritis3m_scale/kritis3m_control/types"
	"github.com/rs/zerolog/log"
	//types
)

type NodeRegisterService interface {
	InitialAssignConfiguration(c *gin.Context)
	InstructedAssignConfiguration(c *gin.Context)
}

type NodeRegisterServiceImpl struct {
	db *db.KSDatabase
}

func NewNodeRegisterServiceImpl(ks_db *db.KSDatabase) NodeRegisterServiceImpl {
	return NodeRegisterServiceImpl{db: ks_db}
}

func (svc NodeRegisterServiceImpl) InitialAssignConfiguration(c *gin.Context) {
	serial_number, err := get_serialnumber(c)
	if err != nil {
		log.Err(err).Msg("cant get identity of request")

		c.Set("caller_error", true)
		c.Error(types.ErrInvalidParam).SetMeta("DistributionService: InitialAssignConfiguration")
		return
	}
	log.Info().Msgf("Initial Configuration Request from Node with Serialnumber %s", serial_number)
	//Get Config:

	node, err := svc.db.GetNodeby_SerialNumber(serial_number)
	if err != nil {
		log.Err(err).Msgf("No Node with Serialnumber %s in Database", serial_number)
		c.Set("internal_error", true)
		c.Error(types.ErrUnauthorized).SetMeta("DistributionService: InitialAssignConfiguration")
		return
	}

	//check if active config exists
	//TODO:check this function
	config, err := svc.db.GetConfigFor_DistributionService(node.ID)
	if err != nil {
		log.Err(err).Msgf("No Configuration available for Node with Serialnumber %s ", serial_number)
		c.Set("internal_error", true)
		c.Error(types.ErrInternalError).SetMeta("DistributionService: InitialAssignConfiguration")
		return
	}
	var endpoint_ids []uint
	//get relevant endpoints
	for _, app := range config.Application {
		endpoint_ids = addToListIfNotPresent(endpoint_ids, app.Ep1ID)
		endpoint_ids = addToListIfNotPresent(endpoint_ids, app.Ep2ID)
	}

	//returnvalue
	var distribution_response types.DistributionResponse
	for _, ep_id := range endpoint_ids {
		ep, err := svc.db.GetEPby_ID(ep_id)
		if err != nil {
			log.Err(err).Msgf("No EP Configuration found for config id  %d ", ep_id)
			c.Set("internal_error", true)
			c.Error(types.ErrInternalError).SetMeta("DistributionService: InitialAssignConfiguration")
			return
		}
		distribution_response.CryptoConfig = append(distribution_response.CryptoConfig, ep)
	}

	var identity_ids []uint
	for _, ep := range distribution_response.CryptoConfig {
		identity_ids = addToListIfNotPresent(identity_ids, ep.IdentityID)
	}
	distribution_response.Identities, err = svc.db.GetIdentities(identity_ids)
	if err != nil {
		log.Err(err).Msgf("Problem with identities ")
		c.Set("internal_error", true)
		c.Error(types.ErrInternalError).SetMeta("DistributionService: InitialAssignConfiguration")
		return
	}

	var lconfig []*types.DBNodeConfig
	lconfig = append(lconfig, config)
	node.Config = lconfig
	distribution_response.Node = *node
	//send out
	c.JSON(200, distribution_response)
}

func (svc NodeRegisterServiceImpl) InstructedAssignConfiguration(c *gin.Context) {
	identity, err := get_identity(c)
	if err != nil {
		log.Err(err).Msg("Identity is not complete. Request is missing either serialnumber or cfg id or version/updatedate")
		c.Set("caller_error", true)
		c.Error(types.ErrUnauthorized).SetMeta("DistributionService: Identity Incomplete")
		return
	}
	log.Info().Msgf("Instructed Configuration Request from Node with Serialnumber %s", identity.Serialnumber)
	//Get Config:
	node, err := svc.db.GetNodeby_SerialNumber(identity.Serialnumber)
	if err != nil {
		log.Err(err).Msgf("No Node with Serialnumber %s in Database", identity.Serialnumber)
		c.Set("internal_error", true)
		c.Error(types.ErrUnauthorized).SetMeta("DistributionService: InstructedAssignConfiguration")
		return
	}
	//check if active config exists
	//TODO:check this function
	//TODO:handle no config available
	config, err := svc.db.GetConfigFor_DistributionService(node.ID)
	if err != nil {
		log.Err(err).Msgf("No Configuration available for Node with Serialnumber %s ", identity.Serialnumber)
		c.Set("internal_error", true)
		c.Error(types.ErrInternalError).SetMeta("DistributionService: InitialAssignConfiguration")
		return
	}
	var endpoint_ids []uint
	//get relevant endpoints
	for _, app := range config.Application {
		endpoint_ids = addToListIfNotPresent(endpoint_ids, app.Ep1ID)
		endpoint_ids = addToListIfNotPresent(endpoint_ids, app.Ep2ID)
	}
	//returnvalue
	var distribution_response types.DistributionResponse
	for _, ep_id := range endpoint_ids {
		ep, err := svc.db.GetEPby_ID(ep_id)
		if err != nil {
			log.Err(err).Msgf("No EP Configuration found for config id  %d ", ep_id)
			c.Set("internal_error", true)
			c.Error(types.ErrInternalError).SetMeta("DistributionService: InitialAssignConfiguration")
			return
		}
		distribution_response.CryptoConfig = append(distribution_response.CryptoConfig, ep)
	}
	var lconfig []*types.DBNodeConfig
	lconfig = append(lconfig, config)
	node.Config = lconfig
	distribution_response.Node = *node
	//send out
	c.JSON(200, distribution_response)

}
func addToListIfNotPresent(list []uint, item uint) []uint {
	if item == 0 {
		return list
	}

	for _, v := range list {
		if v == item {
			return list // Item is already in the list, return unchanged list
		}
	}
	// Item is not in the list, append it
	return append(list, item)
}
