package ptp
/**
Message queue should hold it's own collector, dispatcher and workers.
Dispatcher should have a LOCAL workerQueue, NOT a GLOBAL one.
This way we can run multiple instances of kademlia without mixing up message queues.
 */
type MessageQueue struct{

}
