package karbon

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	// "github.com/rancher/kontainer-engine-driver-karbon/client"
	"github.com/rancher/kontainer-engine-driver-karbon/client"
	v3 "github.com/rancher/kontainer-engine-driver-karbon/client/v3"
	"github.com/rancher/kontainer-engine-driver-karbon/utils"
	"github.com/rancher/kontainer-engine/drivers/options"
	"github.com/rancher/kontainer-engine/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	raw "google.golang.org/api/container/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	runningStatus        = "RUNNING"
	defaultCredentialEnv = "GOOGLE_APPLICATION_CREDENTIALS"
	none                 = "none"
)

var EnvMutex sync.Mutex

// Driver defines the struct of gke driver
type Driver struct {
	driverCapabilities types.Capabilities
}

type state struct {
	Name                string
	Endpoint            string
	DisplayName         string
	Username            string
	Password            string
	Insecure            bool
	AmountOfWorkerNodes int64
	ClusterUUID         string
	KarbonClusterUUID   string
	VMNetworkUUID       string
	VMNetwork           string
	Image               string
	ImageUUID           string
	WorkerCPU           int64
	WorkerDiskMib       int64
	WorkerMemoryMib     int64
	MasterCPU           int64
	MasterDiskMib       int64
	MasterMemoryMib     int64
	EtcdCPU             int64
	EtcdDiskMib         int64
	EtcdMemoryMib       int64
	Version             string
	Cluster             string
	ReclaimPolicy       string
	ClusterUser         string
	ClusterPassword     string
	FileSystem          string
	StorageContainer    string
	KarbonVersion       string
	FlashMode           bool
	ClusterInfo         types.ClusterInfo
}

func NewDriver() types.Driver {
	driver := &Driver{
		driverCapabilities: types.Capabilities{
			Capabilities: make(map[int64]bool),
		},
	}

	driver.driverCapabilities.AddCapability(types.GetVersionCapability)
	driver.driverCapabilities.AddCapability(types.SetVersionCapability)
	driver.driverCapabilities.AddCapability(types.GetClusterSizeCapability)
	driver.driverCapabilities.AddCapability(types.SetClusterSizeCapability)
	return driver
}

// GetDriverCreateOptions implements driver interface
func (d *Driver) GetDriverCreateOptions(ctx context.Context) (*types.DriverFlags, error) {
	logrus.Infof("[DEBUG] GetDriverCreateOptions")
	driverFlag := types.DriverFlags{
		Options: make(map[string]*types.Flag),
	}
	driverFlag.Options["name"] = &types.Flag{
		Type:  types.StringType,
		Usage: "the internal name of the cluster in Rancher",
	}
	driverFlag.Options["endpoint"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Endpoint",
	}
	driverFlag.Options["username"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Username",
	}
	driverFlag.Options["password"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Password",
	}
	driverFlag.Options["display-name"] = &types.Flag{
		Type:  types.StringType,
		Usage: "the name of the cluster that should be displayed to the user",
	}
	driverFlag.Options["workernodes"] = &types.Flag{
		Type:  types.IntType,
		Usage: "Amount of worker nodes",
	}
	driverFlag.Options["insecure"] = &types.Flag{
		Type:  types.BoolType,
		Usage: "Insecure connection",
	}
	driverFlag.Options["image"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Karbon image to be used",
	}
	driverFlag.Options["version"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Kubernetes version to be used",
	}
	driverFlag.Options["cluster"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Nutanix cluster to be used",
	}
	driverFlag.Options["vmnetwork"] = &types.Flag{
		Type:  types.StringType,
		Usage: "VM network to be used",
	}
	driverFlag.Options["workercpu"] = &types.Flag{
		Type:  types.IntType,
		Usage: "Worker CPU",
	}
	driverFlag.Options["workermemorymib"] = &types.Flag{
		Type:  types.IntType,
		Usage: "Worker Memory mb",
	}
	driverFlag.Options["workerdiskmib"] = &types.Flag{
		Type:  types.IntType,
		Usage: "Worker Storage mib",
	}

	driverFlag.Options["mastercpu"] = &types.Flag{
		Type:  types.IntType,
		Usage: "master CPU",
	}
	driverFlag.Options["mastermemorymib"] = &types.Flag{
		Type:  types.IntType,
		Usage: "master Memory mb",
	}
	driverFlag.Options["masterdiskmib"] = &types.Flag{
		Type:  types.IntType,
		Usage: "master Storage mib",
	}

	driverFlag.Options["etcdcpu"] = &types.Flag{
		Type:  types.IntType,
		Usage: "etcd CPU",
	}
	driverFlag.Options["etcdmemorymib"] = &types.Flag{
		Type:  types.IntType,
		Usage: "etcd Memory mb",
	}
	driverFlag.Options["etcddiskmib"] = &types.Flag{
		Type:  types.IntType,
		Usage: "etcd Storage mib",
	}

	driverFlag.Options["clusteruser"] = &types.Flag{
		Type:  types.StringType,
		Usage: "PE user",
	}
	driverFlag.Options["clusterpassword"] = &types.Flag{
		Type:  types.StringType,
		Usage: "PE password",
	}
	driverFlag.Options["storagecontainer"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Storage container",
	}
	driverFlag.Options["filesystem"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Filesystem",
	}
	driverFlag.Options["reclaimpolicy"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Reclaim policy",
	}
	driverFlag.Options["flashmode"] = &types.Flag{
		Type:  types.BoolType,
		Usage: "Flash mode",
	}
	driverFlag.Options["karbonversion"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Karbon version",
	}

	return &driverFlag, nil
}

