# EP18

- server 的 Start 函数逻辑有问题，需要修改，应该监听 peers 的消息才对
- boostrapNodes 函数的逻辑也有问题，err 发生后没有停止处理下文
