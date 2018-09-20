package chains

import (
	"context"

	"github.com/niwho/logs"
)

type DataLoader interface {
	Get(key string) (interface{}, bool)
	Set(key string, val interface{})
	Del(key string)
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
		mode := c.ctx.Value(MODE).(string)
		if mode == GET {
			if data, ok := dl.Get(keystr); ok {
				c.SetData(data)
				// 结束串行调用
				c.Abort()
				return
			}
		}
		// 触发下一级调用
		c.Next()
		switch mode {
		case GET, SET:
			if dat, ok := c.GetData(); ok {
				dl.Set(keystr, dat)
			}
		case DEL:
			dl.Del(keystr)
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

func (dlm *DataLoaderManager) Set(key string, val interface{}) {
	dlm.SetData(context.Background(), key, val)
}

func (dlm *DataLoaderManager) Del(key string) {
	dlm.DelData(context.Background(), key)
}