// GetDriverUpdateOptions implements driver interface
func (d *Driver) GetDriverUpdateOptions(ctx context.Context) (*types.DriverFlags, error) {
	logrus.Infof("[DEBUG] GetDriverUpdateOptions")
	driverFlag := types.DriverFlags{
		Options: make(map[string]*types.Flag),
	}
	driverFlag.Options["workernodes"] = &types.Flag{
		Type:  types.IntType,
		Usage: "The node number for your cluster to update. 0 means no updates",
	}

	return &driverFlag, nil
}

// Create implements driver interface
func (d *Driver) Create(ctx context.Context, opts *types.DriverOptions, _ *types.ClusterInfo) (*types.ClusterInfo, error) {

	utils.PrintToJSON(opts, "[DEBUG] Create OPTS: ")
	utils.PrintToJSON(ctx, "[DEBUG] Create ctx: ")
	state, err := getStateFromOpts(opts)
	if err != nil {
		return nil, err
	}
	utils.PrintToJSON(state, "[DEBUG] Create State: ")
	info := &types.ClusterInfo{}
	utils.PrintToJSON(info, "[DEBUG] Create Info: ")

	karbonManager, err := NewKarbonManager(
		client.Credentials{
			state.Endpoint,
			state.Username,
			state.Password,
			"",
			"",
			true,
			true,
			"",
		}, state.KarbonVersion)
	if err != nil {
		logrus.Debugf("[DEBUG] Error occured during Create after creating KarbonManager %v", err)
		return nil, err
	}
	err = UpdateStateWithUUIDs(karbonManager.GetClient(), &state)
	if err != nil {
		return nil, err
	}
	utils.PrintToJSON(state, "[DEBUG] Create State after searching for UIDs: ")
	//commented for testing purposes
	karbonClusterRequest := &KarbonClusterRequest{
		Name:                  state.DisplayName,
		Description:           state.DisplayName,
		VMNetworkUUID:         state.VMNetworkUUID,
		ServiceClusterIPRange: "172.19.0.0/16",
		NetworkCidr:           "172.20.0.0/16",
		NetworkSubnetLength:   24,
		OSFlavor:              "centos7.5.1804",
		Version:               state.Version,
		ImageUUID:             state.ImageUUID,
		Image:                 state.Image,
		AmountOfWorkerNodes:   state.AmountOfWorkerNodes,
		WorkerCPU:             state.WorkerCPU,
		WorkerDiskMib:         state.WorkerDiskMib,
		WorkerMemoryMib:       state.WorkerMemoryMib,
		MasterCPU:             state.MasterCPU,
		MasterDiskMib:         state.MasterDiskMib,
		MasterMemoryMib:       state.MasterMemoryMib,
		EtcdCPU:               state.EtcdCPU,
		EtcdDiskMib:           state.EtcdDiskMib,
		EtcdMemoryMib:         state.EtcdMemoryMib,
		ReclaimPolicy:         state.ReclaimPolicy,
		ClusterUUID:           state.ClusterUUID,
		ClusterUser:           state.ClusterUser,
		ClusterPassword:       state.ClusterPassword,
		StorageContainer:      state.StorageContainer,
		FileSystem:            state.FileSystem,
		FlashMode:             false,
	}
	KarbonClusterUUID, err := karbonManager.RequestKarbonCluster(karbonClusterRequest, true)
	if err != nil {
		return nil, err
	}
	// KarbonClusterUUID := "37aa1757-c049-4380-6ad4-167951e335dd"
	state.KarbonClusterUUID = KarbonClusterUUID

	err = storeState(info, state)
	if err != nil {
		logrus.Debugf("error storing state %v", err)
		return info, err
	}

	utils.PrintToJSON(info, "[DEBUG] Create END info: ")
	utils.PrintToJSON(state, "[DEBUG] Create END state: ")
	return info, nil
}

