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
