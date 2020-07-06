package v3

//KARBON
///Responses
type KarbonClusterListIntentResponse []KarbonClusterClusterMetadataIntentResponse

// single element in cluster/list not returning the same /cluster/uuid
type KarbonClusterClusterMetadataIntentResponse struct {
	KarbonClusterMetadataResponse *KarbonClusterIntentResponse `json:"cluster_metadata" mapstructure:"cluster_metadata, omitempty"`
	TaskProgressMessage           *string                      `json:"task_progress_message" mapstructure:"task_progress_message, omitempty"`
	TaskProgressPercent           *int64                       `json:"task_progress_percent" mapstructure:"task_progress_percent, omitempty"`
	TaskStatus                    *int64                       `json:"task_status" mapstructure:"task_status, omitempty"`
	TaskType                      *string                      `json:"task_type" mapstructure:"task_type, omitempty"`
}

//return type for /cluster/uuid

type KarbonClusterIntentResponse struct {
	AddonsConfig *KarbonClusterAddonsConfigResponse `json:"addons_config" mapstructure:"addons_config, omitempty"`
	EtcdConfig   *KarbonClusterEtcdConfigResponse   `json:"etcd_config" mapstructure:"etcd_config, omitempty"`
	K8sConfig    *KarbonClusterK8sConfigResponse    `json:"k8s_config" mapstructure:"k8s_config, omitempty"`
	Name         *string                            `json:"name" mapstructure:"name, omitempty"`
	UUID         *string                            `json:"uuid" mapstructure:"uuid, omitempty"`
}

type KarbonClusterAddonsConfigResponse struct {
	KarbonClusterLoggingConfigResponse `json:"logging_config" mapstructure:"logging_config, omitempty"`
}

type KarbonClusterLoggingConfigResponse struct {
	State         *string `json:"state" mapstructure:"state, omitempty"`
	StorageSizeMb *int64  `json:"storage_size_mib" mapstructure:"storage_size_mib, omitempty"`
	Version       *string `json:"version" mapstructure:"version, omitempty"`
}

type KarbonClusterEtcdConfigResponse struct {
	Name         *string                      `json:"name" mapstructure:"name, omitempty"`
	Nodes        []*KarbonClusterNodeResponse `json:"nodes" mapstructure:"nodes, omitempty"`
	NumInstances *int64                       `json:"num_instances" mapstructure:"num_instances, omitempty"`
}

type KarbonClusterNodeResponse struct {
	Health         *string                                  `json:"health" mapstructure:"health, omitempty"`
	Name           *string                                  `json:"name" mapstructure:"name, omitempty"`
	NodePoolName   *string                                  `json:"node_pool_name" mapstructure:"node_pool_name, omitempty"`
	ResourceConfig *KarbonClusterNodeResourceConfigResponse `json:"resource_config" mapstructure:"resource_config, omitempty"`
	UUID           *string                                  `json:"uuid" mapstructure:"uuid, omitempty"`
}

type KarbonClusterNodeResourceConfigResponse struct {
	CPU       *int64  `json:"cpu" mapstructure:"cpu, omitempty"`
	DiskMib   *int64  `json:"disk_mib" mapstructure:"disk_mib, omitempty"`
	Image     *string `json:"image" mapstructure:"image, omitempty"`
	IPAddress *string `json:"ip_address" mapstructure:"ip_address, omitempty"`
	MemoryMib *int64  `json:"memory_mib" mapstructure:"memory_mib, omitempty"`
}

type KarbonClusterK8sConfigResponse struct {
	FQDN                  *string                            `json:"fqdn" mapstructure:"fqdn, omitempty"`
	MasterConfig          *KarbonClusterMasterConfigResponse `json:"master_config" mapstructure:"master_config, omitempty"`
	Masters               []*KarbonClusterNodeResponse       `json:"masters" mapstructure:"masters, omitempty"`
	NetworkCidr           *string                            `json:"network_cidr" mapstructure:"network_cidr, omitempty"`
	NetworkSubnetLength   *int64                             `json:"network_subnet_len" mapstructure:"network_subnet_len, omitempty"`
	OSFlavor              *string                            `json:"os_flavor" mapstructure:"os_flavor, omitempty"`
	ServiceClusterIPRange *string                            `json:"service_cluster_ip_range" mapstructure:"service_cluster_ip_range, omitempty"`
	Version               *string                            `json:"version" mapstructure:"version, omitempty"`
	Workers               []*KarbonClusterNodeResponse       `json:"workers" mapstructure:"workers, omitempty"`
}

