package main

import (
	"log"
	"net/http"

	"github.com/daluu/gorrs/protocol"
	"github.com/divan/gorilla-xmlrpc/xml"
	"github.com/gorilla/rpc"
)

/* add to import list of github.com/daluu/gorrs/protocol/protocol.go,
 * the (exported) go remote (test) library packages
 * to be served by this remote server via reflection. To do that since we
 * have to explicitly reference packages to reflect on and not be able to
 * just pass in package reference at runtime?
 */

/* TODO: also look into whether there's any other alternative to
 * gorilla/rpc and divan/gorilla-xmlrpc/xml in case of issues with XML-RPC
 * support / implementation in go. Or what can be done to extend them to do
 * what we need for a go-based Robot Framework generic remote library server
 *
 * Full spec for said server:
 * https://github.com/robotframework/RemoteInterface
 */

func main() {
	RPC := rpc.NewServer()
	xmlrpcCodec := xml.NewCodec()
	RPC.RegisterCodec(xmlrpcCodec, "text/xml")
	// is there a way to register XML-RPC service such that when XML-RPC client calls the service
	// they refer to service w/o a namespace? e.g. "RunKeyword" instead of "RobotRemoteService.RunKeyword"?
	// see https://github.com/divan/gorilla-xmlrpc/issues/14
	// and https://github.com/gorilla/rpc/issues/48
	RPC.RegisterService(new(protocol.RobotRemoteService), "")
	http.Handle("/RPC2", RPC) //preserve option to use RPC2 endpoint
	http.Handle("/", RPC)     //but not make it required when using with Robot Framework

	//TODO: make port and host/IP address binding be configurable via CLI flags and not fixed to localhost:8270 (the default)
	log.Println("Robot remote server started on localhost:8270 under / and /RPC2 endpoints. Stop server with Ctrl+C, kill, etc. or XML-RPC method 'run_keyword' with parameter 'stop_remote_server'\n")
	log.Fatal(http.ListenAndServe(":8270", nil))
}
