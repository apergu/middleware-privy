package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"middleware/cmd/migration"
	"middleware/infrastructure"
	http_client_infratructure "middleware/infrastructure/http/http_client"
	log "middleware/infrastructure/logger/logrus"
	"middleware/internal/config"
	"middleware/internal/httphandler"
	"middleware/pkg/appemail"
	"middleware/pkg/credential"
	"middleware/pkg/erpprivy"
	"middleware/pkg/pgxdb"
	"middleware/pkg/privy"

	"github.com/go-chi/jwtauth"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"gitlab.com/rteja-library3/rcache"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/remailer"
	"gitlab.com/rteja-library3/rpassword"
	"gitlab.com/rteja-library3/rserver"
	"gitlab.com/rteja-library3/rtoken"
)

func Execute() {
	ctx, canc := context.WithCancel(context.Background())
	cfg := config.Init()

	logrus.Infof("Service %s run on port %s", cfg.Application.Name, cfg.Application.Port)

	// create db pool
	pool := pgxdb.InitDatabase(ctx, cfg.Database)

	// jwt
	jwtAuth := jwtauth.New("HS256", []byte(cfg.Jwt.SignatureKey), nil)
	defaultToken := rtoken.NewJWTLestrrat(jwtAuth, logrus.StandardLogger(), rtoken.NewJWTProperty().
		SetAudience(cfg.Jwt.Audience).
		SetIssuers(cfg.Jwt.Issuers).
		SetExpiredDuration(cfg.Jwt.Expiration),
	)

	defaultRefreshToken := rtoken.NewJWT(logrus.StandardLogger(), nil, rtoken.NewJWTProperty().
		SetAudience(cfg.RefreshJWT.Audience).
		SetIssuers(cfg.RefreshJWT.Issuers).
		SetExpiredDuration(cfg.RefreshJWT.Expiration),
	)

	// cache
	defaultCache := rcache.NewMemoryCache(cfg.Jwt.Expiration)
	if cfg.Application.IsRedis {
		defaultCache = rcache.NewRedisCache(rcache.NewRedisClient(
			&redis.Options{
				Addr:     "localhost:6379",
				Password: "rahmanteja",
				DB:       0,
			},
		), cfg.Jwt.Expiration)
	}

	// Create Password Encryptor
	var defaultPwdEncryptor rpassword.Encryptor = rpassword.NewBcryptPassword(logrus.StandardLogger(), 0)

	var defaultEmailSender remailer.Remail = appemail.AppDummyEmail{}

	// create credential privy
	credPrivy := credential.NewCredentialPrivy(credential.CredentialPrivyProperty{
		Host:     cfg.CredentialPrivy.Host,
		Client:   http.DefaultClient,
		Username: cfg.CredentialPrivy.Username,
		Password: cfg.CredentialPrivy.Password,
	})

	credErpPrivy := erpprivy.NewCredentialERPPrivy(erpprivy.ErpPrivyProperty{
		Host:           cfg.ErpPrivy.Host,
		Username:       cfg.ErpPrivy.Username,
		Password:       cfg.ErpPrivy.Password,
		ApplicationKey: cfg.ErpPrivy.ApplicationKey,
		RequestId:      cfg.ErpPrivy.RequestId,
	})

	defaultPrivy := privy.NewPrivyGeneral(privy.PrivyProperty{
		Host:     cfg.Privy.Host,
		Client:   http.DefaultClient,
		Username: cfg.Privy.Username,
		Password: cfg.Privy.Password,
	})

	httpProperty := httphandler.HTTPHandlerProperty{
		DBPool:              pool,
		DefaultDecoder:      rdecoder.NewJSONEncoder(),
		DefaultToken:        defaultToken,
		DefaultCache:        defaultCache,
		DefaultPwdEncryptor: defaultPwdEncryptor,
		DefaultEmailer:      defaultEmailSender,
		DefaultRefreshToken: defaultRefreshToken,
		DefaultCredential:   credPrivy,
		DefaultPrivy:        defaultPrivy,
		DefaultERPPrivy:     credErpPrivy,
	}

	logger := log.NewLoggerLogrus(cfg)

	httpClient := http_client_infratructure.NewHttpClient(cfg, logger)

	infra := infrastructure.Infrastructure{
		Config:     cfg,
		Logger:     logger,
		HttpClient: httpClient,
	}

	handler := InitHttpHandler(pool, cfg.Cors, httpProperty, jwtAuth, cfg.BasicAuth, &infra)

	// migrate
	err := migration.MigrateUpWithDBName(cfg.Database.Dsn, cfg.Database.DBName)
	if err != nil {
		defaultCache.Close()
		pool.Close()
		canc()

		fmt.Fprintf(os.Stderr, "Unable to MigrateUp database: %v\n", err)
		os.Exit(1)
	}

	// RUN SERVER
	server := rserver.NewServer(
		logrus.StandardLogger(),
		rserver.
			NewOptions().
			SetHandler(handler).
			SetHost("0.0.0.0").
			SetPort(cfg.Application.Port),
	)
	server.Start()

	csignal := make(chan os.Signal, 1)
	signal.Notify(csignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	// WAIT FOR IT
	<-csignal

	logrus.Info("Closing cache")
	defaultCache.Close()

	logrus.Info("Closing pooling")
	pool.Close()

	logrus.Info("Closing server")
	server.Close()
	canc()

	logrus.Infof("Service %s run on port %s stopped", cfg.Application.Name, cfg.Application.Port)
}
