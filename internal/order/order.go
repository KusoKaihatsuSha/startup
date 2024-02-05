package order

type Stages int

const (
	Flag Stages = iota + 1
	File
	Env
)
