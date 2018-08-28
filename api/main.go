//go:generate swagger generate spec
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	ct "burgundy/api/controller"
	conf "burgundy/conf"
	_Repo "burgundy/repository"

	"github.com/juju/errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	banner = `
:::::::::  :::    ::: :::::::::   ::::::::  :::    ::: ::::    ::: :::::::::  :::   ::: 
:+:    :+: :+:    :+: :+:    :+: :+:    :+: :+:    :+: :+:+:   :+: :+:    :+: :+:   :+: 
+:+    +:+ +:+    +:+ +:+    +:+ +:+        +:+    +:+ :+:+:+  +:+ +:+    +:+  +:+ +:+  
+#++:++#+  +#+    +:+ +#++:++#:  :#:        +#+    +:+ +#+ +:+ +#+ +#+    +:+   +#++:   
+#+    +#+ +#+    +#+ +#+    +#+ +#+   +#+# +#+    +#+ +#+  +#+#+# +#+    +#+    +#+    
#+#    #+# #+#    #+# #+#    #+# #+#    #+# #+#    #+# #+#   #+#+# #+#    #+#    #+#    
#########   ########  ###    ###  ########   ########  ###    #### #########     ###    
%s
 => Starting listen %s
`
)

var (
	// BuildDate for Program BuildDate
	BuildDate string
	// Version for Program Version
	Version string
	svrInfo = fmt.Sprintf("burgundy %s(%s)", Version, BuildDate)
)

func init() {
	// use all cpu
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	Burgundy := conf.Burgundy
	if Burgundy.GetBool("v") {
		fmt.Printf("%s\n", svrInfo)
		os.Exit(0)
	}
	Burgundy.SetProfile()

	f := apiLogFile("./burgundy-api.log")
	defer f.Close()
	e := echoInit(Burgundy, f)
	sc := sigInit(e)

	// Prepare Server
	db := _Repo.InitDB(Burgundy)
	defer db.Close()
	if err := ct.InitHandler(Burgundy, e, db); err != nil {
		fmt.Println("InitHandler error : ", errors.Details(err))
		os.Exit(1)
	}

	if !prepareServer(Burgundy, sc, db) {
		os.Exit(1)
	}
	startServer(Burgundy, e)
}

func apiLogFile(logfile string) *os.File {
	// API Logging
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("apiLogFile error : ", err)
		os.Exit(1)
	}
	return f
}

func echoInit(burgundy *conf.ViperConfig, apiLogF *os.File) (e *echo.Echo) {

	// Echo instance
	e = echo.New()
	e.Debug = true

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.POST, echo.GET, echo.PUT, echo.DELETE},
	}))

	// Ping Check
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "burgundy API Alive!\n") })
	e.POST("/", func(c echo.Context) error { return c.String(http.StatusOK, "burgundy API Alive!\n") })

	loggerConfig := middleware.DefaultLoggerConfig
	loggerConfig.Output = apiLogF

	e.Use(middleware.LoggerWithConfig(loggerConfig))
	e.Logger.SetOutput(bufio.NewWriterSize(apiLogF, 1024*16))
	e.Logger.SetLevel(burgundy.APILogLevel())
	e.HideBanner = true

	return e
}

func sigInit(e *echo.Echo) chan os.Signal {

	// Signal
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		sig := <-sc
		e.Logger.Error("Got signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Error(err)
		}
		signal.Stop(sc)
		close(sc)
	}()

	return sc
}

func startServer(burgundy *conf.ViperConfig, e *echo.Echo) {
	// Start Server
	apiServer := fmt.Sprintf("0.0.0.0:%d", burgundy.GetInt("port"))
	log.Printf("%s => Starting server listen %s\n", svrInfo, apiServer)
	fmt.Printf(banner, svrInfo, apiServer)

	if err := e.Start(apiServer); err != nil {
		fmt.Println(err)
		e.Logger.Error(err)
	}
}
