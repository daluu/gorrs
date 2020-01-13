package main

import (
	"log"
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
	log.Fatal("not runnable directly")
}
