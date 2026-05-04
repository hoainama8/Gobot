package main

import (
	"botai/botaigolang/internal/agent"
	"fmt"
	"html/template" // [THEM_MOI]
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
)

type Application struct {
	agent    *agent.Agent
	tmpl     *template.Template // [THEM_MOI]
	order    *agent.OrderInfo
	temp     map[string]*template.Template // [THEM_MOI]
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Nạp file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Lỗi khi nạp file .env")
	}

	// Lấy API Key
	apiKey := os.Getenv("GEMINI_API_KEY")
	fmt.Println("Đã nạp API Key thành công!")
	newTemp, err := newTemplateCache() // [THEM_MOI]
	if err != nil {
		log.Fatal(err)
	}
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := &Application{
		agent: agent.NewAgent(apiKey),
		// tmpl:     template.Must(template.ParseFiles("web/base.html")), // [THEM_MOI]
		order:    agent.NewOrder(),
		temp:     newTemp,
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	server := &http.Server{
		Addr:    ":8081",
		Handler: app.Routes(),
	}
	http2.ConfigureServer(server, &http2.Server{})
	fmt.Printf("Starting server on port 8080\n	http://localhost:8081/chat\n")
	log.Fatal(server.ListenAndServe()) // [SUA]
}