// Update implements driver interface
func (d *Driver) Update(ctx context.Context, info *types.ClusterInfo, opts *types.DriverOptions) (*types.ClusterInfo, error) {
	logrus.Infof("[DEBUG] Update")
	state, err := getState(info)
	if err != nil {
		return nil, err
	}
	newState, err := getStateFromOpts(opts)
	if err != nil {
		return nil, err
	}
	utils.PrintToJSON(ctx, "[DEBUG] Update ctx: ")
	utils.PrintToJSON(info, "[DEBUG] Update info: ")
	utils.PrintToJSON(opts, "[DEBUG] Update opts: ")
	utils.PrintToJSON(state, "[DEBUG] Update state: ")
	utils.PrintToJSON(newState, "[DEBUG] Update newState: ")
	karbonManager, err := NewKarbonManager(
		client.Credentials{
			state.Endpoint,
			state.Username,
			state.Password,
			"",
			"",
			true,
			true,
			"",
		}, state.KarbonVersion)
	// currentAmountOfWorkerNodes := state.AmountOfWorkerNodes
	// state.KarbonClusterUUID
	newAmountOfWorkerNodes := newState.AmountOfWorkerNodes
	karbonClusterInfo := KarbonClusterInfo{
		Name: state.DisplayName,
		UUID: state.KarbonClusterUUID,
	}
	currentAmountOfWorkerNodes, err := karbonManager.GetAmountOfWorkerNodes(karbonClusterInfo)
	if err != nil {
		return nil, err
	}

	logrus.Infof("[DEBUG] update currentAmountOfWorkerNodes %d", currentAmountOfWorkerNodes)
	logrus.Infof("[DEBUG] update newAmountOfWorkerNodes %d", newAmountOfWorkerNodes)
	if currentAmountOfWorkerNodes > newAmountOfWorkerNodes {
		amount := currentAmountOfWorkerNodes - newAmountOfWorkerNodes
		logrus.Infof("[DEBUG] scaling down by nodes %d", amount)
		err = karbonManager.ScaleDownKarbonCluster(karbonClusterInfo, amount)
		if err != nil {
			return nil, err
		}
	}
	if currentAmountOfWorkerNodes < newAmountOfWorkerNodes {
		amount := newAmountOfWorkerNodes - currentAmountOfWorkerNodes
		logrus.Infof("[DEBUG] scaling up by nodes %d", amount)
		err = karbonManager.ScaleUpKarbonCluster(karbonClusterInfo, amount)
		if err != nil {
			return nil, err
		}
	}

	return info, storeState(info, state)
}

