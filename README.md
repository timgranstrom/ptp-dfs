# ptp-dfs

##To-do
* Understand Kademlia & the lab.
* Setup routing table for the existing code provided in the lab.


##Questions for lab supervisors
* How should we format protobuf messages? Should there be only 1 generic
message? or should we format one message for each method? (Ping,findContact,SendData etc..)
* Should we use the same protobuf messages for sending and recieving?
* For lookupContact, should we query the closest relevant contacts, then get a response with new contacts and
then query those contacts if they have not already appeared etc. Is this the recursive
way you had in mind?


##Go-routine structure
* Have one constant goroutine that listens for requests. //listen on separate port?
* Every time you send a request, open up a new goroutine to send.
the request, and keep it open while listening for that response. //listen on separate port for each request
*In what way do we want to use channels? What data to share?