package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"bytes"
	"math/rand"
	"encoding/json"

	"github.com/robfig/cron"
	"github.com/googollee/go-socket.io"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

type Device struct {
	ID string
	Name string
	IP string
	Role string
	Status string
	TimeConnected string
}

func redirectDashboard(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dashboard/", 301)
}

func runCmd(cmd string) string {
	out, err := exec.Command("sh", "-c", cmd).Output()
	_ = err

	return string(out)
}

func arr(array []string, index int) string {
	if len(array) >= index + 1 {
		return array[index]
	}

	return ""
}

func main() {
	config.Load(file.NewSource(
		file.WithPath("config.yaml"),
	))

	cfg := config.Map()

	rand.Seed(time.Now().UTC().UnixNano())

	var devices []Device
	var socket socketio.Socket

	connected := 0

	//DEVELOPMENT - CRON
	if cfg["development"] == "false" {
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
			if (arr(arg, 0) == "new" && len(arg) >= 4) {
				ip_val := arg[1]
				name_val := arg[2]
				role_val := arg[3]
				time_connected_val := time.Now().Format("2006.01.02 15:04:05")

				//GENERATE UID FOR THE DEVICE (PROBABILITY OF REPEATING A UID IS 1/62^[id_length])
				const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

				id_length, _ := cfg["id_length"].(json.Number).Int64()

				b := make([]byte, id_length)

				for i := range b {
					b[i] = letterBytes[rand.Intn(len(letterBytes))]
				}

				id_val := string(b)

				devices = append(devices, Device{ID: id_val, Name: name_val, IP: ip_val, Role: role_val, Status: "UP", TimeConnected: time_connected_val})

				min := fmt.Sprintf("%d", 40000)
				max := fmt.Sprintf("%d", 41000)

				port := runCmd("./get_port " + min + " " + max)

				w.Write([]byte(id_val + "!" + port))

				deviceB, _ := json.Marshal(devices[len(devices) - 1])
				deviceStr := string(deviceB)

				fmt.Println(deviceStr)

				if connected > 0 {
					socket.Emit("message", command + "!" + deviceStr)
				}

				var cmd string

				if role_val == "CLIENT" {
					if cfg["development"] == true {
						fmt.Println("INS ROLE_VAL: " + role_val)
						cmd = "gnome-terminal -x sh -c 'netcat -l " + port + "'"
					}
				} else if role_val == "SERVER" {
					if cfg["development"] == true {
						cmd = "gnome-terminal -x sh -c 'sleep 2.5; netcat " + ip_val + " " + port + "'"
					}
				}

				w.(http.Flusher).Flush()

				fmt.Println(cmd)

				runCmd(cmd)
			} else if (arr(arg, 0) == "status" && len(arg) >= 3) {
				var i int

				for i = range devices {
					if devices[i].ID == arg[1] {
						devices[i].Status = arg[2];

						w.Write([]byte("OK"))

						if(connected > 0) {
							b, _ := json.Marshal(devices[i])
							changeText := string(b)

							fmt.Println("change!" + changeText)

							socket.Emit("message", "change!" + changeText)
						}
					}
				}
			}
		}
	})

	//CONTROL ON /socket.io/
	srvSio, _ := socketio.NewServer(nil)

	srvSio.On("connection", func(sio socketio.Socket) {
		connected = connected + 1
		socket = sio

		socket.On("message", func(msg string) {
			category := strings.Split(msg, "/")[0]

			var arg []string

			if len(strings.Split(msg, "/")) > 1 {
				arg = strings.Split(msg, "/")[1:]
			}

			if category == "info" {
				switch arg[0] {
				case "temperature":
					if arr(arg, 1) == "all" {
						//DEVELOPMENT
						if cfg["development"] == "false" {
							socket.Emit("message", runCmd("cat ./data/temperature.data"))
						}
					} else {
						//DEVELOPMENT
						if cfg["development"] == "false" {
							socket.Emit("message", runCmd("cat /sys/devices/virtual/thermal/thermal_zone0/temp"))
						}
					}
				case "uptime":
					//DEVELOPMENT
					if cfg["development"] == "false" {
						sio.Emit("message", strings.Split(runCmd("cat /proc/uptime"), ".")[0])
					}
				}
			} else if category == "led" {
				var cmd string

				switch arg[0] {
				case "on":
					//DEVELOPMENT
					if cfg["development"] == "false" {
						cmd = "echo 1 > /sys/class/leds/red_led/brightness"
					}

				case "off":
					//DEVELOPMENT
					if cfg["development"] == "false" {
						cmd = "echo 0 > /sys/class/leds/red_led/brightness"
					}
				}

				_ = runCmd(cmd)
			} else if category == "devices" {
				if arr(arg, 0) == "all" && len(arg) == 1 {
					var text bytes.Buffer

					text.WriteString(msg)
					text.WriteString("!")

					for _, element := range devices {
						text.WriteString(element.ID + "\n")
						text.WriteString(element.Name + "\n")
						text.WriteString(element.IP + "\n")
						text.WriteString(element.Role + "\n")
						text.WriteString(element.Status + "\n")
						text.WriteString(element.TimeConnected + "\n")
					}

					socket.Emit("message", text.String())
				} else if arr(arg, 0) == "all" && arr(arg, 1) == "json" {
					var text bytes.Buffer

					text.WriteString(msg)
					text.WriteString("!")

					for _, element := range devices {
						if element != (Device{"", "", "", "", "", ""}) {
							deviceB, _ := json.Marshal(element)
							deviceStr := string(deviceB)

							text.WriteString(deviceStr + "\n")
						}
					}

					socket.Emit("message", text.String())
				}
			}
		})

		socket.On("disconnection", func() {
			connected = connected - 1
		})
	})

	srvSio.On("error", func(so socketio.Socket, err error) {
		fmt.Println("Error: ", err)
	})

	http.Handle("/socket.io/", srvSio)

	http.HandleFunc("/", redirectDashboard)

	//DEVELOPMENT
	if cfg["development"] == "false" {
		runCmd("kill -9 " + os.Args[2])
		runCmd("echo 1 > /sys/class/leds/green_led/brightness")
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}