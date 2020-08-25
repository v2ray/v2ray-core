package control

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"

	"v2ray.com/core/common"
	"v2ray.com/core/common/platform"
)

const content = "H4sIAAAAAAAC/4SVMaskNwzH+/kUW6izcSthMGrcqLhVk0rdQS5cSMg7Xu4S0vizB8meZd57M3ta2GHX/ukvyZZmY2ZKDMzCzJyY5yOlxKII1omsf+qkBiiC6WhbYsbkjDAfySQsJqD3jtrD0EBM3sBHzG3kUsrglIQREXonpd47kYIi4AHmgI9Wcq2jlJITC6JZJ+v3ECYzBMAHyYm392yuY4zWsjACmHZSh6l3A0JETzGlWZqDsnArpTg62mhJONhOdO90p97V1BAnteoaOcuummtrrtuERQwUiJwP8a4KGKcyxdOCw1spOY+WHueFqmakAIgUSSuhwKNgobxKXSLbtg6r5cFmBiAeF6yCkYycmv+BiCIiW8ScHa3DgxAuZQbRhFNrLTFo96RBmx9jKWWG5nBsjyJzuIkftUblonppZU5t5LzwIks5L1a4lijagQxLokbIYwxfytNDC+XQqrWW9fzAunhqh5/Tg8PuaMw0d/Tcw3iDO81bHfWM/AnutMh2xqSUntMzd3wHDy9iHMQz8bmUZYvqedTJ5GgOnrNt7FIbSlwXE3wDI19n/KA38MsLaP4l89b5F8AV3ESOMIEhIBgezHBc0H6xV9KbaXwMvPcNvIHcC0C7UPZQx4JVTb35/AneSQq+bAYXsBmY7TCRupF2NTdVm/+ch22xa0pvRERKqt1oxj9DUbXzU84Gvj5hc5a81SlAUwMwgEs4T9+7sg9lb9h+908MWiKV8xtWciVTmnB3tivRjNerfXdxpfEBbq2NUvLMM5R9NLuyQg8nXT0PIh1xPd/wrcV49oJ6zbZaPlj2V87IY9T3F2XCOcW2MbZyZd49H+9m81E1N9SxlU+ff/1y+/f3719vf7788+Ugv/ffbMIH7ZNj0dsT4WMHHwLPu/Rp2O75uh99AK+N2xn7ZHq1OK6gczkN+9ngdOl1Qvki5xwSR8vFX6D+9vXA97B/+fr5rz9u/738uP328urP19vfP759e3n9Xs6jamvqlfJ/AAAA//+YAMZjDgkAAA=="

type LoveCommand struct{}

func (*LoveCommand) Name() string {
	return "lovevictoria"
}

func (*LoveCommand) Hidden() bool {
	return false
}

func (c *LoveCommand) Description() Description {
	return Description{
		Short: "",
		Usage: []string{""},
	}
}

func (*LoveCommand) Execute([]string) error {
	c, err := base64.StdEncoding.DecodeString(content)
	common.Must(err)
	reader, err := gzip.NewReader(bytes.NewBuffer(c))
	common.Must(err)
	b := make([]byte, 4096)
	nBytes, _ := reader.Read(b)

	bb := bytes.NewBuffer(b[:nBytes])
	scanner := bufio.NewScanner(bb)
	for scanner.Scan() {
		s := scanner.Text()
		fmt.Print(s + platform.LineSeparator())
	}

	return nil
}

func init() {
	common.Must(RegisterCommand(&LoveCommand{}))
}
