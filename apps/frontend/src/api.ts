import type { Board, Card, Comment } from './types'

const API = '/api'

export async function getBoards(): Promise<Board[]> {
  const res = await fetch(`${API}/boards`)
  return res.json()
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
  return res.json()
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
  return res.json()
}

export async function createComment(cardId: string, author: string, body: string): Promise<Comment> {
  const res = await fetch(`${API}/cards/${cardId}/comments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ author, body }),
  })
  return res.json()
}
