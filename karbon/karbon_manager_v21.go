package karbon

import (
	"fmt"
	"strings"

	v3 "github.com/rancher/kontainer-engine-driver-karbon/client/v3"
	"github.com/rancher/kontainer-engine-driver-karbon/utils"
	"gopkg.in/yaml.v2"
)

type karbonManagerv21 struct {
	Client         *v3.Client
	KarbonClusters map[string]*v3.KarbonCluster21IntentResponse
}

func (km karbonManagerv21) GetKarbonCluster(karbonClusterName string) (*v3.KarbonCluster21IntentResponse, error) {
	if karbonCluster, ok := km.KarbonClusters[karbonClusterName]; ok {
		return karbonCluster, nil
	}
	karbonCluster, err := km.Client.V3.GetKarbonCluster21(karbonClusterName)
	if err != nil {
		return nil, fmt.Errorf("Error occured getting karbon cluster: %s", err)
	}
	if km.KarbonClusters == nil {
		km.KarbonClusters = map[string]*v3.KarbonCluster21IntentResponse{
			karbonClusterName: karbonCluster,
		}
	} else {
		km.KarbonClusters[karbonClusterName] = karbonCluster
	}
	return karbonCluster, nil
}

func (km karbonManagerv21) ScaleUpKarbonCluster(karbonClusterInfo KarbonClusterInfo, amountOfNodes int64) error {
	karbonCluster, err := km.GetKarbonCluster(karbonClusterInfo.Name)
	if err != nil {
		return err
	}
	karbonClusterActionResponse, err := km.Client.V3.ScaleUpKarbonCluster(
		karbonCluster.UUID,
		&v3.KarbonClusterScaleUpIntentInput{
			NodePoolName: karbonCluster.WorkerConfig.NodePools[0],
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

func (km karbonManagerv21) ScaleDownKarbonCluster(karbonClusterInfo KarbonClusterInfo, amountOfNodes int64) error {
	karbonCluster, err := km.GetKarbonCluster(karbonClusterInfo.Name)
	if err != nil {
		return err
	}
	fmt.Printf("got cluster! scaling down")
	nodePool, err := km.Client.V3.GetKarbonCluster21NodePool(karbonClusterInfo.Name, karbonCluster.WorkerConfig.NodePools[0])
	if err != nil {
		return err
	}
	workers := make([]string, 0)
	for i := 0; i < int(amountOfNodes); i++ {
		// workers := append(workers, "x")
		workers := append(workers, *nodePool.Nodes[i].Hostname)
		karbonClusterActionResponseList, err := km.Client.V3.ScaleDownKarbonCluster(karbonCluster.UUID, workers, amountOfNodes)
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
func (km karbonManagerv21) GetKubernetesVersion(karbonClusterInfo KarbonClusterInfo) (string, error) {
	karbonCluster, err := km.GetKarbonCluster(karbonClusterInfo.Name)
	if err != nil {
		return "", err
	}
	return karbonCluster.Version, nil
}

func (km karbonManagerv21) GetAmountOfWorkerNodes(karbonClusterInfo KarbonClusterInfo) (int64, error) {
	karbonCluster, err := km.GetKarbonCluster(karbonClusterInfo.Name)
	if err != nil {
		return 0, err
	}
	nodePool, err := km.Client.V3.GetKarbonCluster21NodePool(karbonClusterInfo.Name, karbonCluster.WorkerConfig.NodePools[0])
	if err != nil {
		return 0, err
	}
	return int64(len(nodePool.Nodes)), nil
}

func (km karbonManagerv21) GetClient() *v3.Client {
	return km.Client
}

func (km karbonManagerv21) DeleteKarbonCluster(karbonClusterInfo KarbonClusterInfo) error {
	clusterDeleteResponse, err := km.Client.V3.DeleteKarbonCluster21(karbonClusterInfo.Name)
	if err != nil {
		return fmt.Errorf("Error occured getting karbon cluster")
	}
	err = WaitForCluster(km.Client, clusterDeleteResponse.TaskUUID)
	if err != nil {
		return fmt.Errorf("Failure waiting for cluster deletion")
	}
	return nil
}

func (km karbonManagerv21) RequestKarbonCluster(karbonClusterRequest *KarbonClusterRequest, WaitCompletion bool) (string, error) {

	nodeOSVersion := strings.Replace(karbonClusterRequest.Image, "karbon-", "", -1)
	karbon_cluster := &v3.KarbonCluster21IntentInput{
		Name:    karbonClusterRequest.Name,
		Version: karbonClusterRequest.Version,
		CNIConfig: v3.KarbonCluster21CNIConfigIntentInput{
			FlannelConfig:    v3.KarbonCluster21FlannelConfigIntentInput{},
			NodeCIDRMaskSize: karbonClusterRequest.NetworkSubnetLength,
			PodIPv4CIDR:      karbonClusterRequest.NetworkCidr,
			ServiceIPv4CIDR:  karbonClusterRequest.ServiceClusterIPRange,
		},
		ETCDConfig: v3.KarbonCluster21ETCDConfigIntentInput{
			NodePools: []v3.KarbonCluster21NodePoolIntentInput{
				v3.KarbonCluster21NodePoolIntentInput{
					Name:          "etcd-nodepool",
					NodeOSVersion: nodeOSVersion,
					NumInstances:  1,
					AHVConfig: v3.KarbonCluster21NodePoolAHVConfigIntentInput{
						CPU:                     karbonClusterRequest.EtcdCPU,
						DiskMib:                 karbonClusterRequest.EtcdDiskMib,
						MemoryMib:               karbonClusterRequest.EtcdMemoryMib,
						NetworkUUID:             karbonClusterRequest.VMNetworkUUID,
						PrismElementClusterUUID: karbonClusterRequest.ClusterUUID,
					},
				},
			},
		},
		MastersConfig: v3.KarbonCluster21MasterConfigIntentInput{
			NodePools: []v3.KarbonCluster21NodePoolIntentInput{
				v3.KarbonCluster21NodePoolIntentInput{
					Name:          "master-nodepool",
					NodeOSVersion: nodeOSVersion,
					NumInstances:  1,
					AHVConfig: v3.KarbonCluster21NodePoolAHVConfigIntentInput{
						CPU:                     karbonClusterRequest.MasterCPU,
						DiskMib:                 karbonClusterRequest.MasterDiskMib,
						MemoryMib:               karbonClusterRequest.MasterMemoryMib,
						NetworkUUID:             karbonClusterRequest.VMNetworkUUID,
						PrismElementClusterUUID: karbonClusterRequest.ClusterUUID,
					},
				},
			},
		},
		Metadata: v3.KarbonCluster21MetadataIntentInput{
			APIVersion: "2.0.0",
		},
		StorageClassConfig: v3.KarbonCluster21StorageClassConfigIntentInput{
			DefaultStorageClass: true,
			Name:                "default-storageclass",
			ReclaimPolicy:       karbonClusterRequest.ReclaimPolicy,
			VolumesConfig: v3.KarbonCluster21VolumesConfigIntentInput{
				FileSystem:              karbonClusterRequest.FileSystem,
				FlashMode:               false,
				Password:                karbonClusterRequest.ClusterPassword,
				PrismElementClusterUUID: karbonClusterRequest.ClusterUUID,
				StorageContainer:        karbonClusterRequest.StorageContainer,
				Username:                karbonClusterRequest.ClusterUser,
			},
		},
		WorkersConfig: v3.KarbonCluster21WorkerConfigIntentInput{
			NodePools: []v3.KarbonCluster21NodePoolIntentInput{
				v3.KarbonCluster21NodePoolIntentInput{
					Name:          "worker-nodepool",
					NodeOSVersion: nodeOSVersion,
					NumInstances:  karbonClusterRequest.AmountOfWorkerNodes,
					AHVConfig: v3.KarbonCluster21NodePoolAHVConfigIntentInput{
						CPU:                     karbonClusterRequest.WorkerCPU,
						DiskMib:                 karbonClusterRequest.WorkerDiskMib,
						MemoryMib:               karbonClusterRequest.WorkerMemoryMib,
						NetworkUUID:             karbonClusterRequest.VMNetworkUUID,
						PrismElementClusterUUID: karbonClusterRequest.ClusterUUID,
					},
				},
			},
		},
	}
	createClusterResponse, err := km.Client.V3.CreateKarbonCluster21(karbon_cluster)
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

func (km karbonManagerv21) GetKubeConfigForCluster(karbonClusterInfo KarbonClusterInfo) (*v3.KarbonClusterKubeconfig, error) {
	kubeconfig, err := km.Client.V3.GetKubeConfigForKarbonCluster21(karbonClusterInfo.Name)
	if err != nil {
		return nil, err
	}
	karbonClusterKubeconfig := v3.KarbonClusterKubeconfig{}
	err = yaml.Unmarshal([]byte(kubeconfig.KubeConfig), &karbonClusterKubeconfig)
	if err != nil {
		return nil, err
	}
	utils.PrintToJSON(karbonClusterKubeconfig, "[karbonClusterKubeconfig]")
	return &karbonClusterKubeconfig, nil
}
