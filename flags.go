package vector

const (
	// FlagInit points vector is properly initialized.
	FlagInit = iota
	// FlagNoClear disables clear step.
	FlagNoClear
	// FlagExtraBool enables YAML style bool check [On, Off] in addition to [true, false].
	FlagExtraBool
)
