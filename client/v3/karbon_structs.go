package v3

//KARBON v2.0 api
///Responses
type KarbonCluster20ListIntentResponse []KarbonCluster20ClusterMetadataIntentResponse

// single element in cluster/list not returning the same /cluster/uuid
type KarbonCluster20ClusterMetadataIntentResponse struct {
	KarbonClusterMetadataResponse *KarbonCluster20IntentResponse `json:"cluster_metadata" mapstructure:"cluster_metadata, omitempty"`
	TaskProgressMessage           *string                        `json:"task_progress_message" mapstructure:"task_progress_message, omitempty"`
	TaskProgressPercent           *int64                         `json:"task_progress_percent" mapstructure:"task_progress_percent, omitempty"`
	TaskStatus                    *int64                         `json:"task_status" mapstructure:"task_status, omitempty"`
	TaskType                      *string                        `json:"task_type" mapstructure:"task_type, omitempty"`
}

//return type for /cluster/uuid

type KarbonCluster20IntentResponse struct {
	AddonsConfig *KarbonCluster20AddonsConfigResponse `json:"addons_config" mapstructure:"addons_config, omitempty"`
	EtcdConfig   *KarbonCluster20EtcdConfigResponse   `json:"etcd_config" mapstructure:"etcd_config, omitempty"`
	K8sConfig    *KarbonCluster20K8sConfigResponse    `json:"k8s_config" mapstructure:"k8s_config, omitempty"`
	Name         *string                              `json:"name" mapstructure:"name, omitempty"`
	UUID         *string                              `json:"uuid" mapstructure:"uuid, omitempty"`
}

type KarbonCluster20AddonsConfigResponse struct {
	KarbonCluster20LoggingConfigResponse `json:"logging_config" mapstructure:"logging_config, omitempty"`
}

type KarbonCluster20LoggingConfigResponse struct {
	State         *string `json:"state" mapstructure:"state, omitempty"`
	StorageSizeMb *int64  `json:"storage_size_mib" mapstructure:"storage_size_mib, omitempty"`
	Version       *string `json:"version" mapstructure:"version, omitempty"`
}

type KarbonCluster20EtcdConfigResponse struct {
	Name         *string                        `json:"name" mapstructure:"name, omitempty"`
	Nodes        []*KarbonCluster20NodeResponse `json:"nodes" mapstructure:"nodes, omitempty"`
	NumInstances *int64                         `json:"num_instances" mapstructure:"num_instances, omitempty"`
}

type KarbonCluster20NodeResponse struct {
	Health         *string                                    `json:"health" mapstructure:"health, omitempty"`
	Name           *string                                    `json:"name" mapstructure:"name, omitempty"`
	NodePoolName   *string                                    `json:"node_pool_name" mapstructure:"node_pool_name, omitempty"`
	ResourceConfig *KarbonCluster20NodeResourceConfigResponse `json:"resource_config" mapstructure:"resource_config, omitempty"`
	UUID           *string                                    `json:"uuid" mapstructure:"uuid, omitempty"`
}

type KarbonCluster20NodeResourceConfigResponse struct {
	CPU       *int64  `json:"cpu" mapstructure:"cpu, omitempty"`
	DiskMib   *int64  `json:"disk_mib" mapstructure:"disk_mib, omitempty"`
	Image     *string `json:"image" mapstructure:"image, omitempty"`
	IPAddress *string `json:"ip_address" mapstructure:"ip_address, omitempty"`
	MemoryMib *int64  `json:"memory_mib" mapstructure:"memory_mib, omitempty"`
}

type KarbonCluster20K8sConfigResponse struct {
	FQDN                  *string                              `json:"fqdn" mapstructure:"fqdn, omitempty"`
	MasterConfig          *KarbonCluster20MasterConfigResponse `json:"master_config" mapstructure:"master_config, omitempty"`
	Masters               []*KarbonCluster20NodeResponse       `json:"masters" mapstructure:"masters, omitempty"`
	NetworkCidr           *string                              `json:"network_cidr" mapstructure:"network_cidr, omitempty"`
	NetworkSubnetLength   *int64                               `json:"network_subnet_len" mapstructure:"network_subnet_len, omitempty"`
	OSFlavor              *string                              `json:"os_flavor" mapstructure:"os_flavor, omitempty"`
	ServiceClusterIPRange *string                              `json:"service_cluster_ip_range" mapstructure:"service_cluster_ip_range, omitempty"`
	Version               *string                              `json:"version" mapstructure:"version, omitempty"`
	Workers               []*KarbonCluster20NodeResponse       `json:"workers" mapstructure:"workers, omitempty"`
}

