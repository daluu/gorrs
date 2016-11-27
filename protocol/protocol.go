package protocol

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"time"
)

/* add to import list above, the (exported) go remote (test) library packages
 * to be served by this remote server via reflection. To do that since we
 * have to explicitly reference packages to reflect on and not be able to
 * just pass in package reference at runtime?
 */

/* TODO: look into what can be reflected by go in terms of finding stuff in the
 * imported packages namespace, execute an exported function, optionally
 * extract out or lookup the arguments (#, name, type) for an exported function,
 * and optionally extract the go documentation for an exported function,
 * all via reflection (or something equivalent in go).
 * If not feasible, then go users will have to statically "serve" a chosen
 * package rather than dynamically serve via reflection for what's in the
 * imported namespace. Test library imports would be done in the server main program "../main.go"
 *
 * Just search online for go reflection for resources, or here's some from a search:
 * http://www.jerf.org/iri/post/2945
 * https://blog.golang.org/laws-of-reflection
 * http://blog.ralch.com/tutorial/golang-reflection/
 * http://merbist.com/2011/06/27/golang-reflection-exampl/
 * https://blog.gopheracademy.com/birthday-bash-2014/advanced-reflection-with-go-at-hashicorp/
 * https://jimmyfrasche.github.io/go-reflection-codex/
 * https://gist.github.com/drewolson/4771479
 * https://golang.org/pkg/reflect/
 * https://godoc.org/?q=reflect
 *
 * also, how to redirect stdout & stderr (here & in reflected code),
 * such that we pipe a copy of that data into variables
 * for sending back with XML-RPC call for RunKeyword method?
 */

type RobotRemoteService struct{}

type KeywordNamesReturnValue struct {
	Keywords []interface{}
}

//sample XML-RPC input: <methodCall><methodName>GetKeywordNames</methodName><params></params></methodCall>
/* sample XML-RPC output:
 * <methodResponse><params><param><value><array><data>
 *   <value><string>TruthOfLife</string></value>
 *   <value><string>StringsShouldBeEqual</string></value>
 *   <value><string>StopRemoteServer</string></value>
 * </data></array></value></param></params></methodResponse>
 */
func (h *RobotRemoteService) GetKeywordNames(r *http.Request, args *struct{}, reply *KeywordNamesReturnValue) error {
	//TODO: use reflection to generate array of keywords (found in the imported namespace) to return in reply
	//maybe rather than all imported packages, restrict to a specific one, etc. as specified at server startup?

	//add special keyword built-in to the server:
	reply.Keywords = append(reply.Keywords, "StopRemoteServer")
	return nil
}

func (h *RobotRemoteService) StopRemoteServer() {
	//TODO: no need to call this function with goroutine if we make stopping the server more idiomatic with proper "shutdown"
	//perhaps make use of channels, and have the stop server task wait on channel and only pass to channel
	//when this XML-RPC method is called? And/or other ways to stop the server...

	time.Sleep(5 * time.Second) //let's arbitrarily set delay at 5 seconds
	log.Printf("Remote server/library shut down at %v\n", time.Now())
	_stopRemoteServer()
}

func _stopRemoteServer() {
	os.Exit(1)
}

type RunKeywordReturnValue struct {
	Return    interface{} `xml:"return"`
	Status    string      `xml:"status"`
	Stdout    string      `xml:"output"`
	Stderr    string      `xml:"error"`
	Traceback string      `xml:"traceback"`
}

type KeywordAndArgsInput struct {
	KeywordName     string
	KeywordAguments []interface{}
}

/* e.g. sample XML-RPC input
 * <methodCall><methodName>RunKeyword</methodName>
 *   <params>
 *     <param><value><string>KeywordName</string></value></param>
 *     <param><value><array><data>
 *       <value><string>keyword_arg1</string></value>
 *       <value><string>keyword_arg2</string></value>
 *     </data></array></value></param>
 *   </params></methodCall>
 *
 * sample XML-RPC output
 * <methodResponse>
 *   <params>
 *     <param>
 *       <value><struct>
 *         <member>
 *           <name>return</name>
 *           <value><int>42</int></value>
 *         </member>
 *         <member>
 *           <name>status</name>
 *           <value><string>PASS</string></value>
 *        </member>
 *        <member>
 *          <name>output</name>
 *          <value><string></string></value>
 *        </member>
 *        <member>
 *          <name>error</name>
 *          <value><string></string></value>
 *         </member>
 *         <member>
 *           <name>traceback</name>
 *           <value><string></string></value>
 *         </member>
 *       </struct></value>
 *     </param>
 *   </params>
 * </methodResponse>
 */
