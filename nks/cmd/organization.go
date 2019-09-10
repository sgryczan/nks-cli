package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	nks "github.com/NetApp/nks-sdk-go/nks"
	"github.com/spf13/cobra"
)

// organizationCmd represents the organization command
var organizationCmd = &cobra.Command{
	Use:   "organization",
	Short: "A brief description of your command",
	Long:  ``,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("organization called")
	//},
}

func printOrgs(o *[]nks.Organization) {
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fmt.Fprintf(w, "NAME\tID\t\n")
	for _, org := range *o {
		fmt.Fprintf(w, "%s\t%d\t\n", org.Name, org.ID)
	}
	w.Flush()
}

func GetOrgs() (*[]nks.Organization, error) {
	c := newClient()
	data, err := c.GetOrganizations()
	check(err)

	return &data, err
}

func getDefaultOrg() (nks.Organization, error) {
	o, err := GetOrgs()
	if err != nil {
		o = &[]nks.Organization{}
	}

	// First Org returned is the default
	return (*o)[0], err
}

// organizationCmd represents the organization command
var getOrganizationsCmd = &cobra.Command{
	Use:     "orgs",
	Aliases: []string{"o", "organizations"},
	Short:   "list organizations",
	Long:    ``,
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
	c := newClient()
	o, err := c.GetOrganization(id)
	check(err)

	orgs := []nks.Organization{
		*o,
	}

	printOrgs(&orgs)
}

var getOrgsId string

func init() {
	getCmd.AddCommand(getOrganizationsCmd)
	getOrganizationsCmd.Flags().StringVarP(&getOrgsId, "id", "", "", "ID of organization")
}
