package karbon

import (
	"github.com/rancher/kontainer-engine-driver-karbon/client"
	v3 "github.com/rancher/kontainer-engine-driver-karbon/client/v3"
)

type KarbonClusterRequest struct {
	Name                  string
	Description           string
	VMNetworkUUID         string
	ServiceClusterIPRange string
	NetworkCidr           string
	Image                 string
	ImageUUID             string
	AmountOfWorkerNodes   int64
	WorkerCPU             int64
	WorkerDiskMib         int64
	WorkerMemoryMib       int64
	MasterCPU             int64
	MasterDiskMib         int64
	MasterMemoryMib       int64
	EtcdCPU               int64
	EtcdDiskMib           int64
	EtcdMemoryMib         int64
	OSFlavor              string
	NetworkSubnetLength   int64
	Version               string
	ClusterUUID           string
	ReclaimPolicy         string
	ClusterUser           string
	ClusterPassword       string
	FileSystem            string
	StorageContainer      string
	FlashMode             bool
	CNIProvider           string
	Deployment            string
	AmountOfMasterNodes   int64
	AmountOfETCDNodes     int64
	MasterVIPIP           string
	MasterIP1             string
	MasterIP2             string
	MasterIP3             string
	MasterIP4             string
	MasterIP5             string
}

type KarbonClusterInfo struct {
	Name string
	UUID string
}

type KarbonManager interface {
	// GetKarbonCluster(karbonClusterInfo KarbonClusterInfo) (*v3.KarbonCluster20IntentResponse, error)
	GetClient() *v3.Client
	GetAmountOfWorkerNodes(karbonClusterInfo KarbonClusterInfo) (int64, error)
	DeleteKarbonCluster(karbonClusterInfo KarbonClusterInfo) error
	RequestKarbonCluster(karbonClusterRequest *KarbonClusterRequest, WaitCompletion bool) (string, error)
	GetKubeConfigForCluster(karbonClusterInfo KarbonClusterInfo) (*v3.KarbonClusterKubeconfig, error)
	ScaleDownKarbonCluster(karbonClusterInfo KarbonClusterInfo, amountOfNodes int64) error
	ScaleUpKarbonCluster(karbonClusterInfo KarbonClusterInfo, amountOfNodes int64) error
	GetKubernetesVersion(karbonClusterInfo KarbonClusterInfo) (string, error)
}

func NewKarbonManager(credentials client.Credentials) (KarbonManager, error) {
	client, err := v3.NewV3Client(credentials)
	if err != nil {
		return nil, err
	}
	return &karbonManagerv21{
		Client: client,
	}, nil
}
