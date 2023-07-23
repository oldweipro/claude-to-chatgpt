package model

type ServerConfig struct {
	Claude Claude `mapstructure:"claude" json:"claude" yaml:"claude"`
	Proxy  Proxy  `mapstructure:"proxy" json:"proxy" yaml:"proxy"`
}

type Claude struct {
	BaseUrl          string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	SessionKey       string `mapstructure:"session-key" json:"session-key" yaml:"session-key"`
	OrganizationUuid string `mapstructure:"organization-uuid" json:"organization-uuid" yaml:"organization-uuid"`
}

type Proxy struct {
	Protocol string `mapstructure:"protocol" json:"protocol" yaml:"protocol"`
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     string `mapstructure:"port" json:"port" yaml:"port"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}
