import type { Board, Card, Comment } from './types'

const API = '/api'

export async function getBoards(): Promise<Board[]> {
  const res = await fetch(`${API}/boards`)
  const data = await res.json()
  return data || []
}

export async function createBoard(title: string, description: string): Promise<Board> {
  const res = await fetch(`${API}/boards`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title, description }),
  })
  return res.json()
}

export async function getCards(boardId: string): Promise<Card[]> {
  const res = await fetch(`${API}/boards/${boardId}/cards`)
  const data = await res.json()
  return data || []
}

export async function createCard(boardId: string, card: Partial<Card>): Promise<Card> {
  const res = await fetch(`${API}/boards/${boardId}/cards`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(card),
  })
  return res.json()
}

export async function updateCard(id: string, card: Partial<Card>): Promise<Card> {
  const res = await fetch(`${API}/cards/${id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(card),
  })
  return res.json()
}

export async function deleteCard(id: string): Promise<void> {
  await fetch(`${API}/cards/${id}`, { method: 'DELETE' })
}

export async function getComments(cardId: string): Promise<Comment[]> {
  const res = await fetch(`${API}/cards/${cardId}/comments`)
  const data = await res.json()
  return data || []
}

export async function createComment(cardId: string, author: string, body: string): Promise<Comment> {
  const res = await fetch(`${API}/cards/${cardId}/comments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ author, body }),
  })
  return res.json()
}
