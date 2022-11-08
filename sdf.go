package main

type Sdf struct {
	target string

	basePath  string
	seedPath  string
	casePath  string
	crashPath string

	iterCount  int
	crashCount int

	dbg        *debugger
	cmd        string
	statsDAddr string
}

var _ fuzzer = (*Sdf)(nil)

func NewSdf(target string) *Sdf {
	return &Sdf{
		target: target,

		basePath:  "",
		seedPath:  "",
		casePath:  "",
		crashPath: "",

		iterCount:  0,
		crashCount: 0,

		dbg:        nil,
		cmd:        "",
		statsDAddr: "",
	}
}

func (s *Sdf) SetDebugger(dbg debugger) error {
	return nil
}

func (s *Sdf) SetCommand(cmd string) error {
	return nil
}

func (s *Sdf) SetStatsD(addr string) error {
	return nil
}

func (s Sdf) pick() error {
	return nil
}
func (s Sdf) mutate() error {
	return nil
}
func (s Sdf) execute() error {
	return nil
}
func (s Sdf) monitor() error {
	return nil
}
func (s Sdf) report() error {
	return nil
}
func (s Sdf) clean() error {
	return nil
}

func (s Sdf) Run() {
}
