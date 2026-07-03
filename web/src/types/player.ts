import type { Role } from './role-type'

export type Player = {
  id: string
  name: string
  role: Role
  isSecretWordSubmitted: boolean
}
