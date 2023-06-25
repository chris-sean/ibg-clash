package history

import (
	"github.com/Dreamacro/clash/constant"
	"sync"
	"time"
)

var cache = map[string]struct{}{}
var cacheLock = sync.Mutex{}

func Add(proxy constant.Proxy, metadata *constant.Metadata) {
	if proxy.Type() == constant.Reject {
		return
	}

	cacheLock.Lock()
	defer cacheLock.Unlock()
	cache[metadata.Host] = struct{}{}
	//fmt.Printf("HISTORY %#v\n", metadata)
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
	defer cacheLock.Unlock()

	if len(cache) == 0 {
		return
	}

}
