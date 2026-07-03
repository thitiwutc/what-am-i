import type { Player } from './player'
import type { RoomState } from './room-state'

export type Room = {
  id: string
  state: RoomState
  players: Player[]
}
