package main

import (
	"log"
	"net/http"

	"github.com/simrantanwani226/compete-finder/gen/compete/v1/competev1connect"
	"github.com/simrantanwani226/compete-finder/internal/handler"
	"github.com/simrantanwani226/compete-finder/internal/provider/yc"
)

func main() {
	ycProvider := yc.New("https://yc-oss.github.io/api/companies/all.json")
	h := handler.NewHandler(ycProvider)
	mux := http.NewServeMux()
	path, connectHandler := competev1connect.NewCompeteServiceHandler(h)
	mux.Handle(path, connectHandler)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
