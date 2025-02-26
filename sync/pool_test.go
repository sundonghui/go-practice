package sync

import (
	"bytes"
	"sync"

	"github.com/google/uuid"
)

var bufferPool = sync.Pool{
	New: func() interface{} { return bytes.NewBuffer(nil) },
}

func process(data []byte) {
	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)

	// 错误假设：buf 的底层数组是空的
	buf.Write(data) // 可能将数据追加到已有内容之后！
}

type Session struct {
	ID string
}

var pool = sync.Pool{
	New: func() interface{} { return &Session{ID: uuid.NewString()} },
}
var cache = make(map[string]*Session)

// 错误用法：依赖对象指针的“唯一性”
func processRequest() {
	session := pool.Get().(*Session)
	defer pool.Put(session)

	// 假设 session 的指针是唯一的身份标识
	cache[session.ID] = session // 危险！session.ID 可能被覆盖

	// 后续操作假设 session 的 ID 是唯一的
}
