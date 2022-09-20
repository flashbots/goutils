package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/flashbots/go-utils/httplogger"
	"go.uber.org/zap"
)

var (
	listenAddr = "localhost:8124"
	// logJSON    = os.Getenv("LOG_JSON") == "1"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
	w.WriteHeader(http.StatusOK)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	log.Error("this is an error", "err", errors.New("testError"))
	http.Error(w, "this is an error", http.StatusInternalServerError)
}

func PanicHandler(w http.ResponseWriter, r *http.Request) {
	panic("foo!")
}

func getZapLogger() *zap.SugaredLogger {
	// logger, _ := zap.NewDevelopment()
	logger, _ := zap.NewProduction()
	log := logger.Sugar()
	return log
}

func main() {
	// logFormat := log.TerminalFormat(true)
	// if logJSON {
	// 	logFormat = log.JSONFormat()
	// }

	// log.Root().SetHandler(log.LvlFilterHandler(log.LvlDebug, log.StreamHandler(os.Stderr, logFormat)))

	fmt.Println("Webserver running at", listenAddr)
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", HelloHandler)
	mux.HandleFunc("/error", ErrorHandler)
	mux.HandleFunc("/panic", PanicHandler)
	loggedRouter := httplogger.LoggingMiddlewareZap(getZapLogger(), mux)
	if err := http.ListenAndServe(listenAddr, loggedRouter); err != nil {
		log.Crit("webserver failed", "err", err)
	}

}
