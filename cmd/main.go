package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	btcapi "github.com/A-Danylevych/btc-api"
	"github.com/A-Danylevych/btc-api/pkg/handler"
	"github.com/A-Danylevych/btc-api/pkg/repository"
	"github.com/A-Danylevych/btc-api/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Starting and finishing the API
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	filename := viper.GetString("userdata")

	err := checkFile(filename)
	if err != nil {
		logrus.Fatalf("error occured while opening json file: %s", err.Error())
	}

	repos := repository.NewRepository(filename)
	services := service.NewService(repos, viper.GetString("api"))
	handlers := handler.NewHandler(services)

	srv := new(btcapi.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRouters()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Print("BTC API Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("BTC API Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}

//Check the existence of the file. And creating otherwise
func checkFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

//Path to config file
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
