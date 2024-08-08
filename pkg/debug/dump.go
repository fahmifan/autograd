package debug

import (
	"encoding/json"
	"fmt"
)

func DumpJSON(label string, args any) {
	buf, _ := json.Marshal(args)
	fmt.Println("DEBUG >>> ", label, string(buf))
}
