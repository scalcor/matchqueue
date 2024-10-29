package matchqueue

type (
	PlayerID uint64
	Player   struct {
		ID    PlayerID
		Score float64
	}

	GroupID uint64
	Group   struct {
		ID           GroupID
		Players      [2][]*Player
		CreatedRound uint64
	}
)

func (id PlayerID) IsValid() bool {
	return id > 0
}

func (id GroupID) IsValid() bool {
	return id > 0
}
