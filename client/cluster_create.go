package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

type CreateManagedClusterRequest struct {
	OrganizationID string
	ProjectID      string
	NetworkId      string `json:"networkId"`
	Name           string `json:"description"`
	Topology       string `json:"topology"`
	InstanceType   string `json:"instanceType"`
	DiskSizeGB     int32  `json:"diskSizeGb"`
	DiskType       string `json:"diskType"`
	ServerVersion  string `json:"serverVersion"`
}

type CreateManagedClusterResponse struct {
	ClusterID string `json:"id"`
}

func (c *Client) ManagedClusterCreate(ctx context.Context, req *CreateManagedClusterRequest) (*CreateManagedClusterResponse, error) {
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	requestURL := *c.apiURL
	requestURL.Path = path.Join("mesdb", "v1", "organizations", req.OrganizationID, "projects", req.ProjectID, "clusters")

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL.String(), bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error constructing request: %w", err)
	}
	request.Header.Add("Content-Type", "application/json")
	if err := c.addAuthorizationHeader(request); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, translateStatusCode(resp.StatusCode, "creating managed cluster")
	}

	decoder := json.NewDecoder(resp.Body)
	result := CreateManagedClusterResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &result, nil
}