package cache

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestAdd(t *testing.T) {
	lru := New(int64(0), nil)

	lru.Add("key1", String("1"))
	if lru.currBytes != int64(len("key1")+len("1")) {
		t.Fatal("expected 4 but got", lru.currBytes)
	}

	if v, ok := lru.Get("key1"); string(v.(String)) != "1" || !ok {
		t.Fatal("cache hit key1=1 failed")
	}

	lru.Add("key1", String("111"))
	if lru.currBytes != int64(len("key1")+len("111")) {
		t.Fatal("expected 6 but got", lru.currBytes)
	}

	if v, ok := lru.Get("key1"); string(v.(String)) != "111" || !ok {
		t.Fatal("cache hit key1=111 failed")
	}

	lru.Add("key2", String("222"))
	if lru.currBytes != int64(len("key1")+len("111")+len("key2")+len("222")) {
		t.Fatal("expected 6 but got", lru.currBytes)
	}

	if v, ok := lru.Get("key2"); string(v.(String)) != "222" || !ok {
		t.Fatal("cache hit key2=222 failed")
	}
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if v, ok := lru.Get("key2"); !ok || string(v.(String)) != "value2" {
		t.Fatal("Fail to store key2 and value2")
	}

	if v, ok := lru.Get("k3"); !ok || string(v.(String)) != "v3" {
		t.Fatal("Fail to store k3 and v3")
	}

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Remove oldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
