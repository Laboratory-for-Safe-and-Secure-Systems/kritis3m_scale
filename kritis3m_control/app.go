package kritis3m_control

import (
	// "io"
	// "log"

	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	asl "github.com/Laboratory-for-Safe-and-Secure-Systems/go-asl"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	// "os"

	"net/http"

	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/controller"
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/db"
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/node_server"
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/service"
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/types"
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/util"
)

//tasks: start server
//       start db
//       start ip allocator

// Headscale represents the base app of the service.
type Kritis3m_Scale struct {
	cfg             *types.Config
	db              *db.KSDatabase
	log_db          *db.KSDatabase
	server          http.Server
	insecure_server http.Server
}

func NewKritis3m_scale(cfg *types.Config) (*Kritis3m_Scale, error) {
	log.Info().Msg("In function new kritis3m_scale")

	app := Kritis3m_Scale{
		cfg: cfg,
		// noisePrivateKey: noisePrivateKey,
	}
	var err error
	gormlog := log.With().Str("service", "gorm").Logger().Level(cfg.Database.Level)
	app.db, err = db.NewKritis3mScaleDatabase(cfg.Database, gormlog)
	if err != nil {
		log.Err(err).Msg("error initing db")
		return nil, err
	}
	// app.log_db, err = db.NewLogDatabase(cfg.Log_Database)

	serv_log := service.NewLogServiceImpl(app.log_db)
	serv_reg := service.NewNodeRegisterServiceImpl(app.db)
	ctrl_log := controller.NewLogControllerImpl(serv_log)

	ctrl_reg := controller.NewNodeRegisterControllerImpl(serv_reg)

	ginLogger := log.With().Str("service", "gin").Logger().Level(app.cfg.NodeServer.Log.Level)
	router := node_server.Init(ctrl_log, ctrl_reg, ginLogger, cfg.NodeServer.GinMode)

	err = asl.ASLinit(&cfg.ASLConfig)
	if err != nil {
		log.Err(err).Msg("err asl init")
		panic(err)
	}
	app.server = http.Server{
		Handler: router,
	}
	app.insecure_server = http.Server{
		Handler: router}
	return &app, nil
}

func (ks *Kritis3m_Scale) Import() {
	path := "./db_startup_data/startup.json"
	path = util.AbsolutePathFromConfigPath(path)
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println("error occured")
	}
	defer jsonFile.Close()
	var parsed types.ImportStructure
	bytevalue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(
		[]byte(bytevalue),
		&parsed,
	)
	if err != nil {
		log.Err(err)
		return
	}

	_, err = ks.db.AddIdentitys(parsed.Identites)
	if err != nil {
		log.Err(err).Msg("error occured during import")
		return
	}

	_, err = ks.db.AddEPs(parsed.CryptoConfig)
	if err != nil {
		log.Err(err).Msg("error occured during import")
		return
	}

	for _, node := range parsed.Node {
		err = ks.db.AddNode(node)
		if err != nil {
			log.Err(err).Msg("cant add node to db")
		}
		for _, config := range node.Config {
			// Ensure config.Application is a slice (if it's not, adjust accordingly)
			for _, app := range config.Application {
				_, err = ks.db.AddTrustedClientsto_Application(app, config.Whitelist.TrustedClients)
				if err != nil {
					log.Err(err).Msg("can't update trusted clients")
				}
			}
		}
	}
}

func (ks *Kritis3m_Scale) Listconfigs(cfg_id int, includeAppls bool) {
	var configs []*types.DBNodeConfig
	var err error
	if cfg_id == -1 {
		configs, err = ks.db.GetAllConfigs()
		if err != nil {
			log.Err(err).Msg("error getting nodes")
			return
		}
	} else {
		config, err := ks.db.GetConfigby_ID(uint(cfg_id))
		if err != nil {
			log.Err(err).Msg("error getting node")
			return
		}
		configs = append(configs, config)
	}

	// Print the table header
	if includeAppls {
		fmt.Printf("%-10s %-20s %-10s %-10s %-10s %-20s %-20s %-10s %-10s %-10s\n",
			"ID", "Version", "Whitelist ID", "Appl ID", "State", "Type", "Listening IP:Port", "Client IP:Port", "Ep1 ID", "Ep2 ID")
	} else {
		fmt.Printf("%-10s %-20s %-10s\n", "ID", "Version", "Whitelist ID")
	}

	// Iterate over each configuration
	for _, config := range configs {
		if includeAppls {
			// Retrieve applications related to this config
			appls, err := ks.db.GetApplicationsByCfgID(config.ID)
			if err != nil {
				log.Err(err).Msgf("couldn't get applications of config with ID %d", config.ID)
				return
			}

			// Print each application for the current configuration
			if len(appls) > 0 {
				for _, appl := range appls {
					fmt.Printf("%-10d %-20d %-10d %-10d %-10t %-20s %-20s %-20s %-10d %-10d\n",
						config.ID, config.Version, config.Whitelist.ID,
						appl.ID, appl.State, appl.Type.String(), appl.ServerEndpointAddr, appl.ClientEndpointAddr, appl.Ep1ID, appl.Ep2ID)
				}
			} else {
				// No applications, print config without app details
				fmt.Printf("%-10d %-20d %-10d %-10s %-10s %-10s %-20s %-20s %-10s %-10s\n",
					config.ID, config.Version, config.Whitelist.ID,
					"N/A", "N/A", "N/A", "N/A", "N/A", "N/A", "N/A")
			}
		} else {
			// Print configuration without application details
			fmt.Printf("%-10d %-20d %-10d\n", config.ID, config.Version, config.Whitelist.ID)
		}
	}
}

