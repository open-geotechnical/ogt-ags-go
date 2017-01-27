package main

// This is WIP and some R+D into AGS4 data format
// This is not official = unofficial playground
// but hopefully in due course a reliable lib in foss
//
// Important Note: Its r+d so not production ready.. yet ;-)

import (

	"flag"

	"github.com/open-geotechnical/ogt-ags-go/ogtags"
	"github.com/open-geotechnical/ogt-ags-go/server"
)

func main() {

	// TODO check listen is a valid address/port etc
	listen := flag.String("server", "0.0.0.0:13777", "HTTP server address and port")

	ags_data_dict := flag.String("data-dict", "/home/ags/ags-data-dict", "Path to data dict dir")

	flag.Parse()

	// Initialise the datadict stores
	ogtags.InitLoad(*ags_data_dict)

	// TODO make server a flag, for now its on for fun
	if true {
		server.Start(*listen)
	}


}
