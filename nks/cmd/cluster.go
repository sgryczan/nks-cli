package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	nks "github.com/NetApp/nks-sdk-go/nks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ClusterConfig struct {
	AdminConf      string `json:"b64_admin_conf`
	InstallMethod  string `json:"install_method`
	CACrt          string `json:"b64_ca_crt"`
	CAKey          string `json:"b64_ca_key"`
	ProxyCert      string `json:"b64_front_proxy_ca_cert"`
	ProxyKey       string `json:"b64_front_proxy_ca_key"`
	SAPub          string `json:"b64_sa_pub"`
	SAKey          string `json:"b64_sa_key"`
	EtcDClientCert string `json:"b64_etcd_client_crt"`
	EtdDClientKey  string `json:"b64_etcd_client_key"`
	CASHA256Hash   string `json:"ca_sha256_hash"`
}
type Cluster struct {
	ID         int           `json:"pk"`
	Name       string        `json:"name"`
	Org        string        `json:"org"`
	Provider   string        `json:"provider"`
	Workspace  Workspace     `json:"workspace"`
	K8sVersion string        `json:"k8s_version"`
	NodeCount  int           `json:"node_count"`
	Config     ClusterConfig `json:"config"`
}

type createClusterInputGCE struct {
	Name             string `json:"name"`
	Provider         string `json:"provider"`
	WorkspaceID      int    `json:"workspace"`
	ProviderKeysetID string `json:"provider_keyset"`

	MasterCount           int    `json:"master_count"`
	MasterSize            string `json:"master_size"`
	MasterRootDiskSize    int    `json:"master_root_disk_size"`
	MasterGPUInstanceSize string `json:"master_gpu_instance_size"`
	MasterGPUCoreCount    int    `json:"master_gpu_core_count"`

	WorkerCount           int    `json:"worker_count"`
	WorkerSize            string `json:"worker_size"`
	WorkerGPUInstanceSize string `json:"worker_gpu_instance_size"`
	WorkerGPUCoreCount    int    `json:"worker_gpu_core_count"`
	WorkerRootDiskSize    int    `json:"worker_root_disk_size"`

	K8sVersion          string `json:"k8s_version`
	K8sDashboardEnabled bool   `json:"k8s_dashboard_enabled"`
	K8sRBACEnabled      bool   `json:"k8s_rbac_enabled"`
	K8sPodCIDR          string `json:"k8s_pod_cidr"`
	K8sServiceCIDR      string `json:"k8s_service_cidr"`

	ProjectID string `json:"project_id"`

	UserSSHKeyset string `json:"user_ssh_keyset"`

	EtcdType string `json:"etcd_type"`
	Platform string `json:"platform"`
	Channel  string `json:"channel"`
	Region   string `json:"region"`
	Zone     string `json:"zone"`

	Config string `json:"config"`

	MinNodeCount int `json:"min_node_count"`
	MaxNodeCount int `json:"max_node_count"`
	Owner        int `json:"owner"`
}

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "clusters",
	Short: "list, create, destroy clusters",
	Long:  ``,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("cluster called")
	//},
}

var clusterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "deploy a new cluster",
	Long:  ``,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("cluster called")
	//},
}

var getClustersCmd = &cobra.Command{
	Use:     "cluster",
	Aliases: []string{"cls", "clusters"},
	Short:   "list clusters",
	Run: func(cmd *cobra.Command, args []string) {
		if getClusterAllf {
			getAllClusters()
		} else if getClusterId != "" {
			i, err := strconv.Atoi(getClusterId)
			check(err)
			cl, err := getClusterByID(i)
			check(err)
			cls := []nks.Cluster{
				*cl,
			}
			printClusters(cls)
		} else {
			cs, err := getClusters()
			if err != nil {
				fmt.Printf("There was an error retrieving items:\n\t%s\n\n", err)
				cs = &[]nks.Cluster{}
			}
			printClusters(*cs)
		}
	},
}

func printClusters(cs []nks.Cluster) {
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fmt.Fprintf(w, "NAME\tID\tPROVIDER\tNODES\tK8s_VERSION\t\n")
	for _, c := range cs {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t\n", c.Name, c.ID, c.Provider, c.NodeCount, c.KubernetesVersion)
	}
	w.Flush()
}

func getClusters() (*[]nks.Cluster, error) {
	o, err := strconv.Atoi(viper.GetString("org_id"))
	check(err)
	c := newClient()
	cls, err := c.GetClusters(o)

	check(err)

	return &cls, err
}

func getAllClusters() (*[]nks.Cluster, error) {
	o, err := strconv.Atoi(viper.GetString("org_id"))
	check(err)
	c := newClient()
	cls, err := c.GetAllClusters(o)

	check(err)

	return &cls, err
}

func getClusterByID(clusterId int) (*nks.Cluster, error) {
	o, err := strconv.Atoi(viper.GetString("org_id"))
	check(err)
	c := newClient()
	cl, err := c.GetCluster(o, clusterId)

	check(err)

	return cl, err
}

func createCluster() (string, error) {
	orgID := viper.GetString("org_id")
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%s/clusters", orgID)
	res, err := httpRequest("POST", url)

	data := []Cluster{}

	_ = json.Unmarshal(res, &data)
	//check(err)

	return "", err
}

var getClusterId string
var getClusterAllf bool

func init() {
	getCmd.AddCommand(getClustersCmd)
	getClustersCmd.Flags().StringVarP(&getClusterId, "id", "", "", "ID of cluster")
	getClustersCmd.Flags().BoolVarP(&getClusterAllf, "all", "a", false, "Get everything (incl. Service clusters)")
}
