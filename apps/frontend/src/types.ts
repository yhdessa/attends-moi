export interface Board {
  id: string
  title: string
  description: string
  created_at: string
  updated_at: string
}

export interface Card {
  id: string
  board_id: string
  title: string
  description: string
  status: CardStatus
  priority: CardPriority
  labels: string[]
  assignee: string
  created_at: string
  updated_at: string
}

export type CardStatus = 'backlog' | 'todo' | 'in_progress' | 'review' | 'done'

export type CardPriority = 'low' | 'medium' | 'high' | 'critical'

export interface Comment {
  id: string
  card_id: string
  author: string
  body: string
  created_at: string
}

export const COLUMNS: { key: CardStatus; label: string; color: string }[] = [
  { key: 'backlog', label: 'Backlog', color: 'bg-gray-500' },
  { key: 'todo', label: 'To Do', color: 'bg-blue-500' },
  { key: 'in_progress', label: 'In Progress', color: 'bg-yellow-500' },
  { key: 'review', label: 'Review', color: 'bg-purple-500' },
  { key: 'done', label: 'Done', color: 'bg-green-500' },
]

export const PRIORITY_COLORS: Record<CardPriority, string> = {
  low: 'badge-ghost',
  medium: 'badge-info',
  high: 'badge-warning',
  critical: 'badge-error',
}
