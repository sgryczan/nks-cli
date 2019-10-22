package cmd

import (
	"fmt"
	"os"

	"text/tabwriter"

	nks "github.com/NetApp/nks-sdk-go/nks"
	"github.com/spf13/cobra"
	vpr "github.com/spf13/viper"
)

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

func init() {
	keysetsCmd.AddCommand(getKeysetsCmd)
	rootCmd.AddCommand(keysetsCmd)
}
