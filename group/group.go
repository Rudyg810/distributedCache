package group

import (
	"distributed_cache/cache"
	"distributed_cache/peer"
	"fmt"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache *cache.SafeCache
	peers     peer.PeerPicker
	mu        sync.RWMutex
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache.NewSafeCache(cacheBytes),
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (cache.ByteView, error) {
	if key == "" {
		return cache.ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.Get(key); ok {
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value cache.ByteView, err error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if value, err = g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
		}
	}

	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (cache.ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return cache.ByteView{}, err
	}
	value := cache.NewByteView(bytes)
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value cache.ByteView) {
	g.mainCache.Add(key, value)
}

func (g *Group) RegisterPeers(peers peer.PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) getFromPeer(peer peer.PeerGetter, key string) (cache.ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return cache.ByteView{}, err
	}
	return cache.NewByteView(bytes), nil
}
