package connections

type Client interface {
	Call(method, path string, body interface{}, response interface{}) error
	Close() error
}

type Config struct {
	Type    string
	BaseURL string
	Timeout int
}

func NewClient(cfg Config) (Client, error) {
	switch cfg.Type {
	case "http":
		return NewHTTPClient(cfg)
	default:
		return NewHTTPClient(cfg)
	}
}
