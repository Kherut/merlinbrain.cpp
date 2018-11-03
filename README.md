# merlinbrain.go
An open-source Merlin Brain client in Go.

Firstly install Go CRON library and clone this repository:
```
$ go get -v github.com/robfig/cron
$ git clone https://github.com/kherut-io/merlinbrain.go
```

Compile `get_port.c`:
```
$ g++ get_port.c -o get_port
```

Make `run.sh` executable:
```
# chmod +x run.sh
```

Then run with
```
$ ./run.sh
```
