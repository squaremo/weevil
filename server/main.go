package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	etcd "github.com/coreos/etcd/client"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

const (
	ServicePrefix = "/weave/service"
)

var prefixLen = len(ServicePrefix) + 1

func main() {
	var (
		etcdAddress string
	)

	flag.StringVar(&etcdAddress, "etcd", "http://etcd:4001", "Address of etcd server")
	flag.Parse()

	api, err := newWeevilAPI(etcdAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to etcd at %s", etcdAddress)

	router := mux.NewRouter()

	router.HandleFunc("/", homePage)
	router.HandleFunc("/index.html", homePage)
	router.PathPrefix("/res/").HandlerFunc(handleResource)

	router.HandleFunc("/api/{service}/", api.listInstances)
	router.HandleFunc("/api/", api.listServices)

	http.ListenAndServe("0.0.0.0:7070", router)
}

func handleResource(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[1:]
	http.ServeFile(w, r, file)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

//=== API handlers

type api struct {
	keysAPI etcd.KeysAPI
}

func newWeevilAPI(etcdAddress string) (*api, error) {
	etcdCfg := etcd.Config{
		Endpoints:               []string{etcdAddress},
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	client, err := etcd.New(etcdCfg)
	if err != nil {
		return nil, err
	}

	keysAPI := etcd.NewKeysAPI(client)
	return &api{
		keysAPI: keysAPI,
	}, nil
}

func (api *api) listServices(w http.ResponseWriter, r *http.Request) {
	opts := &etcd.GetOptions{
		Recursive: true,
	}
	res, err := api.keysAPI.Get(context.Background(), ServicePrefix, opts)
	if err != nil {
		http.Error(w, "Error from etcd: "+err.Error(), http.StatusInternalServerError)
		return
	}
	services, err := servicesFromNode(res.Node)
	if err != nil {
		http.Error(w, "Error in value for key: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(services)
}

func (api *api) listInstances(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	key := ServicePrefix + "/" + args["service"]
	// FIXME _details?
	opts := &etcd.GetOptions{
		Recursive: true,
	}
	res, err := api.keysAPI.Get(context.Background(), key, opts)
	if err != nil {
		http.Error(w, "Error from etcd: "+err.Error(), http.StatusInternalServerError)
		return
	}
	instances, err := instancesFromNode(res.Node)
	if err != nil {
		http.Error(w, "Error in value for key: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(instances)
}

func servicesFromNode(node *etcd.Node) ([]interface{}, error) {
	result := []interface{}{}
	for _, s := range node.Nodes {
		children := 0
		service := map[string]interface{}{}
		service["name"] = s.Key[prefixLen:]
		result = append(result, service)
		for _, i := range s.Nodes {
			if strings.HasSuffix(i.Key, "_details") {
				details, err := detailsFromNode(i)
				if err != nil {
					return nil, err
				}
				service["details"] = details
			} else {
				children++
			}
		}
		service["children"] = children
	}
	return result, nil
}

func instancesFromNode(node *etcd.Node) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	result["name"] = node.Key[prefixLen:]
	children := map[string]interface{}{}
	result["children"] = children
	serviceNameLen := len(node.Key) + 1
	for _, i := range node.Nodes {
		details, err := detailsFromNode(i)
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(i.Key, "_details") {
			result["details"] = details
		} else {
			children[i.Key[serviceNameLen:]] = details
		}
	}
	return result, nil
}

func detailsFromNode(node *etcd.Node) (map[string]interface{}, error) {
	var details map[string]interface{}
	err := json.Unmarshal([]byte(node.Value), &details)
	if err != nil {
		return nil, err
	}
	return details, nil
}