type KarbonCluster20MasterConfigResponse struct {
	DeploymentType *string `json:"deployment_type" mapstructure:"deployment_type, omitempty"`
	ExternalIP     *string `json:"external_ip" mapstructure:"external_ip, omitempty"`
}

//Inputs
type KarbonCluster20IntentInput struct {
	Name               string                                       `json:"name" mapstructure:"name, omitempty"`
	Description        string                                       `json:"description" mapstructure:"description, omitempty"`
	VMNetwork          string                                       `json:"vm_network" mapstructure:"vm_network, omitempty"`
	K8sConfig          KarbonCluster20K8sConfigIntentInput          `json:"k8s_config" mapstructure:"k8s_config, omitempty"`
	ClusterRef         string                                       `json:"cluster_ref" mapstructure:"cluster_ref, omitempty"`
	LoggingConfig      KarbonCluster20LoggingConfigIntentInput      `json:"logging_config" mapstructure:"logging_config, omitempty"`
	StorageClassConfig KarbonCluster20StorageClassConfigIntentInput `json:"storage_class_config" mapstructure:"storage_class_config, omitempty"`
	EtcdConfig         KarbonCluster20EtcdConfigIntentInput         `json:"etcd_config" mapstructure:"etcd_config, omitempty"`
}

type KarbonCluster20K8sConfigIntentInput struct {
	ServiceClusterIPRange string                           `json:"service_cluster_ip_range" mapstructure:"service_cluster_ip_range, omitempty"`
	NetworkCidr           string                           `json:"network_cidr" mapstructure:"network_cidr, omitempty"`
	FQDN                  string                           `json:"fqdn" mapstructure:"fqdn, omitempty"`
	Workers               []KarbonCluster20NodeIntentInput `json:"workers" mapstructure:"workers, omitempty"`
	Masters               []KarbonCluster20NodeIntentInput `json:"masters" mapstructure:"masters, omitempty"`
	OSFlavor              string                           `json:"os_flavor" mapstructure:"os_flavor, omitempty"`
	NetworkSubnetLength   int64                            `json:"network_subnet_len" mapstructure:"network_subnet_len, omitempty"`
	Version               string                           `json:"version" mapstructure:"version, omitempty"`
}

type KarbonCluster20NodeIntentInput struct {
	Name           string                                       `json:"name" mapstructure:"name, omitempty"`
	NodePoolName   string                                       `json:"node_pool_name" mapstructure:"node_pool_name, omitempty"`
	ResourceConfig KarbonCluster20NodeResourceConfigIntentInput `json:"resource_config" mapstructure:"resource_config, omitempty"`
	UUID           string                                       `json:"uuid" mapstructure:"uuid, omitempty"`
}

type KarbonCluster20NodeResourceConfigIntentInput struct {
	CPU     int64  `json:"cpu" mapstructure:"cpu, omitempty"`
	DiskMib int64  `json:"disk_mib" mapstructure:"disk_mib, omitempty"`
	Image   string `json:"image" mapstructure:"image, omitempty"`
	// IPAddress string `json:"ip_address" mapstructure:"ip_address, omitempty"`
	MemoryMib int64 `json:"memory_mib" mapstructure:"memory_mib, omitempty"`
}

type KarbonCluster20LoggingConfigIntentInput struct {
	EnableAppLogging bool `json:"enable_app_logging" mapstructure:"enable_app_logging, omitempty"`
}

type KarbonCluster20StorageClassConfigIntentInput struct {
	Metadata KarbonCluster20StorageClassConfigMetadataIntentInput `json:"metadata" mapstructure:"metadata, omitempty"`
	Spec     KarbonCluster20StorageClassConfigSpecIntentInput     `json:"spec" mapstructure:"spec, omitempty"`
}

type KarbonCluster20StorageClassConfigMetadataIntentInput struct {
	Name string `json:"name" mapstructure:"name, omitempty"`
}

type KarbonCluster20StorageClassConfigSpecIntentInput struct {
	ReclaimPolicy string                                                  `json:"reclaim_policy" mapstructure:"reclaim_policy, omitempty"`
	SCVolumeSpec  KarbonCluster20StorageClassConfigVolumesSpecIntentInput `json:"sc_volumes_spec" mapstructure:"sc_volumes_spec, omitempty"`
}

type KarbonCluster20StorageClassConfigVolumesSpecIntentInput struct {
	ClusterRef       string `json:"cluster_ref" mapstructure:"cluster_ref, omitempty"`
	User             string `json:"user" mapstructure:"user, omitempty"`
	Password         string `json:"password" mapstructure:"password, omitempty"`
	StorageContainer string `json:"storage_container" mapstructure:"storage_container, omitempty"`
	FileSystem       string `json:"file_system" mapstructure:"file_system, omitempty"`
	FlashMode        bool   `json:"flash_mode" mapstructure:"flash_mode, omitempty"`
}

