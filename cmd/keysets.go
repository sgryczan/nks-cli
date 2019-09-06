package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"

	"text/tabwriter"

	"github.com/spf13/cobra"
)

type Key struct {
	ID          string `json:"pk"`
	Keyset      string `json:"keyset"`
	Type        string `json:"key_type"`
	Fingerprint string `json:"fingerprint"`
	User        string `json:"user"`
}

type Keyset struct {
	ID         int         `json:"pk"`
	Name       string      `json:"name"`
	Category   string      `json:"category"`
	Entity     string      `json:"entity"`
	Org        string      `json:"org"`
	Workspaces []Workspace `json:"workspaces"`
	metadata   map[string]string
	User       int    `json:"user"`
	Keys       []Key  `json:"keys"`
	Created    string `json:"created"`
}

// keysetsCmd represents the keysets command
var keysetsCmd = &cobra.Command{
	Use:   "keysets",
	Short: "add or edit keysets",
	Long:  ``,
	/* Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("keysets called")
	}, */
}

var getKeysetsCmd = &cobra.Command{
	Use:   "keysets",
	Short: "list keysets",
	Run: func(cmd *cobra.Command, args []string) {
		ks, err := getKeySets()
		if err != nil {
			fmt.Printf("There was an error retrieving items:\n\t%s\n\n", err)
			ks = &[]Keyset{
				Keyset{},
			}
		}
		printKeysets(*ks)
	},
}

func printKeysets(ks []Keyset) {
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fmt.Fprintf(w, "NAME\tID\tCATEGORY\tENTITY\t\n")
	for _, k := range ks {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", k.Name, k.ID, k.Category, k.Entity)
	}
	w.Flush()
}

func getKeySets() (*[]Keyset, error) {
	orgID := viper.GetString("org_id")
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%s/keysets", orgID)
	res, err := httpRequest("GET", url)

	data := []Keyset{}

	_ = json.Unmarshal(res, &data)
	//check(err)

	return &data, err
}

func init() {
	getCmd.AddCommand(getKeysetsCmd)
}
