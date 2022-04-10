package eval

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

func (cc *CacheCluster) QueryAllNodeLength() (map[string]int, error) {
	lenMap := make(map[string]int)
	for _, node := range cc.Nodes {
		length, err := node.GetLen()
		if err != nil {
			return nil, err
		}
		lenMap[node.Url] = length
	}
	return lenMap, nil
}

func (cc *CacheCluster) PutKeywords(n int) error {
	// read keywords from json
	jsonFile, err := os.Open("keywords.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	bytes, _ := io.ReadAll(jsonFile)

	var keywordsMap map[string][]string
	json.Unmarshal([]byte(bytes), &keywordsMap)
	keywords := keywordsMap["keywords"][:n]

	// put them into cache
	fmt.Println("start to put " + strconv.Itoa(n) + " keywords...")
	for i, keyword := range keywords {
		err = cc.Master.Put(keyword, "dummy")
		if err != nil {
			return err
		}
		if (i+1)%1000 == 0 {
			fmt.Printf("finished %d keys\n", i+1)
		}
	}
	fmt.Println("all finished putting...")

	return nil
}