func (d *Driver) PostCheck(ctx context.Context, info *types.ClusterInfo) (*types.ClusterInfo, error) {
	logrus.Infof("[DEBUG] PostCheck")
	state, err := getState(info)
	if err != nil {
		return nil, err
	}

	utils.PrintToJSON(state, "[DEBUG] PostCheckSTATE: ")
	karbonManager, err := NewKarbonManager(
		client.Credentials{
			state.Endpoint,
			state.Username,
			state.Password,
			"",
			"",
			true,
			true,
			"",
		}, state.KarbonVersion)
	karbonClusterInfo := KarbonClusterInfo{
		Name: state.DisplayName,
		UUID: state.KarbonClusterUUID,
	}
	kubeconfig, err := karbonManager.GetKubeConfigForCluster(karbonClusterInfo)
	if err != nil {
		return nil, err
	}

	amountOfWorkerNodes, err := karbonManager.GetAmountOfWorkerNodes(karbonClusterInfo)
	if err != nil {
		return nil, err
	}
	version, err := karbonManager.GetKubernetesVersion(karbonClusterInfo)
	if err != nil {
		return nil, err
	}
	info.Endpoint = kubeconfig.Clusters[0].Cluster.Server
	info.Version = version
	info.RootCaCertificate = kubeconfig.Clusters[0].Cluster.CertificateAuthorityData
	info.NodeCount = amountOfWorkerNodes
	serviceAccountToken, err := GenerateServiceAccountToken(kubeconfig)
	if err != nil {
		return nil, err
	}
	info.ServiceAccountToken = serviceAccountToken
	utils.PrintToJSON(info, "[DEBUG] CLUSTERINFO: ")
	return info, nil
}

// Remove implements driver interface
func (d *Driver) Remove(ctx context.Context, info *types.ClusterInfo) error {
	logrus.Infof("[DEBUG]remove")
	state, err := getState(info)
	if err != nil {
		logrus.Infof("[DEBUG]Remove Error occured: %s", err)
		return err
	}
	utils.PrintToJSON(info, "[DEBUG]Remove Info:	")
	utils.PrintToJSON(state, "[DEBUG] Remove STATE: ")
	karbonManager, err := NewKarbonManager(
		client.Credentials{
			state.Endpoint,
			state.Username,
			state.Password,
			"",
			"",
			true,
			true,
			"",
		}, state.KarbonVersion)
	karbonClusterInfo := KarbonClusterInfo{
		Name: state.DisplayName,
		UUID: state.KarbonClusterUUID,
	}
	logrus.Infof("[DEBUG]Deleting cluster ")
	karbonManager.DeleteKarbonCluster(karbonClusterInfo)
	logrus.Infof("[DEBUG]Done cluster ")
	return nil
}

func (d *Driver) GetVersion(ctx context.Context, info *types.ClusterInfo) (*types.KubernetesVersion, error) {
	logrus.Info("[DEBUG] GetVersion")
	cluster, err := d.getClusterStats(ctx, info)

	if err != nil {
		return nil, err
	}

	version := &types.KubernetesVersion{Version: cluster.CurrentMasterVersion}

	return version, nil
}

func (d *Driver) SetVersion(ctx context.Context, info *types.ClusterInfo, version *types.KubernetesVersion) error {
	logrus.Info("updating master version")

	err := d.updateAndWait(ctx, info, &raw.UpdateClusterRequest{
		Update: &raw.ClusterUpdate{
			DesiredMasterVersion: version.Version,
		}})

	if err != nil {
		return err
	}

	logrus.Info("master version updated successfully")
	logrus.Info("updating node version")

	err = d.updateAndWait(ctx, info, &raw.UpdateClusterRequest{
		Update: &raw.ClusterUpdate{
			DesiredNodeVersion: version.Version,
		},
	})

	if err != nil {
		return err
	}

	logrus.Info("node version updated successfully")

	return nil
}

func (d *Driver) GetClusterSize(ctx context.Context, info *types.ClusterInfo) (*types.NodeCount, error) {
	logrus.Info("[DEBUG] GetClusterSize")
	cluster, err := d.getClusterStats(ctx, info)

	if err != nil {
		return nil, err
	}

	version := &types.NodeCount{Count: int64(cluster.NodePools[0].InitialNodeCount)}

	return version, nil
}

