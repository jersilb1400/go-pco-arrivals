package config

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	PCO      PCOConfig      `json:"pco"`
	Auth     AuthConfig     `json:"auth"`
	Redis    RedisConfig    `json:"redis"`
	Realtime RealtimeConfig `json:"realtime"`
}

type ServerConfig struct {
	Port        int      `json:"port"`
	Host        string   `json:"host"`
	CORSOrigins []string `json:"cors_origins"`
	TrustProxy  bool     `json:"trust_proxy"`
}

type DatabaseConfig struct {
	URL             string `json:"url"`
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"`
}

type PCOConfig struct {
	BaseURL      string `json:"base_url"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scopes       string `json:"scopes"`
	AccessToken  string `json:"access_token"`
	AccessSecret string `json:"access_secret"`
}

type AuthConfig struct {
	SessionTTL            int      `json:"session_ttl"`
	RememberMeDays        int      `json:"remember_me_days"`
	AuthorizedUsers       []string `json:"authorized_users"`
	SessionSecret         string   `json:"session_secret"`
	JWTSecret             string   `json:"jwt_secret"`
	TokenRefreshThreshold int      `json:"token_refresh_threshold"`
}

type RedisConfig struct {
	URL      string `json:"url"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type RealtimeConfig struct {
	WebSocketEnabled     bool `json:"websocket_enabled"`
	PollingFallback      bool `json:"polling_fallback"`
	PollingInterval      int  `json:"polling_interval"`
	LocationPollInterval int  `json:"location_poll_interval"`
	MaxConnections       int  `json:"max_connections"`
	HeartbeatInterval    int  `json:"heartbeat_interval"`
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnvInt("PORT", 3000),
			Host:        getEnv("HOST", "0.0.0.0"),
			CORSOrigins: strings.Split(getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173"), ","),
			TrustProxy:  getEnvBool("TRUST_PROXY", false),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "file:./data/pco_billboard.db?cache=shared&mode=rwc"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", 300),
		},
		PCO: PCOConfig{
			ClientID:     getEnv("PCO_CLIENT_ID", ""),
			ClientSecret: getEnv("PCO_CLIENT_SECRET", ""),
			RedirectURI:  getEnv("PCO_REDIRECT_URI", ""),
			BaseURL:      getEnv("PCO_BASE_URL", "https://api.planningcenteronline.com"),
			Scopes:       getEnv("PCO_SCOPES", "people check_ins"),
		},
		Auth: AuthConfig{
			SessionTTL:            getEnvInt("SESSION_TTL", 3600),
			RememberMeDays:        getEnvInt("REMEMBER_ME_DAYS", 30),
			AuthorizedUsers:       strings.Split(getEnv("AUTHORIZED_USERS", ""), ","),
			SessionSecret:         getEnv("SESSION_SECRET", generateSessionSecret()),
			JWTSecret:             getEnv("JWT_SECRET", generateJWTSecret()),
			TokenRefreshThreshold: getEnvInt("TOKEN_REFRESH_THRESHOLD", 300),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", ""),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Realtime: RealtimeConfig{
			WebSocketEnabled:     getEnvBool("WEBSOCKET_ENABLED", true),
			PollingFallback:      getEnvBool("POLLING_FALLBACK", true),
			PollingInterval:      getEnvInt("POLLING_INTERVAL", 10),
			LocationPollInterval: getEnvInt("LOCATION_POLL_INTERVAL", 60),
			MaxConnections:       getEnvInt("MAX_CONNECTIONS", 1000),
			HeartbeatInterval:    getEnvInt("HEARTBEAT_INTERVAL", 30),
		},
	}

	// Validate required fields
	if cfg.PCO.ClientID == "" {
		return nil, fmt.Errorf("PCO_CLIENT_ID is required")
	}
	if cfg.PCO.ClientSecret == "" {
		return nil, fmt.Errorf("PCO_CLIENT_SECRET is required")
	}
	if cfg.PCO.RedirectURI == "" {
		return nil, fmt.Errorf("PCO_REDIRECT_URI is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func generateSessionSecret() string {
	return randomString(64)
}

func generateJWTSecret() string {
	return randomString(64)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist, we'll use system environment variables
	}
}
