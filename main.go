package main

import (
	"fmt"
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

		category := strings.Split(command, "/")[0]

		if category == "led" {
			if len(strings.Split(command, "/")) > 1 {
				arg := strings.Split(command, "/")[1]

				var cmd string

				switch arg {
				case "on":
					cmd = "echo 1 > /sys/class/leds/red_led/brightness"

				case "off":
					cmd = "echo 0 > /sys/class/leds/red_led/brightness"
				}

				out, err := exec.Command("sh", "-c", cmd).Output()
				_ = err

				if len(out) > 0 {
					w.Write([]byte("\n\nOutput: "))
					w.Write(out)
					fmt.Println(out)
				}
			}
		}
	})

	http.HandleFunc("/", redirectDashboard)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
