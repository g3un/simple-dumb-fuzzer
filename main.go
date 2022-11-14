package main

func main() {
	sdf := NewSdf("/Users/g3un/Desktop/test")
	lldb := NewLldb("/usr/bin/lldb")

	if err := sdf.SetDebugger(lldb); err != nil {
		panic(err)
	}
	if err := sdf.SetStatsD("eff.g3un.com:8125"); err != nil {
		panic(err)
	}

	sdf.Run()
}
