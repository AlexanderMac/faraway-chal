package main

import (
	"flag"
	"fmt"

	_net "github.com/AlexanderMac/faraway-chal/internal/net"
)

func main() {
	var host = flag.String("host", "localhost", "The hostname or IP to connect to (defaults to \"localhost\").")
	var port = flag.Int("port", 3000, "The port to connect to (defaults to 3000).")
	flag.Parse()

	srvAddr := fmt.Sprintf("%s:%d", *host, *port)
	client := _net.NewClient(srvAddr)
	client.Start()
}
