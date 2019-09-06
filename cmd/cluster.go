package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

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

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("cluster called")
	//},
}

var getClustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "list clusters",
	Run: func(cmd *cobra.Command, args []string) {
		cs, err := getClusters()
		if err != nil {
			fmt.Printf("There was an error retrieving items:\n\t%s\n\n", err)
			cs = &[]Cluster{
				Cluster{},
			}
		}
		printClusters(*cs)
	},
}

func printClusters(cs []Cluster) {
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fmt.Fprintf(w, "NAME\tID\tPROVIDER\tNODES\tK8s_VERSION\t\n")
	for _, c := range cs {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t\n", c.Name, c.ID, c.Provider, c.NodeCount, c.K8sVersion)
	}
	w.Flush()
}

func getClusters() (*[]Cluster, error) {
	orgID := viper.GetString("org_id")
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%s/clusters", orgID)
	res, err := httpRequest("GET", url)

	data := []Cluster{}

	_ = json.Unmarshal(res, &data)
	//check(err)

	return &data, err
}

func init() {
	getCmd.AddCommand(getClustersCmd)
}
