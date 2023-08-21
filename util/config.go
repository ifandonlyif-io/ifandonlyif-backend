package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DBDriver                 string        `mapstructure:"DB_DRIVER"`
	DBSource                 string        `mapstructure:"DB_SOURCE"`
	MigrationURL             string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress        string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress        string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	AccessTokenSymmetricKey  string        `mapstructure:"ACCESS_TOKEN_SYMMETRIC_KEY"`
	RefreshTokenSymmetricKey string        `mapstructure:"REFRESH_TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration      time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration     time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EnableProfiler           bool          `mapstructure:"ENABLE_PROFILER"`
	AlchemyApiUrl            string        `mapstructure:"ALCHEMY_API_URL"`
	AlchemyNftApiUrl         string        `mapstructure:"ALCHEMY_NFT_API_URL"`
	IFFNftContractAddress    string        `mapstructure:"IFFNFT_CONTRACT_ADDRESS"`
	MoralisApiUrl            string        `mapstructure:"MORALIS_API_URL"`
	MoralisApiKey            string        `mapstructure:"MORALIS_API_KEY"`
	MoralisEthNetwork        string        `mapstructure:"MORALIS_ETH_NETWORK"`
	BlackholeAddress         string        `mapstructure:"BLACKHOLE_ADDRESS"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
