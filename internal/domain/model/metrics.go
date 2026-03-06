package model

// Metric は共起強度の算出指標を表す。
type Metric int

const (
	MetricPMI       Metric = iota
	MetricNPMI
	MetricJaccard
	MetricFrequency
)
