package main

import (
	"awesomeProject/internal/orders"
	"awesomeProject/internal/orders/pb"
	"awesomeProject/pkg/db"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})

	listener, err := net.Listen("tcp", ":9100")
	if err != nil {
		log.WithFields(log.Fields{
			"new error": err,
		}).Info("Listener failed")
	}
	mongo := db.NewMongoStore()
	s := grpc.NewServer()
	server := orders.NewOrderServer(&log.Logger{}, mongo)
	reflection.Register(s)
	pb.RegisterOrdersServer(s, server)
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt)
	// Block main routine until a signal is received
	<-c
	log.Warning("Shutting down server")
	s.Stop()
	log.Println("Shutting down mongo")
	closeErr := mongo.Client.Disconnect(mongo.Context)
	if closeErr != nil {
		log.Warning("Mongo close connection failed ", closeErr)
	}
	lisErr := listener.Close()
	if lisErr != nil {
		log.Println("Something happened", err)
	}
}