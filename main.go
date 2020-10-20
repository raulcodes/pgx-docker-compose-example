package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func WaitForPostgres(service string, timeOut time.Duration) error {
	var pgChan = make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
			go func(s string) {
				defer wg.Done()
				for {
					_, err := net.Dial("tcp", service)
					if err == nil {
						return
					}
					time.Sleep(1 * time.Second)
				}
			}(service)
		wg.Wait()
		close(pgChan)
	}()

	select {
	case <-pgChan:
		return nil
	case <-time.After(timeOut):
		return fmt.Errorf("postgres isn't ready in %s", timeOut)
	}
}

func main() {
	if err := WaitForPostgres("db:5432", 30 * time.Second); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}