package karbon

import (
	"encoding/base64"
	"fmt"

	v3 "github.com/rancher/kontainer-engine-driver-karbon/client/v3"
	"github.com/rancher/kontainer-engine-driver-karbon/utils"
	"gopkg.in/yaml.v2"
)

type karbonManagerv20 struct {
	Client         *v3.Client
	KarbonClusters map[string]*v3.KarbonClusterIntentResponse
}

func (km karbonManagerv20) ScaleDownKarbonCluster(karbonClusterInfo KarbonClusterInfo, amountOfNodes int64) error {
	karbonCluster, err := km.GetKarbonCluster(karbonClusterInfo.UUID)
	if err != nil {
		return err
	}
	fmt.Printf("got cluster! scaling down")
	workers := make([]string, 0)
	for i := 0; i < int(amountOfNodes); i++ {
		workers := append(workers, *karbonCluster.K8sConfig.Workers[i].Name)
		karbonClusterActionResponseList, err := km.Client.V3.ScaleDownKarbonCluster(karbonClusterInfo.UUID, workers, amountOfNodes)
		fmt.Printf("looping")
		for _, r := range *karbonClusterActionResponseList {
			fmt.Printf("waiting for cluster")
			err = WaitForCluster(km.Client, r.TaskUUID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (km karbonManagerv20) GetKubernetesVersion(karbonClusterInfo KarbonClusterInfo) (string, error) {
	karbonCluster, err := km.GetKarbonCluster(karbonClusterInfo.UUID)
	if err != nil {
		return "", err
	}
	return *karbonCluster.K8sConfig.Version, nil
}

func (km karbonManagerv20) ScaleUpKarbonCluster(karbonClusterInfo KarbonClusterInfo, amountOfNodes int64) error {
	karbonCluster, err := km.GetKarbonCluster(karbonClusterInfo.UUID)
	if err != nil {
		return err
	}
	lenWorkers := len(karbonCluster.K8sConfig.Workers)
	if lenWorkers == 0 {
		return fmt.Errorf("Amount of worker nodes was %d but should be >0", lenWorkers)
	}
	karbonClusterActionResponse, err := km.Client.V3.ScaleUpKarbonCluster(
		karbonClusterInfo.UUID,
		&v3.KarbonClusterScaleUpIntentInput{
			NodePoolName: *karbonCluster.K8sConfig.Workers[0].NodePoolName,
			WorkerCount:  amountOfNodes,
		},
	)
	if err != nil {
		return err
	}
	err = WaitForCluster(km.Client, karbonClusterActionResponse.TaskUUID)
	if err != nil {
		return err
	}
	return nil
}

func (km karbonManagerv20) GetKarbonCluster(karbonClusterUUID string) (*v3.KarbonClusterIntentResponse, error) {
	if karbonCluster, ok := km.KarbonClusters[karbonClusterUUID]; ok {
		return karbonCluster, nil
	}
	karbonCluster, err := km.Client.V3.GetKarbonCluster(karbonClusterUUID)
	if err != nil {
		return nil, fmt.Errorf("Error occured getting karbon cluster")
	}
	if km.KarbonClusters == nil {
		km.KarbonClusters = map[string]*v3.KarbonClusterIntentResponse{
			karbonClusterUUID: karbonCluster,
		}
	} else {
		km.KarbonClusters[karbonClusterUUID] = karbonCluster
	}
	return karbonCluster, nil
}

func (km karbonManagerv20) GetAmountOfWorkerNodes(karbonClusterInfo KarbonClusterInfo) (int64, error) {
	karbonCluster, err := km.GetKarbonCluster(karbonClusterInfo.UUID)
	if err != nil {
		return 0, err
	}
	return int64(len(karbonCluster.K8sConfig.Workers)), nil
}

func (km karbonManagerv20) GetClient() *v3.Client {
	return km.Client
}

func (km karbonManagerv20) DeleteKarbonCluster(karbonClusterInfo KarbonClusterInfo) error {
	clusterDeleteResponse, err := km.Client.V3.DeleteKarbonCluster(karbonClusterInfo.UUID)
	if err != nil {
		return fmt.Errorf("Error occured getting karbon cluster")
	}
	err = WaitForCluster(km.Client, clusterDeleteResponse.TaskUUID)
	if err != nil {
		return fmt.Errorf("Failure waiting for cluster deletion")
	}
	return nil
}

func (km karbonManagerv20) RequestKarbonCluster(karbonClusterRequest *KarbonClusterRequest, WaitCompletion bool) (string, error) {
	fmt.Printf("karbonClusterRequest.Name")
	fmt.Printf(karbonClusterRequest.Name)
	workerNodes, err := km.GenerateNodeSlice(karbonClusterRequest.AmountOfWorkerNodes, karbonClusterRequest.ImageUUID, karbonClusterRequest.WorkerCPU, karbonClusterRequest.WorkerDiskMib, karbonClusterRequest.WorkerMemoryMib)
	if err != nil {
		fmt.Printf("Error RequestKarbonCluster when generating Node slice %s", err)
		return "", err
	}

	utils.PrintToJSON(workerNodes, "Workernodes ")
	karbon_cluster := &v3.KarbonClusterIntentInput{
		Name:        karbonClusterRequest.Name,
		Description: karbonClusterRequest.Description,
		VMNetwork:   karbonClusterRequest.VMNetworkUUID,
		K8sConfig: v3.KarbonClusterK8sConfigIntentInput{
			ServiceClusterIPRange: karbonClusterRequest.ServiceClusterIPRange,
			NetworkCidr:           karbonClusterRequest.NetworkCidr,
			Workers:               workerNodes,
			Masters: []v3.KarbonClusterNodeIntentInput{
				v3.KarbonClusterNodeIntentInput{
					ResourceConfig: v3.KarbonClusterNodeResourceConfigIntentInput{
						CPU:       karbonClusterRequest.MasterCPU,
						DiskMib:   karbonClusterRequest.MasterDiskMib,
						Image:     karbonClusterRequest.ImageUUID,
						MemoryMib: karbonClusterRequest.MasterMemoryMib,
					},
				},
			},
			OSFlavor:            karbonClusterRequest.OSFlavor,
			NetworkSubnetLength: karbonClusterRequest.NetworkSubnetLength,
			Version:             karbonClusterRequest.Version,
		},
		ClusterRef: karbonClusterRequest.ClusterUUID,
		LoggingConfig: v3.KarbonClusterLoggingConfigIntentInput{
			EnableAppLogging: true,
		},
		StorageClassConfig: v3.KarbonClusterStorageClassConfigIntentInput{
			Metadata: v3.KarbonClusterStorageClassConfigMetadataIntentInput{Name: "default-storageclass"},
			Spec: v3.KarbonClusterStorageClassConfigSpecIntentInput{
				ReclaimPolicy: karbonClusterRequest.ReclaimPolicy,
				SCVolumeSpec: v3.KarbonClusterStorageClassConfigVolumesSpecIntentInput{
					ClusterRef:       karbonClusterRequest.ClusterUUID,
					User:             karbonClusterRequest.ClusterUser,
					Password:         karbonClusterRequest.ClusterPassword,
					StorageContainer: karbonClusterRequest.StorageContainer,
					FileSystem:       karbonClusterRequest.FileSystem,
					FlashMode:        karbonClusterRequest.FlashMode,
				},
			},
		},
		EtcdConfig: v3.KarbonClusterEtcdConfigIntentInput{
			NumInstances: 1,
			Nodes: []v3.KarbonClusterNodeIntentInput{
				v3.KarbonClusterNodeIntentInput{
					ResourceConfig: v3.KarbonClusterNodeResourceConfigIntentInput{
						CPU:       karbonClusterRequest.EtcdCPU,
						DiskMib:   karbonClusterRequest.EtcdDiskMib,
						Image:     karbonClusterRequest.ImageUUID,
						MemoryMib: karbonClusterRequest.EtcdMemoryMib,
					},
				},
			},
		},
	}
	createClusterResponse, err := km.Client.V3.CreateKarbonCluster(karbon_cluster)
	if err != nil {
		return "", fmt.Errorf("Error occured during cluster creation:\n %s", err)
	}
	if createClusterResponse.TaskUUID == "" {
		return "", fmt.Errorf("Did not retrieve Task UUID exiting!")
	}
	if WaitCompletion == true {
		err = WaitForCluster(km.Client, createClusterResponse.TaskUUID)
		if err != nil {
			return "", err
		}
	}
	fmt.Printf("Cluster uuid: %s", createClusterResponse.ClusterUUID)
	fmt.Printf("Task uuid: %s", createClusterResponse.TaskUUID)
	return createClusterResponse.ClusterUUID, nil
}

func (km karbonManagerv20) GetKubeConfigForCluster(karbonClusterInfo KarbonClusterInfo) (*v3.KarbonClusterKubeconfig, error) {
	kubeconfig, _ := km.GetKubeConfigStringForCluster(km.Client, karbonClusterInfo.UUID)
	karbonClusterKubeconfig := v3.KarbonClusterKubeconfig{}
	err := yaml.Unmarshal([]byte(kubeconfig), &karbonClusterKubeconfig)
	if err != nil {
		return nil, err
	}
	utils.PrintToJSON(karbonClusterKubeconfig, "[karbonClusterKubeconfig]")
	return &karbonClusterKubeconfig, nil
}

func (km karbonManagerv20) GetKubeConfigStringForCluster(client *v3.Client, ClusterUUID string) (string, error) {
	kubeconfigResponse, err := client.V3.GetKubeConfigForKarbonCluster(ClusterUUID)
	if err != nil {
		return "", err
	}
	kubeconfig, _ := base64.StdEncoding.DecodeString(kubeconfigResponse.YmlConfig)
	return string(kubeconfig), nil
}

func (km karbonManagerv20) GenerateNodeSlice(AmountOfNodes int64, Image string, CPU int64, DiskMib int64, MemoryMib int64) ([]v3.KarbonClusterNodeIntentInput, error) {
	var nodeList []v3.KarbonClusterNodeIntentInput
	if AmountOfNodes < 1 {
		return nil, fmt.Errorf("Amount of Nodes must be >0")
	}
	for i := 0; i < int(AmountOfNodes); i++ {
		nodeList = append(nodeList, v3.KarbonClusterNodeIntentInput{
			ResourceConfig: v3.KarbonClusterNodeResourceConfigIntentInput{
				CPU:       CPU,
				DiskMib:   DiskMib,
				Image:     Image,
				MemoryMib: MemoryMib,
			},
		})
	}
	return nodeList, nil
}