func (d *Driver) SetClusterSize(ctx context.Context, info *types.ClusterInfo, count *types.NodeCount) error {

	logrus.Info("[DEBUG] updating cluster size")

	logrus.Info("[DEBUG] cluster size updated successfully")

	return nil
}

func (d *Driver) GetCapabilities(ctx context.Context) (*types.Capabilities, error) {
	logrus.Info("[DEBUG] GetCapabilities")
	return &d.driverCapabilities, nil
}

func (d *Driver) RemoveLegacyServiceAccount(ctx context.Context, info *types.ClusterInfo) error {
	logrus.Info("[DEBUG] RemoveLegacyServiceAccount")

	return nil
}

func (d *Driver) ETCDSave(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) error {
	return fmt.Errorf("ETCD backup operations are not implemented")
}

func (d *Driver) ETCDRestore(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) error {
	return fmt.Errorf("ETCD backup operations are not implemented")
}

func (d *Driver) GetK8SCapabilities(ctx context.Context, options *types.DriverOptions) (*types.K8SCapabilities, error) {
	logrus.Info("[DEBUG] GetK8SCapabilities")

	capabilities := &types.K8SCapabilities{
		L4LoadBalancer: &types.LoadBalancerCapabilities{
			Enabled:              true,
			Provider:             "GCLB",
			ProtocolsSupported:   []string{"TCP", "UDP"},
			HealthCheckSupported: true,
		},
	}
	return capabilities, nil
}

// SetDriverOptions implements driver interface
func getStateFromOpts(driverOptions *types.DriverOptions) (state, error) {
	utils.PrintToJSON(driverOptions, "[DEBUG] getStateFromOpts driverOptions:")
	d := state{}
	d.Name = options.GetValueFromDriverOptions(driverOptions, types.StringType, "name").(string)
	d.Endpoint = options.GetValueFromDriverOptions(driverOptions, types.StringType, "endpoint").(string)
	d.DisplayName = options.GetValueFromDriverOptions(driverOptions, types.StringType, "display-name", "displayName").(string)
	d.Username = options.GetValueFromDriverOptions(driverOptions, types.StringType, "username").(string)
	d.Password = options.GetValueFromDriverOptions(driverOptions, types.StringType, "password").(string)
	d.Insecure = options.GetValueFromDriverOptions(driverOptions, types.BoolType, "insecure").(bool)
	d.FlashMode = options.GetValueFromDriverOptions(driverOptions, types.BoolType, "flashmode").(bool)
	d.AmountOfWorkerNodes = options.GetValueFromDriverOptions(driverOptions, types.IntType, "workernodes").(int64)

	d.WorkerCPU = options.GetValueFromDriverOptions(driverOptions, types.IntType, "workercpu").(int64)
	d.WorkerDiskMib = options.GetValueFromDriverOptions(driverOptions, types.IntType, "workerdiskmib").(int64)
	d.WorkerMemoryMib = options.GetValueFromDriverOptions(driverOptions, types.IntType, "workermemorymib").(int64)
	d.MasterCPU = options.GetValueFromDriverOptions(driverOptions, types.IntType, "mastercpu").(int64)
	d.MasterDiskMib = options.GetValueFromDriverOptions(driverOptions, types.IntType, "masterdiskmib").(int64)
	d.MasterMemoryMib = options.GetValueFromDriverOptions(driverOptions, types.IntType, "mastermemorymib").(int64)
	d.EtcdCPU = options.GetValueFromDriverOptions(driverOptions, types.IntType, "etcdcpu").(int64)
	d.EtcdDiskMib = options.GetValueFromDriverOptions(driverOptions, types.IntType, "etcddiskmib").(int64)
	d.EtcdMemoryMib = options.GetValueFromDriverOptions(driverOptions, types.IntType, "etcdmemorymib").(int64)
	d.Version = options.GetValueFromDriverOptions(driverOptions, types.StringType, "version").(string)
	d.ReclaimPolicy = options.GetValueFromDriverOptions(driverOptions, types.StringType, "reclaimpolicy").(string)
	d.ClusterUser = options.GetValueFromDriverOptions(driverOptions, types.StringType, "clusteruser").(string)
	d.ClusterPassword = options.GetValueFromDriverOptions(driverOptions, types.StringType, "clusterpassword").(string)
	d.FileSystem = options.GetValueFromDriverOptions(driverOptions, types.StringType, "filesystem").(string)
	d.StorageContainer = options.GetValueFromDriverOptions(driverOptions, types.StringType, "storagecontainer").(string)
	d.VMNetwork = options.GetValueFromDriverOptions(driverOptions, types.StringType, "vmnetwork").(string)
	d.Image = options.GetValueFromDriverOptions(driverOptions, types.StringType, "image").(string)
	d.Cluster = options.GetValueFromDriverOptions(driverOptions, types.StringType, "cluster").(string)
	d.KarbonVersion = options.GetValueFromDriverOptions(driverOptions, types.StringType, "karbonversion").(string)

	utils.PrintToJSON(d, "[DEBUG] getStateFromOpts: ")
	return d, d.validate()
}

