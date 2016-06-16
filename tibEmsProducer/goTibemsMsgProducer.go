package main

/*
#cgo CFLAGS: -I../../include
#cgo LDFLAGS: -L../../lib/64 -L../../lib -ltibems64 -ltibemslookup64 -ltibemsufo64 -ltibemsadmin64 -lldap -llber -lxml2 -lssl -lcrypto -lz -lpthread -ldl
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "tibemsUtilities.h"
*/
import "C"
import "fmt"
import "os"


/*-----------------------------------------------------------------------
 * Parameters
 *----------------------------------------------------------------------*/

var serverUrl    	*C.char
var userName     	*C.char
var password     	*C.char
var pk_password  	*C.char
var name 					*C.char
var factory 			C.tibemsConnectionFactory
var connection		C.tibemsConnection
var session				C.tibemsSession
var msgProducer		C.tibemsMsgProducer
var destination		C.tibemsDestination
var queue		      C.tibemsQueue
var sslParams			C.tibemsSSLParams
var errorContext	C.tibemsErrorContext


func onCompletion(msg C.tibemsMsg,  status C.tibems_status) {
    var text *C.char

    if (C.tibemsTextMsg_GetText(msg, &text) != C.TIBEMS_OK) {
        fmt.Printf("Error retrieving message!\n");
        return;
    }

    if (status == C.TIBEMS_OK){
        fmt.Printf("Successfully sent message %s.\n", text);
    } else {
        fmt.Printf("Error sending message %s.\n", text);
        fmt.Printf("Error:  %s.\n", C.tibemsStatus_GetText(status));
    }
}

func fail(message string,  errContext C.tibemsErrorContext){
    var status C.tibems_status
    var str *C.char
    status=C.TIBEMS_OK

    fmt.Printf("ERROR: %s\n",message);

    status = C.tibemsErrorContext_GetLastErrorString(errContext, &str);
    fmt.Printf("\nLast error message =\n%s\n", C.GoString(str));
    status = C.tibemsErrorContext_GetLastErrorStackTrace(errContext, &str);
    fmt.Printf("\nStack trace = \n%s\n", C.GoString(str));
    _ = status
    os.Exit(1);
}


func main() {
 name := "test"
 serverUrl=C.CString("tcp://localhost:7222")
 userName=C.CString("admin")
 password=C.CString("admin")
 var status C.tibems_status
 var msg C.tibemsTextMsg
 status=  C.TIBEMS_OK
 fmt.Printf("Publishing to destination '%s'\n",name);
 status = C.tibemsErrorContext_Create(&errorContext);
 if (status != C.TIBEMS_OK){
     fmt.Printf("ErrorContext create failed: %s\n", C.tibemsStatus_GetText(status));
     os.Exit(1);
 }
 factory = C.tibemsConnectionFactory_Create();
 if (factory == nil) {
     fail("Error creating tibemsConnectionFactory", errorContext);
 }

 status = C.tibemsConnectionFactory_SetServerURL(factory,serverUrl);
 if (status != C.TIBEMS_OK) {
     fail("Error setting server url", errorContext);
 }
 status = C.tibemsConnectionFactory_CreateConnection(factory, &connection, userName, password);
 if (status != C.TIBEMS_OK) {
     fail("Error creating tibemsConnection", errorContext);
 }
 status = C.tibemsQueue_Create(&queue,C.CString(name));
 if (status != C.TIBEMS_OK) {
   fail("Error creating tibemsDestination", errorContext);
 }
 destination = (C.tibemsDestination)(queue)
 /* create the session */
 status = C.tibemsConnection_CreateSession(connection, &session,C.TIBEMS_FALSE,C.TIBEMS_AUTO_ACKNOWLEDGE);
 if (status != C.TIBEMS_OK){
     fail("Error creating tibemsSession", errorContext);
 }
 /* create the producer */
 status = C.tibemsSession_CreateProducer(session, &msgProducer,destination);
 if (status != C.TIBEMS_OK){
     fail("Error creating tibemsMsgProducer", errorContext);
 }
 status = C.tibemsTextMsg_Create(&msg);
 if (status != C.TIBEMS_OK){
     fail("Error creating tibemsTextMsg", errorContext);
 }
 /* set the message text */
 status = C.tibemsTextMsg_SetText(msg,C.CString("Go Test Message"));
 if (status != C.TIBEMS_OK){
     fail("Error setting tibemsTextMsg text", errorContext);
 }
 status = C.tibemsMsgProducer_Send(msgProducer,msg);
 if (status != C.TIBEMS_OK){
     fail("Error publishing tibemsTextMsg", errorContext);
 }

}
