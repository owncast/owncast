package cache

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpiration(t *testing.T) {
	ttl := 500 * time.Millisecond
	c := New()
	defer c.Stop()

	inst := c.GetOrCreate("expiring", ttl)
	inst.Set("k", []byte("v"))

	if got := inst.Get("k"); string(got) != "v" {
		t.Fatalf("immediate read: got %q want %q", got, "v")
	}

	time.Sleep(ttl + 250*time.Millisecond)
	if got := inst.Get("k"); got != nil {
		t.Errorf("after TTL: expected nil, got %q", got)
	}
}

func TestGetOrCreateReusesInstance(t *testing.T) {
	c := New()
	defer c.Stop()

	a := c.GetOrCreate("named", time.Minute)
	b := c.GetOrCreate("named", time.Minute)
	assert.Same(t, a, b, "GetOrCreate must return the same Instance for the same name")
}

func TestMissingKeyReturnsNil(t *testing.T) {
	c := New()
	defer c.Stop()

	inst := c.GetOrCreate("empty", time.Minute)
	assert.Nil(t, inst.Get("never-set"))
}

func TestConcurrentSetAndGet(t *testing.T) {
	c := New()
	defer c.Stop()

	inst := c.GetOrCreate("concurrent", 5*time.Second)
	const goroutines = 10

	done := make(chan struct{})
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer func() { done <- struct{}{} }()

			key := "key" + strconv.Itoa(idx)
			val := "val" + strconv.Itoa(idx)
			inst.Set(key, []byte(val))
			time.Sleep(50 * time.Millisecond)
			assert.Equal(t, val, string(inst.Get(key)))
		}(i)
	}
	for i := 0; i < goroutines; i++ {
		<-done
	}
}

func TestStopIsIdempotent(t *testing.T) {
	c := New()
	c.GetOrCreate("a", time.Minute)
	c.GetOrCreate("b", time.Minute)
	c.Stop()
	// Calling Stop on an empty container after Stop should not panic.
	c.Stop()
}
