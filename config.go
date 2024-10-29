package matchqueue

type Config struct {
	// match window
	InitMatchWindow      float64 `json:"init_match_window"`
	MinMatchWindow       float64 `json:"min_match_window"`
	MaxMatchWindow       float64 `json:"max_match_window"`
	WindowAdjustStep     float64 `json:"window_adjust_step"`
	MinRateToKeepWindow  float64 `json:"min_rate_to_keep_window"`
	MaxRateToKeepWindow  float64 `json:"max_rate_to_keep_window"`
	WindowAdjustPerRetry float64 `json:"window_adjust_per_retry"`

	// filter
	ScoreBoundFilter     string    `json:"score_bound_filter"`
	ScoreModRatio        []float64 `json:"score_mod_ratio"`
	MatchingWindowFilter string    `json:"matching_window_filter"`

	// group
	MinNumToCreateGroup    int `json:"min_num_to_create_group"`
	MaxNumToCreateGroup    int `json:"max_num_to_create_group"`
	NumPlayerToCreateGroup int `json:"num_player_to_create_group"`
	NumRoundToCreateGroup  int `json:"num_round_to_create_group"`

	// notification
	NotifyThresholdWait int `json:"notify_threshold_wait"`
}

var defaultModRatio = []float64{1.0, 1.0, 1.0, 1.0, 1.0, 0.8, 0.6, 0.4, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2}

// DefaultConfig returns the predefined configuration.
func DefaultConfig() *Config {
	return &Config{
		InitMatchWindow:        10.0,
		MinMatchWindow:         5.0,
		MaxMatchWindow:         50.0,
		WindowAdjustStep:       0.1,
		MinRateToKeepWindow:    0.85,
		MaxRateToKeepWindow:    0.95,
		WindowAdjustPerRetry:   0.5,
		ScoreBoundFilter:       "curve",
		ScoreModRatio:          defaultModRatio,
		MatchingWindowFilter:   "calculated",
		MinNumToCreateGroup:    10,
		MaxNumToCreateGroup:    16,
		NumPlayerToCreateGroup: 40,
		NumRoundToCreateGroup:  2,
		NotifyThresholdWait:    10,
	}
}
