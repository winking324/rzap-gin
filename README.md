# rzap-gin

Replace the log output of gin with zap

## How to install

`go get github.com/winking324/rzap-gin`

## How to use

``` go
rzap.NewGlobalLogger([]zapcore.Core{
    rzap.NewCore(&lumberjack.Logger{
        Filename: "/your/log/path/app.log",
    }, zap.InfoLevel),
})

r := gin.New()
r.Use(rzap_gin.Logger(nil), rzap_gin.Recovery(nil, false))
```