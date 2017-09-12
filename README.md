# ptp-dfs

## To-do
* Understand Kademlia & the lab.
* Setup routing table for the existing code provided in the lab.

## Messages
* There will be 4 messages of type ping, lookupContact, lookupData, store.
* Each message will have their own protobuf template.
* Each message will have an ID that matches them to the correct goRoutine.
* All messages will have "isReply" bool to flag if the message is
a request or a response.
* All messages will have an optional response data containing the response of the
message, if any.

## Go-routine structure
### Worker-Queue
* Have a constant goroutine (Listen) that recieves messages and
adds them to a work queue.
* Have a dispatcher that take messages from the work queue and go through a
worker queue to find the worker that the message belong to.
* When making a new request, register a new channel for the request 
and add it to the worker queue to get notified when recieving the correct respone.


## Questions for lab supervisors
*(done) How should we format protobuf messages? Should there be only 1 generic
message? or should we format one message for each method? (Ping,findContact,SendData etc..)
*(done) Should we use the same protobuf messages for sending and recieving?
*(done) For lookupContact, should we query the closest relevant contacts, then get a response with new contacts and
then query those contacts if they have not already appeared etc. Is this the recursive
way you had in mind?

## Useful links and docs
#### Worker Queue
http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html
