# go-demo-grpc


### HTTP/2

- 流（stream）：在一个已建立的连接上的双向字节流。一个流可以携带一条或多条消息。
- 帧（frame）：HTTP/2 中最小的通信单元。每一帧都包含一个帧头，它至少要标记该帧所属的流。
- 消息（message）：完整的帧序列，映射为一条逻辑上的 HTTP 消息，由一帧或多帧组成。这样的话，允许消息进行多路复用，客户端和服务器端能够将消息分解成独立的帧，交叉发送它们，然后在另一端进行重新组合。



