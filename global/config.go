package global

type Config struct {
	Name        string
	Version     string
	Env         string
	HttpAddress string
	Log         struct {
		Level    string
		Path     string
		SaveDays int
	}
	Mysql struct {
		DataSourceName    string
		ConnMaxLifeSecond int
		MaxOpenConns      int
		MaxIdleConns      int
		SqlShow           bool
	}
	Redis struct {
		Addr     string
		DB       int
		Password string
	}
	RuntimeDataPath string
	Jwt             struct {
		Secret     string
		ExpireHour int
	}
	SentinelConfigPath string
	Jaeger             struct {
		Way          string
		Endpoint     string
		Token        string
		SamplerType  string
		SamplerParam float64
	}
	Plugins []Plugin
}

type Plugin struct {
	Name    string
	Setting string
	Open    bool
}
