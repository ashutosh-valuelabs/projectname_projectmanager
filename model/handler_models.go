package model

type Conf struct {
	R []Routes `yaml:"Routes"`
}
type Routes struct {
	Path          string `yaml:"Path"`
	Callback      string `yaml:"Callback"`
	Method        string `yaml:"Method"`
	Authorization string `yaml:"Authorization"`
}
