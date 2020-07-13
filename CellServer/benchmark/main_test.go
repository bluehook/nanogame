package benchmark

import (
	"sync"
	"testing"
)

type TestData struct {
	i int
	f float64
	s string
	d int64
	l []string
}

func BenchmarkNew(b *testing.B) {
	tmp := map[string]*TestData{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b := &TestData{
			l: []string{},
		}
		tmp["abc"] = b
	}
}

func BenchmarkPool(b *testing.B) {
	pool := &sync.Pool{
		New: func() interface{} {
			return &TestData{
				l: []string{},
			}
		},
	}

	init := pool.Get()
	pool.Put(init)

	tmp := map[string]*TestData{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := pool.Get().(*TestData)
		tmp["abc"] = x
		pool.Put(x)
	}
}
