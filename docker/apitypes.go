package docker

import "time"

type Catalog struct {
	Repositories []string `json:"repositories"`
}

type Tags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type Blob struct {
	Architecture string `json:"architecture"`
	Config       struct {
		Hostname     string                 `json:"Hostname"`
		Domainname   string                 `json:"Domainname"`
		User         string                 `json:"User"`
		AttachStdin  bool                   `json:"AttachStdin"`
		AttachStdout bool                   `json:"AttachStdout"`
		AttachStderr bool                   `json:"AttachStderr"`
		ExposedPorts map[string]interface{} `json:"ExposedPorts"`
		Tty          bool                   `json:"Tty"`
		OpenStdin    bool                   `json:"OpenStdin"`
		StdinOnce    bool                   `json:"StdinOnce"`
		Env          []string               `json:"Env"`
		Cmd          []string               `json:"Cmd"`
		Image        string                 `json:"Image"`
		Volumes      map[string]interface{} `json:"Volumes"`
		WorkingDir   string                 `json:"WorkingDir"`
		Entrypoint   []string               `json:"Entrypoint"`
		OnBuild      interface{}            `json:"OnBuild"`
		Labels       map[string]string      `json:"Labels"`
	} `json:"config"`
	Container       string `json:"container"`
	ContainerConfig struct {
		Hostname     string            `json:"Hostname"`
		Domainname   string            `json:"Domainname"`
		User         string            `json:"User"`
		AttachStdin  bool              `json:"AttachStdin"`
		AttachStdout bool              `json:"AttachStdout"`
		AttachStderr bool              `json:"AttachStderr"`
		Tty          bool              `json:"Tty"`
		OpenStdin    bool              `json:"OpenStdin"`
		StdinOnce    bool              `json:"StdinOnce"`
		Env          []string          `json:"Env"`
		Cmd          []string          `json:"Cmd"`
		Image        string            `json:"Image"`
		Volumes      interface{}       `json:"Volumes"`
		WorkingDir   string            `json:"WorkingDir"`
		Entrypoint   []string          `json:"Entrypoint"`
		OnBuild      interface{}       `json:"OnBuild"`
		Labels       map[string]string `json:"Labels"`
	} `json:"container_config"`
	Created       time.Time `json:"created"`
	DockerVersion string    `json:"docker_version"`
	History       []struct {
		Created    time.Time `json:"created"`
		CreatedBy  string    `json:"created_by"`
		EmptyLayer bool      `json:"empty_layer,omitempty"`
	} `json:"history"`
	Os     string `json:"os"`
	Rootfs struct {
		Type    string   `json:"type"`
		DiffIds []string `json:"diff_ids"`
	} `json:"rootfs"`
}

type ManifestV1 struct {
	SchemaVersion int    `json:"schemaVersion"`
	Name          string `json:"name"`
	Tag           string `json:"tag"`
	Architecture  string `json:"architecture"`
	FsLayers      []struct {
		BlobSum string `json:"blobSum"`
	} `json:"fsLayers"`
	History []struct {
		V1Compatibility string `json:"v1Compatibility"`
	} `json:"history"`
	Signatures []struct {
		Header struct {
			Jwk struct {
				Crv string `json:"crv"`
				Kid string `json:"kid"`
				Kty string `json:"kty"`
				X   string `json:"x"`
				Y   string `json:"y"`
			} `json:"jwk"`
			Alg string `json:"alg"`
		} `json:"header"`
		Signature string `json:"signature"`
		Protected string `json:"protected"`
	} `json:"signatures"`
}

type ManifestV2 struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"config"`
	Layers []struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"layers"`
}
