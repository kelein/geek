package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group of SingleFlight
type Group struct {
	m map[string]*call
	l sync.Mutex
}

// Do of SingleFlight group
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.l.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.l.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)  // 发送请求前加锁
	g.m[key] = c // 将 call 添加到正在处理的请求中
	g.l.Unlock()

	c.val, c.err = fn() // 调用 fn 发送请求
	c.wg.Done()

	g.l.Lock()
	delete(g.m, key) // 更新已处理的 call
	g.l.Unlock()
	return c.val, c.err
}
