package main

import (
	"context"
	"os"
	"os/signal"
	"time"
	"user_posts/business_logic"
	"user_posts/datasource"
	"user_posts/service"
)

func main() {
	// Probably here a better approach would be having a dependency injection framework like UberFX
	postDS := datasource.NewPostClient()
	userDS := datasource.NewUserClient()
	logic := business_logic.NewUserPostLogic(userDS, postDS)
	server := service.NewServer(logic)
	server.Start(":8080")

	// Getting interrupt signal to stop
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Giving time to gracefully shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)
	defer cancel()
	server.Shutdown(ctx)
	os.Exit(0)
}