package main

import (
	"WBABEProject-04/controller"
	"WBABEProject-04/model"
	"WBABEProject-04/router"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
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
		if err := g.Wait(); err != nil {
			panic(fmt.Errorf("error > %v", err))
		}

	}
}
