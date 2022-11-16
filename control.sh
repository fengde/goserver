#!/bin/bash

echo '
########################################################

control.sh 使用提示:
  ----for 本地----
  运行 ./control.sh               console运行
  运行 ./control.sh docker        容器运行
  运行 ./control.sh pack          编译打包，生成app.tar.gz 用于当前环境
  运行 ./control.sh pack linux    编译打包，生成app.tar.gz 用于linux下运行
  运行 ./control.sh pack windows  编译打包，生成app.tar.gz 用于windows下运行
  运行 ./control.sh pack darwin   编译打包，生成app.tar.gz 用于mac下运行

  ----for 线上 linux环境----
  运行 ./control.sh start   启动服务-非容器运行 (日志存于/var/log)
  运行 ./control.sh stop    停止服务-非容器运行
  运行 ./control.sh restart 重启服务-非容器运行
  运行 ./control.sh status  查看服务运行状态-非容器运行

########################################################
'

# 这里project根据情况自行修改
project='server'
image="$project:latest"
maingo="main.go"
app=$project
pidfile="/var/run/$app.pid"
logfile="/var/log/$app.log"
port=8080

function debug() {
  if [ -f $maingo ]; then
    go run $maingo
  fi

  return 0
}

function docker_() {
  docker ps >> /dev/null
  if [ $? -gt 0 ]; then
    return 1
  fi

  if [ $(docker ps | grep $app | wc -l | tr -s ' ') -gt 0 ]; then
      docker rm -f $app
  fi

  docker build -t $image .
  if [ $? -gt 0 ]; then
    echo "docker build failed..."
    return 1
  fi

  docker run -d -p $port:$port --name $app $image
  if [ $? -gt 0 ]; then
    echo "docker run failed..."
    return 1
  fi

  return 0
}

function pack() {
  if [ -f $maingo ]; then
    echo '正在编译...'

    mkdir -p app

    case $1 in
      linux)
        CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod vendor -o app/$app $maingo
        ;;
      windows)
        CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -mod vendor -o app/$app $maingo
        ;;
      darwin)
        CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -mod vendor -o app/$app $maingo
        ;;
      *)
        go build -mod vendor -o app/$app $maingo
        ;;
    esac
    
    if [ $? -gt 0 ]; then
      echo "go build failed..."
      return 1
    fi

    echo '正在打包...'

    cp -fr ./conf app/conf
    cp -fr ./.env_online app/.env
    cp -fr ./control.sh app/control.sh
    tar zcf app.tar.gz app
    rm -fr ./app

    if [ -f ./app.tar.gz ]; then
      echo '已打包完成...'
      echo '文件位置：./app.tar.gz'
      echo `md5 ./app.tar.gz`
    else
      echo '编译打包失败...'
      return 1
    fi
  fi

  return 0
}

function check_pid() {
  if [ -f $pidfile ]; then
    pid=`cat $pidfile`
    if [ $pid -gt 0 ]; then
      running=`ps -p $pid | grep -v "PID TTY" | wc -l`
      return $running
    fi
  else
    echo "$pidfile not found..."
  fi

  return 1
}

function start() {
  check_pid
  running=$?
  if [ $? -gt 0 ]; then
    echo -n "$app now is running already, pid="
    cat $pidfile
    return 1
  fi

  nohup ./$app >> $logfile 2>&1 &

  echo $! > $pidfile

  echo "$app started... pid=$!"

  return 0
}

function stop() {
  pid=`cat $pidfile`
  if [ $pid -gt 0 ]; then
    kill -15 $pid >/dev/null 2>&1
    if [ $? -eq 0 ]; then
      echo "$app stoped... pid=$pid"
      return 0
    fi
  fi

  return 1
}

function restart() {
    stop
    sleep 1
    start
}

function status() {
    check_pid
    if [ $? -gt 0 ]; then
      echo "$app is running..."
    else
      echo "$app is stoped..."
    fi
}

case $1 in
  docker)
    docker_
    ;;
  pack)
    pack $2
    ;;
  start)
    start
    ;;
  stop)
    stop
    ;;
  restart)
    restart
    ;;
  status)
    status
    ;;
  help)
    ;;
  --help)
    ;;
  -help)
    ;;   
  -h)
    ;;        
  *)
    debug
    ;;
esac