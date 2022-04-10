package eval

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func queryAllNodeLength(urls []string) []int {
	var lengths []int
	for _, url := range urls {
		resp, _ := http.Get(url + "/getlen")
		defer resp.Body.Close()
		byteArr, _ := io.ReadAll(resp.Body)
		respStr := string(byteArr)
		length, _ := strconv.Atoi(strings.Split(respStr, ":")[1])
		lengths = append(lengths, length)
	}
	return lengths
}

func putKeywords(masterUrl string, n int) error {
	jsonFile, err := os.Open("keywords.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	bytes, _ := io.ReadAll(jsonFile)

	var keywordsMap map[string][]string
	json.Unmarshal([]byte(bytes), &keywordsMap)

	keywords := keywordsMap["keywords"][:n]

	fmt.Println("start to put " + strconv.Itoa(n) + " keywords...")
	for _, keyword := range keywords {
		keywordUrl := url.QueryEscape(keyword)
		_, err = http.Get(masterUrl + "/put?key=" + keywordUrl + "&value=dummy")
		if err != nil {
			return err
		}
	}
	fmt.Println("finished putting...")

	return nil
}
