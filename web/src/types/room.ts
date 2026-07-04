import type { Question } from './question'
import type { Player } from './player'
import type { RoomState } from './room-state'
import type { Guess } from './guess'

export type Room = {
  id: string
  state: RoomState
  players: Player[]
  active_player: Player
  last_question: Question | null
  last_guess: Guess | null
}
