## 目录结构定义

- conf 配置文件
- consts 常量定义
- db   db表结构，基础查询维护
- doc  文档集合
- global 全局调用对象
- http web api入口逻辑
- i10n 国际化
- service 核心业务逻辑
- sql sql变更记录
- test 业务测试
- util 工程内的脚手架工具方法
- vendor 依赖包
- .env 多环境配置维护


## 依赖
```
要求 go1.18+
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

## 常用操作-升级包
```
直接修改go.mod版本
然后执行
    go mod tidy
    go mod vendor 
    or 
    go mod download
```

## 常用操作-调试打包
```
本地调试：
./control.sh
二进制打包：
./control.sh pack 
```