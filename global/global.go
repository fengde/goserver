package global

import (
	"context"
	"flag"
	"goserver/conf"
	"goserver/consts"
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
	Conf      conf.Config
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
	var err error
	DB, err = mysqlx.NewCluster([]string{Conf.Mysql.DataSourceName}, time.Second*time.Duration(Conf.Mysql.ConnMaxLifeSecond), Conf.Mysql.MaxOpenConns, Conf.Mysql.MaxIdleConns, !Conf.Mysql.SqlShow)
	if err != nil {
		return err
	}

	Cache, err = redisx.NewClient(Conf.Redis.Addr, Conf.Redis.DB, Conf.Redis.Password)
	if err != nil {
		return err
	}

	Locker = redisx.NewLockerV2(Cache)

	return nil
}

// 当前运行环境是否为dev
func IsDevEnv() bool {
	return Conf.Env == consts.ENV_DEV
}

// 当前运行环境是否为qa
func IsQaEnv() bool {
	return Conf.Env == consts.ENV_QA
}

// 当前运行环境是否为online
func IsOnlineEnv() bool {
	return Conf.Env == consts.ENV_ONLINE
}

// 系统是否还在运行，带sleep
func Continue(sleep ...time.Duration) bool {
	if len(sleep) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), sleep[0])
		defer cancel()

		select {
		case <-exit:
			return false
		case <-ctx.Done():
			return true
		}
	}

	select {
	case <-exit:
		return false
	default:
		return true
	}
}

// 设置系统关闭
func Exist() {
	close(exit)
}
