package master

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ServerMaster struct {
	port  int
	lib *LibMasterCH
}

func NewServerMaster(port int) *ServerMaster {
	return &ServerMaster{port: port, lib: NewLibMasterCH()}
}

func (svr ServerMaster) getMembersHandler() string {
	return strings.Join(svr.lib.GetMembers(), ",")
}

func (svr ServerMaster) addMemberHandler(r *http.Request) string {
	queryMap, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return "Wrong Request Format"
	}

	host, ok := queryMap["host"]
	if !ok {
		return "Wrong Request Format"
	}

	port, ok := queryMap["port"]
	if !ok {
		return "Wrong Request Format"
	}

	hostport := host[0] + ":" + port[0]

	vidStrs, ok := queryMap["vid"]
	if !ok {
		return "Wrong Request Format"
	}
	
	var vids []uint64

	for i := range vidStrs {
		vid, err := strconv.ParseUint(vidStrs[i], 10, 64)
		if err != nil {
			return "vid must be a uint64 number"
		}
		vids = append(vids, vid)
	}

	if err := svr.lib.AddMember(hostport, vids); err != nil {
		return err.Error()
	}
	return "Done"
}

func (svr ServerMaster) Serve() {
	http.HandleFunc("/getMembers", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, svr.getMembersHandler())
	})

	http.HandleFunc("/addMember", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, svr.addMemberHandler(r))
	})

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(svr.port), nil))
}