package debug

import _ "net/http/pprof"
import "net/http"

func init() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()
}
