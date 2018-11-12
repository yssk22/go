package flow

// Options is to controll the generactor behavior
type Options struct {
	RootDir        string `json:"root_dir"`
	GenerateJS     bool   `json:"generate_js"`
	GenerateLibDef bool   `json:"generate_lib_def"`
}

// NewOptions returns a new *Option with default values.
func NewOptions() *Options {
	return &Options{
		RootDir:        ".",
		GenerateJS:     true,
		GenerateLibDef: true,
	}
}
