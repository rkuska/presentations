package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/server"
)

func main() {
	r := server.NewRouter()
	server := server.NewServer("127.0.0.1:8080", r)

	errs := make(chan error)
	go func() {
		fmt.Println("listening on :8080")
		errs <- server.ListenAndServe()
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	exit := <-errs
	fmt.Println("exit error", exit.Error())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("shutdown error", err.Error())
	}

}
