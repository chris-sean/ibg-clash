package history

import (
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"sync"
	"time"
)

var cache = make(map[string]struct{}, 100)
var cacheLock = sync.Mutex{}

func Add(proxy constant.Proxy, metadata *constant.Metadata) {
	if proxy.Type() == constant.Reject {
		return
	}

	cacheLock.Lock()
	defer cacheLock.Unlock()
	cache[metadata.Host] = struct{}{}
}

func init() {
	go func() {
		for {
			uploadHistory()
			time.Sleep(time.Minute)
		}
	}()
}

func uploadHistory() {
	cacheLock.Lock()
	defer func() {
		cacheLock.Unlock()
		if err := recover(); err != nil {
			log.Warnln("[uploadHistory] failed: %v", err)
		}
	}()

	if len(cache) == 0 {
		return
	}

	// do upload

	cache = make(map[string]struct{}, 100)
}
