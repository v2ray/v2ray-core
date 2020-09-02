package v2board

import "fmt"

func (v *V2Board) APIEntry() string {
	return fmt.Sprintf("%s/api/v1/server/battleroach", v.config.Server)
}

func (v *V2Board) ConfigUri() string {
	return fmt.Sprintf("%s/config?node_id=%d&token=%s", v.APIEntry(), v.config.Node, v.config.Token)
}

func (v *V2Board) UserListUri() string {
	return fmt.Sprintf("%s/user?node_id=%d&token=%s", v.APIEntry(), v.config.Node, v.config.Token)
}

func (v *V2Board) ReportTrafficUri() string {
	return fmt.Sprintf("%s/submit?node_id=%d&token=%s", v.APIEntry(), v.config.Node, v.config.Token)
}
