package main

import (
	"jbndlr/example/api/esid"
	"jbndlr/example/api/grpc"
	"jbndlr/example/api/rest"
	"jbndlr/example/conf"
	"log"

	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group

	g.Go(func() error { return rest.Serve(conf.P.API.RESTPort) })
	g.Go(func() error { return esid.Serve(conf.P.API.ESIDPort) })
	g.Go(func() error { return grpc.Serve(conf.P.API.GRPCPort) })

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
