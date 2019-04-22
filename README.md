# event-sourcing

## Concept

### Publisher

Publisher are responsible for messaging using Amazon SNS.

### Subscriber

Subscribers are responsible for polling for queuing.

### Butler

Butler does chores such as creation, destruction of SNS/SQS.  
FYI: This has nothing to do with the concept of the Event sourcing architecture. 
