package main

import (
	"html/template"
	"net/http"
	"sync"

	"pomodoro.tracker/internal/controller"
	"pomodoro.tracker/internal/usecase"
)

var (
	mu *sync.Mutex
)

func main() {
	// Роут для статики (CSS, картинки)
	http.Handle("/static/", http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("static")),
	))

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	tasksUC, _ := usecase.NewTaskManager()

	httpServ, _ := controller.NewHttpWithRoutes(mu, tasksUC, tmpl)

	httpServ.Run()
}
