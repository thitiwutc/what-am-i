import type { Question } from './question'
import type { Role } from './role-type'

export type Player = {
  id: string
  name: string
  role: Role
  // word is the secret word each player has to guess.
  secret_word: string
  // isws is true if the player submitted secret word.
  isws: boolean
  questions: Question[]
  is_word_guessed: boolean
}
