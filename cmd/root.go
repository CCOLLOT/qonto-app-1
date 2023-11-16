package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/CCOLLOT/appnametochange/internal/app"
	"github.com/CCOLLOT/appnametochange/internal/logger"
	"github.com/spf13/cobra"
)

type GracefulService interface {
	Name() string
	Start() error
	Shutdown() error
}

func InitAndRunCommand() error {
	rootCmd := &cobra.Command{
		Use:   "root",
		Short: "Run the main process",
	}
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Run the main process",
		Run: func(cmd *cobra.Command, args []string) {
			Run()
		},
	}
	rootCmd.AddCommand(startCmd)
	return rootCmd.Execute()
}

func Run() error {
	log, err := logger.New()
	if err != nil {
		return err
	}
	services := []GracefulService{}
	app, err := app.New(log)
	if err != nil {
		return err
	}
	go app.Start()

	services = append(services, app)

	waitGracefulShutdown := make(chan any)
	var wg sync.WaitGroup
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
		<-s
		for _, svc := range services {
			wg.Add(1)
			go func(svc GracefulService) {
				defer wg.Done()
				fmt.Println(fmt.Sprintf("gracefully shutting down %s...", svc.Name()))
				err := svc.Shutdown()
				if err != nil {
					fmt.Println(fmt.Sprintf("failed to gracefully shut down %s, err: %s", svc.Name(), err.Error()))
					return
				}
				fmt.Println(fmt.Sprintf("successfully shut down %s", svc.Name()))
				fmt.Println("bye.")
			}(svc)
		}
		wg.Wait()
		close(waitGracefulShutdown)
	}()
	<-waitGracefulShutdown
	return nil
}
