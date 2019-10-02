package models

import (
	nks "github.com/NetApp/nks-sdk-go/nks"
)

type CreateClusterInputGCE struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	//WorkspaceID      int    `json:"workspace"`
	ProviderKey int `json:"provider_keyset"`

	MasterCount int    `json:"master_count"`
	MasterSize  string `json:"master_size"`
	//MasterRootDiskSize int    `json:"master_root_disk_size"`
	//MasterGPUInstanceSize string `json:"master_gpu_instance_size"`
	//MasterGPUCoreCount    int    `json:"master_gpu_core_count"`

	WorkerCount int    `json:"worker_count"`
	WorkerSize  string `json:"worker_size"`
	//WorkerGPUInstanceSize string `json:"worker_gpu_instance_size"`
	//WorkerGPUCoreCount    int    `json:"worker_gpu_core_count"`
	//WorkerRootDiskSize    int    `json:"worker_root_disk_size"`

	KubernetesVersion string `json:"k8s_version"`
	DashboardEnabled  bool   `json:"k8s_dashboard_enabled"`
	RbacEnabled       bool   `json:"k8s_rbac_enabled"`
	//K8sPodCIDR          string `json:"k8s_pod_cidr"`
	//K8sServiceCIDR      string `json:"k8s_service_cidr"`

	//ProjectID string `json:"project_id"`

	SSHKeySet int `json:"user_ssh_keyset"`

	EtcdType  string         `json:"etcd_type"`
	Platform  string         `json:"platform"`
	Channel   string         `json:"channel"`
	Region    string         `json:"region"`
	Solutions []nks.Solution `json:"solutions"`

	//Config string `json:"config"`

	//MinNodeCount int `json:"min_node_count"`
	//MaxNodeCount int `json:"max_node_count"`
	//Owner        int `json:"owner"`
	//ProviderSubnetID    string
	//ProviderSubnetCidr  string
	//ProviderNetworkID   string
	//ProviderNetworkCIDR string
}

type CreateClusterInput struct {
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	WorkspaceID int    `json:"workspace,omitempty"`
	ProviderKey int    `json:"provider_keyset,omitempty"`

	MasterCount           int    `json:"master_count,omitempty"`
	MasterSize            string `json:"master_size,omitempty"`
	MasterRootDiskSize    int    `json:"master_root_disk_size,omitempty"`
	MasterGPUInstanceSize string `json:"master_gpu_instance_size,omitempty"`
	MasterGPUCoreCount    *int   `json:"master_gpu_core_count,omitempty"`

	WorkerCount           int    `json:"worker_count,omitempty"`
	WorkerSize            string `json:"worker_size,omitempty"`
	WorkerGPUInstanceSize string `json:"worker_gpu_instance_size,omitempty"`
	WorkerGPUCoreCount    *int   `json:"worker_gpu_core_count,omitempty"`
	WorkerRootDiskSize    int    `json:"worker_root_disk_size,omitempty"`

	KubernetesVersion string `json:"k8s_version,omitempty"`
	DashboardEnabled  bool   `json:"k8s_dashboard_enabled,omitempty"`
	RbacEnabled       bool   `json:"k8s_rbac_enabled,omitempty"`
	K8sPodCIDR        string `json:"k8s_pod_cidr,omitempty"`
	K8sServiceCIDR    string `json:"k8s_service_cidr,omitempty"`

	ProjectID string `json:"project_id,omitempty"`

	SSHKeySet int `json:"user_ssh_keyset,omitempty"`

	EtcdType  string         `json:"etcd_type,omitempty"`
	Platform  string         `json:"platform,omitempty"`
	Channel   string         `json:"channel,omitempty"`
	Region    string         `json:"region,omitempty"`
	Solutions []nks.Solution `json:"solutions,omitempty"`

	Config *map[string]bool `json:"config,omitempty"`

	MinNodeCount        *int     `json:"min_node_count,omitempty"`
	MaxNodeCount        *int     `json:"max_node_count,omitempty"`
	Owner               int      `json:"owner,omitempty"`
	ProviderSubnetID    string   `json:"provider_subnet_id,omitempty"`
	ProviderSubnetCidr  string   `json:"provider_subnet_cidr,omitempty"`
	ProviderNetworkID   string   `json:"provider_network_id,omitempty"`
	ProviderNetworkCIDR string   `json:"provider_network_cidr,omitempty"`
	Features            []string `json:"features,omitempty"`
}
