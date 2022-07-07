package concurrent

import (
	"sync"
)

//SyncSlice仅支持并发写数据，等所有协程写完后再统一读取数据，例子见Test_SyncSlice
type SyncSlice struct {
	items interface{} //interface约定是[]interface{}形式
	ml    *sync.RWMutex
}

func NewSyncSlice() *SyncSlice {
	return &SyncSlice{ml: new(sync.RWMutex), items: []interface{}{}}
}

func (syncSlice *SyncSlice) Append(items interface{}) { //放入的必须是切片类型
	syncSlice.ml.Lock()
	defer syncSlice.ml.Unlock()
	syncSlice.items = append(syncSlice.items.([]interface{}), items)
}

//获取结果数据
func (syncSlice *SyncSlice) GetSlice() interface{} {
	return syncSlice.items
}
