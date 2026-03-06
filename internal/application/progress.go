package application

// Stage は処理パイプラインの段階を表す。
type Stage string

const (
	StageReading     Stage = "reading"
	StageSplitting   Stage = "splitting"
	StageTokenizing  Stage = "tokenizing"
	StageFiltering   Stage = "filtering"
	StageVocabulary  Stage = "vocabulary"
	StageCooccur     Stage = "cooccurrence"
	StageStatistics  Stage = "statistics"
	StageGraphBuild  Stage = "graph_build"
	StageCentrality  Stage = "centrality"
	StageCommunity   Stage = "community"
	StageComplete    Stage = "complete"
)

// Progress は処理進捗を表す。
type Progress struct {
	Stage      Stage   `json:"stage"`
	Percentage float64 `json:"percentage"`
	Message    string  `json:"message"`
}

// ProgressListener は進捗通知を受け取るコールバック型。
type ProgressListener func(Progress)
