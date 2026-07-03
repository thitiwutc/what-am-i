import type { ActionType } from './action-type'

export type Action = {
  id: string
  type: ActionType
  value: string
  yes_count: number
  no_count: number
  unsure_count: number
  created_at: string
}
