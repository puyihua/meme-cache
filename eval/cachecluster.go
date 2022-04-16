package eval

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/puyihua/meme-cache/internal/master"
	"github.com/puyihua/meme-cache/internal/node"
)

const localhost = "http://127.0.0.1:"

type CacheCluster struct {
	Master *MasterServer
	Nodes  []*NodeServer
}

func NewCacheCluster(masterPort int, nodePorts []int, numVidPerNode int, storeType int) *CacheCluster {
	// master
	masterUrl := localhost + strconv.Itoa(masterPort)
	ms := &MasterServer{masterUrl}
	go func() {
		masterSrv := master.NewServerMaster(masterPort)
		masterSrv.Serve()
	}()

	var nodes []*NodeServer
	// nodes
	for _, port := range nodePorts {
		go func() {
			fmt.Printf("port: %d\n", port)
			nodeSrv := node.NewServerWithType(port, storeType)
			nodeSrv.Serve()
		}()
		vids := randVids(numVidPerNode)
		ms.AddMember(port, vids)

		nodeServer := &NodeServer{
			Url:  localhost + strconv.Itoa(port),
			Vids: vids,
		}
		nodes = append(nodes, nodeServer)
	}

	return &CacheCluster{ms, nodes}
}



func randVids(n int) []uint64 {
	vids := make([]uint64, n)
	for i := range vids {
		vids[i] = rand.Uint64()
	}
	return vids
}

type MasterServer struct {
	Url string
}

func (ms *MasterServer) AddMember(port int, vids []uint64) error {
	vidsStr := make([]string, len(vids))
	for i, vid := range vids {
		vidsStr[i] = strconv.FormatUint(vid, 10)
	}
	vidsStrJoin := strings.Join(vidsStr, ",")

	_, err := http.Head(ms.Url + "/addMember?host=127.0.0.1&port=" + strconv.Itoa(port) + "&vid=" + vidsStrJoin)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MasterServer) Get(key string) (string, error) {
	resp, err := http.Get(ms.Url + "/kv/get?key=" + key)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	value := string(bytes)
	return value, nil
}

func (ms *MasterServer) Put(key string, value string) error {
	keyUrl := url.QueryEscape(key)
	valueUrl := url.QueryEscape(value)
	resp, err := http.Get(ms.Url + "/kv/put?key=" + keyUrl + "&value=" + valueUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Client->Master: Put failed with status code %d", resp.StatusCode)
	}

	return nil
}

type NodeServer struct {
	Url  string
	Vids []uint64
}

func (ns *NodeServer) GetLen() (int, error) {
	url := ns.Url
	resp, err := http.Get(url + "/getlen")
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	byteArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	length, err := strconv.Atoi(string(byteArr))
	if err != nil {
		return -1, err
	}
	return length, nil
}
