package pgxdb

import (
	"context"
	"fmt"
	"os"

	"middleware/internal/config"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type tracer struct{}

func (t tracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	fmt.Println(middleware.GetReqID(ctx), data)
	return ctx
}

func (t tracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	// fmt.Println(data.Err)
	// fmt.Println(data.CommandTag)
	// fmt.Println(ctx.Value("x"))
}

func InitDatabase(ctx context.Context, cfg config.Database) *pgxpool.Pool {
	logrus.Info("[Database] START ", cfg.Dsn)
	defer logrus.Info("[Database] END ", cfg.Dsn)

	config, err := pgxpool.ParseConfig(cfg.Dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Database] Unable to parse dsn: %v\n", err)
		os.Exit(1)
	}

	config.ConnConfig.Tracer = tracer{}

	// pgx.Logger = pgxlogrus.NewLogger(logrus.StandardLogger())

	config.MaxConnIdleTime = cfg.MaxConnIdleTime
	config.MaxConnLifetime = cfg.MaxConnLifetime
	config.MaxConns = cfg.MaxOpenConn

	logrus.Info("[Database] Connecting to database ", cfg.Dsn)
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Database] Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	logrus.Info("[Database] Connected ", cfg.Dsn)

	logrus.Info("[Database] Ping database ", cfg.Dsn)
	err = pool.Ping(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Database] Unable to ping to database: %v\n", err)
		os.Exit(1)
	}
	logrus.Info("[Database] Success ping database ", cfg.Dsn)

	return pool
}
