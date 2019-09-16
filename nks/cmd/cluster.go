package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"reflect"
	"text/tabwriter"
	"io/ioutil"

	homedir "github.com/mitchellh/go-homedir"
	nks "github.com/NetApp/nks-sdk-go/nks"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getClusterId int
var getClusterAllf bool

var createClusterNamef string
var createClusterNumWorkers int
var createClusterWorkerSize string
var createClusterMasterSize string
var deleteClusterIDf int

type createClusterInputGCE struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	//WorkspaceID      int    `json:"workspace"`
	ProviderKey int `json:"provider_keyset"`

	MasterCount int    `json:"master_count"`
	MasterSize  string `json:"master_size"`
	//MasterRootDiskSize int    `json:"master_root_disk_size"`
	//MasterGPUInstanceSize string `json:"master_gpu_instance_size"`
	//MasterGPUCoreCount    int    `json:"master_gpu_core_count"`

	WorkerCount int    `json:"worker_count"`
	WorkerSize  string `json:"worker_size"`
	//WorkerGPUInstanceSize string `json:"worker_gpu_instance_size"`
	//WorkerGPUCoreCount    int    `json:"worker_gpu_core_count"`
	//WorkerRootDiskSize    int    `json:"worker_root_disk_size"`

	KubernetesVersion string `json:"k8s_version"`
	DashboardEnabled  bool   `json:"k8s_dashboard_enabled"`
	RbacEnabled       bool   `json:"k8s_rbac_enabled"`
	//K8sPodCIDR          string `json:"k8s_pod_cidr"`
	//K8sServiceCIDR      string `json:"k8s_service_cidr"`

	//ProjectID string `json:"project_id"`

	SSHKeySet int `json:"user_ssh_keyset"`

	EtcdType  string         `json:"etcd_type"`
	Platform  string         `json:"platform"`
	Channel   string         `json:"channel"`
	Region    string         `json:"region"`
	Solutions []nks.Solution `json:"solutions"`

	//Config string `json:"config"`

	//MinNodeCount int `json:"min_node_count"`
	//MaxNodeCount int `json:"max_node_count"`
	//Owner        int `json:"owner"`
	//ProviderSubnetID    string
	//ProviderSubnetCidr  string
	//ProviderNetworkID   string
	//ProviderNetworkCIDR string
}

var gceDefaults = map[string]interface{}{
	"Provider":          &CurrentConfig.Provider,
	"ProviderKey":       &CurrentConfig.ProviderKeySetID,
	"MasterCount":       1,
	"MasterSize":        "n1-standard-1",
	"WorkerCount":       2,
	"WorkerSize":        "n1-standard-1",
	"Region":            "us-west1-c",
	"KubernetesVersion": "v1.13.2",
	"RbacEnabled":       true,
	"DashboardEnabled":  true,
	"EtcdType":          "self_hosted",
	"Platform":          "coreos",
	"Channel":           "stable",
	"SSHKeySet":         &CurrentConfig.SSHKeySetId,
	"Solutions":         []nks.Solution{nks.Solution{Solution: "helm_tiller"}},
	//"ProviderSubnetID":    "__new__",
	//"ProviderSubnetCidr":  "172.23.1.0/24",
	//"ProviderNetworkID":   "__new__",
	//"ProviderNetworkCIDR": "172.23.0.0/16",
}

var clusterCmd = &cobra.Command{
	Use:   "clusters",
	Aliases: []string{"cl", "clus", "clu", "clusters"},
	Short: "manage cluster resources",
	Long:  ``,
	//Run: func(cmd *cobra.Command, args []string) {
	//	
	//},
}

var createClusterCmd = &cobra.Command{
	Use:   "create",
	Short: "deploy a new cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		template := generateClusterTemplate()
		template.Name = createClusterNamef
		if flagDebug {
			fmt.Printf("Template:\n%+v", template)
		}
		
		if createClusterMasterSize != "" {
			template.MasterSize = createClusterMasterSize
		}

		if createClusterNumWorkers != 2 {
			template.WorkerCount = createClusterNumWorkers
		}

		if createClusterWorkerSize != "" {
			template.WorkerSize = createClusterWorkerSize
		}

		if flagDebug {
			fmt.Printf("Template:\n \t%+v", template)
		} 

		newCluster, err := createCluster(template)
		check(err)
		printClusters([]nks.Cluster{newCluster})

		if CurrentConfig.ClusterId == 0 {
			setCluster(newCluster.ID)	
		}

	},
}

var deleteClusterCmd = &cobra.Command{
	Use:   "delete",
	Aliases: []string{"rm", "del"},
	Short: "delete a cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := deleteClusterByID(deleteClusterIDf)
		check(err)
	},
}

var listClustersCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "li"},
	Short:   "list clusters",
	Run: func(cmd *cobra.Command, args []string) {
		c := &([]nks.Cluster{})
		var err error

		if getClusterAllf {
			c, err = getAllClusters()
			if err != nil {
				fmt.Printf("There was an error retrieving items:\n\t%s\n\n", err)
			}
		} else {
			c, err = getClusters()
			if err != nil {
				fmt.Printf("There was an error retrieving items:\n\t%s\n\n", err)
			}
			
		}
		if len(*c) > 0 {
			if CurrentConfig.ClusterId == 0 {
				setCluster((*c)[0].ID)	
			}
		}
		
		printClusters(*c)
	},
}

var getClustersCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"cls", "clusters"},
	Short:   "get cluster details",
	Run: func(cmd *cobra.Command, args []string) {
		i := CurrentConfig.ClusterId
		if getClusterId != 0 {
			i = getClusterId
		}

		cl, err := getClusterByID(i)
		check(err)
			
		s := reflect.ValueOf(cl).Elem()
		typeOfT := s.Type()

		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			fmt.Printf("%s %s = %v\n",
				typeOfT.Field(i).Name, f.Type(), f.Interface())
		}
	},
}

func printClusters(cs []nks.Cluster) {
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fmt.Fprintf(w, "NAME\tID\tPROVIDER\tNODES\tK8s_VERSION\tSTATE\t\n")
	for _, c := range cs {
		if c.ID == CurrentConfig.ClusterId {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v (default)\t\n", c.Name, c.ID, c.Provider, c.NodeCount, c.KubernetesVersion, c.State)
		} else {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t\n", c.Name, c.ID, c.Provider, c.NodeCount, c.KubernetesVersion, c.State)
		}
		
	}
	w.Flush()
}

func getClusters() (*[]nks.Cluster, error) {
	o := CurrentConfig.OrgID

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

func setClusterKubeConfig(clusterId int)  {
	c := newClient()
	kubeConfig, err := c.GetKubeConfig(CurrentConfig.OrgID, CurrentConfig.ClusterId)
	if err != nil {
		fmt.Printf("There was an error retrieving config for cluster %d: \n\t%v\n", CurrentConfig.ClusterId, err)
	}
	home, err := homedir.Dir()
	//fmt.Printf("Setting kubeconfig to cluster %d", clusterId)
	b := []byte(kubeConfig)
	if err != nil {
		fmt.Printf("There was an error finding your home directory: %v", err)
	}


	err = ioutil.WriteFile(fmt.Sprintf("%s/.kube/config", home), b, 0644)
	if err != nil {
		fmt.Printf("There was an error writing the kubeconfig: %v", err)
	}
}

func deleteClusterByID(clusterId int) error {
	o := CurrentConfig.OrgID
	c := newClient()
	var err error

	if flagForce {
		err = c.ForceDeleteCluster(o, clusterId)
	} else {
		err = c.DeleteCluster(o, clusterId)
	}
	check(err)
	
	if clusterId == CurrentConfig.ClusterId {
		setCluster(0)	
	}

	cl, err := getAllClusters()
	check(err)
	printClusters(*cl)
	return nil
}

func createCluster(cl createClusterInputGCE) (nks.Cluster, error) {
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%d/clusters", CurrentConfig.OrgID)
	b, err := json.Marshal(cl)
	check(err)
	if flagDebug {
		fmt.Printf("\nSending request with body:\n%s\n", b)
	}

	res, err := httpRequestPost("POST", url, b)
	if flagDebug {
		fmt.Printf("Got response: \n%s", string(res))
	}
	check(err)

	data := nks.Cluster{}

	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Println(err)
	}

	return data, nil
}

func generateClusterTemplate() createClusterInputGCE {
	c := createClusterInputGCE{}
	mapstructure.Decode(gceDefaults, &c)
	return c
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(getClustersCmd)
	clusterCmd.AddCommand(createClusterCmd)
	clusterCmd.AddCommand(deleteClusterCmd)
	clusterCmd.AddCommand(listClustersCmd)


	listClustersCmd.Flags().IntVarP(&getClusterId, "id", "", 0, "ID of cluster")
	getClustersCmd.Flags().IntVarP(&getClusterId, "id", "", 0, "ID of cluster")
	getClustersCmd.Flags().BoolVarP(&getClusterAllf, "all", "a", false, "Get everything (incl. Service clusters)")

	createClusterCmd.Flags().StringVarP(&createClusterNamef, "name", "n", "", "ID of cluster")
	createClusterCmd.Flags().StringVarP(&createClusterMasterSize, "master-size", "", "", "Instance size of master nodes")

	createClusterCmd.Flags().StringVarP(&createClusterWorkerSize, "worker-size", "", "", "Instance size of worker nodes")
	createClusterCmd.Flags().IntVarP(&createClusterNumWorkers, "num-workers", "w", 2, "Number of worker nodes (default: 2)")
	e := createClusterCmd.MarkFlagRequired("name")
	check(e)


	deleteClusterCmd.Flags().BoolVarP(&flagForce, "force", "f", false, "ID of cluster to delete")
	deleteClusterCmd.Flags().IntVarP(&deleteClusterIDf, "id", "i", 0, "ID of cluster to delete")
	e = deleteClusterCmd.MarkFlagRequired("id")
	check(e)
}