type KarbonCluster20EtcdConfigIntentInput struct {
	NumInstances int64                            `json:"num_instances" mapstructure:"num_instances, omitempty"`
	Name         string                           `json:"name" mapstructure:"name, omitempty"`
	Nodes        []KarbonCluster20NodeIntentInput `json:"nodes" mapstructure:"nodes, omitempty"`
}

//KARBON 2.1

type KarbonCluster21ListIntentResponse []KarbonCluster21IntentResponse
type KarbonCluster21IntentResponse struct {
	Name                     string `json:"name" mapstructure:"name, omitempty"`
	UUID                     string `json:"uuid" mapstructure:"uuid, omitempty"`
	Status                   string `json:"status" mapstructure:"status, omitempty"`
	Version                  string `json:"version" mapstructure:"version, omitempty"`
	KubeApiServerIPv4Address string `json:"kubeapi_server_ipv4_address" mapstructure:"kubeapi_server_ipv4_address, omitempty"`
	ETCDConfig               struct {
		NodePools []string `json:"node_pools" mapstructure:"node_pools, omitempty"`
	} `json:"etcd_config" mapstructure:"etcd_config, omitempty"`
	MasterConfig struct {
		DeploymentType string   `json:"deployment_type" mapstructure:"deployment_type, omitempty"`
		NodePools      []string `json:"node_pools" mapstructure:"node_pools, omitempty"`
	} `json:"master_config" mapstructure:"master_config, omitempty"`
	WorkerConfig struct {
		NodePools []string `json:"node_pools" mapstructure:"node_pools, omitempty"`
	} `json:"worker_config" mapstructure:"worker_config, omitempty"`
}

type KarbonCluster21NodePoolIntentResponse struct {
	AHVConfig     KarbonCluster21NodePoolAHVConfigIntentResponse `json:"ahv_config" mapstructure:"ahv_config, omitempty"`
	Name          string                                         `json:"name" mapstructure:"name, omitempty"`
	NodeOSVersion string                                         `json:"node_os_version" mapstructure:"node_os_version, omitempty"`
	NumInstances  int64                                          `json:"num_instances" mapstructure:"num_instances, omitempty"`
	Nodes         []KarbonCluster21NodeIntentResponse            `json:"nodes" mapstructure:"nodes, omitempty"`
}

type KarbonCluster21NodePoolAHVConfigIntentResponse struct {
	CPU                     *int64  `json:"cpu" mapstructure:"cpu, omitempty"`
	DiskMib                 *int64  `json:"disk_mib" mapstructure:"disk_mib, omitempty"`
	MemoryMib               *int64  `json:"memory_mib" mapstructure:"memory_mib, omitempty"`
	NetworkUUID             *string `json:"network_uuid" mapstructure:"network_uuid, omitempty"`
	PrismElementClusterUUID *string `json:"prism_element_cluster_uuid" mapstructure:"prism_element_cluster_uuid, omitempty"`
}

type KarbonCluster21NodeIntentResponse struct {
	Hostname    *string `json:"hostname" mapstructure:"hostname, omitempty"`
	IPv4Address *string `json:"ipv4_address" mapstructure:"ipv4_address, omitempty"`
}

type KarbonCluster21KubeconfigResponse struct {
	KubeConfig string `json:"kube_config" mapstructure:"kube_config, omitempty"`
}

//inputs
type KarbonCluster21IntentInput struct {
	Name               string                                       `json:"name" mapstructure:"name, omitempty"`
	Version            string                                       `json:"version" mapstructure:"version, omitempty"`
	CNIConfig          KarbonCluster21CNIConfigIntentInput          `json:"cni_config" mapstructure:"cni_config, omitempty"`
	ETCDConfig         KarbonCluster21ETCDConfigIntentInput         `json:"etcd_config" mapstructure:"etcd_config, omitempty"`
	MastersConfig      KarbonCluster21MasterConfigIntentInput       `json:"masters_config" mapstructure:"masters_config, omitempty"`
	Metadata           KarbonCluster21MetadataIntentInput           `json:"metadata" mapstructure:"metadata, omitempty"`
	StorageClassConfig KarbonCluster21StorageClassConfigIntentInput `json:"storage_class_config" mapstructure:"storage_class_config, omitempty"`
	WorkersConfig      KarbonCluster21WorkerConfigIntentInput       `json:"workers_config" mapstructure:"workers_config, omitempty"`
}
type KarbonCluster21MetadataIntentInput struct {
	APIVersion string `json:"api_version" mapstructure:"api_version, omitempty"`
}

