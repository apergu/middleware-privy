package migration

import (
	"fmt"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"gitlab.com/mohamadikbal/project-privy/internal/config"
)

func Execute() {
	args := os.Args

	if len(args) < 2 {
		logrus.Error("[Migration] arguments is missing")
		os.Exit(1)
	}

	cfg := config.Init()

	switch args[1] {
	case "up":
		err := MigrateUp(cfg.Database.Dsn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to MigrateUp: %v\n", err)
			os.Exit(1)
		}
	case "down":
		err := MigrateDown(cfg.Database.Dsn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to MigrateDown: %v\n", err)
			os.Exit(1)
		}
	case "drop":
		err := MigrateDrop(cfg.Database.Dsn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to MigrateDrop: %v\n", err)
			os.Exit(1)
		}
	case "force":
		if len(args) < 3 {
			logrus.Error("[Migration] arguments force is missing")
			os.Exit(1)
		}

		version, _ := strconv.Atoi(args[2])
		err := MigrateForce(cfg.Database.Dsn, version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to MigrateForce: %v\n", err)
			os.Exit(1)
		}
	case "step":
		if len(args) < 3 {
			logrus.Error("[Migration] arguments force is missing")
			os.Exit(1)
		}

		version, _ := strconv.Atoi(args[2])
		err := MigrateStep(cfg.Database.Dsn, version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to MigrateStep: %v\n", err)
			os.Exit(1)
		}
	default:
		logrus.Error("[Migration] invalid command")
		os.Exit(1)
	}
}

type migrationLogger struct {
	logger *logrus.Logger
}

// Printf is like fmt.Printf
func (m migrationLogger) Printf(format string, v ...interface{}) {
	m.logger.Printf(format, v...)
}

// Verbose should return true when verbose logging output is wanted
func (m migrationLogger) Verbose() bool {
	return true
}

func InitMigration(dsn string) *migrate.Migrate {
	db := &pgx.Postgres{}

	driver, err := db.Open(dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migration",
		"jatis",
		driver,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create database instance: %v\n", err)
		os.Exit(1)
	}

	m.Log = migrationLogger{logrus.StandardLogger()}

	return m
}

func MigrateUp(url string) error {
	logrus.Info("[migration] Start Migrate Up")
	defer logrus.Info("[migration] Finish Migrate Up")

	err := InitMigration(url).Up()
	if err == migrate.ErrNoChange {
		return nil
	}

	if err != nil {
		logrus.
			Error("[migration] Migrate Up Error : ", err)

		return err
	}

	return nil
}

func MigrateDown(url string) error {
	logrus.Info("[migration] Start Migrate Down")
	defer logrus.Info("[migration] Finish Migrate Down")

	err := InitMigration(url).Down()
	if err == migrate.ErrNoChange {
		return nil
	}

	if err != nil {
		logrus.
			Error("[migration] Migrate Up Down : ", err)

		return err
	}

	return nil
}

func MigrateDrop(url string) error {
	logrus.Info("[migration] Start Migrate Drop")
	defer logrus.Info("[migration] Finish Migrate Drop")

	err := InitMigration(url).Drop()
	if err == migrate.ErrNoChange {
		return nil
	}

	if err != nil {
		logrus.
			Error("[migration] Migrate Drop Error : ", err)

		return err
	}

	return nil
}

func MigrateForce(url string, version int) error {
	logrus.Info("[migration] Start Migrate Force")
	defer logrus.Info("[migration] Finish Migrate Force")

	err := InitMigration(url).Force(version)
	if err == migrate.ErrNoChange {
		return nil
	}

	if err != nil {
		logrus.
			Error("[migration] Migrate Force Error : ", err)

		return err
	}

	return nil
}

func MigrateStep(url string, version int) error {
	logrus.Info("[migration] Start Migrate Step")
	defer logrus.Info("[migration] Finish Migrate Step")

	err := InitMigration(url).Steps(version)
	if err == migrate.ErrNoChange {
		return nil
	}

	if err != nil {
		logrus.
			Error("[migration] Migrate Force Error : ", err)

		return err
	}

	return nil
}
