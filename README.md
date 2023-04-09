# go-demo-grpc


### HTTP/2

- 流（stream）：在一个已建立的连接上的双向字节流。一个流可以携带一条或多条消息。
- 帧（frame）：HTTP/2 中最小的通信单元。每一帧都包含一个帧头，它至少要标记该帧所属的流。
- 消息（message）：完整的帧序列，映射为一条逻辑上的 HTTP 消息，由一帧或多帧组成。这样的话，允许消息进行多路复用，客户端和服务器端能够将消息分解成独立的帧，交叉发送它们，然后在另一端进行重新组合。


### gRPC 错误码

错误码 | 数字 | 描述
--- | --- | ---
OK | 0 | 成功状态
CANCELLED | 1 | 操作已被（调用者）取消
UNKNOWN | 2 | 未知错误
INVALID_ARGUMENT | 3 | 客户端指定了非法参数
DEADLINE_EXCEEDED | 4 | 在操作完成前，就已超过了截止时间
NOT_FOUND | 5 | 某些请求实体没有找到
ALREADY_EXISTS | 6 | 客户端试图创建的实体已存在
PERMISSION_DENIED | 7 | 调用者没有权限执行特定的操作
RESOURCE_EXHAUSTED | 8 | 某些资源已被耗尽
FAILED_PRECONDITION | 9 | 操作被拒绝，系统没有处于执行操作所需的状态
ABORTED | 10|  操作被中止
OUT_OF_RANGE | 11 | 尝试进行的操作超出了合法的范围
UNIMPLEMENTED | 12 | 在该服务中，未实现或不支持（未启用）本操作
INTERNAL | 13 | 内部错误
UNAVAILABLE | 14 | 该服务当前不可用
DATA_LOSS | 15 | 不可恢复的数据丢失或损坏
UNAUTHENTICATED | 16 | 客户端没有进行操作的合法认证凭证

