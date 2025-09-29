package main

import (
	"flag"
	"fmt"
	"log"
	"rtk/api-mocker/internal/app"
	"rtk/api-mocker/internal/config"
	"rtk/api-mocker/pkg/logger"
	"runtime"
	"runtime/debug"

	"github.com/joho/godotenv"
)

const (
	flagLocal  = "local"
	flagDocker = "docker"
)

func init() {

	info, _ := debug.ReadBuildInfo()
	fmt.Println("Go version:", info.GoVersion, runtime.GOARCH)

	// without default value
	env := flag.String("env", "", "specify .env filename for flag")
	flag.Parse()
	// for dynamic Load(".env." + *env);

	if *env == flagLocal {
		if err := godotenv.Load(".env.local", ".env"); err != nil {
			log.Fatal(err)
		}
		log.Printf("⚙️ loaded .env.%s, .env\n", *env)
	}

	if *env == flagDocker {
		log.Printf("🐳 app runs in %q mod", flagDocker)
	}

}

func main() {

	config := config.New()
	zap := logger.New(config)
	defer zap.Sync()

	app := app.New(config, zap)
	app.Run()
}
