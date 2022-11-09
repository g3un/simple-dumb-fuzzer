package main

func main() {
	sdf := NewSdf("/Users/g3un/Desktop/test")
	lldb := NewLldb("/usr/bin/lldb")

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
