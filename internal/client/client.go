package client

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// Config holds the configuration for the Dockhand API client
type Config struct {
	Endpoint      string
	Cookie        string
	Timeout       int
	TLSSkipVerify bool
}

// Client manages communication with the Dockhand API
type Client struct {
	httpClient *resty.Client
	config     *Config
}

// NewClient creates a new Dockhand API client
func NewClient(config *Config) *Client {
	httpClient := resty.New().
		SetBaseURL(config.Endpoint).
		SetHeader("Cookie", config.Cookie).
		SetHeader("Content-Type", "application/json").
		SetTimeout(time.Duration(config.Timeout) * time.Second)

	if config.TLSSkipVerify {
		httpClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return &Client{
		httpClient: httpClient,
		config:     config,
	}
}

// GetHTTPClient returns the underlying HTTP client
func (c *Client) GetHTTPClient() *resty.Client {
	return c.httpClient
}

// Container operations

// ListContainers retrieves all containers
func (c *Client) ListContainers(environmentID string) ([]Container, error) {
	var containers []Container
	resp, err := c.httpClient.R().
		SetResult(&containers).
		Get(fmt.Sprintf("/api/environments/%s/containers", environmentID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to list containers: %d %s", resp.StatusCode(), resp.String())
	}

	return containers, nil
}

// GetContainer retrieves a specific container
func (c *Client) GetContainer(environmentID, containerID string) (*Container, error) {
	var container Container
	resp, err := c.httpClient.R().
		SetResult(&container).
		Get(fmt.Sprintf("/api/environments/%s/containers/%s", environmentID, containerID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to get container: %d %s", resp.StatusCode(), resp.String())
	}

	return &container, nil
}

// CreateContainer creates a new container
func (c *Client) CreateContainer(environmentID string, container *Container) (*Container, error) {
	var result Container
	resp, err := c.httpClient.R().
		SetBody(container).
		SetResult(&result).
		Post(fmt.Sprintf("/api/environments/%s/containers", environmentID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to create container: %d %s", resp.StatusCode(), resp.String())
	}

	return &result, nil
}

// UpdateContainer updates a container
func (c *Client) UpdateContainer(environmentID, containerID string, container *Container) (*Container, error) {
	var result Container
	resp, err := c.httpClient.R().
		SetBody(container).
		SetResult(&result).
		Put(fmt.Sprintf("/api/environments/%s/containers/%s", environmentID, containerID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to update container: %d %s", resp.StatusCode(), resp.String())
	}

	return &result, nil
}

// DeleteContainer deletes a container
func (c *Client) DeleteContainer(environmentID, containerID string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/environments/%s/containers/%s", environmentID, containerID))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("failed to delete container: %d %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// StartContainer starts a container
func (c *Client) StartContainer(environmentID, containerID string) error {
	resp, err := c.httpClient.R().
		Post(fmt.Sprintf("/api/environments/%s/containers/%s/start", environmentID, containerID))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("failed to start container: %d %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// StopContainer stops a container
func (c *Client) StopContainer(environmentID, containerID string) error {
	resp, err := c.httpClient.R().
		Post(fmt.Sprintf("/api/environments/%s/containers/%s/stop", environmentID, containerID))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("failed to stop container: %d %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// RestartContainer restarts a container
func (c *Client) RestartContainer(environmentID, containerID string) error {
	resp, err := c.httpClient.R().
		Post(fmt.Sprintf("/api/environments/%s/containers/%s/restart", environmentID, containerID))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("failed to restart container: %d %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// PauseContainer pauses a container
func (c *Client) PauseContainer(environmentID, containerID string) error {
	resp, err := c.httpClient.R().
		Post(fmt.Sprintf("/api/environments/%s/containers/%s/pause", environmentID, containerID))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("failed to pause container: %d %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// UnpauseContainer unpauses a container
func (c *Client) UnpauseContainer(environmentID, containerID string) error {
	resp, err := c.httpClient.R().
		Post(fmt.Sprintf("/api/environments/%s/containers/%s/unpause", environmentID, containerID))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("failed to unpause container: %d %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Compose Stack operations

// ListComposeStacks retrieves all compose stacks
func (c *Client) ListComposeStacks(environmentID string) ([]ComposeStack, error) {
	var stacks []ComposeStack
	resp, err := c.httpClient.R().
		SetResult(&stacks).
		Get(fmt.Sprintf("/api/environments/%s/compose-stacks", environmentID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to list compose stacks: %d %s", resp.StatusCode(), resp.String())
	}

	return stacks, nil
}

// GetComposeStack retrieves a specific compose stack
func (c *Client) GetComposeStack(environmentID, stackID string) (*ComposeStack, error) {
	var stack ComposeStack
	resp, err := c.httpClient.R().
		SetResult(&stack).
		Get(fmt.Sprintf("/api/environments/%s/compose-stacks/%s", environmentID, stackID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to get compose stack: %d %s", resp.StatusCode(), resp.String())
	}

	return &stack, nil
}

// CreateComposeStack creates a new compose stack
func (c *Client) CreateComposeStack(environmentID string, stack *ComposeStack) (*ComposeStack, error) {
	var result ComposeStack
	resp, err := c.httpClient.R().
		SetBody(stack).
		SetResult(&result).
		Post(fmt.Sprintf("/api/environments/%s/compose-stacks", environmentID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to create compose stack: %d %s", resp.StatusCode(), resp.String())
	}

	return &result, nil
}

// UpdateComposeStack updates a compose stack
func (c *Client) UpdateComposeStack(environmentID, stackID string, stack *ComposeStack) (*ComposeStack, error) {
	var result ComposeStack
	resp, err := c.httpClient.R().
		SetBody(stack).
		SetResult(&result).
		Put(fmt.Sprintf("/api/environments/%s/compose-stacks/%s", environmentID, stackID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to update compose stack: %d %s", resp.StatusCode(), resp.String())
	}

	return &result, nil
}

// DeleteComposeStack deletes a compose stack
func (c *Client) DeleteComposeStack(environmentID, stackID string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/environments/%s/compose-stacks/%s", environmentID, stackID))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("failed to delete compose stack: %d %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// StartComposeStack starts a compose stack
func (c *Client) StartComposeStack(environmentID, stackID string) error {
	resp, err := c.httpClient.R().
		Post(fmt.Sprintf("/api/environments/%s/compose-stacks/%s/start", environmentID, stackID))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("failed to start compose stack: %d %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// StopComposeStack stops a compose stack
func (c *Client) StopComposeStack(environmentID, stackID string) error {
	resp, err := c.httpClient.R().
		Post(fmt.Sprintf("/api/environments/%s/compose-stacks/%s/stop", environmentID, stackID))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("failed to stop compose stack: %d %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Environment operations

// ListEnvironments retrieves all environments
func (c *Client) ListEnvironments() ([]Environment, error) {
	var environments []Environment
	resp, err := c.httpClient.R().
		SetResult(&environments).
		Get("/api/environments")

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to list environments: %d %s", resp.StatusCode(), resp.String())
	}

	return environments, nil
}

// GetEnvironment retrieves a specific environment
func (c *Client) GetEnvironment(environmentID string) (*Environment, error) {
	var env Environment
	resp, err := c.httpClient.R().
		SetResult(&env).
		Get(fmt.Sprintf("/api/environments/%s", environmentID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to get environment: %d %s", resp.StatusCode(), resp.String())
	}

	return &env, nil
}

// CreateEnvironment creates a new environment
func (c *Client) CreateEnvironment(env *Environment) (*Environment, error) {
	var result Environment
	resp, err := c.httpClient.R().
		SetBody(env).
		SetResult(&result).
		Post("/api/environments")

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to create environment: %d %s", resp.StatusCode(), resp.String())
	}

	return &result, nil
}

// UpdateEnvironment updates an environment
func (c *Client) UpdateEnvironment(environmentID string, env *Environment) (*Environment, error) {
	var result Environment
	resp, err := c.httpClient.R().
		SetBody(env).
		SetResult(&result).
		Put(fmt.Sprintf("/api/environments/%s", environmentID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to update environment: %d %s", resp.StatusCode(), resp.String())
	}

	return &result, nil
}

// DeleteEnvironment deletes an environment
func (c *Client) DeleteEnvironment(environmentID string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/environments/%s", environmentID))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("failed to delete environment: %d %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Network operations

// ListNetworks retrieves all networks
func (c *Client) ListNetworks(environmentID string) ([]Network, error) {
	var networks []Network
	resp, err := c.httpClient.R().
		SetResult(&networks).
		Get(fmt.Sprintf("/api/environments/%s/networks", environmentID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to list networks: %d %s", resp.StatusCode(), resp.String())
	}

	return networks, nil
}

// GetNetwork retrieves a specific network
func (c *Client) GetNetwork(environmentID, networkID string) (*Network, error) {
	var network Network
	resp, err := c.httpClient.R().
		SetResult(&network).
		Get(fmt.Sprintf("/api/environments/%s/networks/%s", environmentID, networkID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to get network: %d %s", resp.StatusCode(), resp.String())
	}

	return &network, nil
}

// CreateNetwork creates a new network
func (c *Client) CreateNetwork(environmentID string, network *Network) (*Network, error) {
	var result Network
	resp, err := c.httpClient.R().
		SetBody(network).
		SetResult(&result).
		Post(fmt.Sprintf("/api/environments/%s/networks", environmentID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to create network: %d %s", resp.StatusCode(), resp.String())
	}

	return &result, nil
}

// DeleteNetwork deletes a network
func (c *Client) DeleteNetwork(environmentID, networkID string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/environments/%s/networks/%s", environmentID, networkID))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("failed to delete network: %d %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Volume operations

// ListVolumes retrieves all volumes
func (c *Client) ListVolumes(environmentID string) ([]Volume, error) {
	var volumes []Volume
	resp, err := c.httpClient.R().
		SetResult(&volumes).
		Get(fmt.Sprintf("/api/environments/%s/volumes", environmentID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to list volumes: %d %s", resp.StatusCode(), resp.String())
	}

	return volumes, nil
}

// GetVolume retrieves a specific volume
func (c *Client) GetVolume(environmentID, volumeID string) (*Volume, error) {
	var volume Volume
	resp, err := c.httpClient.R().
		SetResult(&volume).
		Get(fmt.Sprintf("/api/environments/%s/volumes/%s", environmentID, volumeID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to get volume: %d %s", resp.StatusCode(), resp.String())
	}

	return &volume, nil
}

// CreateVolume creates a new volume
func (c *Client) CreateVolume(environmentID string, volume *Volume) (*Volume, error) {
	var result Volume
	resp, err := c.httpClient.R().
		SetBody(volume).
		SetResult(&result).
		Post(fmt.Sprintf("/api/environments/%s/volumes", environmentID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to create volume: %d %s", resp.StatusCode(), resp.String())
	}

	return &result, nil
}

// DeleteVolume deletes a volume
func (c *Client) DeleteVolume(environmentID, volumeID string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/environments/%s/volumes/%s", environmentID, volumeID))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("failed to delete volume: %d %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Image operations

// ListImages retrieves all images
func (c *Client) ListImages(environmentID string) ([]Image, error) {
	var images []Image
	resp, err := c.httpClient.R().
		SetResult(&images).
		Get(fmt.Sprintf("/api/environments/%s/images", environmentID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to list images: %d %s", resp.StatusCode(), resp.String())
	}

	return images, nil
}

// GetImage retrieves a specific image
func (c *Client) GetImage(environmentID, imageID string) (*Image, error) {
	var image Image
	resp, err := c.httpClient.R().
		SetResult(&image).
		Get(fmt.Sprintf("/api/environments/%s/images/%s", environmentID, imageID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to get image: %d %s", resp.StatusCode(), resp.String())
	}

	return &image, nil
}

// PullImage pulls an image
func (c *Client) PullImage(environmentID string, pullReq *ImagePullRequest) error {
	resp, err := c.httpClient.R().
		SetBody(pullReq).
		Post(fmt.Sprintf("/api/environments/%s/images/pull", environmentID))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("failed to pull image: %d %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// DeleteImage deletes an image
func (c *Client) DeleteImage(environmentID, imageID string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/environments/%s/images/%s", environmentID, imageID))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("failed to delete image: %d %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Health check operations

// HealthCheck performs a health check on the API
func (c *Client) HealthCheck() (bool, error) {
	resp, err := c.httpClient.R().
		Get("/api/health")

	if err != nil {
		return false, err
	}

	return resp.IsSuccess(), nil
}
