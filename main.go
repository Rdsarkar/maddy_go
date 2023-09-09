package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/godemo/service"
)


var (
	_version = "default"
)
// @title ABLOOMAPI
// @version 2.0
// @description An API to perform AbloomERP
// @termsOfService http://swagger.io/terms/

// @contact.name Sumon Sarker
// @contact.email suman@satcombd.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	fmt.Println("Starting Go Demo Service ", _version)
	defer fmt.Println("Done....")
	port := flag.Int("p", 8081, "Service listen port")
	bindAddress := flag.String("b", "0.0.0.0", "Bind address")
	verbose := flag.Bool("v", false, "Verbose output")
	flag.Parse()
	if *verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Config file missing")
		fmt.Println("hrms [flags] <path to config file> ")
		flag.Usage()
		os.Exit(1)
	}
	//Read the config file
	configBytes, err := ioutil.ReadFile(args[0])
	if err != nil {
		fmt.Println("Unable to read config file ", err)
		os.Exit(1)
	}
	if jotService := service.NewJOTRestService(configBytes, *verbose); jotService != nil {
		stopSignal := make(chan bool)
		termination := make(chan os.Signal)
		signal.Notify(termination, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-termination
			fmt.Println("SIGTERM/SIGINT received from os")
			stopSignal <- true
		}()
		jotService.Serve(*bindAddress, *port, stopSignal)
	} else {
		fmt.Println("Unable to start the service ...")
		os.Exit(2)
	}

}