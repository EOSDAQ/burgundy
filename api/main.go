//go:generate swagger generate spec
package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	ct "burgundy/api/controller"
	mw "burgundy/api/middleware"
	conf "burgundy/conf"
	_Repo "burgundy/repository"
	"burgundy/util"

	"github.com/juju/errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
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
	mlog    *zap.SugaredLogger
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
	mlog, _ = util.InitLog("main", Burgundy.GetString("loglevel"))

	e := echoInit(Burgundy)
	sc := sigInit(e)

	// Prepare Server
	db := _Repo.InitDB(Burgundy)
	defer db.Close()
	if err := ct.InitHandler(Burgundy, e, db); err != nil {
		mlog.Errorw("InitHandler", "err", errors.Details(err))
		os.Exit(1)
	}

	if !prepareServer(Burgundy, sc, db) {
		os.Exit(1)
	}
	startServer(Burgundy, e)
}

func echoInit(burgundy *conf.ViperConfig) (e *echo.Echo) {

	// Echo instance
	e = echo.New()

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

	e.Use(mw.ZapLogger(mlog))
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
	mlog.Infow("Starting server", "info", svrInfo, "listen", apiServer)
	fmt.Printf(banner, svrInfo, apiServer)

	if err := e.Start(apiServer); err != nil {
		mlog.Errorw("End server", "err", err)
	}
}
