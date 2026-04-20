package config

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Xendit   XenditConfig   `mapstructure:"xendit"`
	Email    EmailConfig    `mapstructure:"email"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Fee      FeeConfig      `mapstructure:"fee"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
	URL  string `mapstructure:"url"` // base URL for email links, e.g. http://localhost:8080
}

type DatabaseConfig struct {
	Host     string     `mapstructure:"host"`
	Port     int        `mapstructure:"port"`
	User     string     `mapstructure:"user"`
	Password string     `mapstructure:"password"`
	Name     string     `mapstructure:"name"`
	Pool     PoolConfig `mapstructure:"pool"`
}

type PoolConfig struct {
	Idle            int `mapstructure:"idle"`
	Max             int `mapstructure:"max"`
	LifetimeSeconds int `mapstructure:"lifetime_seconds"`
}

type RedisConfig struct {
	Addr string `mapstructure:"addr"`
	DB   int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	AccessTTL  string `mapstructure:"access_ttl"`
	RefreshTTL string `mapstructure:"refresh_ttl"`
}

type XenditConfig struct {
	SecretKey    string `mapstructure:"secret_key"`
	WebhookToken string `mapstructure:"webhook_token"`
	CallbackURL  string `mapstructure:"callback_url"`
}

type EmailConfig struct {
	Provider string `mapstructure:"provider"`
	APIKey   string `mapstructure:"api_key"`
	From     string `mapstructure:"from"`
	SMTPHost string `mapstructure:"smtp_host"`
	SMTPPort int    `mapstructure:"smtp_port"`
}

type StorageConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	PublicURL string `mapstructure:"public_url"`
}

// FeeConfig — semua nominal dalam rupiah (integer, bukan float)
type FeeConfig struct {
	PlatformPercent int   `mapstructure:"platform_percent"` // 5 = 5%
	PlatformFlat    int64 `mapstructure:"platform_flat"`    // Rp 500
	MinDonation     int64 `mapstructure:"min_donation"`     // Rp 1.000
	MinWithdrawal   int64 `mapstructure:"min_withdrawal"`   // Rp 50.000
}
