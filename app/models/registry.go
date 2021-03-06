package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	client "github.com/zhujingfa/docker-registry-client/registry"
)

// AllRegistries contains a list of added registries using their hostnames
// access granted via mutex locks/unlocks
var AllRegistries Registries

func init() {
	AllRegistries.Registries = map[string]*Registry{}
}

// Registries contains a map of all active registries identified by their name, locked when necessary
type Registries struct {
	Registries map[string]*Registry
	sync.Mutex
}

const STATUS_DOWN="DOWN"
const STATUS_UP="UP"

type Registry struct {
	*client.Registry
	Repositories map[string]*Repository
	TTL          time.Duration
	Ticker       *time.Ticker
	Name         string
	Host         string
	Scheme       string
	Version      string
	Port         int
	sync.Mutex
	status       string
	ip           string
}

func (r *Registry) IP() string {
	return r.ip
}

// Refresh is called with the configured TTL time for the given registry
func (r *Registry) Refresh() {
	//原来这里更新的都是副本。
	// Copy the registry information to a new object, and update it
	ur := *r

	err := r.Ping();
	if err != nil {
		ur.status = STATUS_DOWN
	} else {
		ur.status = STATUS_UP
	}

	ip, _ := net.LookupHost(r.Host)
	if len(ip) > 0 {
		r.ip = ip[0]
	}

	logrus.Info("Refreshing " + r.URL)
	// Get the list of repositories
	repos, err := ur.Registry.Repositories()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
		}).Error("Failed to retrieve an updated list of repositories for " + ur.URL)
	}
	// Get the repository information
	ur.Repositories = make(map[string]*Repository)
	for _, repoName := range repos {

		// Get the list of tags for the repository
		tags, err := ur.Tags(repoName)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Error":           err.Error(),
				"Repository Name": repoName,
			}).Error("Failed to retrieve an updated list of tags for " + ur.URL)
			continue
		}

		repo := Repository{Name: repoName, Tags: make(map[string]*Tag)}
		// Get the manifest for each of the tags
		for _, tagName := range tags {

			// Using v2 required getting the manifest then retrieving the blob
			// for the config digest
			man, err := ur.Manifest(repoName, tagName)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"Error":           err.Error(),
					"Repository Name": repoName,
					"Tag Name":        tagName,
				}).Error("Failed to retrieve manifest information for " + ur.URL)
				continue
			}

			// Get the v1 config information
			v1Bytes, err := ur.ManifestMetadata(repoName, man.Config.Digest)
			if err != nil {
				logrus.Error(err)
				continue
			}
			var v1 V1Compatibility
			err = json.Unmarshal(v1Bytes, &v1)
			if err != nil {
				logrus.Error(err)
				continue
			}

			// add the pointer for the history to its layer using its index
			layerIndex := 0
			for i, history := range v1.History {
				if !history.EmptyLayer {
					v1.History[i].ManifestLayer = &man.Layers[layerIndex]
					layerIndex++
				}
				sh := strings.Split(history.CreatedBy, "/bin/sh -c")
				if len(sh) > 1 {
					v1.History[i].ShellType = "/bin/sh -c"
					commands := strings.SplitAfter(sh[1], "&&")
					for _, cmd := range commands {
						//Keywords: Keywords(cmd)
						v1.History[i].Commands = append(v1.History[i].Commands, Command{Cmd: cmd, Keywords: []string{}})
					}
				}
			}

			// Get the tag size information
			size, err := ur.TagSizeByObj(man)
			if err != nil {
				logrus.Error(err)
			}

			repo.Tags[tagName] = &Tag{Name: tagName, V1Compatibility: &v1, Size: int64(size), DeserializedManifest: man}
		}
		ur.Repositories[repoName] = &repo

		//time.Sleep(10*time.Millisecond)
	}

	AllRegistries.Lock()
	AllRegistries.Registries[ur.Name] = &ur
	AllRegistries.Unlock()
}

// TagCount returns the total number of tags across all repositories
func (r *Registry) TagCount() (count int) {
	for _, repo := range r.Repositories {
		count += len(repo.Tags)
	}
	return count
}

// TagCount returns the total number of layers across all repositories
func (r *Registry) LayerCount() int {
	layerDigests := make(map[string]struct{})
	for _, repo := range r.Repositories {
		for _, tag := range repo.Tags {
			for _, layer := range tag.Layers {
				layerDigests[layer.Digest.String()] = struct{}{}
			}
		}
	}
	return len(layerDigests)
}

// Pushes returns the number of pushes recorded by passing the forwarded registry events
func (r *Registry) Pushes() (pushes int) {
	AllEvents.Lock()
	defer AllEvents.Unlock()
	if _, ok := AllEvents.Events[r.Name]; !ok {
		return 0
	}

	for _, e := range AllEvents.Events[r.Name] {
		// TODO: really need to find a better way to exclude the managers queries
		if e.Action == "push" && e.Request.Useragent != "Go-http-client/1.1" && e.Request.Method != "HEAD" {
			pushes++
		}
	}
	return pushes
}

// Pulls returns the number of pulls recorded by passing the forwarded registry events
func (r *Registry) Pulls() (pulls int) {
	AllEvents.Lock()
	defer AllEvents.Unlock()
	if _, ok := AllEvents.Events[r.Name]; !ok {
		return 0
	}

	for _, e := range AllEvents.Events[r.Name] {
		// exclude heads since thats the method the manager uses for getting meta info
		// TODO: really need to find a better way to exclude the managers queries
		if e.Action == "pull" && e.Request.Useragent != "Go-http-client/1.1" && e.Request.Method != "HEAD" {
			pulls++
		}
	}
	return pulls
}

// Status returns the text representation of whether the registry is reachable
func (r *Registry) Status() string {
	return r.status
}


// AddRegistry adds the new registry for viewing in the interface and sets up
// the go routine for automatic refreshes
func AddRegistry(scheme, host, user, password string, port int, ttl time.Duration, skipTLS bool) (*Registry, error) {
	switch {
	case scheme == "":
		return nil, errors.New("Invalid scheme: " + scheme)
	case host == "":
		return nil, errors.New("Invalid host: " + host)
	case port == 0:
		return nil, errors.New("Invalid port: " + strconv.Itoa(port))
	}

	var url string
	if scheme == "https" && port == 443 {
		url = fmt.Sprintf("%s://%s", scheme, host)
	} else {
		url = fmt.Sprintf("%s://%s:%v", scheme, host, port)
	}


	var hub *client.Registry
	var err error
	if skipTLS {
		hub, err = client.NewInsecure(url, user, password)
		if err != nil {
			logrus.Error("Failed to connect to unvalidated TLS registry: " + err.Error())
			return nil, err
		}
	} else {
		hub, err = client.New(url, user, password)
		if err != nil {
			logrus.Error("Failed to connect to validated registry: " + err.Error())
			return nil, err
		}
	}

	//fuck, 这里是对象，不是指针
	r := &Registry{
		Registry: hub,
		TTL:      ttl,
		Ticker:   time.NewTicker(ttl),
		Host:     host,
		Scheme:   scheme,
		Port:     port,
		Version:  "v2",
		Name:     host + ":" + strconv.Itoa(port),
		status:STATUS_DOWN,
	}

	r.Refresh()


	go func() {
		for range r.Ticker.C {
			r.Refresh()
		}
	}()

	return AllRegistries.Registries[r.Name], nil
}
