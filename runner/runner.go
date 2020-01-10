package runner

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

func RunRemoteServer(library interface{}) {
	RPC := rpc.NewServer()
	xmlrpcCodec := xml.NewCodec()
	//map XML-RPC methods to the go implemented functions
	//CamelCase mapping
	xmlrpcCodec.RegisterAlias("GetKeywordNames", "RobotRemoteService.GetKeywordNames")
	xmlrpcCodec.RegisterAlias("GetKeywordArguments", "RobotRemoteService.GetKeywordArguments")
	xmlrpcCodec.RegisterAlias("GetKeywordDocumentation", "RobotRemoteService.GetKeywordDocumentation")
	xmlrpcCodec.RegisterAlias("RunKeyword", "RobotRemoteService.RunKeyword")
	//pythonic mapping
	xmlrpcCodec.RegisterAlias("get_keyword_names", "RobotRemoteService.GetKeywordNames")
	xmlrpcCodec.RegisterAlias("get_keyword_arguments", "RobotRemoteService.GetKeywordArguments")
	xmlrpcCodec.RegisterAlias("get_keyword_documentation", "RobotRemoteService.GetKeywordDocumentation")
	xmlrpcCodec.RegisterAlias("run_keyword", "RobotRemoteService.RunKeyword")

	//set server to handle both XML MIME types
	RPC.RegisterCodec(xmlrpcCodec, "application/xml")
	RPC.RegisterCodec(xmlrpcCodec, "text/xml")
	protocolType := new(protocol.RobotRemoteService)
	protocolType.InitilizeRemoteLibrary(library)
	RPC.RegisterService(protocolType, "")
	http.Handle("/RPC2", RPC) //preserve option to use RPC2 endpoint
	http.Handle("/", RPC)     //but not make it required when using with Robot Framework

	//TODO: make port and host/IP address binding be configurable via CLI flags and not fixed to localhost:8270 (the default)
	log.Println("Robot remote server started on localhost:8270 under / and /RPC2 endpoints. Stop server with Ctrl+C, kill, etc. or XML-RPC method 'run_keyword' with parameter 'stop_remote_server'")
	log.Fatal(http.ListenAndServe(":8270", nil))
}
