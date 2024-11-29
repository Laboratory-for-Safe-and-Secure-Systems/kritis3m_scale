package service

import (
	"strconv"
	"time"

	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/types"
	"github.com/gin-gonic/gin"
)

type Identity struct {
	Serialnumber string
	Config_id    uint
	Updated_at   time.Time
}

// not production ready
func get_identity(c *gin.Context) (Identity, error) {
	var identity Identity
	node := c.Param("serialnumber")
	config := c.Param("configID")
	updated_at := c.Param("time_stamp")

	if node == "" || config == "" || updated_at == "" {
		return identity, types.ErrNotFound
	}

	config_id, err := strconv.ParseUint(config, 10, 64)
	if err != nil {
		return Identity{}, err
	}
	identity.Config_id = uint(config_id)
	identity.Serialnumber = node

	updatedAtNumber, err := strconv.ParseInt(updated_at, 10, 64)
	if err != nil {
		return Identity{}, err
	}
	// Convert the timestamp to time.Time (assuming it's in seconds)
	identity.Updated_at = time.Unix(updatedAtNumber, 0)

	return identity, nil

}

func get_serialnumber(c *gin.Context) (string, error) {
	node := c.Param("serialnumber")

	if node == "" {
		return node, types.ErrNotFound
	}
	return node, nil

}
