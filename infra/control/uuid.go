package control

import (
	"fmt"

	"v2ray.com/core/common"
	"v2ray.com/core/common/uuid"
)

type UUIDCommand struct{}

func (c *UUIDCommand) Name() string {
	return "uuid"
}

func (c *UUIDCommand) Description() Description {
	return Description{
		Short: "Generate new UUIDs",
		Usage: []string{"v2ctl uuid"},
	}
}

func (c *UUIDCommand) Execute([]string) error {
	u := uuid.New()
	fmt.Println(u.String())
	return nil
}

func init() {
	common.Must(RegisterCommand(&UUIDCommand{}))
}
