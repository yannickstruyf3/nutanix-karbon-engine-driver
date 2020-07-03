package karbon

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	v3 "github.com/rancher/kontainer-engine-driver-karbon/client/v3"
	"github.com/rancher/kontainer-engine-driver-karbon/utils"
	"github.com/rancher/kontainer-engine/drivers/util"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KarbonClusterRequest struct {
	Name                  string
	Description           string
	VMNetworkUUID         string
	ServiceClusterIPRange string
	NetworkCidr           string
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
}

func GetKarbonCluster(client *v3.Client, karbonClusterUUID string) (*v3.KarbonClusterIntentResponse, error) {
	cluster, err := client.V3.GetKarbonCluster(karbonClusterUUID)
	if err != nil {
		return nil, fmt.Errorf("Error occured getting karbon cluster")
	}
	return cluster, nil
}

func DeleteKarbonCluster(client *v3.Client, karbonClusterUUID string) error {
	clusterDeleteResponse, err := client.V3.DeleteKarbonCluster(karbonClusterUUID)
	if err != nil {
		return fmt.Errorf("Error occured getting karbon cluster")
	}
	err = WaitForCluster(client, clusterDeleteResponse.TaskUUID)
	if err != nil {
		return fmt.Errorf("Failure waiting for cluster deletion")
	}
	return nil
}

func RequestKarbonCluster(client *v3.Client, karbonClusterRequest *KarbonClusterRequest, WaitCompletion bool) (string, error) {
	fmt.Printf("karbonClusterRequest.Name")
	fmt.Printf(karbonClusterRequest.Name)
	workerNodes, err := GenerateNodeSlice(karbonClusterRequest.AmountOfWorkerNodes, karbonClusterRequest.ImageUUID, karbonClusterRequest.WorkerCPU, karbonClusterRequest.WorkerDiskMib, karbonClusterRequest.WorkerMemoryMib)
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
	createClusterResponse, err := client.V3.CreateKarbonCluster(karbon_cluster)
	if err != nil {
		return "", fmt.Errorf("Error occured during cluster creation:\n %s", err)
	}
	if createClusterResponse.TaskUUID == "" {
		return "", fmt.Errorf("Did not retrieve Task UUID exiting!")
	}
	if WaitCompletion == true {
		err = WaitForCluster(client, createClusterResponse.TaskUUID)
		if err != nil {
			return "", err
		}
	}
	fmt.Printf("Cluster uuid: %s", createClusterResponse.ClusterUUID)
	fmt.Printf("Task uuid: %s", createClusterResponse.TaskUUID)
	return createClusterResponse.ClusterUUID, nil
}

func WaitForCluster(client *v3.Client, taskUUID string) error {
	log.Printf("Starting wait")
	sleepTime := 30
	var status string = "QUEUED"

	for status == "QUEUED" || status == "RUNNING" {
		time.Sleep(time.Duration(sleepTime) * time.Second)
		v, err := client.V3.GetTask(taskUUID)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "INVALID_UUID") {
				return fmt.Errorf("INVALID_UUID retrieved!")
			}
			return err
		}
		status = *v.Status
		log.Printf("Status: %s", status)
		if status == "INVALID_UUID" || status == "FAILED" {
			return fmt.Errorf("error_detail: %s, progress_message: %s", utils.StringValue(v.ErrorDetail), utils.StringValue(v.ProgressMessage))
		}

	}
	if status == "SUCCEEDED" {
		return nil
	}
	return fmt.Errorf("End state was NOT succeeded! %s", status)
}

func GetKubeConfigForCluster(client *v3.Client, ClusterUUID string) (*v3.KarbonClusterKubeconfig, error) {
	kubeconfig, _ := GetKubeConfigStringForCluster(client, ClusterUUID)
	karbonClusterKubeconfig := v3.KarbonClusterKubeconfig{}
	err := yaml.Unmarshal([]byte(kubeconfig), &karbonClusterKubeconfig)
	if err != nil {
		return nil, err
	}
	utils.PrintToJSON(karbonClusterKubeconfig, "[karbonClusterKubeconfig]")
	return &karbonClusterKubeconfig, nil
}

func GetKubeConfigStringForCluster(client *v3.Client, ClusterUUID string) (string, error) {
	kubeconfigResponse, err := client.V3.GetKubeConfigForKarbonCluster(ClusterUUID)
	if err != nil {
		return "", err
	}
	kubeconfig, _ := base64.StdEncoding.DecodeString(kubeconfigResponse.YmlConfig)
	return string(kubeconfig), nil
}

func ScaleUpKarbonCluster(client *v3.Client, karbonCluster *v3.KarbonClusterIntentResponse, amountOfNodes int64) error {

	lenWorkers := len(karbonCluster.K8sConfig.Workers)
	if lenWorkers == 0 {
		return fmt.Errorf("Amount of worker nodes was %d but should be >0", lenWorkers)
	}
	karbonClusterScaleUpIntentInput := v3.KarbonClusterScaleUpIntentInput{
		NodePoolName: *karbonCluster.K8sConfig.Workers[0].NodePoolName,
		WorkerCount:  amountOfNodes,
	}
	utils.PrintToJSON(karbonClusterScaleUpIntentInput, "test")
	karbonClusterActionResponse, err := client.V3.ScaleUpKarbonCluster(
		karbonCluster,
		&karbonClusterScaleUpIntentInput,
	)
	if err != nil {
		return err
	}
	err = WaitForCluster(client, karbonClusterActionResponse.TaskUUID)
	if err != nil {
		return err
	}
	return nil
}

func ScaleDownKarbonCluster(client *v3.Client, karbonCluster *v3.KarbonClusterIntentResponse, amountOfNodes int64) error {

	fmt.Printf("got cluster! scaling down")
	karbonClusterActionResponseList, err := client.V3.ScaleDownKarbonCluster(karbonCluster, amountOfNodes)
	fmt.Printf("looping")
	for _, r := range *karbonClusterActionResponseList {
		fmt.Printf("waiting for cluster")
		err = WaitForCluster(client, r.TaskUUID)
		if err != nil {
			return err
		}
	}

	return nil
}

func GenerateNodeSlice(AmountOfNodes int64, Image string, CPU int64, DiskMib int64, MemoryMib int64) ([]v3.KarbonClusterNodeIntentInput, error) {
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

func GetClientsetFromKubeconfig(kubeconfig *v3.KarbonClusterKubeconfig) (*kubernetes.Clientset, error) {
	capem, err := base64.StdEncoding.DecodeString(kubeconfig.Clusters[0].Cluster.CertificateAuthorityData)
	if err != nil {
		return nil, err
	}
	host := kubeconfig.Clusters[0].Cluster.Server
	if !strings.HasPrefix(host, "https://") {
		host = fmt.Sprintf("https://%s", host)
	}

	// in here we have to use http basic auth otherwise we can't get the permission to create cluster role
	config := &rest.Config{
		Host: host,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: capem,
		},
		BearerToken: kubeconfig.Users[0].User.Token,
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func GenerateServiceAccountToken(karbonClusterKubeconfig *v3.KarbonClusterKubeconfig) (string, error) {
	clientSet, err := GetClientsetFromKubeconfig(karbonClusterKubeconfig)
	if err != nil {
		return "", err
	}

	return util.GenerateServiceAccountToken(clientSet)
}
