# golang 笔试题答卷

## 1. 简述异步与同步 io 的区别与实际运用中的注意事项
#### 区别：
- 同步IO：调用线程会阻塞直到操作完成（如read()必须等到数据就绪）
- 异步IO：调用立即返回，内核通过回调/事件通知结果（如epoll+非阻塞文件描述符）
#### 注意事项：
- 资源泄漏：异步IO需要显式关闭文件描述符
- 回调地狱：使用协程+channel结构化异步代码
- 缓冲区管理：注意数据生命周期（避免异步读写时缓冲区被回收）
- 错误处理：异步场景下错误可能延迟出现
#### 实际应用
- 同步IO更适合简单逻辑开发（代码直观）
- 异步IO适合高并发场景
## 2. 简述在多线程情景下如何在线程（协程）之间同步数据
#### 核心方法：
- Channel（首选方案）：make(chan T)+select
- sync包：Mutex/RWMutex+WaitGroup
- 原子操作：atomic.Value/sync/atomic
#### 注意事项
- Channel 关闭后写入会 panic（使用recover()保护）
- Mutex未释放会导致死锁（建议用defer mu.Unlock()）
- 共享内存尽量限制在小作用域（如结构体内部）  
## 3. 简述 golang 语言中 channel 的作用和利弊
#### 作用
- 协程间通信管道
- 替代锁实现CSP并发模型
- 事件通知机制
#### 优点
- 自带同步机制（无需额外锁）
- 类型安全（编译时检查）
- 支持多路复用（select）
#### 缺点
- 过度使用会增加GC压力（大缓冲channel）
- 错误使用易导致死锁（如无缓冲channel未配对使用）
#### 注意事项
- 优先使用无缓冲channel实现精准控制
- close后需要处理接收的零值（配合ok语法）
- nil channel会永久阻塞（初始化必须make）
## 4. 设计一种在您的理解下最小消耗的内存缓冲
#### 设计要点
- 对象池：使用sync.Pool复用内存
- 内存对齐：避免false sharing
- 分块管理：按不同size分级管理
#### 代码示例
```go
// 按需选择合适尺寸的buffer
type BufferPool struct {
    pools [4]*sync.Pool
}

func NewBufferPool() *BufferPool {
    return &BufferPool{
        pools: [4]*sync.Pool{
            {New: func() any { return make([]byte, 128) }},
            {New: func() any { return make([]byte, 512) }},
            {New: func() any { return make([]byte, 2048) }},
            {New: func() any { return make([]byte, 8192) }},
        }
    }
}
```
#### 注意事项
- 池中的对象具有不可预测性（不能假设对象始终存在）
- 大对象不要放回池中（影响GC效率）
- 清空复用前的数据（避免脏数据残留）
## 5. 简述网络编程中，在调用 socket 或您使用过的网络框架下 send/write 的具体含义
#### 本质理解：
- 系统调用：将用户态缓冲区数据复制到内核协议栈
- 非阻塞模式下可能只发送部分数据（需检查返回值）
- TCP层会自动处理分包/组包
#### Golang实现细节：
- net.Conn的Write内部使用非阻塞IO
- 实际通过runtime调度实现异步
- 连接级写缓冲区管理（每个conn有自己的写缓冲）
#### 注意事项：
- 必须检查返回值（实际写入字节数）
- 大数据发送要循环调用Write（直到全部发送）
- 配合SetWriteDeadline避免阻塞
## 6. 设计一个 http/websocket 框架，写出简要思路或伪代码
#### 架构要点：
- 路由层：前缀树路由+中间件链
- 协议升级：兼容HTTP/WebSocket
- 连接管理：心跳机制+优雅关闭
#### 伪代码示例
```go
type Context struct {
    Writer  http.ResponseWriter
    Request *http.Request
    // 中间件数据传递...
}

type HandlerFunc func(*Context)

// WebSocket升级逻辑
func UpgradeWS(c *Context) {
    conn, err := websocket.Upgrade(c.Writer, c.Request...)
    go handleWSConnection(conn)
}

// 主处理逻辑
func handleWSConnection(conn *websocket.Conn) {
    defer conn.Close()
    for {
        msgType, msg, err := conn.ReadMessage()
        // 使用channel分发消息
        select {
        case writeChan <- msg:
        case <-closeChan:
            return
        }
    }
}
```
#### 注意事项：
- 限制单个连接内存占用（防止内存耗尽）
- 使用sync.Pool复用协议解析器
- 监控连接状态（如使用prometheus统计）
## 7. 将一百万份 10kb 的小文件合并成一个文件，不考虑顺序/错误情况，只要求最高效率，请给出伪代码
#### 优化要点：
- 并行读取：worker池并发处理
- 批量写入：合并小IO为大块写入
- 零拷贝：使用sendfile系统调用（Linux）
#### 伪代码示例：
```go
func mergeFiles(output string, files []string) {
    dst, _ := os.Create(output)
    defer dst.Close()

    sem := make(chan struct{}, runtime.NumCPU()*2) // 控制并发度
    var wg sync.WaitGroup

    for _, f := range files {
        wg.Add(1)
        go func(path string) {
            defer wg.Done()
            sem <- struct{}{}
            data, _ := os.ReadFile(path)
            
            // 批量写入锁
            mu.Lock()
            dst.Write(data)
            mu.Unlock()
            
            <-sem
        }(f)
    }
    wg.Wait()
}
```  
#### 注意事项：
- 文件描述符限制（ulimit -n调大）
- 磁盘IO模式（SSD vs HDD策略不同）
- 内存压力控制（使用流式处理替代全加载）