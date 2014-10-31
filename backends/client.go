package backends

import (
	"errors"
	"strings"

	"github.com/kelseyhightower/confd/backends/consul"
	"github.com/kelseyhightower/confd/backends/env"
	"github.com/kelseyhightower/confd/backends/etcd"
	"github.com/kelseyhightower/confd/log"
)

// The StoreClient interface is implemented by objects that can retrieve
// key/value pairs from a backend store.
type StoreClient interface {
	GetValues(keys []string) (map[string]string, error)
	WatchPrefix(prefix string, waitIndex uint64, stopChan chan bool) (uint64, error)
}

// NewBackendClient is the function that custom backends must implement.  It takes
// an array of backend nodes.
type NewBackendClient func([]string) (StoreClient, error)

var backends = map[string]NewBackendClient{}

// New is used to create a storage client based on our configuration.
func New(config Config) (StoreClient, error) {
	if config.Backend == "" {
		config.Backend = "etcd"
	}
	backendNodes := config.BackendNodes
	log.Notice("Backend nodes set to " + strings.Join(backendNodes, ", "))
	switch config.Backend {
	case "consul":
		return consul.NewConsulClient(backendNodes)
	case "etcd":
		// Create the etcd client upfront and use it for the life of the process.
		// The etcdClient is an http.Client and designed to be reused.
		return etcd.NewEtcdClient(backendNodes, config.ClientCert, config.ClientKey, config.ClientCaKeys)
	case "env":
		return env.NewEnvClient()
	default:
		if backend := backends[config.Backend]; backend != nil {
			return backend(backendNodes)
		}
		return nil, errors.New("Invalid backend")
	}
}

// Register can be used to extend confd. Backends must implement the StoreClient interface
func Register(name string, constructor NewBackendClient) {
	backends[name] = constructor
}
