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
	port int
	lib  *LibMasterCH
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

	vidStrSingleStr, ok := queryMap["vid"]
	if !ok {
		return "Wrong Request Format"
	}
	vidStrs := strings.Split(vidStrSingleStr[0], ",")

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

func (svr ServerMaster) getHandler(u *url.URL) string {
	queryMap, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "Wrong Request Format"
	}
	keys, ok := queryMap["key"]
	if !ok {
		return "Wrong Request Format"
	}

	val, err := svr.lib.Get(keys[0])

	if err != nil {
		return err.Error()
	}

	return "{" + keys[0] + ": " + val + "}"
}

func (svr ServerMaster) putHandler(u *url.URL) string {
	queryMap, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "Wrong Request Format"
	}

	keys, ok := queryMap["key"]
	if !ok {
		return "Wrong Request Format"
	}

	values, ok := queryMap["value"]
	if !ok {
		return "Wrong Request Format"
	}
	err2 := svr.lib.Put(keys[0], values[0])

	if err2 != nil {
		return err.Error()
	}

	return "Done"
}

func (svr ServerMaster) routerHandler(u *url.URL) string {
	queryMap, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "Wrong Request Format"
	}

	keys, ok := queryMap["key"]
	if !ok {
		return "Wrong Request Format"
	}

	hostport, errRouter := svr.lib.Router(keys[0])

	if errRouter != nil {
		return errRouter.Error()
	}

	return hostport
}

func (svr ServerMaster) Serve() {
	http.HandleFunc("/getMembers", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, svr.getMembersHandler())
	})

	http.HandleFunc("/addMember", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, svr.addMemberHandler(r))
	})

	http.HandleFunc("/kv/get", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, svr.getHandler(r.URL))
	})

	http.HandleFunc("/kv/put", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, svr.putHandler(r.URL))
	})

	http.HandleFunc("/kv/router", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, svr.routerHandler(r.URL))
	})

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(svr.port), nil))
}
