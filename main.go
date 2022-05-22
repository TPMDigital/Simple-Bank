package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	// "runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	// "github.com/tpmdigital/simplebank/api"
	db "github.com/tpmdigital/simplebank/db/sqlc"
	"github.com/tpmdigital/simplebank/gapi"
	"github.com/tpmdigital/simplebank/pb"
	"github.com/tpmdigital/simplebank/util"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configuration file:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	// go runGinServer(config, store)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {

	// create a new instance of the simplebank (gapi) server
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// create a grpc server
	grpcServer := grpc.NewServer()

	// register our server with this grpc server
	pb.RegisterSimpleBankServer(grpcServer, server)

	// turn on relection so we can see the endpoints using Evans
	reflection.Register(grpcServer)

	// create a grpc listener
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create grpc listener:", err)
	}

	// tell the server about the listener and start the grpc server
	log.Printf("starting gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server:", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {

	// create a new instance of the simplebank (gapi) server
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: false, // false = camel case, true is snake case 
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// create a http listener
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create http listener:", err)
	}

	// tell the server about the listener and start the grpc server
	log.Printf("starting HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start http gateway server:", err)
	}
}

/*
func runGinServer(config util.Config, store db.Store) {

	// create a new instance of the simplebank (gin) server
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create http server:", err)
	}

	// start this server
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start http server:", err)
	}
}
*/
