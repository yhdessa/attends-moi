<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import type { Board, Card, CardStatus, CardPriority } from '../types'
import { COLUMNS } from '../types'
import { getCards, createCard, updateCard, deleteCard } from '../api'
import KanbanCard from './KanbanCard.vue'
import CreateCardModal from './CreateCardModal.vue'

const props = defineProps<{ board: Board }>()
const emit = defineEmits<{ 'edit-board': [] }>()

const cards = ref<Card[]>([])
const showCreateModal = ref(false)
const draggedCardId = ref<string | null>(null)

const search = ref('')
const assigneeFilter = ref('')
const labelFilter = ref('')
const priorityFilter = ref<CardPriority | 'all'>('all')

const hasFilters = computed(() => {
  return (
    search.value.trim() ||
    assigneeFilter.value.trim() ||
    labelFilter.value.trim() ||
    priorityFilter.value !== 'all'
  )
})

onMounted(async () => {
  const fetched = await getCards(props.board.id)
  cards.value = fetched || []
})

function getPosition(card: Card) {
  if (typeof card.position === 'number' && !Number.isNaN(card.position)) {
    return card.position
  }
  const fallback = Date.parse(card.created_at)
  return Number.isNaN(fallback) ? 0 : fallback
}

function matchesFilters(card: Card) {
  const q = search.value.trim().toLowerCase()
  if (q) {
    const haystack = [
      card.title,
      card.description,
      card.assignee,
      ...(card.labels || []),
    ]
      .filter(Boolean)
      .join(' ')
      .toLowerCase()
    if (!haystack.includes(q)) return false
  }

  const assigneeQuery = assigneeFilter.value.trim().toLowerCase()
  if (assigneeQuery && !card.assignee.toLowerCase().includes(assigneeQuery)) return false

  const labelQuery = labelFilter.value.trim().toLowerCase()
  if (labelQuery) {
    const hasLabel = (card.labels || []).some((label) => label.toLowerCase().includes(labelQuery))
    if (!hasLabel) return false
  }

  if (priorityFilter.value !== 'all' && card.priority !== priorityFilter.value) return false

  return true
}

function getStatusCards(status: CardStatus, excludeId?: string) {
  return cards.value
    .filter((c) => c.status === status && c.id !== excludeId)
    .sort((a, b) => getPosition(a) - getPosition(b))
}

function cardsByStatus(status: CardStatus) {
  return getStatusCards(status).filter(matchesFilters)
}

function nextPositionForStatus(status: CardStatus, excludeId?: string) {
  const statusCards = getStatusCards(status, excludeId)
  const last = statusCards[statusCards.length - 1]
  return last ? getPosition(last) + 1 : Date.now()
}

function betweenPositions(prev?: Card, next?: Card) {
  if (prev && next) return (getPosition(prev) + getPosition(next)) / 2
  if (!prev && next) return getPosition(next) - 1
  if (prev && !next) return getPosition(prev) + 1
  return Date.now()
}

function clearFilters() {
  search.value = ''
  assigneeFilter.value = ''
  labelFilter.value = ''
  priorityFilter.value = 'all'
}

async function handleCreate(card: Partial<Card>) {
  const status = (card.status ?? 'backlog') as CardStatus
  const position = nextPositionForStatus(status)
  const newCard = await createCard(props.board.id, { ...card, status, position })
  cards.value.push(newCard)
  showCreateModal.value = false
}

async function handleStatusChange(cardId: string, status: CardStatus) {
  const position = nextPositionForStatus(status, cardId)
  await handleMove(cardId, status, position)
}

async function handleUpdate(cardId: string, patch: Partial<Card>) {
  const current = cards.value.find((c) => c.id === cardId)
  const nextPatch = { ...patch }

  if (current && patch.status && patch.status !== current.status && nextPatch.position === undefined) {
    nextPatch.position = nextPositionForStatus(patch.status as CardStatus, cardId)
  }

  const updated = await updateCard(cardId, nextPatch)
  const idx = cards.value.findIndex((c) => c.id === cardId)
  if (idx !== -1) cards.value[idx] = updated
}

