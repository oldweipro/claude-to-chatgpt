package model

import (
	"strings"
	"sync"
)

type ServerConfig struct {
	Claude  Claude `mapstructure:"claude" json:"claude" yaml:"claude"`
	Proxy   Proxy  `mapstructure:"proxy" json:"proxy" yaml:"proxy"`
	BaseUrl string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
}

type Claude struct {
	SessionKey []string   `mapstructure:"session-key" json:"session-key" yaml:"session-key"`
	Lock       sync.Mutex `json:"-"`
}

type Proxy struct {
	Protocol string `mapstructure:"protocol" json:"protocol" yaml:"protocol"`
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     string `mapstructure:"port" json:"port" yaml:"port"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

func (c *Claude) GetSessionKey() (sessionKey string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	if len(c.SessionKey) == 0 {
		return
	}
	sessionKey = c.SessionKey[0]
	if !strings.HasPrefix(sessionKey, "sessionKey=") {
		sessionKey = "sessionKey=" + sessionKey
	}
	c.SessionKey = append(c.SessionKey[1:], sessionKey)
	return
}
