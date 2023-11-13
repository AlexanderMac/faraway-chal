package main

import (
	"flag"
	"fmt"

	_net "github.com/AlexanderMac/faraway-chal/internal/net"
)

func main() {
	var host = flag.String("host", "", "The host to listen to (default is \"\" (all interfaces)).")
	var port = flag.Int("port", 3000, "The port to listen on (default is 3000).")
	flag.Parse()

	srvAddr := fmt.Sprintf("%s:%d", *host, *port)
	server := _net.NewServer(srvAddr)
	server.Start()
}
