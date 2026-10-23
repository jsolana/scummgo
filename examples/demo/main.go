package main

import (
	"github.com/jsolana/scummgo"
)

func main() {
	loader := scummgo.NewResourceFileLoader("./")

	app := scummgo.New(loader, scummgo.WithScreenZoom(4), scummgo.WithDebugMode())
	app.RunCommand(scummgo.ScriptRun{ScriptRef: scummgo.NewResourceRef("resources", "scripts/boot")})
	app.Run()
}
