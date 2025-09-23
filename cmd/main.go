package main

import (
	"flag"
	"fmt"
	"log"
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

	zap, err := logger.New(config)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	zap.Debug("debug message")
	zap.Info("info message")
	zap.Warn("warn message")
	zap.Error("error message")
	// logger := logger.New(config.Environment)

	// app := app.NewApp(config, logger, db, r2store)
	// app.Run()
}
