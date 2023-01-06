package main

import (
	"WBABEProject-04/conf"
	"WBABEProject-04/controller"
	"WBABEProject-04/logger"
	"WBABEProject-04/model"
	"WBABEProject-04/router"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	var configFlag = flag.String("config", "./conf/config.toml", "toml file to use for configuration")
	flag.Parse()
	cf := conf.NewConfig(*configFlag)
	if err := logger.InitLogger(cf); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	if model, err := model.NewModel(); err != nil {
		panic(fmt.Errorf("model.NewMode > %v", err))
	} else if controller, err := controller.NewController(model); err != nil {
		panic(fmt.Errorf("controller.NewController > %v", err))
	} else if router, err := router.NewRouter(controller); err != nil {
		panic(fmt.Errorf("router.NewRouter > %v", err))
	} else {
		mapi := &http.Server{
			Addr:           ":8080",
			Handler:        router.Index(),
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		g.Go(func() error {
			return mapi.ListenAndServe()
		})
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Warn("Shutdown Server ...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mapi.Shutdown(ctx); err != nil {
			logger.Error("Server Shutdown:", err)
		}
		select {
		case <-ctx.Done():
			logger.Info("timeout of 5 seconds.")
		}

		logger.Info("Server exiting")
	}
	if err := g.Wait(); err != nil {
		logger.Error(err)
	}
}
