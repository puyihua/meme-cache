package eval

import (
	"fmt"
	"testing"
	"time"
)

const masterPort = 8080

func TestVids(t *testing.T) {
	numNode := 8
	numKeys := 1000
	err := startClusterAndPutKeys(numNode, 8, numKeys)
	if err != nil {
		t.Error(err)
	}
}

func startClusterAndPutKeys(numNode, numVidsPerNode, numKeys int) error {
	nodePorts := make([]int, numNode)
	for i := 0; i < numNode; i++ {
		nodePorts[i] = masterPort + i + 1
	}

	cc := NewCacheCluster(masterPort, nodePorts[:], numVidsPerNode)

	time.Sleep(1 * time.Second)

	err := cc.PutKeywords(numKeys)
	if err != nil {
		fmt.Println(err.Error())
	}

	lenMap, err := cc.QueryAllNodeLength()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(lenMap)
	return nil
}
