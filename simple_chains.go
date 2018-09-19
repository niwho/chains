package chains

import (
	"context"

	"github.com/niwho/logs"
)

type DataLoader interface {
	Get(key string) (interface{}, bool)
	Set(key string, val interface{})
}

type DataLoaderManager struct {
	*ChainManager
}

func NewDataLoaderManager() *DataLoaderManager {
	return &DataLoaderManager{
		ChainManager: NewChainManager(),
	}
}

func (dlm *DataLoaderManager) wrapper(dl DataLoader) ChainMethod {
	return func(c *ChainRequest) {
		key := c.GetKey()
		keystr, ok := key.(string)
		if !ok {
			logs.Log(logs.F{"key": key}).Error("DataLoaderManager:invalid key")
			c.Abort()
			return
		}
		if data, ok := dl.Get(keystr); ok {
			c.SetData(data)
			// 结束串行调用
			c.Abort()
			return
		}
		// 触发下一级调用
		c.Next()
		if dat, ok := c.GetData(); ok {
			dl.Set(keystr, dat)
		}
	}
}

func (dlm *DataLoaderManager) Use(dl DataLoader) *DataLoaderManager {
	dlm.AddChains(dlm.wrapper(dl))
	return dlm
}

func (dlm *DataLoaderManager) Get(key string) interface{} {
	return dlm.GetData(context.Background(), key)
}