func (ks *Kritis3m_Scale) ListActive() {
	selectedCfgs, err := ks.db.GetActiveConfigs()
	if err != nil {
		log.Err(err).Msg("can't get active configs")
		return
	}

	// Print the table header
	fmt.Printf("%-10s %-10s\n", "Node ID", "Config ID")
	fmt.Println(strings.Repeat("-", 20)) // Print a separator line

	// Iterate over each selected configuration and print it
	for _, cfg := range selectedCfgs {
		fmt.Printf("%-10d %-10d\n", cfg.NodeID, cfg.ConfigID)
	}
}

func (ks *Kritis3m_Scale) ActivateConfig(node_id uint, cfg_id uint) {
	cfg, err := ks.db.GetConfigby_ID(cfg_id)
	if err != nil {
		log.Err(err).Msgf("couldnt fetch matching configuration of cfg_id %d", cfg_id)
		return
	}

	if cfg.NodeID != node_id {
		log.Error().Msgf("cfg with id %d, does not belong to node with id %d", cfg_id, node_id)
		return
	}

	node, err := ks.db.GetNodeby_ID(node_id)
	if err != nil {
		return
	}

	ks.db.ActivateConfig_byCfgID(cfg_id, node.SerialNumber)
}

// launches a GIN server with Kritis3m_api
func (ks *Kritis3m_Scale) ListNodes(id int, includeConfig bool) {
	var nodes []*types.DBNode
	var err error
	if id == -1 {
		nodes, err = ks.db.GetAllNodes()
		if err != nil {
			log.Err(err).Msg("error getting nodes")
			return
		}
	} else {
		node, err := ks.db.GetNodeby_ID(uint(id))
		if err != nil {
			log.Err(err).Msg("error getting node")
			return
		}
		nodes = append(nodes, node)
	}
	// Print the table header
	if includeConfig {
		fmt.Printf("%-10s %-20s %-15s %-10s\n", "ID", "Serial Number", "Network Index", "Config IDs")
	} else {
		fmt.Printf("%-10s %-20s %-15s\n", "ID", "Serial Number", "Network Index")
	}

	// Print each node
	for _, node := range nodes {
		configs, err := ks.db.GetAllConfigsOfNodeby_ID(node.ID)
		if err != nil {
			log.Err(err).Msgf("couldn't get configs of node with id %d", node.ID)
			return
		}

		if includeConfig {
			// Check if there are any configurations for the node
			if len(configs) > 0 {
				// Print the node details with the first config ID
				fmt.Printf("%-10d %-20s %-15d %-10d\n", node.ID, node.SerialNumber, node.NodeNetworkIndex, configs[0].ID)
				// Print any additional config IDs on new lines with aligned columns
				for _, config := range configs[1:] {
					fmt.Printf("%-10s %-20s %-15s %-10d\n", "", node.SerialNumber, "", config.ID)
				}
			} else {
				// If no config IDs exist for the node
				fmt.Printf("%-10d %-20s %-15d %-10s\n", node.ID, node.SerialNumber, node.NodeNetworkIndex, "N/A")
			}
		} else {
			fmt.Printf("%-10d %-20s %-15d\n", node.ID, node.SerialNumber, node.NodeNetworkIndex)
		}
	}
}

// launches a GIN server with Kritis3m_api
func (ks *Kritis3m_Scale) Serve() {

	go func() {
		serverEndpoint := asl.ASLsetupServerEndpoint(&ks.cfg.NodeServer.ASL_Endpoint)
		if serverEndpoint == nil {
			fmt.Println("Error setting up server endpoint")
			os.Exit(1)
		}
		defer asl.ASLFreeEndpoint(serverEndpoint)
		addr, err := net.ResolveTCPAddr("tcp", ks.cfg.NodeServer.Address)
		if err != nil {
			log.Err(err).Msg("cant parse ip address correctly")
			log.Fatal()
		}
		tcpListener, _ := net.ListenTCP("tcp", addr)

		aslListener := node_server.ASLListener{
			L:  tcpListener,
			Ep: serverEndpoint,
		}
		ks.server.Serve(aslListener)
	}()

	select {}
	log.Info().Msg("server down")

}
func CustomLoggerMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process the request
		c.Next()

		// Use the provided logger to log the request details
		duration := time.Since(start)
		logger.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Str("client_ip", c.ClientIP()).
			Dur("duration", duration).
			Msg("Request handled")
	}
}

/*
This function enablesthe app.
This is the point where the state and the tooling is obtained
In this current state, just the config is loaded.
TODO:
- Database init
- Ip allocator init
- Rest API init
*/