async function handleMove(cardId: string, status: CardStatus, position: number) {
  const idx = cards.value.findIndex((c) => c.id === cardId)
  if (idx === -1) return

  cards.value[idx] = { ...cards.value[idx], status, position }
  const updated = await updateCard(cardId, { status, position })
  cards.value[idx] = updated
}

async function handleDelete(cardId: string) {
  await deleteCard(cardId)
  cards.value = cards.value.filter((c) => c.id !== cardId)
}

function handleDragStart(cardId: string) {
  draggedCardId.value = cardId
}

async function handleDropOnColumn(event: DragEvent, status: CardStatus) {
  event.preventDefault()
  if (!draggedCardId.value) return

  const position = nextPositionForStatus(status, draggedCardId.value)
  await handleMove(draggedCardId.value, status, position)
  draggedCardId.value = null
}

async function handleDropOnCard(event: DragEvent, status: CardStatus, targetId: string) {
  event.preventDefault()
  event.stopPropagation()
  if (!draggedCardId.value || draggedCardId.value === targetId) return

  const statusCards = getStatusCards(status, draggedCardId.value)
  const targetIndex = statusCards.findIndex((c) => c.id === targetId)
  const prev = targetIndex > 0 ? statusCards[targetIndex - 1] : undefined
  const next = targetIndex >= 0 ? statusCards[targetIndex] : undefined
  const position = betweenPositions(prev, next)

  await handleMove(draggedCardId.value, status, position)
  draggedCardId.value = null
}
</script>

<template>
  <div>
    <div class="flex items-center gap-2 mb-6">
      <router-link to="/" class="btn btn-ghost btn-sm">← Back</router-link>
      <h2 class="text-2xl font-bold">{{ board.title }}</h2>
      <button class="btn btn-ghost btn-sm" @click="emit('edit-board')">Edit</button>
      <button class="btn btn-primary btn-sm ml-auto" @click="showCreateModal = true">
        + Add Card
      </button>
    </div>

    <div class="flex flex-wrap gap-2 mb-4">
      <input
        v-model="search"
        type="text"
        placeholder="Search cards..."
        class="input input-bordered input-sm w-full md:w-64"
      />
      <input
        v-model="assigneeFilter"
        type="text"
        placeholder="Assignee"
        class="input input-bordered input-sm w-full md:w-40"
      />
      <input
        v-model="labelFilter"
        type="text"
        placeholder="Label"
        class="input input-bordered input-sm w-full md:w-32"
      />
      <select v-model="priorityFilter" class="select select-bordered select-sm w-full md:w-40">
        <option value="all">All priorities</option>
        <option value="low">Low</option>
        <option value="medium">Medium</option>
        <option value="high">High</option>
        <option value="critical">Critical</option>
      </select>
      <button v-if="hasFilters" class="btn btn-ghost btn-sm" @click="clearFilters">
        Clear
      </button>
    </div>

    <div class="flex gap-4 overflow-x-auto pb-4">
      <div
        v-for="col in COLUMNS"
        :key="col.key"
        class="min-w-[280px] flex-1 bg-base-200 rounded-xl p-3"
        @dragover.prevent
        @drop="handleDropOnColumn($event, col.key)"
      >
        <div class="flex items-center gap-2 mb-3">
          <span :class="['w-3 h-3 rounded-full', col.color]"></span>
          <h3 class="font-semibold">{{ col.label }}</h3>
          <span class="badge badge-sm">{{ cardsByStatus(col.key).length }}</span>
        </div>

        <div class="space-y-2">
          <div
            v-for="card in cardsByStatus(col.key)"
            :key="card.id"
            @dragover.prevent
            @drop="handleDropOnCard($event, col.key, card.id)"
          >
            <KanbanCard
              :card="card"
              @status-change="handleStatusChange"
              @delete="handleDelete"
              @drag-start="handleDragStart"
              @update="handleUpdate"
            />
          </div>
        </div>
      </div>
    </div>

    <CreateCardModal
      :open="showCreateModal"
      @close="showCreateModal = false"
      @create="handleCreate"
    />
  </div>
</template>
