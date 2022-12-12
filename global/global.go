package global

import (
	"flag"
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/fengde/gocommon/confx"
	"github.com/fengde/gocommon/logx"
	"github.com/fengde/gocommon/redisx"
	"github.com/fengde/gocommon/timex"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Conf      Config
	Cache     *redisx.Client
	exit      = make(chan int)
	StartUnix int64
	Locker    *redisx.LockerV2
	Enforcer  *casbin.Enforcer
	DB        *gorm.DB
)

func Init() error {

	var err error

	StartUnix = timex.NowUnix()

	// 初始化Conf
	{
		flag.Parse()
		confx.MustLoad(*flag.String("f", "conf/config.yaml", "the config file"), &Conf, confx.UseEnv())
	}
	// 初始化Log
	{
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
	}
	// 初始化DB
	{
		logx.Info("mysql connecting...")

		DB, err = gorm.Open(mysql.New(mysql.Config{
			DSN: Conf.Mysql.DataSourceName,
		}), &gorm.Config{
			SkipDefaultTransaction: true,
		})
		if err != nil {
			return err
		}
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		sqlDB.SetMaxIdleConns(Conf.Mysql.MaxIdleConns)
		sqlDB.SetMaxOpenConns(Conf.Mysql.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(Conf.Mysql.ConnMaxLifeSecond))

		logx.Info("mysql connected")
	}
	// 初始化Redis
	{
		logx.Info("redis connecting...")

		Cache, err = redisx.NewClient(Conf.Redis.Addr, Conf.Redis.DB, Conf.Redis.Password)
		if err != nil {
			return err
		}

		logx.Info("redis connected")
	}
	// 分布式锁
	{
		Locker = redisx.NewLockerV2(Cache)
	}
	// 初始化RBAC对象
	{
		a, err := gormadapter.NewAdapter("mysql", Conf.Mysql.DataSourceName, true)
		if err != nil {
			return err
		}
		// 如果DataSourceName不存在db，则会创建casbin库；如果存在db，则会自动创建casbin_rule表
		Enforcer, err = casbin.NewEnforcer("conf/rbac_model.conf", a)
		if err != nil {
			return err
		}
	}

	return nil
}