func (s *state) validate() error {
	logrus.Infof("[DEBUG] validate")

	if s.Name == "" {
		return fmt.Errorf("Karbon cluster name is required")
	}
	//Check endpoint
	if s.Endpoint == "" {
		return fmt.Errorf("Prism Central endpoint is required")
	}
	matchedEndpoint, err := regexp.MatchString("^.*:(\\d{2,4})$", s.Endpoint)
	if err != nil || matchedEndpoint == false {
		return fmt.Errorf("Endpoint must be formatted: IP/FQDN:Port")
	}

	if s.Username == "" {
		return fmt.Errorf("Username is required")
	}
	if s.Password == "" {
		return fmt.Errorf("Password is required")
	}

	if s.AmountOfWorkerNodes < 1 {
		return fmt.Errorf("AmountOfWorkerNodes must be >= 1")
	}
	if s.WorkerCPU < 1 {
		return fmt.Errorf("WorkerCPU must be >= 1")
	}
	if s.WorkerDiskMib < 1 {
		return fmt.Errorf("WorkerDiskMib must be >= 1")
	}
	if s.WorkerMemoryMib < 1 {
		return fmt.Errorf("WorkerMemoryMib must be >= 1")
	}
	if s.MasterCPU < 1 {
		return fmt.Errorf("MasterCPU must be >= 1")
	}
	if s.MasterDiskMib < 1 {
		return fmt.Errorf("MasterDiskMib must be >= 1")
	}
	if s.MasterMemoryMib < 1 {
		return fmt.Errorf("MasterMemoryMib must be >= 1")
	}
	if s.EtcdCPU < 1 {
		return fmt.Errorf("EtcdCPU must be >= 1")
	}
	if s.EtcdDiskMib < 1 {
		return fmt.Errorf("EtcdDiskMib must be >= 1")
	}
	if s.EtcdMemoryMib < 1 {
		return fmt.Errorf("EtcdMemoryMib must be >= 1")
	}
	if s.Version == "" {
		return fmt.Errorf("Version is required")
	}
	if s.ReclaimPolicy == "" {
		return fmt.Errorf("ReclaimPolicy is required")
	}
	if s.ReclaimPolicy != "Retain" && s.ReclaimPolicy != "Delete" {
		return fmt.Errorf("ReclaimPolicy must be Retain or Delete")
	}
	if s.ClusterUser == "" {
		return fmt.Errorf("ClusterUser is required")
	}
	if s.ClusterPassword == "" {
		return fmt.Errorf("ClusterPassword is required")
	}
	if s.FileSystem == "" {
		return fmt.Errorf("FileSystem is required")
	}
	if s.FileSystem != "xfs" && s.FileSystem != "ext4" {
		return fmt.Errorf("FileSystem is must be ext4 or xfs")
	}
	if s.StorageContainer == "" {
		return fmt.Errorf("StorageContainer is required")
	}
	if s.VMNetwork == "" {
		return fmt.Errorf("VMNetwork is required")
	}
	if s.Image == "" {
		return fmt.Errorf("Image is required")
	}
	if s.Cluster == "" {
		return fmt.Errorf("Cluster is required")
	}
	if s.KarbonVersion == "" {
		return fmt.Errorf("KarbonVersion is required")
	}
	if s.KarbonVersion != "2.0" && s.KarbonVersion != "2.1" {
		return fmt.Errorf("KarbonVersion must be 2.0 or 2.1")
	}

	return nil
}

