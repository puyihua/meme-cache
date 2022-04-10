package eval

import (
	"fmt"
	"testing"
	"time"
)

func TestVids(t *testing.T) {
	masterPort := 8080
	nodePorts := [2]int{8081, 8082}

	cc := NewCacheCluster(masterPort, nodePorts[:], 1)

	time.Sleep(1 * time.Second)

	err := cc.PutKeywords(1000)
	if err != nil {
		fmt.Println(err.Error())
	}

	lenMap, err := cc.QueryAllNodeLength()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(lenMap)
}
