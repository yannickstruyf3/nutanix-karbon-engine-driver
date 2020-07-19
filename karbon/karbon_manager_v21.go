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
	workerNodePoolName := genUUID()
	etcdNodePoolName := genUUID()
	masterNodePoolName := genUUID()
	// workerNodePoolName := fmt.Sprintf("%s-%s", karbonClusterRequest.Name, "worker-nodepool")
	// etcdNodePoolName := fmt.Sprintf("%s-%s", karbonClusterRequest.Name, "etcd-nodepool")
	// masterNodePoolName := fmt.Sprintf("%s-%s", karbonClusterRequest.Name, "master-nodepool")
	karbon_cluster := &v3.KarbonCluster21IntentInput{
		Name:    karbonClusterRequest.Name,
		Version: karbonClusterRequest.Version,
		CNIConfig: v3.KarbonCluster21CNIConfigIntentInput{
			// FlannelConfig:    v3.KarbonCluster21FlannelConfigIntentInput{},
			NodeCIDRMaskSize: karbonClusterRequest.NetworkSubnetLength,
			PodIPv4CIDR:      karbonClusterRequest.NetworkCidr,
			ServiceIPv4CIDR:  karbonClusterRequest.ServiceClusterIPRange,
		},
		ETCDConfig: v3.KarbonCluster21ETCDConfigIntentInput{
			NodePools: []v3.KarbonCluster21NodePoolIntentInput{
				v3.KarbonCluster21NodePoolIntentInput{
					Name:          etcdNodePoolName,
					NodeOSVersion: nodeOSVersion,
					NumInstances:  karbonClusterRequest.AmountOfETCDNodes,
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
					Name:          masterNodePoolName,
					NodeOSVersion: nodeOSVersion,
					NumInstances:  karbonClusterRequest.AmountOfMasterNodes,
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
					Name:          workerNodePoolName,
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

	switch depl := strings.ToLower(karbonClusterRequest.Deployment); depl {
	case "production - active/active":
		if karbonClusterRequest.AmountOfMasterNodes < 2 || karbonClusterRequest.AmountOfMasterNodes > 5 {
			return "", fmt.Errorf("AmountOfMasterNodes must be between 2 and 5 when creating an Active/Active Karbon cluster")
		}
		if karbonClusterRequest.AmountOfETCDNodes != 3 && karbonClusterRequest.AmountOfETCDNodes != 5 {
			return "", fmt.Errorf("AmountOfETCDNodes must be 3 or 5 when creating an Active/Active Karbon cluster")
		}
		if karbonClusterRequest.MasterVIPIP == "" {
			return "", fmt.Errorf("MasterVIPIP must be set when creating an Active/Passive Karbon cluster")
		}
		masterNodeConfigList := make([]v3.KarbonCluster21MasterNodeMasterConfigIntentInput, 0)
		masterIPs, err := km.ParseMasterIP(karbonClusterRequest)
		if err != nil {
			return "", err
		}
		for i := 0; i < int(karbonClusterRequest.AmountOfMasterNodes); i++ {
			masterNodeConfigList = append(masterNodeConfigList, v3.KarbonCluster21MasterNodeMasterConfigIntentInput{
				IPv4Address:  masterIPs[i],
				NodePoolName: masterNodePoolName,
			})
		}
		karbon_cluster.MastersConfig.ExternalLBConfig = &v3.KarbonCluster21ExternalLBMasterConfigIntentInput{
			ExternalIPv4Address: karbonClusterRequest.MasterVIPIP,
			MasterNodesConfig:   masterNodeConfigList,
		}
	case "production - active/passive":
		if karbonClusterRequest.AmountOfMasterNodes != 2 {
			return "", fmt.Errorf("AmountOfMasterNodes must be 2 when creating an Active/Passive Karbon cluster")
		}
		if karbonClusterRequest.AmountOfETCDNodes != 3 && karbonClusterRequest.AmountOfETCDNodes != 5 {
			return "", fmt.Errorf("AmountOfETCDNodes must be 3 or 5 when creating an Active/Passive Karbon cluster")
		}
		if karbonClusterRequest.MasterVIPIP == "" {
			return "", fmt.Errorf("MasterVIPIP must be set when creating an Active/Passive Karbon cluster")
		}
		karbon_cluster.MastersConfig.ActivePassiveConfig = &v3.KarbonCluster21ActivePassiveMasterConfigIntentInput{
			ExternalIPv4Address: karbonClusterRequest.MasterVIPIP,
		}
	case "development":
		if karbonClusterRequest.AmountOfMasterNodes != 1 {
			return "", fmt.Errorf("AmountOfMasterNodes must be 1 when creating an Active/Passive Karbon cluster")
		}
		if karbonClusterRequest.AmountOfETCDNodes != 1 {
			return "", fmt.Errorf("AmountOfETCDNodes must be 1 when creating a Development Karbon cluster")
		}
		karbon_cluster.MastersConfig.SingleMasterConfig = &v3.KarbonCluster21SingleMasterConfigIntentInput{}
	default:
		return "", fmt.Errorf("Unsupported deployment type: %s.\n", depl)
	}

	if strings.ToLower(karbonClusterRequest.CNIProvider) != "flannel" && strings.ToLower(karbonClusterRequest.CNIProvider) != "calico" {
		return "", fmt.Errorf("CNIProvider must be Flannel or Calico")
	}
	if strings.ToLower(karbonClusterRequest.CNIProvider) == "calico" {
		karbon_cluster.CNIConfig.CalicoConfig = &v3.KarbonCluster21CalicoConfigIntentInput{
			IpPoolConfigs: []v3.KarbonCluster21CalicoConfigIpPoolConfigIntentInput{
				v3.KarbonCluster21CalicoConfigIpPoolConfigIntentInput{
					CIDR: karbonClusterRequest.NetworkCidr,
				},
			},
		}
	} else {
		karbon_cluster.CNIConfig.FlannelConfig = &v3.KarbonCluster21FlannelConfigIntentInput{}
	}

	utils.PrintToJSON(karbon_cluster, "[DEBUG karbon_cluster: ")
	createClusterResponse, err := km.Client.V3.CreateKarbonCluster21(karbon_cluster)
	if err != nil {
		return "", fmt.Errorf("Error occured during cluster creation:\n %s", err)
	}
	utils.PrintToJSON(createClusterResponse, "[DEBUG createClusterResponse: ")
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

func (km karbonManagerv21) ParseMasterIP(karbonClusterRequest *KarbonClusterRequest) ([]string, error) {
	masterIpList := make([]string, 0)
	if karbonClusterRequest.MasterIP1 != "" {
		masterIpList = append(masterIpList, karbonClusterRequest.MasterIP1)
	}
	if karbonClusterRequest.MasterIP2 != "" {
		masterIpList = append(masterIpList, karbonClusterRequest.MasterIP2)
	}
	if karbonClusterRequest.MasterIP3 != "" {
		masterIpList = append(masterIpList, karbonClusterRequest.MasterIP3)
	}
	if karbonClusterRequest.MasterIP4 != "" {
		masterIpList = append(masterIpList, karbonClusterRequest.MasterIP4)
	}
	if karbonClusterRequest.MasterIP5 != "" {
		masterIpList = append(masterIpList, karbonClusterRequest.MasterIP5)
	}
	if len(masterIpList) < int(karbonClusterRequest.AmountOfMasterNodes) {
		return nil, fmt.Errorf("Not enough Master IPs were passed!")
	}
	return masterIpList, nil
}
