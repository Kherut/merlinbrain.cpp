package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"bytes"

	"github.com/robfig/cron"
	"github.com/googollee/go-socket.io"
)

type Device struct {
	Name string
	IP string
	Role string
	Status string
	ConnectedAt string
}

func redirectDashboard(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dashboard/", 301)
}

func runCmd(cmd string) string {
	out, err := exec.Command("sh", "-c", cmd).Output()
	_ = err

	return string(out)
}

func main() {
	development := true

	var devices []Device
	var socket socketio.Socket

	//DEVELOPMENT
	if !development {
		c := cron.New()
		counter := 0

		c.AddFunc("0 0 * * * *", func() {
			currentTime := time.Now()
	
			cmd := "echo " + currentTime.Format("2006.01.02-15:04:05") + "#" + strings.Split(runCmd("cat /sys/devices/virtual/thermal/thermal_zone0/temp"), "\n")[0] + " >> ./data/temperature.data"
	
			runCmd("tail -n 48 data/temperature.data > data/temperature.tmp && mv data/temperature.tmp data/temperature.data")
			runCmd(cmd)
	
			counter += 1
		})

		c.Start()
	}

	//DASHBOARD AT /dashboard
	http.Handle("/dashboard/", http.StripPrefix("/dashboard/", http.FileServer(http.Dir("www"))))

	//CONTROL AT /control (THERE'S ONLY devices/new IN THERE)
	http.HandleFunc("/control/", func(w http.ResponseWriter, r *http.Request) {
		command := strings.Join(strings.Split(r.URL.Path[1:], "/")[1:], "/")

		category := strings.Split(command, "/")[0]

		var arg []string

		if len(strings.Split(command, "/")) > 1 {
			arg = strings.Split(command, "/")[1:]
		}

		if category == "devices" {
			if len(arg) > 0 {
				if (arg[0] == "new" && len(arg) >= 4) {
					ip := arg[1]
					name := arg[2]
					role := arg[3]
					connectedat := time.Now().Format("2006.01.02-15:04:05")

					devices = append(devices, Device{Name: name, IP: ip, Role: role, Status: "UP", ConnectedAt: connectedat})

					min := fmt.Sprintf("%d", 40000)
					max := fmt.Sprintf("%d", 41000)

					port := runCmd("./get_port " + min + " " + max)

					w.Write([]byte(port))

					socket.Emit("message", command + "///" + "{\"name\": \"" + name + "\", \"ip\": \"" + ip + "\", \"role\": \"" + role + "\", \"status\": \"" + "UP" + "\", \"connectedat\": \"" + connectedat + "\"}")

					var cmd string

					if role == "CLIENT" {
						if development {
							cmd = "gnome-terminal -x sh -c 'netcat -l " + port + "'"
						}
					} else if role == "SERVER" {
						if development {
							cmd = "gnome-terminal -x sh -c 'sleep 2.5; netcat " + ip + " " + port + "'"
						}
					}

					w.(http.Flusher).Flush()
					runCmd(cmd)
				}
			}
		}
	})

	//CONTROL ON /socket.io/
	srvSio, _ := socketio.NewServer(nil)

	srvSio.On("connection", func(sio socketio.Socket) {
		socket = sio

		socket.On("message", func(msg string) {
			category := strings.Split(msg, "/")[0]

			var arg []string

			if len(strings.Split(msg, "/")) > 1 {
				arg = strings.Split(msg, "/")[1:]
			}

			fmt.Print(category + " -> ")
			fmt.Println(arg)

			if category == "info" {
				switch arg[0] {
				case "temperature":
					if len(arg) > 1 {
						if arg[1] == "all" {
							//DEVELOPMENT
							if !development {
								sio.Emit("message", runCmd("cat ./data/temperature.data"))
							}
						}
					} else {
						//DEVELOPMENT
						if !development {
							sio.Emit("message", runCmd("cat /sys/devices/virtual/thermal/thermal_zone0/temp"))
						}
					}
				case "uptime":
					//DEVELOPMENT
					if !development {
						sio.Emit("message", strings.Split(runCmd("cat /proc/uptime"), ".")[0])
					}
				}
			} else if category == "led" {
				var cmd string

				switch arg[0] {
				case "on":
					//DEVELOPMENT
					if !development {
						cmd = "echo 1 > /sys/class/leds/red_led/brightness"
					}

				case "off":
					//DEVELOPMENT
					if !development {
						cmd = "echo 0 > /sys/class/leds/red_led/brightness"
					}
				}

				_ = runCmd(cmd)
			} else if category == "devices" {
				if len(arg) == 1 {
					var text bytes.Buffer

					text.WriteString(msg)
					text.WriteString("///")

					for _, element := range devices {
						text.WriteString(element.Name + "\n")
						text.WriteString(element.IP + "\n")
						text.WriteString(element.Role + "\n")
						text.WriteString(element.Status + "\n")
						text.WriteString(element.ConnectedAt + "\n")
					}

					socket.Emit("message", text.String())
				} else if len(arg) >= 2 && arg[1] == "json" {
					var text bytes.Buffer

					text.WriteString(msg)
					text.WriteString("///")

					for _, element := range devices {
						if element != (Device{"", "", "", "", ""}) {
							text.WriteString("{ \"name\": \"" + element.Name + "\", \"ip\": \"" + element.IP + "\", \"role\": \"" + element.Role + "\", \"status\": \"" + element.Status + "\", \"connectedat\": \"" + element.ConnectedAt + "\" }\n")
						}
					}

					fmt.Println(text.String())
					socket.Emit("message", text.String())
				}
			}
		})
	})

	srvSio.On("error", func(so socketio.Socket, err error) {
		fmt.Println("Error: ", err)
	})

	http.Handle("/socket.io/", srvSio)

	http.HandleFunc("/", redirectDashboard)

	//DEVELOPMENT
	if !development {
		runCmd("kill -9 " + os.Args[2])
		runCmd("echo 1 > /sys/class/leds/green_led/brightness")
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}