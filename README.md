# merlinbrain.go
An open-source Merlin Brain client in Go.

## Compile
Firstly install Go CRON library and clone this repository:
```
$ go get -v github.com/robfig/cron
$ git clone https://github.com/kherut-io/merlinbrain.go
```

Compile `get_port.cpp`:
```
$ g++ get_port.cpp -o get_port
```

## Run
### Run *Brain* normally (with LED support)
Make `run.sh` executable:
```
# chmod +x run.sh
```

Then run with
```
$ ./run.sh
```

### Run *Brain* in development mode
In `main.go` find line
```
development := false
```
and change it to
```
development := true
```

Now run using
```
$ go run main.go
```