package main

import (
	"net/http"
	"os/exec"
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

		switch command {
		case "led/on":
			cmd := exec.Command("echo", "1 > /sys/class/leds/red_led/brightness")
			err := cmd.Run()

			w.Write([]byte("\n\n" + "Status: " + err.Error()))
		case "led/off":
			cmd := exec.Command("echo", "0 > /sys/class/leds/red_led/brightness")
			err := cmd.Run()

			w.Write([]byte("\n\n" + "Status: " + err.Error()))
		}
	})

	http.HandleFunc("/", redirectDashboard)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
