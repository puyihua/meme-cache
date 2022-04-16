package eval

import (
	"fmt"
	"github.com/puyihua/meme-cache/internal/node/store"
	"testing"
	"time"
)

func TestMigration(t *testing.T) {
	cc := NewCacheCluster(masterPort, []int{8081}, 2, store.TypeFineGrained)
	err := cc.PutKeywords(1000)
	if err != nil {
		fmt.Println(err.Error())
	}

	lenMap1, err := cc.QueryAllNodeLength()
	if err != nil {
		fmt.Println(err.Error())
	}

	total1 := sumValue(lenMap1)
	if total1 != 1000 {
		t.Fatalf("Wrong number of total keys %d != %d", total1, 1000)
	}

	cc.AddNode(8082, 2, store.TypeFineGrained)

	lenMap2, err := cc.QueryAllNodeLength()
	if err != nil {
		fmt.Println(err.Error())
	}

	time.Sleep(1 * time.Second)

	fmt.Println(lenMap2)

	total2 := sumValue(lenMap2)
	if total2 != 1000 {
		t.Fatalf("Wrong number of total keys %d != %d", total2, 1000)
	}

}

func sumValue(m map[string]int) int {
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}