type KarbonCluster21MasterConfigIntentInput struct {
	SingleMasterConfig KarbonCluster21SingleMasterConfigIntentInput `json:"single_master_config" mapstructure:"single_master_config, omitempty"`
	NodePools          []KarbonCluster21NodePoolIntentInput         `json:"node_pools" mapstructure:"node_pools, omitempty"`
}

type KarbonCluster21SingleMasterConfigIntentInput struct {
}
type KarbonCluster21WorkerConfigIntentInput struct {
	NodePools []KarbonCluster21NodePoolIntentInput `json:"node_pools" mapstructure:"node_pools, omitempty"`
}
type KarbonCluster21ETCDConfigIntentInput struct {
	NodePools []KarbonCluster21NodePoolIntentInput `json:"node_pools" mapstructure:"node_pools, omitempty"`
}

type KarbonCluster21CNIConfigIntentInput struct {
	NodeCIDRMaskSize int64                                   `json:"node_cidr_mask_size" mapstructure:"node_cidr_mask_size, omitempty"`
	PodIPv4CIDR      string                                  `json:"pod_ipv4_cidr" mapstructure:"pod_ipv4_cidr, omitempty"`
	ServiceIPv4CIDR  string                                  `json:"service_ipv4_cidr" mapstructure:"service_ipv4_cidr, omitempty"`
	FlannelConfig    *KarbonCluster21FlannelConfigIntentInput `json:"flannel_config" mapstructure:"flannel_config, omitempty"`
	CalicoConfig     *KarbonCluster21CalicoConfigIntentInput   `json:"calico_config" mapstructure:"calico_config, omitempty"`
}

type KarbonCluster21CalicoConfigIntentInput struct{
	IpPoolConfigs []KarbonCluster21CalicoConfigIpPoolConfigIntentInput `json:"ip_pool_configs" mapstructure:"ip_pool_configs,omitempty"`
}

type KarbonCluster21CalicoConfigIpPoolConfigIntentInput struct{
	CIDR string  `json:"cidr" mapstructure:"cidr"`
}

type KarbonCluster21FlannelConfigIntentInput struct{}

type KarbonCluster21NodePoolIntentInput struct {
	AHVConfig     KarbonCluster21NodePoolAHVConfigIntentInput `json:"ahv_config" mapstructure:"ahv_config, omitempty"`
	Name          string                                      `json:"name" mapstructure:"name, omitempty"`
	NodeOSVersion string                                      `json:"node_os_version" mapstructure:"node_os_version, omitempty"`
	NumInstances  int64                                       `json:"num_instances" mapstructure:"num_instances, omitempty"`
}

type KarbonCluster21NodePoolAHVConfigIntentInput struct {
	CPU                     int64  `json:"cpu" mapstructure:"cpu, omitempty"`
	DiskMib                 int64  `json:"disk_mib" mapstructure:"disk_mib, omitempty"`
	MemoryMib               int64  `json:"memory_mib" mapstructure:"memory_mib, omitempty"`
	NetworkUUID             string `json:"network_uuid" mapstructure:"network_uuid, omitempty"`
	PrismElementClusterUUID string `json:"prism_element_cluster_uuid" mapstructure:"prism_element_cluster_uuid, omitempty"`
}

type KarbonCluster21StorageClassConfigIntentInput struct {
	DefaultStorageClass bool                                    `json:"default_storage_class" mapstructure:"default_storage_class, omitempty"`
	Name                string                                  `json:"name" mapstructure:"name, omitempty"`
	ReclaimPolicy       string                                  `json:"reclaim_policy" mapstructure:"reclaim_policy, omitempty"`
	VolumesConfig       KarbonCluster21VolumesConfigIntentInput `json:"volumes_config" mapstructure:"volumes_config, omitempty"`
}

type KarbonCluster21VolumesConfigIntentInput struct {
	FileSystem              string `json:"file_system" mapstructure:"file_system, omitempty"`
	FlashMode               bool   `json:"flash_mode" mapstructure:"flash_mode, omitempty"`
	Password                string `json:"password" mapstructure:"password, omitempty"`
	PrismElementClusterUUID string `json:"prism_element_cluster_uuid" mapstructure:"prism_element_cluster_uuid, omitempty"`
	StorageContainer        string `json:"storage_container" mapstructure:"storage_container, omitempty"`
	Username                string `json:"username" mapstructure:"username, omitempty"`
}

//KARBON shared

type KarbonClusterActionResponse struct {
	ClusterName string `json:"cluster_name" mapstructure:"cluster_name, omitempty"`
	ClusterUUID string `json:"cluster_uuid" mapstructure:"cluster_uuid, omitempty"`
	TaskUUID    string `json:"task_uuid" mapstructure:"task_uuid, omitempty"`
}

type KarbonCluster20KubeconfigResponse struct {
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
