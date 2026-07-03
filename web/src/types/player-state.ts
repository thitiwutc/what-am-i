import type { Role } from './role-type'

export type PlayerState = {
  id: string
  name: string
  role: Role
  // isws is true if the player submitted secret word.
  isws: boolean
}
