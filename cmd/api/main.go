package main

import (
	"context"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/2zqa/ssot-specs-server/internal/data"
	"github.com/charmbracelet/log"
	oidc "github.com/coreos/go-oidc"
	"github.com/peterbourgon/ff/v3"
	"golang.org/x/oauth2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	oidc struct {
		clientID, secret, redirectURL, issuer string
	}
	apiKey string
	cors   struct {
		trustedOrigins []string
	}
	disableAuth bool
}

type application struct {
	config   config
	logger   *log.Logger
	models   data.Models
	oidcInfo oidcInfo
}

type oidcInfo struct {
	verifier *oidc.IDTokenVerifier
	provider *oidc.Provider
	config   oauth2.Config
}

func main() {
	var cfg config
	logger := log.New(os.Stdout)

	// Read the flags into the config struct
	fs := flag.NewFlagSet("ssot-specs-server", flag.ExitOnError)

	fs.IntVar(&cfg.port, "port", 4000, "API server port")
	fs.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	fs.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("SSOT_DB_DSN"), "PostgreSQL DSN")

	fs.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	fs.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	fs.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	fs.StringVar(&cfg.oidc.clientID, "oidc-client-id", os.Getenv("SSOT_OIDC_CLIENT_ID"), "OIDC client ID")
	fs.StringVar(&cfg.oidc.secret, "oidc-secret", os.Getenv("SSOT_OIDC_SECRET"), "OIDC client secret")
	fs.StringVar(&cfg.oidc.redirectURL, "oidc-redirect-url", os.Getenv("SSOT_OIDC_REDIRECT_URL"), "OIDC redirect URL")
	fs.StringVar(&cfg.oidc.issuer, "oidc-issuer", os.Getenv("SSOT_OIDC_ISSUER"), "OIDC issuer (authorization server)")

	fs.StringVar(&cfg.apiKey, "api-key", os.Getenv("SSOT_API_KEY"), "API key for authenticating collector clients")

	fs.BoolVar(&cfg.disableAuth, "disable-auth", false, "Disable authentication middleware")

	fs.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("SSOT"), ff.WithConfigFile(".env"), ff.WithConfigFileParser(ff.EnvParser)); err != nil {
		logger.Fatal(err)
	}

	if cfg.disableAuth {
		logger.Warn("Authentication middleware disabled")
	}

	db, err := openDB(cfg, logger)
	if err != nil {
		logger.Fatal("Could not set up database connection", "err", err)
	}
	logger.Info("Database connection pool established")

	oidcInfo, err := configureOIDC(cfg)
	if err != nil {
		logger.Fatal("Could not set up OIDC", "err", err)
	}
	logger.Info("OIDC configured")

	app := &application{
		config:   cfg,
		logger:   logger,
		models:   data.NewModels(db),
		oidcInfo: oidcInfo,
	}

	err = app.serve()
	if err != nil {
		logger.Fatal(err)
	}
}

func configureOIDC(cfg config) (oidcInfo, error) {
	provider, err := oidc.NewProvider(context.Background(), cfg.oidc.issuer)
	if err != nil {
		return oidcInfo{}, err
	}

	oauthCfg := oauth2.Config{
		ClientID:     cfg.oidc.clientID,
		ClientSecret: cfg.oidc.secret,
		RedirectURL:  cfg.oidc.redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}

	// Configure an OpenID Connect aware OAuth2 client.
	return oidcInfo{
		provider: provider,
		verifier: provider.Verifier(&oidc.Config{ClientID: cfg.oidc.clientID}),
		config:   oauthCfg,
	}, nil
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config, appLogger *log.Logger) (*gorm.DB, error) {

	newLogger := logger.New(
		appLogger,
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,        // Don't include params in the SQL log
		},
	)

	db, err := gorm.Open(postgres.Open(cfg.db.dsn), &gorm.Config{
		Logger:               newLogger,
		FullSaveAssociations: true,
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&data.Device{})
	db.AutoMigrate(&data.Specs{})
	db.AutoMigrate(&data.Motherboard{})
	db.AutoMigrate(&data.CPU{})
	db.AutoMigrate(&data.Disk{})
	db.AutoMigrate(&data.Partition{})
	db.AutoMigrate(&data.Network{})
	db.AutoMigrate(&data.NetworkInterface{})
	db.AutoMigrate(&data.NetworkInterfaceDriver{})
	db.AutoMigrate(&data.Bios{})
	db.AutoMigrate(&data.Memory{})
	db.AutoMigrate(&data.SwapDevice{})
	db.AutoMigrate(&data.Kernel{})
	db.AutoMigrate(&data.Release{})
	db.AutoMigrate(&data.DIMM{})
	db.AutoMigrate(&data.OEM{})
	db.AutoMigrate(&data.Virtualization{})

	db.AutoMigrate(&data.Search{})

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.db.maxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.db.maxIdleConns)
	sqlDB.SetConnMaxIdleTime(duration)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database.
	err = sqlDB.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
