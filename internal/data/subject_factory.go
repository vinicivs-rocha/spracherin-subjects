package data

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Params   map[string]string
}

type MySQLPoolConfig struct {
	MaxOpen int
	MaxIdle int
}

func NewMySQLSubjectRepositoryFromEnv(pool MySQLPoolConfig) (*MySQLSubjectRepository, error) {
	cfg, err := MySQLConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return NewMySQLSubjectRepositoryFromConfig(cfg, pool)
}

func MySQLConfigFromEnv() (MySQLConfig, error) {
	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	if host == "" {
		return MySQLConfig{}, errors.New("DB_HOST is required")
	}

	port := 3306
	parsed, err := strconv.Atoi(strings.TrimSpace(os.Getenv("DB_PORT")))
	if err != nil {
		return MySQLConfig{}, errors.New("DB_PORT must be a number")
	}
	port = parsed

	user := strings.TrimSpace(os.Getenv("DB_USER"))
	if user == "" {
		return MySQLConfig{}, errors.New("DB_USER is required")
	}

	password := os.Getenv("DB_PASSWORD")
	if strings.TrimSpace(password) == "" {
		return MySQLConfig{}, errors.New("DB_PASSWORD is required")
	}

	name := strings.TrimSpace(os.Getenv("DB_NAME"))
	if name == "" {
		return MySQLConfig{}, errors.New("DB_NAME is required")
	}

	params := parseDBParams(os.Getenv("DB_PARAMS"))

	return MySQLConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     name,
		Params:   params,
	}, nil
}

func NewMySQLSubjectRepositoryFromConfig(cfg MySQLConfig, pool MySQLPoolConfig) (*MySQLSubjectRepository, error) {
	dsn, err := buildMySQLDSN(cfg)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if pool.MaxOpen > 0 {
		db.SetMaxOpenConns(pool.MaxOpen)
	}
	if pool.MaxIdle > 0 {
		db.SetMaxIdleConns(pool.MaxIdle)
	}

	return NewMySQLSubjectRepository(db), nil
}

func buildMySQLDSN(cfg MySQLConfig) (string, error) {
	if strings.TrimSpace(cfg.Host) == "" {
		return "", errors.New("host is required")
	}
	if cfg.Port <= 0 {
		return "", errors.New("port must be positive")
	}
	if strings.TrimSpace(cfg.User) == "" {
		return "", errors.New("user is required")
	}
	if strings.TrimSpace(cfg.Password) == "" {
		return "", errors.New("password is required")
	}
	if strings.TrimSpace(cfg.Name) == "" {
		return "", errors.New("database name is required")
	}

	params := url.Values{}
	for key, value := range cfg.Params {
		if strings.TrimSpace(key) == "" {
			continue
		}
		params.Set(key, value)
	}

	query := params.Encode()
	if query == "" {
		query = "parseTime=true"
	} else if !strings.Contains(query, "parseTime=") {
		query = query + "&parseTime=true"
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		query,
	), nil
}

func parseDBParams(raw string) map[string]string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return map[string]string{}
	}

	params := make(map[string]string)
	pairs := strings.Split(raw, "&")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		key := strings.TrimSpace(parts[0])
		if key == "" {
			continue
		}
		value := ""
		if len(parts) == 2 {
			value = strings.TrimSpace(parts[1])
		}
		params[key] = value
	}

	return params
}
