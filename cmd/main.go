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
	"github.com/vladoatanasov/logrus_amqp"
)

// Starting and finishing the API
func main() {
	log := logrus.New()
	log.SetFormatter(new(logrus.JSONFormatter))

	hook := logrus_amqp.NewAMQPHook("localhost:5672", "guest", "guest", "BTC-API", "")

	log.Hooks.Add(hook)

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	filename := viper.GetString("userdata")

	if err := checkFile(filename); err != nil {
		log.Fatalf("error occured while opening json file: %s", err.Error())
	}

	repos := repository.NewRepository(filename)
	services := service.NewService(repos, viper.GetString("microservice"))
	handlers := handler.NewHandler(services, log)

	srv := new(btcapi.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRouters()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	log.Debug("BTC API Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Debug("BTC API Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Errorf("error occured on server shutting down: %s", err.Error())
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
