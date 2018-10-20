package main

import (
	"net/http"
	"os"
	"strings"
)

func redirectDashboard(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dashboard/", 301)
}

func main() {
	//DASHBOARD AT /dashboard
	http.Handle("/dashboard/", http.StripPrefix("/dashboard/", http.FileServer(http.Dir("template"))))

	//CONTROL AT /control
	http.HandleFunc("/control/", func(w http.ResponseWriter, r *http.Request) {
		command := strings.Join(strings.Split(r.URL.Path[1:], "/")[1:], "/")

		w.Write([]byte(command))

		file, err := os.Open("/sys/class/leds/red_led/brightness")

		w.Write([]byte("\n\nStatus: " + err.Error()))

		switch command {
		case "led/on":
			file.Write([]byte("1"))

		case "led/off":
			file.Write([]byte("0"))
		}

		file.Close()
	})

	http.HandleFunc("/", redirectDashboard)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
