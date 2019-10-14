package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	vpr "github.com/spf13/viper"
	models "gitlab.com/sgryczan/nks-cli/nks/models"
)

var notificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "View notifications",

	Run: func(cmd *cobra.Command, args []string) {
		PrintNotifications(GetNotifications(), 10)
	},
}

func init() {
	rootCmd.AddCommand(notificationsCmd)
}

func GetNotifications() []models.Notification {
	if FlagDebug {
		fmt.Printf("Debug - GetNotifications()\n")
	}
	url := fmt.Sprintf("%s/user/notifications", vpr.GetString("api_url"))
	if FlagDebug {
		fmt.Printf("Debug - GetNotifications() - url: %s\n", url)
	}
	ns := models.Notifications{}

	res, err := httpRequest("GET", url)
	if FlagDebug {
		fmt.Printf("Debug - GetNotifications() - response: \n%s\n", res)
	}
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	err = json.Unmarshal(res, &ns)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	return ns
}

func PrintNotifications(ns []models.Notification, n int) {
	if FlagDebug {
		fmt.Printf("Debug - printNotifications()\n")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fmt.Fprintf(w, "NAME\tMESSAGE\tCATEGORY\tTOPIC\tLEVEL\tCLUSTER\tORG\t\n")
	for i := 0; i < n; i++ {
		n := ns[i]
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t\n", n.ID, n.Message, n.Category, n.Topic, n.Level, n.ExtraData.Cluster.Name, n.ExtraData.Org.Name)
	}
	w.Flush()
}