type KarbonClusterMasterConfigResponse struct {
	DeploymentType *string `json:"deployment_type" mapstructure:"deployment_type, omitempty"`
	ExternalIP     *string `json:"external_ip" mapstructure:"external_ip, omitempty"`
}

//Inputs
type KarbonClusterIntentInput struct {
	Name               string                                     `json:"name" mapstructure:"name, omitempty"`
	Description        string                                     `json:"description" mapstructure:"description, omitempty"`
	VMNetwork          string                                     `json:"vm_network" mapstructure:"vm_network, omitempty"`
	K8sConfig          KarbonClusterK8sConfigIntentInput          `json:"k8s_config" mapstructure:"k8s_config, omitempty"`
	ClusterRef         string                                     `json:"cluster_ref" mapstructure:"cluster_ref, omitempty"`
	LoggingConfig      KarbonClusterLoggingConfigIntentInput      `json:"logging_config" mapstructure:"logging_config, omitempty"`
	StorageClassConfig KarbonClusterStorageClassConfigIntentInput `json:"storage_class_config" mapstructure:"storage_class_config, omitempty"`
	EtcdConfig         KarbonClusterEtcdConfigIntentInput         `json:"etcd_config" mapstructure:"etcd_config, omitempty"`
}

type KarbonClusterK8sConfigIntentInput struct {
	ServiceClusterIPRange string                         `json:"service_cluster_ip_range" mapstructure:"service_cluster_ip_range, omitempty"`
	NetworkCidr           string                         `json:"network_cidr" mapstructure:"network_cidr, omitempty"`
	FQDN                  string                         `json:"fqdn" mapstructure:"fqdn, omitempty"`
	Workers               []KarbonClusterNodeIntentInput `json:"workers" mapstructure:"workers, omitempty"`
	Masters               []KarbonClusterNodeIntentInput `json:"masters" mapstructure:"masters, omitempty"`
	OSFlavor              string                         `json:"os_flavor" mapstructure:"os_flavor, omitempty"`
	NetworkSubnetLength   int64                          `json:"network_subnet_len" mapstructure:"network_subnet_len, omitempty"`
	Version               string                         `json:"version" mapstructure:"version, omitempty"`
}

type KarbonClusterNodeIntentInput struct {
	Name           string                                     `json:"name" mapstructure:"name, omitempty"`
	NodePoolName   string                                     `json:"node_pool_name" mapstructure:"node_pool_name, omitempty"`
	ResourceConfig KarbonClusterNodeResourceConfigIntentInput `json:"resource_config" mapstructure:"resource_config, omitempty"`
	UUID           string                                     `json:"uuid" mapstructure:"uuid, omitempty"`
}

type KarbonClusterNodeResourceConfigIntentInput struct {
	CPU     int64  `json:"cpu" mapstructure:"cpu, omitempty"`
	DiskMib int64  `json:"disk_mib" mapstructure:"disk_mib, omitempty"`
	Image   string `json:"image" mapstructure:"image, omitempty"`
	// IPAddress string `json:"ip_address" mapstructure:"ip_address, omitempty"`
	MemoryMib int64 `json:"memory_mib" mapstructure:"memory_mib, omitempty"`
}

type KarbonClusterLoggingConfigIntentInput struct {
	EnableAppLogging bool `json:"enable_app_logging" mapstructure:"enable_app_logging, omitempty"`
}