//this function doesn't fully work yet, see
//https://github.com/divan/gorilla-xmlrpc/issues/ #16 and 18
func (h *RobotRemoteService) RunKeyword(r *http.Request, args *KeywordAndArgsInput, reply *RunKeywordReturnValue) error {
	//use reflection to run function "keyword name" out of 1st arg
	//with 2nd arg (array) containing the args for the keyword function
	//sample debug/test code for now...
	log.Printf("keyword: %+v\n", args.KeywordName)
	if args.KeywordName == "StopRemoteServer" {
		go h.StopRemoteServer()
	}
	log.Printf("args: %+v\n", args.KeywordAguments)
	for _, a := range args.KeywordAguments {
		log.Printf("arg: %+v\n", a)
		switch reflect.TypeOf(a).Kind() {
		case reflect.Slice:
			log.Printf("args:\n")
			s := reflect.ValueOf(a)
			for i := 0; i < s.Len(); i++ {
				log.Printf("%v: %+v\n", i, s.Index(i))
			}
		default:
			log.Println("Somehow didn't get an array of arguments for keyword.")
		}
	}
	//and return the results in struct below (sample static output for now):
	var retval interface{}
	retval = 42 //truth of life
	reply.Return = retval
	reply.Status = "FAIL"
	reply.Stdout = "TODO: stdout from keyword execution gets piped into this"
	reply.Stderr = "TODO: stderr from keyword execution gets piped into this"
	reply.Traceback = "TODO: stack trace info goes here, if any..."
	return nil
}

//the below functions & structs are optional and since not fully implemented,
//may be commented out if not wish to expose them to Robot Framework via gorrs as keywords

type KeywordInput struct {
	KeywordName string
}

type KeywordArgumentsReturnValue struct {
	KeywordAguments []interface{}
}

//sample XML-RPC input: <methodCall><methodName>GetKeywordArguments</methodName><params><param><value><string>KeywordName</string></value></param></params></methodCall>
//sample XML-RPC output: <methodResponse><params><param><value><array><data><value><string>arg1</string></value>...</data></array></value></param></params></methodResponse>
func (h *RobotRemoteService) GetKeywordArguments(r *http.Request, args *KeywordInput, reply *KeywordArgumentsReturnValue) error {
	//use reflection to get the arguments to keyword function and pass back to reply
	reply.KeywordAguments = make([]interface{}, 0) //if to pass back no arguments
	//http://stackoverflow.com/questions/12990338/cannot-convert-string-to-interface
	//else something like reply.KeywordAguments = append(reply.KeywordAguments,"two","arguments")
	return nil
}

type KeywordDocumentationReturnValue struct {
	KeywordDocumentation string
}

//sample XML-RPC input: <methodCall><methodName>GetKeywordDocumentation</methodName><params><param><value><string>KeywordName</string></value></param></params></methodCall>
//sample XML-RPC output: <methodResponse><params><param><value><string>godoc text</string></value></param></params></methodResponse>
func (h *RobotRemoteService) GetKeywordDocumentation(r *http.Request, args *KeywordInput, reply *KeywordDocumentationReturnValue) error {
	//makes a call shell call to godoc against the source code of the remote library package
	//or equivalent go package exported function (API) if there exists such for godoc
	//to then extract that go documentation for the keyword function and pass back to reply
	//extract off the documentation in source code, or the documentation web endpoint (http://localhost:6060 or http://golang.org if a standard go package)?
	//e.g. godoc -html -q package-name
	reply.KeywordDocumentation = "Unimplemented. TODO: keyword's go documentation goes here..."
	return nil
}
