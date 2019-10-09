package cmd

import (
	"fmt"
	"os"

	"text/tabwriter"

	nks "github.com/NetApp/nks-sdk-go/nks"
	"github.com/spf13/cobra"
	vpr "github.com/spf13/viper"
)

type Key struct {
	ID          string `json:"pk"`
	Keyset      string `json:"keyset"`
	Type        string `json:"key_type"`
	Fingerprint string `json:"fingerprint"`
	User        string `json:"user"`
}

type Keyset struct {
	ID         int             `json:"pk"`
	Name       string          `json:"name"`
	Category   string          `json:"category"`
	Entity     string          `json:"entity"`
	Org        string          `json:"org"`
	Workspaces []nks.Workspace `json:"workspaces"`
	metadata   map[string]string
	User       int    `json:"user"`
	Keys       []Key  `json:"keys"`
	Created    string `json:"created"`
}

var keysetsCmd = &cobra.Command{
	Use:     "keysets",
	Aliases: []string{"ks", "key", "keys"},
	Short:   "add or edit keysets",
	Long:    ``,
	/* Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("keysets called")
	}, */
}

var getKeysetsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "li"},
	Short:   "list keysets",
	Run: func(cmd *cobra.Command, args []string) {
		checkDefaultOrg()

		ks, err := getKeySets()
		if err != nil {
			fmt.Printf("There was an error retrieving items:\n\t%s\n\n", err)
			ks = &[]nks.Keyset{}
		}
		printKeysets(*ks)
	},
}

func printKeysets(ks []nks.Keyset) {
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fmt.Fprintf(w, "NAME\tID\tCATEGORY\tENTITY\t\n")
	for _, k := range ks {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", k.Name, k.ID, k.Category, k.Entity)
	}
	w.Flush()
}

func getKeySets() (*[]nks.Keyset, error) {

	ks, err := SDKClient.GetKeysets(vpr.GetInt("org_id"))

	return &ks, err
}


// This method in the SDK Client returns the first matching keyset for the specified provider
// regardless of what organization it belongs to. This can result in errors such as 
// "Provider credential does not belong to organization: $ORG"
func GetUserProfileKeysetID(up *nks.UserProfile, org int, prov string) (int, error) {
	if up == nil {
		return 0, fmt.Errorf("userprofile object is nil")
	}
	for _, ks := range up.Keysets {
		if (prov == "user_ssh" && ks.Category == "user_ssh" && ks.IsDefault && ks.Org == org) || (ks.Category == "provider" && ks.Entity == prov && ks.Org == org) {
			return ks.ID, nil
		}
	}
	return 0, fmt.Errorf("no %s keyset found in userprofile", prov)
}


func init() {
	keysetsCmd.AddCommand(getKeysetsCmd)
	rootCmd.AddCommand(keysetsCmd)
}
