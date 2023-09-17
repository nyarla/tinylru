package tinylru

import (
	"testing"
)

func TestCache(t *testing.T) {
	c := new(Cache[int, int]).Init(5)

	c.Put(1, 10)
	if c.first.key != 1 || c.first.value != 10 {
		t.Fatalf(`failed to put value: %d => %d`, c.first.key, c.first.value)
	}

	c.Put(2, 20)
	if c.first.key != 2 || c.first.value != 20 {
		t.Fatalf(`failed to put value: %d => %d`, c.first.key, c.first.value)
	}

	c.Put(3, 30)
	if c.first.key != 3 || c.first.value != 30 {
		t.Fatalf(`failed to put value: %d => %d`, c.first.key, c.first.value)
	}

	c.Get(2)
	if c.first.key != 2 || c.first.value != 20 {
		t.Fatalf(`failed to move value: %d, %d != 2, 20`, c.first.key, c.first.value)
	}

	c.Put(4, 40)
	if c.first.key != 4 || c.first.value != 40 {
		t.Fatalf(`failed to put value: %d => %d`, c.first.key, c.first.value)
	}

	c.Put(5, 50)
	if c.first.key != 5 || c.first.value != 50 {
		t.Fatalf(`failed to put value: %d => %d`, c.first.key, c.first.value)
	}

	c.Put(6, 60)
	if c.first.key != 6 || c.first.value != 60 {
		t.Fatalf(`failed to put value: %d => %d`, c.first.key, c.first.value)
	}

	if c.last.key != 3 || c.last.value != 30 {
		t.Fatalf(`failed to move value: %d, %d != 3, 30`, c.last.key, c.last.value)
	}

	c.Get(3)
	if c.last.key != 2 || c.last.value != 20 {
		t.Fatalf(`failed to move value: %d, %d != 3, 30`, c.last.key, c.last.value)
	}
}

func BenchmarkCache(b *testing.B) {
	c := new(Cache[int, int]).Init(128)

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		c.Put(idx, idx)
	}
}
