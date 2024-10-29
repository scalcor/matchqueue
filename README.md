# MatchQueue

MatchQueue is a match-making service with a simple matching alogorithm.

## Example

```go
package main

import "github.com/scalcor/matchqueue"

func main() {
  // create a new queue with default configs
  queue := matchqueue.New(matchqueue.DefaultConfig())

  // add a party of players
  players := []*matchqueue.Player{
    {ID: 1, Score: 10.0},
    {ID: 2, Score: 11.5},
  }
  queue.AddPlayer(players)

  // add more players
  ...
  queue.AddPlayer(players2)

  // proc matching
  groups, err := queue.ProcMatching()
  if err != nil {
    panic(err.Error())
  }

  for _, group := range groups {
    fmt.Println(group)
  }
}

```

## Terms

- `Player`
  - User who is joining the matching.
- `Score`
  - Value which represents the player's skill. (a.k.a. MMP)
  - It is the base factor for the matching.
- `Party`
  - Group of `Player`s. The player joins the matching as party.
  - Base unit of matching.
  - `Score` of the party is the average of all players in the party.
- `Window`
  - Range value of the `Party`.
  - Center of the window is the party's `Score`.
  - Parties with overlapping windows will be matched.
  - Size of window changes as time passes.
- `Group`
  - Matched `Player`s.
  - It has 2 `Team`s.
- `Team`
  - Each `Player` in a `Group` belongs to a `Team`.
  - A group has 2 teams. Number of players in each team cannot differ more than 1.
