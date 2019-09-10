package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	nks "github.com/NetApp/nks-sdk-go/nks"
	"github.com/spf13/cobra"
)

var useHttp bool

type organization struct {
	ID                 int    `json:"pk,Number"`
	Name               string `json:"name"`
	Slug               string `json:"slug"`
	Logo               string `json:"logo"`
	EnableExperimental string `json:"enable_experimental_feature,Bool"`
	CreatedTime        string `json:"created"`
	UpdatedTime        string `json:"updated"`
}

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

	fmt.Printf("Data has %d objects\n", len(data))

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
var organizationGetCmd = &cobra.Command{
	Use:   "orgs",
	Short: "list organizations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		organizationGet()
	},
}

func organizationGet() {
	orgs, err := GetOrgs()
	if err != nil {
		fmt.Printf("Error: There was an errorretrieving items::\n\t%s\n\n", err)
		orgs = &[]nks.Organization{}
	}
	printOrgs(orgs)
}

func init() {
	getCmd.AddCommand(organizationGetCmd)
}
