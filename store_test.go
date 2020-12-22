package concurrent_store

import (
	"math/rand"
	"net/url"
	"testing"
	"time"
)

func TestConcurrentStore_Add(t *testing.T) {
	s := NewStore()
	v := "foo"
	s.Add(v)
	if !s.Contains(v) {
		t.Fatal("store doesn't contain v despite being added")
	}
}

func TestConcurrentStore_Contains(t *testing.T) {
	s := NewStore()
	a := "foo"
	b := "bar"

	s.Add(a)

	if !s.Contains(a) {
		t.Fatal("store doesn't contain a despite being added")
	}
	if s.Contains(b) {
		t.Fatal("store contains b despite not being added")
	}
}

func TestConcurrentStore_Pop(t *testing.T) {
	s := NewStore()
	v := "foo"
	s.Add(v)
	p, err := s.Pop()

	if err != nil {
		t.Fatal(err)
	}
	if p.(string) != v {
		t.Fatal("returned item not equal to entered item")
	}
	if s.Contains(p) {
		t.Fatal("popped item not remove from store")
	}
}

func TestConcurrentStore_All(t *testing.T) {
	s := NewStore()
	a := "foo"
	b := "bar"

	s.Add(a)
	s.Add(b)

	all := s.All()
	if _, ok := all[a]; !ok {
		t.Fatal("inserted item not included in s.All()")
	}
	if _, ok := all[b]; !ok {
		t.Fatal("inserted item not included in s.All()")
	}
}

func BenchmarkConcurrentStore_Contains(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	s := NewStore()
	for i := 0; i < 1000; i++ {
		u, _ := url.Parse("https://" + randStringRunes(4) + ".com")
		s.Add(u)
	}

	u, _ := url.Parse("https://" + randStringRunes(4) + ".com")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Contains(u)
	}
}

func BenchmarkConcurrentStore_Add(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	u, _ := url.Parse("https://" + randStringRunes(8) + ".com")
	s := NewStore()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(u)
	}
}

func BenchmarkConcurrentStore_Pop(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	s := NewStore()
	for i := 0; i < 1000; i++ {
		u, _ := url.Parse("https://" + randStringRunes(8) + ".com")
		s.Add(u)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Pop()
	}
}

func BenchmarkConcurrentStore_All(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	s := NewStore()
	for i := 0; i < 1000; i++ {
		u, _ := url.Parse("https://" + randStringRunes(8) + ".com")
		s.Add(u)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.All()
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
