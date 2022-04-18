package docker

import (
	"encoding/json"
	"fmt"
	"github.com/rcomanne/docker-registry-gui/pkg/configuration"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	client       *http.Client
	dockerConfig *configuration.Docker
}

func NewClient(dockerConfig *configuration.Docker) *Client {
	return &Client{
		client:       &http.Client{},
		dockerConfig: dockerConfig,
	}
}

func (c *Client) Validate() bool {
	valid := true

	// build the request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2", c.dockerConfig.Address), nil)
	handleError(err)

	// execute it
	byteBody := c.doExecute(req)

	// map the result
	if string(byteBody) != "{}" {
		var result ErrorResponse
		err = json.Unmarshal(byteBody, &result)
		for _, e := range result.Errors {
			log.Printf("code = %s, message = %s", e.Code, e.Message)
		}
		valid = false
	}

	return valid
}

func (c *Client) ListRepositories() Catalog {
	// build the request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/_catalog", c.dockerConfig.Address), nil)
	handleError(err)
	// execute it
	byteBody := c.doExecute(req)

	// map the result
	var result Catalog
	err = json.Unmarshal(byteBody, &result)
	handleError(err)

	return result
}

func (c *Client) ListRepositoryTags(repositoryName string) Tags {
	// build request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/%s/tags/list", c.dockerConfig.Address, unescapePath(repositoryName)), nil)
	handleError(err)
	// execute it
	byteBody := c.doExecute(req)

	// map the result
	var result Tags
	err = json.Unmarshal(byteBody, &result)
	handleError(err)

	return result
}

func (c *Client) GetManifestV1(repositoryName, tag string) ManifestV1 {
	// build request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/%s/manifests/%s", c.dockerConfig.Address, unescapePath(repositoryName), tag), nil)
	handleError(err)
	// execute the request
	byteBody := c.doExecute(req)

	// map the result
	var result ManifestV1
	err = json.Unmarshal(byteBody, &result)
	handleError(err)

	return result
}

func (c *Client) GetManifestV2(repositoryName, tag string) ManifestV2 {
	// build request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/%s/manifests/%s", c.dockerConfig.Address, unescapePath(repositoryName), tag), nil)
	handleError(err)
	// add header for a response that contains the actual digest required for other parts of the API
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	// execute it
	byteBody := c.doExecute(req)

	// map the result
	var result ManifestV2
	err = json.Unmarshal(byteBody, &result)
	handleError(err)

	return result
}

func (c *Client) GetBlob(repositoryName, digest string) Blob {
	// build request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/%s/blobs/%s", c.dockerConfig.Address, unescapePath(repositoryName), digest), nil)
	handleError(err)
	// execute it
	byteBody := c.doExecute(req)

	// map the result
	var result Blob
	err = json.Unmarshal(byteBody, &result)
	handleError(err)

	return result
}

func (c *Client) doExecute(req *http.Request) []byte {
	// add basic auth to request
	if c.dockerConfig.HasAuthentication() {
		addAuthentication(req, c.dockerConfig)
	}

	// do request
	resp, err := c.client.Do(req)
	handleError(err)
	defer resp.Body.Close()

	// parse response to byte array
	byteBody, err := ioutil.ReadAll(resp.Body)
	handleError(err)
	return byteBody
}

func addAuthentication(req *http.Request, dockerConfig *configuration.Docker) {
	// add basic auth header from the config
	req.SetBasicAuth(dockerConfig.Username, dockerConfig.Password)
}

func unescapePath(in string) string {
	if out, err := url.PathUnescape(in); err != nil {
		handleError(err)
		return ""
	} else {
		return out
	}
}

func handleError(err error) {
	// TODO: better error handling...
	if err != nil {
		log.Fatalln(err)
	}
}
