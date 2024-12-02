# EP18

- server 的 Start 函数逻辑有问题，需要修改，应该监听 peers 的消息才对
- initRemoteServers 函数的逻辑也有问题，err 发生后没有停止处理下文

# EP19

- initRemoteServers 函数已经被删除
- todo:local_transport_test 的单测需要全部重写
- 目前 REMOTE_NODE_B 尚未被同步到，有 bug，还没排查出来
- REMOTE_NODE_B 的问题已经解决，因为 5000 端口早被占用，将 B 服务器的端口改成 4001 就行了