type KarbonClusterStorageClassConfigIntentInput struct {
	Metadata KarbonClusterStorageClassConfigMetadataIntentInput `json:"metadata" mapstructure:"metadata, omitempty"`
	Spec     KarbonClusterStorageClassConfigSpecIntentInput     `json:"spec" mapstructure:"spec, omitempty"`
}

type KarbonClusterStorageClassConfigMetadataIntentInput struct {
	Name string `json:"name" mapstructure:"name, omitempty"`
}

type KarbonClusterStorageClassConfigSpecIntentInput struct {
	ReclaimPolicy string                                                `json:"reclaim_policy" mapstructure:"reclaim_policy, omitempty"`
	SCVolumeSpec  KarbonClusterStorageClassConfigVolumesSpecIntentInput `json:"sc_volumes_spec" mapstructure:"sc_volumes_spec, omitempty"`
}

type KarbonClusterStorageClassConfigVolumesSpecIntentInput struct {
	ClusterRef       string `json:"cluster_ref" mapstructure:"cluster_ref, omitempty"`
	User             string `json:"user" mapstructure:"user, omitempty"`
	Password         string `json:"password" mapstructure:"password, omitempty"`
	StorageContainer string `json:"storage_container" mapstructure:"storage_container, omitempty"`
	FileSystem       string `json:"file_system" mapstructure:"file_system, omitempty"`
	FlashMode        bool   `json:"flash_mode" mapstructure:"flash_mode, omitempty"`
}

type KarbonClusterEtcdConfigIntentInput struct {
	NumInstances int64                          `json:"num_instances" mapstructure:"num_instances, omitempty"`
	Name         string                         `json:"name" mapstructure:"name, omitempty"`
	Nodes        []KarbonClusterNodeIntentInput `json:"nodes" mapstructure:"nodes, omitempty"`
}

type KarbonClusterActionResponse struct {
	ClusterUUID string `json:"cluster_uuid" mapstructure:"cluster_uuid, omitempty"`
	TaskUUID    string `json:"task_uuid" mapstructure:"task_uuid, omitempty"`
}

type KarbonClusterKubeconfigResponse struct {
	ClusterUUID string `json:"cluster_uuid" mapstructure:"cluster_uuid, omitempty"`
	YmlConfig   string `json:"yml_config" mapstructure:"yml_config, omitempty"`
}

type KarbonClusterKubeconfig struct {
	APIVersion string `yaml:"apiVersion" mapstructure:"apiVersion, omitempty"`
	Kind       string `yaml:"kind" mapstructure:"kind, omitempty"`
	Clusters   []struct {
		Name    string `yaml:"name" mapstructure:"name, omitempty"`
		Cluster struct {
			Server                   string `yaml:"server" mapstructure:"server, omitempty"`
			CertificateAuthorityData string `yaml:"certificate-authority-data" mapstructure:"certificate-authority-data, omitempty"`
		} `yaml:"cluster" mapstructure:"cluster, omitempty"`
	} `yaml:"clusters" mapstructure:"clusters, omitempty"`
	Users []struct {
		Name string `yaml:"name" mapstructure:"name, omitempty"`
		User struct {
			Token string `yaml:"token" mapstructure:"token, omitempty"`
		} `yaml:"user" mapstructure:"user, omitempty"`
	} `yaml:"users" mapstructure:"users, omitempty"`
	Contexts []struct {
		Context struct {
			Cluster string `yaml:"cluster" mapstructure:"cluster, omitempty"`
			User    string `yaml:"user" mapstructure:"user, omitempty"`
		} `yaml:"context" mapstructure:"context, omitempty"`
		Name string `yaml:"name" mapstructure:"name, omitempty"`
	} `yaml:"contexts" mapstructure:"contexts, omitempty"`
	CurrentContext string `yaml:"current-context" mapstructure:"current-context, omitempty"`
}

type KarbonClusterScaleUpIntentInput struct {
	NodePoolName string `json:"node_pool_name" mapstructure:"node_pool_name, omitempty"`
	WorkerCount  int64  `json:"worker_count" mapstructure:"worker_count, omitempty"`
}
