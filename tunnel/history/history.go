package history

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/Dreamacro/clash/win"
	"golang.org/x/exp/slices"
	"net/http"
	"sync"
	"time"
)

var APIAddress string
var cache = make(map[string]time.Time, 100)
var cacheLock = sync.Mutex{}

func Add(proxy constant.Proxy, metadata *constant.Metadata) {
	if proxy.Type() == constant.Reject {
		return
	}

	// 过滤后台服务域名
	if len(metadata.Host) > 0 && !slices.Contains(constant.ServerAPIDomains, metadata.Host) {
		if _, ok := cache[metadata.Host]; !ok {
			now := time.Now().UTC()
			cacheLock.Lock()
			defer cacheLock.Unlock()
			cache[metadata.Host] = now
		}
	}
}

func init() {
	go func() {
		for {
			uploadHistory()
			time.Sleep(time.Minute * 5)
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
	log.Infoln("cache len %d:", len(cache))
	err := upload()
	if err != nil {
		log.Errorln("upload record err: %v:", err)
	}

	cache = make(map[string]time.Time, 100)
}

// Record 上网行为记录
type Record struct {
	MachineId     string    `json:"machine_id"`
	RecordType    int       `json:"record_type"`
	TriggerTime   time.Time `json:"trigger_time"`
	BehaviorValue string    `json:"behavior_value"` // 网站访问url
}

const RecordType = 1 // 上网记录上报

func upload() error {
	records := make([]Record, len(cache))
	i := 0
	for url, time := range cache {
		record := Record{
			MachineId:     win.MachineId,
			RecordType:    RecordType,
			TriggerTime:   time,
			BehaviorValue: url,
		}

		records[i] = record
		i++
	}

	record := make(map[string]interface{})
	record["machine_id"] = win.MachineId
	record["records"] = records
	bytesData, _ := json.Marshal(record)

	_, err := http.Post(
		fmt.Sprintf("%v%v", APIAddress, constant.UploadUrl),
		"application/json;charset=utf-8",
		bytes.NewBuffer(bytesData),
	)

	return err
}
