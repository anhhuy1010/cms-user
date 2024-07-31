package main

import (
	"fmt"
	"log"
	"net"

	"github.com/anhhuy1010/cms-user/config"
	"github.com/anhhuy1010/cms-user/database"
	grpcClient "github.com/anhhuy1010/cms-user/grpc"
	"github.com/anhhuy1010/cms-user/grpc/service"

	pbUser "github.com/anhhuy1010/cms-user/grpc/proto/users"
	"github.com/anhhuy1010/cms-user/routes"
	"github.com/anhhuy1010/cms-user/services/logService"
	"github.com/gin-gonic/gin"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	engine       *gin.Engine
	cfg          *viper.Viper
	logrusLogger *logrus.Logger
	customFunc   grpc_logrus.CodeToLevel
)

func init() {
	engine = gin.New()
	engine.Use(gin.Logger())

	logService.NewLogrus()

	cfg = config.GetConfig()
}

func main() {
	_, err := database.Init()
	if err == nil {
		fmt.Println("\n Database connected!")
	} else {
		fmt.Println("Fatal error database connection", err)
	}

	grpcSV := grpcClient.GrpcService{}
	_, err2 := grpcSV.NewService()
	if err2 == nil {
		fmt.Println("starting HTTP/2 gRPC server")
		fmt.Println()
	} else {
		fmt.Println("Fatal error GRPC connection: ", err2)
	}
	port := cfg.GetString("server.port")
	go func() {
		StartRest(port)
	}()

	GRPCPort := cfg.GetString("server.grpc_port")
	err2 = StartGRPC(GRPCPort)
	if err2 != nil {
		fmt.Println("Fatal error GRPC connection: ", err2)
	}
}

func StartRest(port string) {
	routes.RouteInit(engine)

	if err := engine.Run(":" + port); err != nil {
		log.Fatalln(err)
	}
}

func StartGRPC(port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("GRPC error")
		return err
	}
	var server *grpc.Server

	// register service
	server = grpc.NewServer()
	pbUser.RegisterUserServer(server, service.NewUserServer())

	// start gRPC server
	fmt.Println("starting gRPC server... port: ", port)

	return server.Serve(listen)
}
