/*
*
*
* The main entry of the server
*
*
*
*
*
*/

package main

import(
	"fmt"

	"net/http"

	"flag"

	"log"

	"./service/ticktimer"
)

var(
	// Specify the flag set
	serverAddr = flag.String( "a", "", "Specifies the IP address of server" )
	serverPort = flag.String( "p", "6161", "Specify the port of server which listening the internet connection" )
)

func main(){
	// Parse the flags
	flag.Parse()

	addr := string( *serverAddr ) + ":" + string( *serverPort )
	fmt.Printf( "The address of server: %s\n", addr )

	// Start up the http server
	// A better usage case are as below
	// s := &http.Server{
	// 	Addr:           ":8080",
	// 	Handler:        myHandler,
	// 	ReadTimeout:    10 * time.Second,
	// 	WriteTimeout:   10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }

	// Register the handler that handle the request of static file
	http.Handle( "/", http.FileServer( http.Dir( "./www" ) ) )

	// The ticktimer
	timer := ticktimer.Timer{}
	http.Handle( "/timer", timer )

	// Start the heep server
	log.Fatal( http.ListenAndServe( addr, nil ) )
}