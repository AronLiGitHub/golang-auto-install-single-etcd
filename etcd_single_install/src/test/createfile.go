package main

import "../tools"

func main() {
	name := "./test"
	content := `
		test
	`
	tools.WriteWithIoutil(name,content)
}


