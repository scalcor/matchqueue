package matchqueue

type State struct {
	// snapshot of the latest round
	Round        uint64  // last round
	PlayerQueued int     // number of players who were in the queue when the latest round started.
	MatchWindow  float64 // match window when the latest round started.

	// accumulations
	GroupCreated        int    // number of created groups
	WaitTimeAll         uint64 // sum of wait time (second) of all matched players
	WaitTimeAvg         uint64 // average of wait time (second) of all matched players
	WaitTimeMax         uint64 // maximum wait time (second) among all matched players
	PlayerMatched       int    // number of matched players
	CanceledWaitTimeAll uint64 // sum of wait time (second) of all canceled players
	CanceledWaitTimeAvg uint64 // average of wait time (second) of all canceled players
	CanceledWaitTimeMax uint64 // maximum wait time (second) among all canceled players
	PlayerCanceled      int    // number of canceled players
}

func (s *State) AddMatched(waitTime uint64, cnt int) {
	s.PlayerMatched += cnt
	s.WaitTimeAll += waitTime * uint64(cnt)
	s.WaitTimeAvg = s.WaitTimeAll / uint64(s.PlayerMatched)
	s.WaitTimeMax = max(s.WaitTimeMax, waitTime)
}

func (s *State) AddCanceled(waitTime uint64, cnt int) {
	s.PlayerCanceled += cnt
	s.CanceledWaitTimeAll += waitTime * uint64(cnt)
	s.CanceledWaitTimeAvg = s.CanceledWaitTimeAll / uint64(s.PlayerCanceled)
	s.CanceledWaitTimeMax = max(s.CanceledWaitTimeMax, waitTime)
}
