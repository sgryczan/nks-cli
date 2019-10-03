package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	nks "github.com/NetApp/nks-sdk-go/nks"
	"github.com/spf13/cobra"
	vpr "github.com/spf13/viper"
)

var organizationCmd = &cobra.Command{
	Use:     "organization",
	Aliases: []string{"orgs", "org", "organizations"},
	Short:   "manage organizations",
	Long:    ``,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("organization called")
	//},
}

func checkDefaultOrg() {
	if FlagDebug {
		fmt.Printf("Debug - checkDefaultOrg()\n")
	}
	if vpr.GetInt("org_id") == 0 {
		fmt.Printf("No default organization set. Specify an organization, or set a default one with `nks config set org -i <id>'\n")
		os.Exit(1)
	}
}

func printOrgs(o *[]nks.Organization) {

	if FlagDebug {
		fmt.Printf("Debug - printOrgs()\n")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fmt.Fprintf(w, "NAME\tID\t\n")
	for _, org := range *o {
		if org.ID == vpr.GetInt("org_id") {
			fmt.Fprintf(w, "%s\t%d (current)\t\n", org.Name, org.ID)
		} else {
			fmt.Fprintf(w, "%s\t%d\t\n", org.Name, org.ID)
		}
	}
	w.Flush()
}

func GetOrgs() (*[]nks.Organization, error) {
	if SDKClient.Token == "" {
		if FlagDebug {
			fmt.Println("Debug - GetOrgs(), no configured token in client, initializing")
		}
		initClient()
	}

	if FlagDebug {
		fmt.Println("Debug - GetOrgs(): sending request")

	}
	start := time.Now()
	data, err := SDKClient.GetOrganizations()
	elapsed := time.Since(start)
	check(err)

	if FlagDebug {
		fmt.Printf("Debug - GetOrgs(): request completed in %v\n", elapsed)
	}

	return &data, err
}

func getDefaultOrg() (nks.Organization, error) {
	if FlagDebug {
		fmt.Println("Debug - getDefaultOrg()")
	}
	o, err := GetOrgs()
	if err != nil {
		o = &[]nks.Organization{}
	}

	// First Org returned is the default
	return (*o)[0], err
}

var listOrganizationsCmd = &cobra.Command{
	Use:   "list",
	Short: "list organizations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		getOrganizations()
	},
}

var getOrganizationsCmd = &cobra.Command{
	Use:   "get",
	Short: "get organization details",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if getOrgsId != "" {
			i, err := strconv.Atoi(getOrgsId)
			check(err)
			getOrganizationByID(i)
		} else {
			getOrganizations()
		}
	},
}

func getOrganizations() {
	orgs, err := GetOrgs()
	if err != nil {
		fmt.Printf("Error: There was an error retrieving items::\n\t%s\n\n", err)
		orgs = &[]nks.Organization{}
	}
	printOrgs(orgs)
}

func getOrganizationByID(id int) {

	o, err := SDKClient.GetOrganization(id)
	check(err)

	orgs := []nks.Organization{
		*o,
	}

	printOrgs(&orgs)
}

var getOrgsId string

func init() {
	rootCmd.AddCommand(organizationCmd)
	organizationCmd.AddCommand(getOrganizationsCmd)
	organizationCmd.AddCommand(listOrganizationsCmd)
	getOrganizationsCmd.Flags().StringVarP(&getOrgsId, "id", "", "", "ID of organization")
}
