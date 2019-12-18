package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"text/tabwriter"
	"time"

	nks "github.com/NetApp/nks-sdk-go/nks"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	vpr "github.com/spf13/viper"
	models "gitlab.com/sgryczan/nks-cli/nks/models"
)

var getclusterID int
var flagListAllClusters bool

var flagClusterName string
var flagProviderName string
var createClusterNumWorkers int
var createClusterWorkerSize string
var createClusterMasterSize string
var deleteclusterIDf int

var gceDefaults = map[string]interface{}{
	// Name is implied
	"Provider":          vpr.GetString("provider"),
	"ProviderKey":       vpr.GetString("gce_keyset"),
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
	"SSHKeySet":         vpr.GetInt("ssh_keyset"),
	"Solutions":         []nks.Solution{nks.Solution{Solution: "helm_tiller"}},
	//"ProviderSubnetID":    "__new__",
	//"ProviderSubnetCidr":  "172.23.1.0/24",
	//"ProviderNetworkID":   "__new__",
	//"ProviderNetworkCIDR": "172.23.0.0/16",
}

var hciDefaults = map[string]interface{}{
	// Name is implied
	"Provider":    "hci",
	"ProviderKey": 63207,
	"WorkspaceID": 22022,

	"MasterCount":           1,
	"MasterSize":            "m",
	"MasterRootDiskSize":    50,
	"MasterGPUInstanceSize": "",
	"MasterGPUCoreCount":    nil,

	"WorkerCount":           2,
	"WorkerSize":            "m",
	"WorkerGPUInstanceSize": "",
	"WorkerGPUCoreCount":    nil,
	"WorkerRootDiskSize":    50,

	"Region": "LAB-RTP",

	"KubernetesVersion": "v1.14.3",
	"RbacEnabled":       true,
	"DashboardEnabled":  true,
	"K8sPodCIDR":        "10.2.0.0",
	"K8sServiceCIDR":    "10.3.0.0",
	"EtcdType":          "classic",
	"Platform":          "debian",
	"Channel":           "stable",
	"Zone":              "",
	"Config":            map[string]bool{"enable_experimental_features": true},
	"SSHKeySet":         vpr.GetInt("ssh_keyset"),
	"Solutions":         []nks.Solution{nks.Solution{Solution: "helm_tiller"}},
	"Features":          []string{},
	"MinNodeCount":      nil,
	"MaxNodeCount":      nil,
	"Owner":             26309, // ???
	//"ProviderSubnetID":    "__new__",
	//"ProviderSubnetCidr":  "172.23.1.0/24",
	//"ProviderNetworkID":   "__new__",
	//"ProviderNetworkCIDR": "172.23.0.0/16",
}

var clusterCmd = &cobra.Command{
	Use:     "clusters",
	Aliases: []string{"cl", "clus", "clu", "cluster", "clusters"},
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
		checkDefaultOrg()

		provider := "gce"
		if flagProviderName != "" {
			provider = flagProviderName
		}

		providerKsSetting := fmt.Sprintf("%s_keyset", provider)

		if flagDebug {
			fmt.Printf("Keyset var for provider %s should be: %s\n", provider, providerKsSetting)
		}

		// check user ssh keyset prese
		for _, keyset := range []string{"user_ssh", providerKsSetting} {
			if ks := vpr.GetInt(keyset); ks == 0 {

				if flagDebug {
					fmt.Printf("Debug - No keyset for provider: %s in var: %d. Attempting to set a default\n", provider, ks)
				}

				if keyset == "user_ssh" {
					setDefaultProviderKey(keyset)
				} else {
					setDefaultProviderKey(provider)
				}
			}
		}

		template := generateClusterTemplate(provider)
		template.Name = flagClusterName
		template.Provider = provider
		template.ProviderKey = vpr.GetInt(providerKsSetting)
		template.SSHKeySet = vpr.GetInt("ssh_keyset")
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

		fmt.Printf("Creating cluster '%s'...\n", flagClusterName)
		newCluster, err := createCluster(template)
		check(err)
		printClusters([]nks.Cluster{newCluster})

		if vpr.GetInt("cluster_id") == 0 {
			setClusterAsCurrent(newCluster.ID)
		}

	},
}

var deleteClusterCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "del"},
	Short:   "delete a cluster",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		checkDefaultOrg()

		err := deleteClusterByID(deleteclusterIDf)
		check(err)
	},
}

var listClustersCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "li"},
	Short:   "list clusters",
	Run: func(cmd *cobra.Command, args []string) {
		checkDefaultOrg()

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
			if vpr.GetInt("cluster_id") == 0 {
				setClusterAsCurrent((*c)[0].ID)
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
		checkDefaultOrg()

		i := vpr.GetInt("cluster_id")
		if getclusterID != 0 {
			i = getclusterID
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
	if flagDebug {
		fmt.Printf("Debug - printClusters()\n")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fmt.Fprintf(w, "NAME\tID\tPROVIDER\tNODES\tK8s_VERSION\tSTATE\t\n")
	for _, c := range cs {
		if c.ID == vpr.GetInt("cluster_id") {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v (default)\t\n", c.Name, c.ID, c.Provider, c.NodeCount, c.KubernetesVersion, c.State)
		} else {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t\n", c.Name, c.ID, c.Provider, c.NodeCount, c.KubernetesVersion, c.State)
		}

	}
	w.Flush()
}

func getClusters() (*[]nks.Cluster, error) {
	o := vpr.GetInt("org_id")

	if flagDebug {
		fmt.Printf("Debug - getClusters(%d)\n", o)
	}

	start := time.Now()
	cls, err := SDKClient.GetClusters(o)
	elapsed := time.Since(start)

	if flagDebug {
		fmt.Printf("Debug - getClusters() - completed in %v\n", elapsed)
		fmt.Printf("Debug - getClusters() - response: %v\n", cls)
	}

	check(err)

	return &cls, err
}

func getAllClusters() (*[]nks.Cluster, error) {
	o, err := strconv.Atoi(vpr.GetString("org_id"))
	check(err)

	cls, err := SDKClient.GetAllClusters(o)

	check(err)

	return &cls, err
}

func getClusterByID(clusterID int) (*nks.Cluster, error) {
	o, err := strconv.Atoi(vpr.GetString("org_id"))
	check(err)

	cl, err := SDKClient.GetCluster(o, clusterID)

	check(err)

	return cl, err
}

func setClusterKubeConfig(clusterID int) {
	if flagDebug {
		fmt.Printf("Debug - setClusterKubeConfig(%d)\n", clusterID)
	}

	kubeConfig, err := SDKClient.GetKubeConfig(vpr.GetInt("org_id"), clusterID)
	if err != nil {
		fmt.Printf("There was an error retrieving config for cluster %d: \n\t%v\n", vpr.GetInt("cluster_id"), err)
	}
	home, err := homedir.Dir()
	//fmt.Printf("Setting kubeconfig to cluster %d", clusterID)
	b := []byte(kubeConfig)
	if err != nil {
		fmt.Printf("There was an error finding your home directory: %v", err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/.kube/config", home), b, 0644)
	if err != nil {
		fmt.Printf("There was an error writing the kubeconfig: %v", err)
	}
}

func deleteClusterByID(clusterID int) error {
	o := vpr.GetInt("org_id")

	var err error

	if flagForce {
		err = SDKClient.ForceDeleteCluster(o, clusterID)
	} else {
		err = SDKClient.DeleteCluster(o, clusterID)
	}
	check(err)

	if clusterID == vpr.GetInt("cluster_id") {
		setclusterID(0)
	}

	cl, err := getAllClusters()
	check(err)
	printClusters(*cl)
	return nil
}

func createCluster(cl models.CreateClusterInput) (nks.Cluster, error) {
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%d/clusters", vpr.GetInt("org_id"))
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

	getClustersCmd.Flags().IntVarP(&getclusterID, "id", "", 0, "ID of cluster")

	listClustersCmd.Flags().IntVarP(&getclusterID, "id", "", 0, "ID of cluster")
	listClustersCmd.Flags().BoolVarP(&flagListAllClusters, "all", "a", false, "Get everything (incl. Service clusters)")

	createClusterCmd.Flags().StringVarP(&flagClusterName, "name", "n", "", "ID of cluster")
	createClusterCmd.Flags().StringVarP(&flagProviderName, "provider", "p", "", "Name of provider")
	createClusterCmd.Flags().StringVarP(&createClusterMasterSize, "master-size", "", "", "Instance size of master nodes")
	createClusterCmd.Flags().StringVarP(&createClusterWorkerSize, "worker-size", "", "", "Instance size of worker nodes")
	createClusterCmd.Flags().IntVarP(&createClusterNumWorkers, "num-workers", "w", 2, "Number of worker nodes (default: 2)")
	e := createClusterCmd.MarkFlagRequired("name")
	check(e)

	deleteClusterCmd.Flags().BoolVarP(&flagForce, "force", "f", false, "Force deletion")
	deleteClusterCmd.Flags().IntVarP(&deleteclusterIDf, "id", "i", 0, "ID of cluster to delete")
	e = deleteClusterCmd.MarkFlagRequired("id")
	check(e)
}
