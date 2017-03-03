package client

import "fmt"
import "time"
import "strings"
import "strconv"

var defaultClientHeaders = map[string]string{
	"Content-Type": "application/json",
}

type Config struct {
	hostPort *hostPort
	user     string
	pass     string
	timeout  time.Duration
	headers  map[string]string
}

type hostPort struct {
	Host string
	Port int
}

// NewClientConfig handles the parameters that will be used in the HTTP Client
func NewClientConfig(host string, port int, user string, pass string, timeout int) (*Config, error) {
	err := validateSchema(host)
	if err != nil {
		return nil, err
	}
	hp, err := newHostString(host)
	if err != nil {
		return nil, err
	}

	if hp.Port != port {
		hp.Port = port
	}

	return &Config{
		hostPort: hp,
		user:     user,
		pass:     pass,
		timeout:  time.Duration(timeout) * time.Second,
		headers:  defaultClientHeaders,
	}, nil
}

func validateSchema(host string) error {
	if !(strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://")) {
		return fmt.Errorf("Host doesn't contain a valid HTTP protocol (http|https) => %s", host)
	}
	return nil
}

func newHostString(host string) (*hostPort, error) {
	urlString := strings.Split(host, "/")[2]
	if strings.Contains(urlString, ":") {
		urlStringPort := strings.Split(urlString, ":")[1]
		intedPort, err := strconv.Atoi(urlStringPort)
		if err != nil {
			return nil, fmt.Errorf("Invalid URL:Port combination, Port is not a string => %s", urlStringPort)
		}
		return &hostPort{strings.Join(strings.Split(host, ":")[0:2], ":"), intedPort}, nil
	}
	return &hostPort{strings.Join(strings.Split(host, ":")[0:2], ":"), 9200}, nil
}

// SetHeader that will be sent with the request
func (c *Config) SetHeader(key string, value string) {
	c.headers[key] = value
}

// HttpAddress returns the host and port combination so it can
// be used by the Client http://host:port
func (c *Config) HttpAddress() string {
	return fmt.Sprintf("%s:%d", c.hostPort.Host, c.hostPort.Port)
}

// GetTimeout returns the configured HTTP timeout
func (c *Config) GetTimeout() time.Duration {
	return c.timeout
}

// SetHost modifies the target host
func (c *Config) SetHost(value string) error {
	err := validateSchema(value)
	if err != nil {
		return err
	}
	c.hostPort, err = newHostString(value)
	return err
}

// SetPort modifies the target port
func (c *Config) SetPort(value int) {
	c.hostPort.Port = value
}

// SetUser modifies the user (HTTP Basic Auth)
func (c *Config) SetUser(value string) {
	c.user = value
}

// SetPass modifies the password (HTTP Basic Auth)
func (c *Config) SetPass(value string) {
	c.pass = value
}
