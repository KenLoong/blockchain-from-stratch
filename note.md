# EP8

- gob: type elliptic.p256Curve has no exported fields"错误
- 解决：为 PublicKey 结构体实现 GobEncode 和 GobDecode 接口，将密钥用我们的逻辑去序列化，避免 gob 直接序列化 ecdsa.PublicKey

# EP18

- server 的 Start 函数逻辑有问题，需要修改，应该监听 peers 的消息才对
- initRemoteServers 函数的逻辑也有问题，err 发生后没有停止处理下文

# EP19

- initRemoteServers 函数已经被删除
- todo:local_transport_test 的单测需要全部重写
- 目前 REMOTE_NODE_B 尚未被同步到，有 bug，还没排查出来
- REMOTE_NODE_B 的问题已经解决，因为 5000 端口早被占用，将 B 服务器的端口改成 4001 就行了

# EP20

- err := s.TCPTransport.Start() ; // 修复未感知端口已被占用的漏洞 （我自己修复的，作者自己没有意识到这个问题）
- 还是有 bugs：
  - 因为目前 sendGetStatusMessage 这个函数只会被调用一次，如果第一次同步 blocks 的时候没有同步到最新的 block，那么这个 server 是永远也没有第二次机会追上其他 server 了. sendGetStatusMessage 这个函数的角色类似心跳机制，应该定时调用（看作者在下一集中是否会意识到这个问题）
  - processBlocksMessage 这个函数有问题，因为目前的逻辑是全量同步 blocks，所以 0 这个 block 在被同步时，一定会有 ErrBlockKnown 这个错误，需要忽略这个错误，我已经修复这个问题了(看作者是否会意识到这个问题了)

# EP21

- sendGetStatusMessage 这个函数还是调用一次，但是在处理 GetStatusMessage 时，会开启定时调用 requestBlocksLoop 的逻辑，不停地请求 blocks，也算是不断同步了
- 作者依然没有处理 ErrBlockKnown 的同步问题，但是将 return 改成 continue 了
- 作者修改了 public key 的定义，但是我的实现应该也没有问题，就不跟作者改了

# EP22

- 使用 echo 框架来写了两个 http api 接口
- 项目运行时间较长时，会出现这个错误：peer send error => addr 127.0.0.1:59723 [err: write tcp 127.0.0.1:3000->127.0.0.1:59723: write: broken pipe]

# EP23

- 实现 http 发送 tx
- 目前只有 validator server 才会清理 mempool,其它 server 的 mempool 只会不断膨胀

# EP24

- 实现两种类型 tx 的执行
- 将 gob 换成 json，gob 太多问题了，很复杂
- 现在 late_node 无法同步了，出现 failed to decode message
- Transaction 需要自定义 json 编解码，不知道为何会这么复杂，烦

# fix

- 解决了 failed to decode message 的问题，原因是 TCP 读取消息时的 buf 直接写成 4096，导致消息大于 4096 时，被截断，然后 json 解码失败

# EP25

- 实现基本的账户转账，但是因为 blockchain 实体的 accountState 变量没有在各节点同步，所以只有 LOCAL_NODE_3000 节点能够成功执行 transfer 的交易（LOCAL_NODE_3000 节点初始化时，在 accountState 里放进了自己的 address），其它节点都会失败，出现 address (28621e74135fac8140a7ca4c5d9da32aeba442de) unknown"的失败

# EP26

- 博主的 addBlockWithoutValidation 函数有死锁漏洞，我改成 execTransactions 来避免死锁了

# EP27

- 博主的单测有问题，block 在加完 tx 后，没有 sign 竟然也能成功放进 blockchain，有问题,但我的代码版本可以保证必须 sign 才能放进 blockchain
