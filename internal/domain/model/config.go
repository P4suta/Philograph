package model

// AnalysisConfig は分析パイプラインの設定を表す。
type AnalysisConfig struct {
	WindowSize      int
	MinFrequency    int
	MinCooccurrence int
	TargetPOS       []POS
	Metric          Metric
	MaxNodes        int
	StopWords       []string
	Language        Language
}

// DefaultConfig はゼロ設定で動作するデフォルト設定を返す。
func DefaultConfig() AnalysisConfig {
	return AnalysisConfig{
		WindowSize:      5,
		MinFrequency:    3,
		MinCooccurrence: 2,
		TargetPOS:       []POS{POSNoun, POSVerb, POSAdjective},
		Metric:          MetricNPMI,
		MaxNodes:        150,
		StopWords:       nil,
		Language:        LangUnknown,
	}
}
