package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"text/tabwriter"

	nks "github.com/NetApp/nks-sdk-go/nks"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	models "gitlab.com/sgryczan/nks-cli/nks/models"
)

var getClusterId int
var flagListAllClusters bool

var flagClusterName 			string
var flagProviderName			string
var createClusterNumWorkers 	int
var createClusterWorkerSize 	string
var createClusterMasterSize 	string
var deleteClusterIDf 			int


var gceDefaults = map[string]interface{}{
	// Name is implied
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

var hciDefaults = map[string]interface{}{
	// Name is implied
	"Provider":          		"hci",
	"ProviderKey":       		63207,
	"Workspace":				22022,

	"MasterCount":       		1,
	"MasterSize":        		"m",
	"MasterRootDiskSize": 		50,
	"MasterGPUInstanceSize":	"",
	"MasterGPUCoreCount":		nil,

	"WorkerCount":       		2,
	"WorkerSize":        		"m",
	"WorkerGPUInstanceSize":	"",
	"WorkerGPUCoreCount":		nil,
	"WorkerRootDiskSize":		50,

	"Region":            		"LAB-RTP",

	"KubernetesVersion": 		"v1.14.3",
	"RbacEnabled":       		true,
	"DashboardEnabled":  		true,
	"EtcdType":          		"classic",
	"Platform":          		"debian",
	"Channel":           		"stable",
	"Zone":						"",
	"Config":					map[string]bool{"enable_experimental_features": true},
	"SSHKeySet":         		&CurrentConfig.SSHKeySetId,
	"Solutions":         		[]nks.Solution{nks.Solution{Solution: "helm_tiller"}},
	"Features":					[]string{},
	"MinNodeCount":				nil,
	"MaxNodeCount":				nil,
	"Owner":					26309, // ???
	//"ProviderSubnetID":    "__new__",
	//"ProviderSubnetCidr":  "172.23.1.0/24",
	//"ProviderNetworkID":   "__new__",
	//"ProviderNetworkCIDR": "172.23.0.0/16",
}

var clusterCmd = &cobra.Command{
	Use:     "clusters",
	Aliases: []string{"cl", "clus", "clu", "clusters"},
	Short:   "manage cluster resources",
	Long:    ``,
	//Run: func(cmd *cobra.Command, args []string) {
	//
	//},
}

var createClusterCmd = &cobra.Command{
	Use:   "create",
	Short: "deploy a new cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		provider := "gce"
		if flagProviderName != "" {
			provider = flagProviderName
		}
		template := generateClusterTemplate(provider)
		template.Name = flagClusterName
		if FlagDebug {
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

		if FlagDebug {
			fmt.Printf("Template:\n \t%+v", template)
		}

		fmt.Printf("Creating cluster '%s'...", flagClusterName)
		newCluster, err := createCluster(template)
		check(err)
		printClusters([]nks.Cluster{newCluster})

		if CurrentConfig.ClusterId == 0 {
			setCluster(newCluster.ID)
		}

	},
}

var deleteClusterCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "del"},
	Short:   "delete a cluster",
	Long:    ``,
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

		if flagListAllClusters {
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

	cls, err := SDKClient.GetClusters(o)

	check(err)

	return &cls, err
}

func getAllClusters() (*[]nks.Cluster, error) {
	o, err := strconv.Atoi(viper.GetString("org_id"))
	check(err)

	cls, err := SDKClient.GetAllClusters(o)

	check(err)

	return &cls, err
}

func getClusterByID(clusterId int) (*nks.Cluster, error) {
	o, err := strconv.Atoi(viper.GetString("org_id"))
	check(err)


	cl, err := SDKClient.GetCluster(o, clusterId)

	check(err)

	return cl, err
}

func setClusterKubeConfig(clusterId int) {

	kubeConfig, err := SDKClient.GetKubeConfig(CurrentConfig.OrgID, CurrentConfig.ClusterId)
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

	var err error

	if flagForce {
		err = SDKClient.ForceDeleteCluster(o, clusterId)
	} else {
		err = SDKClient.DeleteCluster(o, clusterId)
	}
	check(err)

	if clusterId == CurrentConfig.ClusterId {
		setClusterID(0)
	}

	cl, err := getAllClusters()
	check(err)
	printClusters(*cl)
	return nil
}

func createCluster(cl models.CreateClusterInput) (nks.Cluster, error) {
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%d/clusters", CurrentConfig.OrgID)
	b, err := json.Marshal(cl)
	check(err)
	if FlagDebug {
		fmt.Printf("\nSending request with body:\n%s\n", b)
	}

	res, err := httpRequestPost("POST", url, b)
	if FlagDebug {
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

func generateClusterTemplate(provider string) models.CreateClusterInput {
	c := models.CreateClusterInput{}

	switch provider {
	case "gce":
		mapstructure.Decode(gceDefaults, &c)
	case "hci":
		mapstructure.Decode(hciDefaults, &c)
	default:
		fmt.Printf("Provider %s is not supported :(", provider)
	}

	return c
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(getClustersCmd)
	clusterCmd.AddCommand(createClusterCmd)
	clusterCmd.AddCommand(deleteClusterCmd)
	clusterCmd.AddCommand(listClustersCmd)

	getClustersCmd.Flags().IntVarP(&getClusterId, "id", "", 0, "ID of cluster")

	listClustersCmd.Flags().IntVarP(&getClusterId, "id", "", 0, "ID of cluster")
	listClustersCmd.Flags().BoolVarP(&flagListAllClusters, "all", "a", false, "Get everything (incl. Service clusters)")

	createClusterCmd.Flags().StringVarP(&flagClusterName, "name", "n", "", "ID of cluster")
	createClusterCmd.Flags().StringVarP(&flagProviderName, "provider", "p", "", "Name of provider")
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
