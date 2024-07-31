package pkg

import "fmt"

var(
	//Package package name 
	Package = "BlockIRC"

	//Version release version 
	Version = "1.6.4"

	//Gitcommit will be overwritten automatically by the build systems 
	Gitcommit = "HEAD"
)

//// FullVersion display the full version and build
func FullVersion() string{
	return fmt.Sprintf("%s-%s@%s", Package, Version, Gitcommit)
}