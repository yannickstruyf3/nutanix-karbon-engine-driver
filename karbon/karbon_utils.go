package karbon

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	v3 "github.com/rancher/kontainer-engine-driver-karbon/client/v3"
	"github.com/rancher/kontainer-engine-driver-karbon/utils"
	"github.com/rancher/kontainer-engine/drivers/util"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

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

func genUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
