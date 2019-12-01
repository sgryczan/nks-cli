package extensions

import (
	"fmt"

	nks "github.com/NetApp/nks-sdk-go/nks"
	"gitlab.com/sgryczan/nks-cli/nks/models"
)

// GetNotifications returns user notifications
func (c *SDK) GetNotifications() (ns []*models.Notification, err error) {
	req := &nks.APIReq{
		Method:       "GET",
		Path:         fmt.Sprintf("/user/notifications"),
		ResponseObj:  &ns,
		WantedStatus: 200,
	}
	err = c.RunRequest(req)
	return
}
