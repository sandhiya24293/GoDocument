GRPC PROTOCAL :

link : https://container-solutions.com/introduction-to-grpc/

What is gRPC?
gRPC is a RPC platform developed by Google which was announced and made open source  in late Feb 2015.  
The letters “gRPC” are a recursive acronym which means, gRPC Remote Procedure Call.

gRPC has two parts, the gRPC protocol, and the data serialization. By default gRPC utilizes Protobuf for serialization, but it is pluggable with any form of serialization you wish to use, with some caveats, which I will get to later.

The Protocol
The protocol itself is based on http2, and exploits many of its benefits. Google’s development of SPDY laid down the first designs for  what eventually became http2.

gRPC supports several built in features inherited from http2, such as compressing headers, persistent single TCP connections, cancellation and timeout contracts between client and server.

The protocol has built in flow control from http2 on data frames.  This is very handy for insuring clients respect the throughput of your system, but does add an extra level of complexity when diagnosing issues in your infrastructure, because either client or server can set their own flow control values.

Load balancing (LB) is normally performed by the client, which chooses the server for a given request from a list provided by a Load Balancing server. The LB server will monitor the health of endpoints and use this and other factors to manage the list provided to clients. Clients will use a simple algorithm such as round-robin internally, but note that the LB server may apply more complex logic when compiling the list for a given client.

RPC Types
gRPC offers two essential types for client server communication.

Unary
Essentially these are synchronous requests made to the gRPC server with a single request  that blocks until a response is received.

Streaming
Streaming is really powerful and can be accomplished in three different configurations:  client pushing messages to a stream; server pushing messages to a stream; or bidirectional, where client and server are both sending data in two streams in the same method.  In all cases the client initiates the RPC method.
Streams don’t provide any acknowledgement of receipt until the stream completes, which can add complexity when the system needs to cope with node failures or network partitions.  This can be mitigated by using a bidirectional stream to return ACKs.  If a server is given a chance to kill a connection gracefully a message will be returned indicating the last received message.

Protobuf
Protobuf is the default serialization format for the data sent between clients and servers.  The encoding allows for small messages and quick decoding and encoding.
Protobuf foregos zero-copy of data like some other data interchange methods (such as Cap’n Proto or Flatbuffers), instead opting for encoding and decoding bytes.  This makes the data smaller at the cost of having to devote CPU to encoding and decoding messages.  Unlike other serialization formats like JSON or XML, protobuf tries to minimize the overhead of encoding by providing strongly-typed fields in an encoded binary format that it can quickly traverse in a predictable manner.

The protobuf file
Protobuf defines how to interpret messages and allows the developer to create stubs that make encoding and decoding these values quick and efficient.  In the below example we only have two fields – name and age.


message Foo {
    string name = 1;
    int32 age = 2;
}
1
2
3
4
message Foo {
    string name = 1;
    int32 age = 2;
}
Each field has a type declaration and field number, the names are only for us mere mortals.  Field numbers are important, they stick with the field and they ensure backward compatibility if someone is using older stubs.  You can always add and remove fields but you should never number the same field as a previous field that was removed, if you have already released stubs.  If you remove a field, you can lock it down to prevent accidental reuse by adding it as reserved.  So if I was to remove age I would change the definition to:


message Foo {
    reserved 2;
    string name = 1;
}
1
2
3
4
message Foo {
    reserved 2;
    string name = 1;
}
While reviewing the protobuf docs you may see that it supports required and optional designators for fields.  This has been removed in protobuf version 3.

Encoding
Before we delve into the mechanizations of encoding and decoding data, I want to cover some behaviors you should be aware of.  Every protobuf encoder/decoder should be able to skip fields it does not know and set defaults for fields it cannot find.  Even though all encoders should write fields in their number order, all decoders should anticipate fields out of order and if duplicates of a field are found they are added or concatenated or merged, all depending on the field definition.

All fields start with the field number, followed by the wire type (which determines how the messages will be decoded) all followed by the actual data contained within.  The decoding strategies would change depending on the field type.

So in a binary stream it would look like this.

Field number – Wire Type – Data

The field number is a Varint, and the wire type only occupies the last three bits after the field number.  Depending on the wire type, other field information may also be included after the field label,  a wire type such as strings will have another varint defining the length of the string.

Every field in protobuf will follow this scheme.  Using this method the encoder can quickly traverse the fields, and get the needed  information to either decode or skip the field quickly.  Since the field carries the wire type in the first byte, as soon as the decoder determines this isn’t the field it is looking for, everything needed to skip the field is known and easily moved forward to the next field.

Protobuf and gRPC
Protobuf and gRPC are completely independent of each other and even though there don’t appear to be any encoding methods adapted to gRPC, in theory you could change out any encoding method.

That being said all the automatic client server side code generation is performed by protoc the protobuf stub generation tool.  So you can change out your encoding but you would lose the ability to generate server and client code stubs for ten different languages.

Stubbing and backward compatibility
Protobuf attempts to put safeguards into its encoding scheme.  In theory, stubs generated by protobuf should be backwards compatible.  I think this compatibility requirement is one of the reasons the removed “required” and “optional” fields.

Some people will still balk at sharing stubs across their microservices and even though protobuf added conventions to ensure backwards compatibility for services that are not updated, that does not mean developers will follow them.  There will be developers that change field name 3 from a uint32 to a string, breaking your application.

Being Cloud Native
The Cloud Native Computing Foundation charter states their goal with microservices to be:

Significantly increase the overall agility and maintainability of applications. The foundation will shape the evolution of the technology to advance the state of the art for application management, and to make the technology ubiquitous and easily available through reliable interfaces.

The chosen language for interservice communication is one of the key components of a microservice architecture. When you pick such a language it should meet some core goals such as; quick to develop, resilient to a changing environment, and be performant on the wire as well as in your application. gRPC seems to meet these goals, which leads to an architecture that is more agile and maintainable, so in that sense it appears to be a strategic choice for being Cloud Native.

Conclusion
Now you are wondering “Should I adopt gRPC?”

You want to be sure any new architecture is vetted and has strong backing, which this project does. It has also been adopted in production environments with companies like Square, Netflix and CoreOS.

You also want to be sure your team can onboard with the protocol quickly.  With the code generation this should be straightforward, but sometimes engineers will need to understand the technical details of http2, which may be more of a burden. Thebackwards-compatibility and code-generation features built into protobuf are a boon to those working in large teams. Pursuing premature optimization can be a pitfall, but gRPC is very robust in its current state, and the speed is only one aspect of what gRPC brings to the table.

There are limitations, such as the fact that there are no front end implementations of the protocol, so you may need to transform the data if it is to be consumed by a front end.  Tools are currently being developed that may ease this process.

So yes, gRPC is not only important, but it should also be very high on the laundry list of things you should learn.  It is no small matter that the CNCF has accepted into their list of hosted projects.
