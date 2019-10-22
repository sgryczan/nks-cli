package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
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

// GetNotifications returns notifications for current user
func GetNotifications() []*models.Notification {
	if flagDebug {
		fmt.Printf("Debug - GetNotifications()\n")
	}
	//ns := models.Notifications{}

	ns, err := SDKClient.GetNotifications()
	if flagDebug {
		fmt.Printf("Debug - GetNotifications() - response: \n%v\n", ns)
	}
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	return ns
}

// PrintNotifications prints a slice of Notifications to stdout
func PrintNotifications(ns []*models.Notification, n int) {
	if flagDebug {
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
