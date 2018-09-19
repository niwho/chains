package chains

import (
	"context"
	"fmt"
	"math"
	"sync"
)

type ChainManager struct {
	sync.Pool
	cacheChains GetDataChain
}

func NewChainManager() *ChainManager {
	cm := &ChainManager{}
	cm.New = func() interface{} {
		return &ChainRequest{}
	}
	return cm
}

//context.WithValue(ctx, userKey, u)
//u, ok := ctx.Value(userKey).(*User)
func (cm *ChainManager) GetData(ctx context.Context, key interface{}) interface{} {
	c := cm.Get().(*ChainRequest)
	c.reset()
	c.ctx = ctx
	c.key = key
	c.cacheChains = cm.cacheChains

	c.Next()

	dat, _ := c.GetData()
	cm.Put(c)

	return dat
}

// 注意按添加的顺序依次执行
func (cm *ChainManager) AddChains(chains ...ChainMethod) {
	cm.cacheChains = append(cm.cacheChains, chains...)
}

type ChainMethod func(*ChainRequest)
type GetDataChain []ChainMethod

const (
	// 有符号，所以是一半
	abortIndex int8 = math.MaxInt8 / 2
)

type ChainRequest struct {
	index       int8
	cacheChains GetDataChain

	key   interface{}
	data  interface{}
	isSet bool

	ctx context.Context
}

func (c *ChainRequest) reset() {
	c.index = -1
	c.cacheChains = nil
	c.key = ""
	c.data = ""
	c.isSet = false
}

func (c *ChainRequest) GetKey() interface{} {
	return c.key
}

func (c *ChainRequest) GetData() (interface{}, bool) {
	return c.data, c.isSet
}

func (c *ChainRequest) SetData(data interface{}) {
	c.data = data
	c.isSet = true
}

func (c *ChainRequest) Next() {
	c.index++
	s := int8(len(c.cacheChains))
	for ; c.index < s; c.index++ {
		fmt.Println("index~~~~~~", c.index)
		c.cacheChains[c.index](c)
	}
}

func (c *ChainRequest) Abort() {
	c.index = abortIndex
}
