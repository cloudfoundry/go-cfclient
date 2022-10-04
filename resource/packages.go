package resource

import (
	"encoding/json"
	"errors"
)

type PackageState string

const (
	AwaitingUpload   PackageState = "AWAITING_UPLOAD"
	ProcessingUpload PackageState = "PROCESSING_UPLOAD"
	Ready            PackageState = "READY"
	Failed           PackageState = "FAILED"
	Copying          PackageState = "COPYING"
	Expired          PackageState = "EXPIRED"
)

type Package struct {
	Type      string          `json:"type,omitempty"` // bits or docker
	Data      json.RawMessage `json:"data,omitempty"` // depends on value of Type
	State     PackageState    `json:"state,omitempty"`
	GUID      string          `json:"guid,omitempty"`
	CreatedAt string          `json:"created_at,omitempty"`
	UpdatedAt string          `json:"updated_at,omitempty"`
	Links     map[string]Link `json:"links,omitempty"`
	Metadata  Metadata        `json:"metadata,omitempty"`
}

// BitsPackage is the data for Packages of type bits.
// It provides an upload link to which a zip file should be uploaded.
type BitsPackage struct {
	Error    string `json:"error,omitempty"`
	Checksum struct {
		Type  string `json:"type,omitempty"`  // eg. sha256
		Value string `json:"value,omitempty"` // populated after the bits are uploaded
	} `json:"checksum,omitempty"`
}

// DockerPackage is the data for Packages of type docker.
// It references a docker image from a registry.
type DockerPackage struct {
	Image    string `json:"image,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type ListPackagesResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []Package  `json:"resources,omitempty"`
}

type DockerCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type DockerPackageData struct {
	Image string `json:"image"`
	*DockerCredentials
}

type CreateDockerPackageRequest struct {
	Type          string                       `json:"type"`
	Relationships map[string]ToOneRelationship `json:"relationships"`
	Data          DockerPackageData            `json:"data"`
}

func (v *Package) BitsData() (BitsPackage, error) {
	var bits BitsPackage
	if v.Type != "bits" {
		return bits, errors.New("this package is not of type bits")
	}

	if err := json.Unmarshal(v.Data, &bits); err != nil {
		return bits, err
	}

	return bits, nil
}

func (v *Package) DockerData() (DockerPackage, error) {
	var docker DockerPackage
	if v.Type != "docker" {
		return docker, errors.New("this package is not of type docker")
	}

	if err := json.Unmarshal(v.Data, &docker); err != nil {
		return docker, err
	}

	return docker, nil
}
