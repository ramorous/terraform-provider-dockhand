package client

// Container represents a Docker container
type Container struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	State   string            `json:"state"`
	Status  string            `json:"status"`
	Ports   []ContainerPort   `json:"ports,omitempty"`
	Mounts  []ContainerMount  `json:"mounts,omitempty"`
	Env     []string          `json:"env,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Memory  int64             `json:"memory,omitempty"`
	CPUs    float64           `json:"cpus,omitempty"`
	Restart string            `json:"restart_policy,omitempty"`
}

// ContainerPort represents a container port mapping
type ContainerPort struct {
	PrivatePort int    `json:"private_port"`
	PublicPort  int    `json:"public_port,omitempty"`
	Type        string `json:"type"`
	IP          string `json:"ip,omitempty"`
}

// ContainerMount represents a container volume mount
type ContainerMount struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	Type        string `json:"type"`
}

// ComposeStack represents a Docker Compose stack
type ComposeStack struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Compose        string                 `json:"compose"`
	Status         string                 `json:"status"`
	DesiredStatus  string                 `json:"desired_status,omitempty"`
	Services       map[string]ComposeService `json:"services,omitempty"`
	Labels         map[string]string      `json:"labels,omitempty"`
	GitRepo        *GitRepository         `json:"git_repo,omitempty"`
	AutoSync       bool                   `json:"auto_sync,omitempty"`
	WebhookToken   string                 `json:"webhook_token,omitempty"`
	CreatedAt      string                 `json:"created_at,omitempty"`
	UpdatedAt      string                 `json:"updated_at,omitempty"`
}

// ComposeService represents a service in a Compose stack
type ComposeService struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
	Count  int    `json:"count,omitempty"`
}

// GitRepository represents a Git repository configuration
type GitRepository struct {
	URL    string `json:"url"`
	Branch string `json:"branch,omitempty"`
	Path   string `json:"path,omitempty"`
	Auth   *GitAuth `json:"auth,omitempty"`
}

// GitAuth represents Git authentication
type GitAuth struct {
	Type  string `json:"type"` // "ssh" or "https"
	Token string `json:"token,omitempty"`
	Key   string `json:"key,omitempty"`
}

// Environment represents a Docker environment/host
type Environment struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"` // "local", "ssh", "docker_socket"
	Host         string                 `json:"host,omitempty"`
	Port         int                    `json:"port,omitempty"`
	Auth         *EnvironmentAuth       `json:"auth,omitempty"`
	Labels       map[string]string      `json:"labels,omitempty"`
	Active       bool                   `json:"active,omitempty"`
	DockerInfo   *DockerInfo            `json:"docker_info,omitempty"`
	CreatedAt    string                 `json:"created_at,omitempty"`
	UpdatedAt    string                 `json:"updated_at,omitempty"`
}

// EnvironmentAuth represents authentication for an environment
type EnvironmentAuth struct {
	Type     string `json:"type"` // "ssh", "tls", "basic"
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Key      string `json:"key,omitempty"`
	CertPath string `json:"cert_path,omitempty"`
}

// DockerInfo represents Docker daemon information
type DockerInfo struct {
	Version        string `json:"version"`
	APIVersion     string `json:"api_version"`
	OS             string `json:"os"`
	Architecture   string `json:"architecture"`
	Containers     int    `json:"containers"`
	ContainersRunning int `json:"containers_running"`
	ContainersPaused int `json:"containers_paused"`
	ContainersStopped int `json:"containers_stopped"`
	Images         int    `json:"images"`
}

// Network represents a Docker network
type Network struct {
	ID     string                 `json:"id"`
	Name   string                 `json:"name"`
	Type   string                 `json:"type"` // "bridge", "overlay", "host", "null"
	Driver string                 `json:"driver"`
	Scope  string                 `json:"scope"`
	Labels map[string]string      `json:"labels,omitempty"`
	IPAM   *NetworkIPAM           `json:"ipam,omitempty"`
	Containers map[string]interface{} `json:"containers,omitempty"`
}

// NetworkIPAM represents IPAM configuration for a network
type NetworkIPAM struct {
	Driver  string                 `json:"driver"`
	Config  []NetworkIPAMConfig    `json:"config,omitempty"`
	Options map[string]string      `json:"options,omitempty"`
}

// NetworkIPAMConfig represents IPAM configuration
type NetworkIPAMConfig struct {
	Subnet  string `json:"subnet,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

// Volume represents a Docker volume
type Volume struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Driver    string                 `json:"driver"`
	Mountpoint string                `json:"mountpoint"`
	Labels    map[string]string      `json:"labels,omitempty"`
	Options   map[string]string      `json:"options,omitempty"`
	Size      int64                  `json:"size,omitempty"`
	Containers []string              `json:"containers,omitempty"`
}

// Image represents a Docker image
type Image struct {
	ID          string            `json:"id"`
	RepoTags    []string          `json:"repo_tags,omitempty"`
	RepoDigests []string          `json:"repo_digests,omitempty"`
	Size        int64             `json:"size"`
	Created     string            `json:"created"`
	Labels      map[string]string `json:"labels,omitempty"`
	Architecture string           `json:"architecture,omitempty"`
	OS          string            `json:"os,omitempty"`
}

// ImagePullRequest represents an image pull request
type ImagePullRequest struct {
	Image      string `json:"image"`
	Registry   string `json:"registry,omitempty"`
	Auth       *ImageAuth `json:"auth,omitempty"`
}

// ImageAuth represents image registry authentication
type ImageAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
