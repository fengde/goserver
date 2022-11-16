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
		ExpireHour int64
	}
	SentinelConfigPath string
	Plugins            []Plugin
}

type Plugin struct {
	Name    string
	Setting string
	Open    bool
}
