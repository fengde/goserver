package global

import (
	"flag"
	"time"

	"github.com/fengde/gocommon/confx"
	"github.com/fengde/gocommon/logx"
	"github.com/fengde/gocommon/storex/mysqlx"
	"github.com/fengde/gocommon/storex/redisx"
	"github.com/fengde/gocommon/timex"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var (
	Conf      Config
	DB        *mysqlx.Cluster
	Cache     *redisx.Client
	exit      = make(chan int)
	StartUnix int64
	Locker    *redisx.LockerV2
)

func Init() error {
	StartUnix = timex.NowUnix()

	flag.Parse()
	confx.MustLoad(*flag.String("f", "conf/config.yaml", "the config file"), &Conf, confx.UseEnv())
	logx.SetLevel(logx.DebugLevel)
	if Conf.Log.Path != "" {
		if Conf.Log.SaveDays < 1 {
			Conf.Log.SaveDays = 7
		}

		logx.SetLogFile(Conf.Log.Path, Conf.Log.SaveDays)
		logx.SetFormatter(logx.TextFormatter)
		if IsOnlineEnv() {
			gin.SetMode(gin.ReleaseMode)
		}
	}
	

	logx.Info("mysql connecting...")

	var err error

	DB, err = mysqlx.NewCluster([]string{Conf.Mysql.DataSourceName}, time.Second*time.Duration(Conf.Mysql.ConnMaxLifeSecond), Conf.Mysql.MaxOpenConns, Conf.Mysql.MaxIdleConns, !Conf.Mysql.SqlShow)
	if err != nil {
		return err
	}

	logx.Info("mysql connected")

	logx.Info("redis connecting...")

	Cache, err = redisx.NewClient(Conf.Redis.Addr, Conf.Redis.DB, Conf.Redis.Password)
	if err != nil {
		return err
	}

	logx.Info("redis connected")

	Locker = redisx.NewLockerV2(Cache)

	return nil
}
