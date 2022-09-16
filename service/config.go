package service

const (
	Transport_REST string = "rest"
	Transport_GRPC string = "grpc"
)

type Config struct {
	PopulateDB bool
	Env        string
	Port       int
	Db         struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	// rate limiting for each ip address
	Limiter struct {
		Rps     float64
		Burst   int
		Enabled bool
	}

	Auth struct {
		Enabled bool
		ApiKey  string
	}
	Transport string
}
