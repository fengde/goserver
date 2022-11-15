## 依赖
```
要求 go1.18+
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

## 升级包
```
直接修改go.mod版本
然后执行
    go mod tidy
    go mod vendor 
```

## 常用操作-调试打包
```
本地调试：
./control.sh
二进制打包：
./control.sh pack 
```