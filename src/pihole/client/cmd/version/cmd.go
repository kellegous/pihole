package version

import (
	"fmt"

	"pihole/build"
)

// Name ...
const Name = "version"

// Main ...
func Main(conf string, args []string) {
	fmt.Printf("sha:    %s\n", build.SHA)
	fmt.Printf("ref:    %s\n", build.Ref)
	fmt.Printf("name:   %s\n", build.Name())
}
