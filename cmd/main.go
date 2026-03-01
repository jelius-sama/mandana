package main

import "C"
import (
    "context"
    "fmt"
    "mandana/api"
    "mandana/api/middleware"
    "mandana/constants"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/jelius-sama/logger"
)

var (
    IS_PROD string
    PORT    string
    UPTIME  time.Time = time.Now()
)

func init() {
    type ResolveEnv struct {
        key      string
        injected *string
        fallback string
    }

    var resolveEnv func(arg ResolveEnv) = func(arg ResolveEnv) {
        if arg.key == "" {
            // no key value passed, misuse of `resolveEnv` function, just return.
            return
        }

        // No compile time injected variable passed or available
        // use fallback value and key to assign the env directly.
        if arg.injected == nil {
            if arg.fallback != "" {
                if err := os.Setenv(arg.key, arg.fallback); err != nil {
                    logger.Panic("failed to set environment variable:\n", err)
                }
            }
            return
        }

        // Compile time injected variable available, assign the env.
        if *arg.injected != "" {
            if err := os.Setenv(arg.key, *arg.injected); err != nil {
                logger.Panic("failed to set environment variable:\n", err)
            }
            return
        }

        // Compile time injected variable available but empty string,
        // assign the env using the fallback value and mutate the
        // `injected` pointer variable value.
        if err := os.Setenv(arg.key, arg.fallback); err != nil {
            logger.Panic("failed to set environment variable:\n", err)
        }
        *arg.injected = arg.fallback
    }

    resolveEnv(ResolveEnv{
        key:      "IS_PROD",
        injected: &IS_PROD,
        fallback: "FALSE",
    })
    resolveEnv(ResolveEnv{
        key:      "PORT",
        injected: &PORT,
        fallback: ":8080",
    })

    logger.Configure(logger.Cnf{
        IsDev: logger.IsDev{
            EnvironmentVariable: logger.StringPtr("IS_PROD"), // OR: nil
            ExpectedValue:       logger.StringPtr("FALSE"),   // OR: nil
            DirectValue:         nil,                         // OR: logger.BoolPtr(IS_PROD != "TRUE"),
        },
        UseSyslog: false,
    })
}

func main() {
    fmt.Println("\n\033[0;36m"+constants.AppName().Japanese, "version", constants.Version, "\033[0m")
    logger.Info("Starting server on port", PORT)

    var quit chan os.Signal = make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    var server *http.Server = &http.Server{
        Addr: PORT,
        Handler: middleware.RecoveryMiddleware(
            middleware.LoggingMiddleware(
                middleware.NotFound(
                    api.Router(),
                ),
            ),
        ),
    }

    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal("Failed to start server on port "+PORT+"\n", err)
        }
    }()

    <-quit
    var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var deadline, _ = ctx.Deadline()
    var done chan struct{} = make(chan struct{})

    var ticker *time.Ticker = time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    go func() {
        if err := server.Shutdown(ctx); err != nil {
            logger.TimedFatal("Server forced to shutdown:", err)
        }
        close(done)
    }()

    for {
        select {
        case <-done:
            logger.TimedInfo("Server stopped.")
            return

        case <-ctx.Done():
            logger.TimedInfo("Timeout reached:", ctx.Err())
            return

        case <-ticker.C:
            var remaining int = int(time.Until(deadline).Seconds())
            if remaining < 0 {
                remaining = 0
            }

            fmt.Printf("\r\033[K\033[0;36m[INFO] Shutting down in %d seconds...\033[0m", remaining)
        }
    }
}

