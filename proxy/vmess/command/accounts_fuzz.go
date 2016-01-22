// +build gofuzz

package command

func Fuzz(data []byte) int {
	cmd := new(SwitchAccount)
	if err := cmd.Unmarshal(data); err != nil {
		return 0
	}
	return 1
}