func UpdateStateWithUUIDs(client *v3.Client, state *state) error {
	var err error
	state.VMNetworkUUID, err = FindSubnetByName(client, state.VMNetwork)
	if err != nil {
		return err
	}
	state.ClusterUUID, err = FindNutanixClusterByName(client, state.Cluster)
	if err != nil {
		return err
	}
	state.ImageUUID, err = FindImageByName(client, state.Image)
	if err != nil {
		return err
	}
	return nil
}

func storeState(info *types.ClusterInfo, state state) error {
	logrus.Infof("[DEBUG] storeState")
	bytes, err := json.Marshal(state)
	if err != nil {
		return err
	}
	if info.Metadata == nil {
		info.Metadata = map[string]string{}
	}
	info.Metadata["state"] = string(bytes)
	return nil
}

func getState(info *types.ClusterInfo) (state, error) {
	logrus.Infof("[DEBUG] getState")
	state := state{}
	// ignore error
	err := json.Unmarshal([]byte(info.Metadata["state"]), &state)
	return state, err
}

func getClientset(cluster *raw.Cluster) (kubernetes.Interface, error) {
	logrus.Infof("[DEBUG] getClientset")
	capem, err := base64.StdEncoding.DecodeString(cluster.MasterAuth.ClusterCaCertificate)
	if err != nil {
		return nil, err
	}
	host := cluster.Endpoint
	if !strings.HasPrefix(host, "https://") {
		host = fmt.Sprintf("https://%s", host)
	}
	// in here we have to use http basic auth otherwise we can't get the permission to create cluster role
	config := &rest.Config{
		Host: host,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: capem,
		},
		Username: cluster.MasterAuth.Username,
		Password: cluster.MasterAuth.Password,
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func (d *Driver) getClusterStats(ctx context.Context, info *types.ClusterInfo) (*raw.Cluster, error) {
	logrus.Infof("[DEBUG] getClusterStats")

	cluster := raw.Cluster{}

	return &cluster, nil
}

func (d *Driver) updateAndWait(ctx context.Context, info *types.ClusterInfo, updateRequest *raw.UpdateClusterRequest) error {
	logrus.Info("[DEBUG] updateAndWait")
	return nil
}

func FindNutanixClusterByName(client *v3.Client, clusterName string) (string, error) {
	filter := &v3.DSMetadata{}
	clusters, err := client.V3.ListCluster(filter)
	if err != nil {
		return "", nil
	}
	for _, c := range clusters.Entities {
		if strings.ToLower(clusterName) == strings.ToLower(*c.Spec.Name) {
			return *c.Metadata.UUID, nil
		}
	}
	return "", fmt.Errorf("Did not find UUID for cluster %s", clusterName)
}

func FindImageByName(client *v3.Client, imageName string) (string, error) {
	// filter := &v3.DSMetadata{}
	karbonPrefix := "karbon-" + imageName
	possibleImageNames := []string{
		karbonPrefix,
		imageName,
	}
	images, err := client.V3.ListAllImage("")
	if err != nil {
		return "", nil
	}
	for _, n := range possibleImageNames {
		for _, i := range images.Entities {
			if strings.ToLower(n) == strings.ToLower(*i.Spec.Name) {
				return *i.Metadata.UUID, nil
			}
		}
	}
	return "", fmt.Errorf("Did not find UUID for image %s", imageName)
}

func FindSubnetByName(client *v3.Client, subnetName string) (string, error) {

	subnets, err := client.V3.ListAllSubnet("")
	if err != nil {
		return "", nil
	}
	for _, i := range subnets.Entities {
		if strings.ToLower(subnetName) == strings.ToLower(*i.Spec.Name) {
			return *i.Metadata.UUID, nil
		}
	}
	return "", fmt.Errorf("Did not find UUID for subnet %s", subnetName)
}
