package chains

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestChainManager(t *testing.T) {
	cha := NewChainManager()
	a1 := func(c *ChainRequest) {
		//own work
		if false {
			c.SetData("a1")
			fmt.Println("a1 set")
			c.Abort()
			return
		}
		c.Next()
		// 从下一级获取到结果
		if dat, ok := c.GetData(); ok {
			//do some work
			fmt.Println("a1 next", dat)

		}
		fmt.Println("a1 next not ok")
	}

	b1 := func(c *ChainRequest) {
		//own work
		if false {
			c.SetData("b1")
			fmt.Println("b1 set")
			c.Abort()
			return
		}
		c.Next()
		// 从下一级获取到结果
		if dat, ok := c.GetData(); ok {
			//do some work
			fmt.Println("b1 next", dat)

		}
		fmt.Println("b1 next not ok")
	}
	b2 := func(c *ChainRequest) {
		//own work
		if false {
			c.SetData("b2")
			fmt.Println("b2 set")
			c.Abort()
			return
		}
		c.Next()
		// 从下一级获取到结果
		if dat, ok := c.GetData(); ok {
			//do some work
			fmt.Println("b2 next", dat)

		}
		fmt.Println("b2 next not ok")
	}
	b3 := func(c *ChainRequest) {
		//own work
		if false {
			c.SetData("b3")
			fmt.Println("b3 set")
			c.Abort()
			return
		}
		c.Next()
		// 从下一级获取到结果
		if dat, ok := c.GetData(); ok {
			//do some work
			fmt.Println("b3 next", dat)

		}
		fmt.Println("b3 next not ok")
	}
	key := "no meaning"
	cha.AddChains(a1, b1, b2, b3)
	dat := cha.GetData(context.Background(), key)
	fmt.Println("final dat", dat)
}

type DA struct{}

func (DA) Get(key string) (interface{}, bool) {
	fmt.Println("da get:", key)
	return nil, false
	return "daaaaaa", true
}

func (DA) Set(key string, val interface{}) {
	//do som worker

	fmt.Println("da set:", key, val)
}

type DB struct{}

func (DB) Get(key string) (interface{}, bool) {
	fmt.Println("dB get:", key)
	return "dBBBb", false
}

func (DB) Set(key string, val interface{}) {
	//do som worker

	fmt.Println("dB set:", key, val)
}

type DC struct{}

func (DC) Get(key string) (interface{}, bool) {
	fmt.Println("dC get:", key)
	return "dCCCC", true
}

func (DC) Set(key string, val interface{}) {
	//do som worker

	fmt.Println("dC set:", key, val)
}

func TestDataLoaderManager(t *testing.T) {
	dlm := NewDataLoaderManager()
	dlm.Use(DA{})
	dlm.Use(DB{})
	dlm.Use(DC{})
	val := dlm.Get("test_key")
	fmt.Println("TestDataLoaderManager", val)
}
