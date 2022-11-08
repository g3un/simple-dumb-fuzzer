package main

func main() {
	sdf := NewSdf("")
	lldb := NewLldb("")

	if err := sdf.SetDebugger(lldb); err != nil {
		panic(err)
	}
	if err := sdf.SetStatsD(""); err != nil {
		panic(err)
	}
	if err := sdf.SetCommand(""); err != nil {
		panic(err)
	}

	sdf.Run()
}
