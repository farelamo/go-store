package config

import (
	"store/utils"
	"strconv"
)

var (
	// App
	AppPort, _ = strconv.Atoi(utils.MustGetVaultEnv("APP_PORT"))
	LogLevel   = utils.MustGetVaultEnv("LOG_LEVEL")

	// Vault
	// VaultAddr       = utils.MustGetVaultEnv("VAULT_ADDR")
	// VaultPort, _    = strconv.Atoi(utils.MustGetVaultEnv("VAULT_PORT"))
	// VaultToken      = utils.MustGetVaultEnv("VAULT_TOKEN")
	// VaultSecretPath = utils.MustGetVaultEnv("VAULT_SECRET_PATH")
	// VaultEnabled    = utils.MustGetVaultEnv("VAULT_ENABLED")

	// PhpMyAdmin
	PhpMyAdminPort, _        = strconv.Atoi(utils.MustGetVaultEnv("PHPMYADMIN_PORT"))
	PhpMyAdminDefaultPort, _ = strconv.Atoi(utils.MustGetVaultEnv("PHPMYADMIN_DEFAULT_PORT"))

	// DB Name
	DbName = utils.MustGetVaultEnv("MYSQL_DB_NAME")
	DbUser = utils.MustGetVaultEnv("MYSQL_USER")
	DbPass = utils.MustGetVaultEnv("MYSQL_PASSWORD")
	DbHost = utils.MustGetVaultEnv("MYSQL_HOST")
	DbPort = utils.MustGetVaultEnv("MYSQL_PORT")

	// Redis
	RedisHost  = utils.MustGetVaultEnv("REDIS_HOST")
	RedisPort  = utils.MustGetVaultEnv("REDIS_PORT")
	RedisDb, _ = strconv.Atoi(utils.MustGetVaultEnv("REDIS_DB"))

	// Jwt
	AccessTokenExpiration, _  = strconv.Atoi(utils.MustGetVaultEnv("JWT_ACCESS_TOKEN_EXPIRATION"))
	RefreshTokenExpiration, _ = strconv.Atoi(utils.MustGetVaultEnv("JWT_REFRESH_TOKEN_EXPIRATION"))
	JwtSecret                 = utils.MustGetVaultEnv("JWT_SECRET")
	RefreshSecret             = utils.MustGetVaultEnv("JWT_REFRESH")
)
