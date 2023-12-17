# Test Task

Create a straightforward TCP client that delivers incremental counters to the given server at the
given message rate.

## Protocol

*NOTE: rate = number of messages per second*

* The client must send packets only at the specified rate.
* The service periodically sends the new rate in JSON. 
    Example:
    `{"rate": 10}` means that the client should send 10 messages per second until a new rate is
    broadcasted.
* The client must send incremental int64 counter starting with `1` encoded as a sequence of bytes
* (varint).
