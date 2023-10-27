package main

import (
	"flag"
	"fmt"

	_net "github.com/AlexanderMac/faraway-chal/internal/net"
)

var host = flag.String("host", "", "The host to listen to (default is \"\" (all interfaces)).")
var port = flag.Int("port", 3000, "The port to listen on (default is 3000).")

func main() {
	flag.Parse()

	srvAddr := fmt.Sprintf("%s:%d", *host, *port)
	_net.RunServer(srvAddr)
}
