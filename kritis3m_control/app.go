package kritis3m_control

import (
	// "io"
	// "log"

	asl "github.com/Laboratory-for-Safe-and-Secure-Systems/go-asl"
	// "github.com/gin-gonic/gin"
	// "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	// "os"

	"net/http"

	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/db"
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/grpc"
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/types"
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

	grpc.NewNorthboundService(app.db.DB, &gormlog)

	// app.log_db, err = db.NewLogDatabase(cfg.Log_Database)

	// ginLogger := log.With().Str("service", "gin").Logger().Level(app.cfg.NodeServer.Log.Level)

	err = asl.ASLinit(&cfg.ASLConfig)
	if err != nil {
		log.Err(err).Msg("err asl init")
		panic(err)
	}
	// app.server = http.Server{
	// 	Handler: router,
	// }
	// app.insecure_server = http.Server{
	// 	Handler: router}

	return &app, nil
}
