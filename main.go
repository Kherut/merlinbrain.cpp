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
	"log"

	"github.com/robfig/cron"
	"github.com/googollee/go-socket.io"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

type Device struct {
	ID string
	Name string
	Alias string
	IP string
	Role string
	Status string
	TimeConnected string
	Approved bool
}

var devices []Device
var keys []string
var cfg map[string]interface{}
var socket socketio.Socket
var connected int
var httpWriter http.ResponseWriter

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

func printMode(str, mode string) {
	if mode == "http" {
		httpWriter.Write([]byte(str))
	} else if mode == "socket" {
		socket.Emit("message", str)
	} else if mode == "console" {
		fmt.Print(str)
	}
}

func randomString(length int64) string {
	//GENERATE RANDOM STRING (PROBABILITY OF REPEATING ONE IS 1/62^[idLength])
	const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	b := make([]byte, length)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
	}
	
    return false
}

func mbCommand(command, mode string) {
	category := strings.Split(command, "/")[0]
	var key string

	if len(strings.Split(command, "@")) > 1 {
		key = strings.Split(command, "@")[1]
		command = strings.Split(command, "@")[0]
	}

	var arg []string

	if len(strings.Split(command, "/")) > 1 {
		arg = strings.Split(command, "/")[1:]
	}

	if category == "devices" {
		if (arr(arg, 0) == "new" && len(arg) == 4) {
			ip_val := arg[1]
			name_val := arg[2]
			role_val := arg[3]
			time_connected_val := time.Now().Format("2006.01.02 15:04:05")

			id_length, _ := cfg["idLength"].(json.Number).Int64()
			id_val := randomString(id_length)

			key_length, _ := cfg["keyLength"].(json.Number).Int64()
			key := randomString(key_length)

			devices = append(devices, Device{ID: id_val, Name: name_val, Alias: "", IP: ip_val, Role: role_val, Status: "UP", TimeConnected: time_connected_val, Approved: false})
			keys = append(keys, key)

			min := fmt.Sprintf("%d", 40000)
			max := fmt.Sprintf("%d", 41000)

			port := runCmd("./get_port " + min + " " + max)

			printMode(id_val + "!" + port + "!" + key, mode)

			deviceB, _ := json.Marshal(devices[len(devices) - 1])
			deviceStr := string(deviceB)

			//LOG
			//fmt.Println(deviceStr)

			if connected > 0 {
				socket.Emit("message", command + "!" + deviceStr)
			}

			var cmd string

			if role_val == "CLIENT" {
				if cfg["development"] == true {
					cmd = "gnome-terminal -x sh -c 'netcat -l " + port + "'"
				}
			} else if role_val == "SERVER" {
				if cfg["development"] == true {
					cmd = "gnome-terminal -x sh -c 'sleep 2.5; netcat " + ip_val + " " + port + "'"
				}
			}

			//LOG
			//fmt.Println(cmd)

			runCmd(cmd)
		} else if arr(arg, 0) == "status" {
			if len(arg) == 2 {
				for i := range devices {
					if devices[i].ID == arg[1] {
						printMode(devices[i].Status, mode)
					}
				}
			} else if len(arg) == 3 {
				if contains(keys, key) {
					var i int

					for i = range devices {
						if devices[i].ID == arg[1] && devices[i].Approved {
							devices[i].Status = arg[2];

							printMode(arg[2], mode)

							if(connected > 0) {
								b, _ := json.Marshal(devices[i])
								changeText := string(b)

								//LOG
								//fmt.Println("change!" + changeText)

								socket.Emit("message", "change!" + changeText)
							}
						}
					}
				} else {
					printMode("Wrong key", mode)
				}
			}
		} else if arr(arg, 0) == "all" && len(arg) == 1 {
			var text bytes.Buffer

			text.WriteString(command)
			text.WriteString("!")

			for _, element := range devices {
				text.WriteString(element.ID + "\n")
				text.WriteString(element.Name + "\n")
				text.WriteString(element.IP + "\n")
				text.WriteString(element.Role + "\n")
				text.WriteString(element.Status + "\n")
				text.WriteString(element.TimeConnected + "\n")
			}

			printMode(text.String(), mode)
		} else if arr(arg, 0) == "all" && arr(arg, 1) == "json" {
			var text bytes.Buffer

			text.WriteString(command)
			text.WriteString("!")

			for _, element := range devices {
				if element != (Device{"", "", "", "", "", "", "", false}) {
					deviceB, _ := json.Marshal(element)
					deviceStr := string(deviceB)

					text.WriteString(deviceStr + "\n")
				}
			}

			printMode(text.String(), mode)
		} else if arr(arg, 0) == "approve" && len(arg) == 2 {
			for i := range devices {
				if devices[i].ID == arg[1] {
					devices[i].Approved = true

					b, _ := json.Marshal(devices[i])
					changeText := string(b)

					socket.Emit("message", "change!" + changeText)
					
					printMode("OK", mode)

					break
				}
			}
		}
	} else if category == "info" {
		switch arg[0] {
		case "temperature":
			if arr(arg, 1) == "all" {
				//DEVELOPMENT
				if cfg["development"] == "false" {
					printMode(runCmd("cat ./data/temperature.data"), mode)
				}
			} else {
				//DEVELOPMENT
				if cfg["development"] == "false" {
					printMode(runCmd("cat /sys/devices/virtual/thermal/thermal_zone0/temp"), mode)
				}
			}
		case "uptime":
			//DEVELOPMENT
			if cfg["development"] == "false" {
				printMode(strings.Split(runCmd("cat /proc/uptime"), ".")[0], mode)
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
	} else if category == "clear" {
		runCmd("clear")
	} else if category == "close" {
		os.Exit(0)
	}
}

func main() {
	firstRun := true

	config.Load(file.NewSource(
		file.WithPath("config.yaml"),
	))

	cfg = config.Map()

	if cfg["logging"] == true {
		logFile, err := os.OpenFile(config.Get("logFile").String("logfile.log"), os.O_CREATE|os.O_APPEND, 0644)

		if err != nil {
			log.Fatal(err)
		}
	
		log.SetOutput(logFile)
		log.Print("Start")
	}

	rand.Seed(time.Now().UTC().UnixNano())

	connected = 0

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

	//CONTROL AT /control
	http.HandleFunc("/control/", func(w http.ResponseWriter, r *http.Request) {
		httpWriter = w

		command := strings.Join(strings.Split(r.URL.Path[1:], "/")[1:], "/")

		mbCommand(command, "http")
	})

	//CONTROL ON /socket.io/
	srvSio, _ := socketio.NewServer(nil)

	srvSio.On("connection", func(sio socketio.Socket) {
		connected = connected + 1
		socket = sio

		socket.On("message", func(msg string) {
			mbCommand(msg, "socket")
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

	go func() {
		var s string

		for true {
			if(!firstRun) {
				fmt.Println()
			}
			
			firstRun = false

			fmt.Print("> ")
			fmt.Scanln(&s)

			mbCommand(s, "console")
		}
	} ()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}