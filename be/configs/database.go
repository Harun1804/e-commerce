package configs

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var DB *gorm.DB

func ConnectionDatabase() error {
	config := NewConfig().Database
	// ensure database exists before opening connection
	if err := checkAndCreateDatabase(config); err != nil {
		return err
	}

	connString := buildPostgresConnString(config.Host, config.Port, config.User, config.Pass, config.Name, config.SSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		zap.L().Error("[Postgress] ConnectionPostgres - 1", zap.Error(err))
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		zap.L().Error("[Postgress] ConnectionPostgres - 2", zap.Error(err))
		return err
	}

	sqlDB.SetMaxOpenConns(config.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(config.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err := sqlDB.Ping(); err != nil {
		zap.L().Error("[Postgress] ConnectionPostgres - ping", zap.Error(err))
		return err
	}

	if err := setupReplicaConnections(DB, config); err != nil {
		return err
	}

	return nil
}

func CloseDatabase() {
	if DB == nil {
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		zap.L().Error("[Postgress] CloseDatabase - 1", zap.Error(err))
		return
	}

	if err := sqlDB.Close(); err != nil {
		zap.L().Error("[Postgress] CloseDatabase - 2", zap.Error(err))
	}
}

// checkAndCreateDatabase connects to the server's default `postgres` database,
// checks whether the configured database exists, and creates it if missing.
func checkAndCreateDatabase(config Database) error {

	// connect to the server using the default 'postgres' database
	connString := buildPostgresConnString(config.Host, config.Port, config.User, config.Pass, "postgres", config.SSLMode)

	tempDB, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		zap.L().Error("[Postgress] CheckAndCreateDatabase - open", zap.Error(err))
		return err
	}

	sqlDB, err := tempDB.DB()
	if err != nil {
		zap.L().Error("[Postgress] CheckAndCreateDatabase - db", zap.Error(err))
		return err
	}
	defer func() { _ = sqlDB.Close() }()

	var exists bool
	row := tempDB.Raw("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)", config.Name).Row()
	if err := row.Scan(&exists); err != nil {
		zap.L().Error("[Postgress] CheckAndCreateDatabase - scan", zap.Error(err))
		return err
	}

	if !exists {
		createSQL := fmt.Sprintf("CREATE DATABASE %s", quoteIdentifier(config.Name))
		if err := tempDB.Exec(createSQL).Error; err != nil {
			zap.L().Error("[Postgress] CheckAndCreateDatabase - create", zap.Error(err))
			return err
		}
		zap.L().Info("[Postgress] CheckAndCreateDatabase - created database", zap.String("database", config.Name))
	}

	return nil
}

func setupReplicaConnections(db *gorm.DB, config Database) error {
	if len(config.SlaveHosts) == 0 {
		return nil
	}

	replicas := make([]gorm.Dialector, 0, len(config.SlaveHosts))
	for _, target := range config.SlaveHosts {
		host, port, ok := parseReplicaTarget(target, config.Port)
		if !ok {
			continue
		}

		replicaConnString := buildPostgresConnString(host, port, config.User, config.Pass, config.Name, config.SSLMode)
		replicas = append(replicas, postgres.Open(replicaConnString))
	}

	if len(replicas) == 0 {
		return nil
	}

	err := db.Use(
		dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}).SetMaxOpenConns(config.DBMaxOpenConns).
			SetMaxIdleConns(config.DBMaxIdleConns).
			SetConnMaxLifetime(time.Hour),
	)
	if err != nil {
		zap.L().Error("[Postgress] setupReplicaConnections", zap.Error(err))
		return err
	}

	zap.L().Info("[Postgress] replica connection enabled", zap.Int("replica_count", len(replicas)))
	return nil
}

func buildPostgresConnString(host string, port int, user, pass, dbName, sslMode string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user,
		pass,
		host,
		port,
		dbName,
		sslMode,
	)
}

func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func parseReplicaTarget(raw string, defaultPort int) (string, int, bool) {
	target := strings.TrimSpace(raw)
	if target == "" {
		return "", 0, false
	}

	parts := strings.Split(target, ":")
	if len(parts) == 1 {
		return parts[0], defaultPort, true
	}

	if len(parts) != 2 {
		zap.L().Warn("[Postgress] invalid replica target format", zap.String("target", raw))
		return "", 0, false
	}

	host := strings.TrimSpace(parts[0])
	if host == "" {
		zap.L().Warn("[Postgress] invalid replica host", zap.String("target", raw))
		return "", 0, false
	}

	port, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || port <= 0 {
		zap.L().Warn("[Postgress] invalid replica port", zap.String("target", raw))
		return "", 0, false
	}

	return host, port, true
}
