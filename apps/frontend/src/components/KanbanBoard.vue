<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Board, Card, CardStatus } from '../types'
import { COLUMNS } from '../types'
import { getCards, createCard, updateCard, deleteCard } from '../api'
import KanbanCard from './KanbanCard.vue'
import CreateCardModal from './CreateCardModal.vue'

const props = defineProps<{ board: Board }>()
const emit = defineEmits<{ back: [] }>()

const cards = ref<Card[]>([])
const showCreateModal = ref(false)

onMounted(async () => {
  cards.value = await getCards(props.board.id)
})

function cardsByStatus(status: CardStatus) {
  return cards.value.filter((c) => c.status === status)
}

async function handleCreate(card: Partial<Card>) {
  const newCard = await createCard(props.board.id, card)
  cards.value.push(newCard)
  showCreateModal.value = false
}

async function handleStatusChange(cardId: string, status: CardStatus) {
  const updated = await updateCard(cardId, { status })
  const idx = cards.value.findIndex((c) => c.id === cardId)
  if (idx !== -1) cards.value[idx] = updated
}

async function handleDelete(cardId: string) {
  await deleteCard(cardId)
  cards.value = cards.value.filter((c) => c.id !== cardId)
}
</script>

<template>
  <div>
    <div class="flex items-center gap-4 mb-6">
      <button class="btn btn-ghost btn-sm" @click="emit('back')">← Back</button>
      <h2 class="text-2xl font-bold">{{ board.title }}</h2>
      <button class="btn btn-primary btn-sm ml-auto" @click="showCreateModal = true">
        + Add Card
      </button>
    </div>

    <div class="flex gap-4 overflow-x-auto pb-4">
      <div
        v-for="col in COLUMNS"
        :key="col.key"
        class="min-w-[280px] flex-1 bg-base-200 rounded-xl p-3"
      >
        <div class="flex items-center gap-2 mb-3">
          <span :class="['w-3 h-3 rounded-full', col.color]"></span>
          <h3 class="font-semibold">{{ col.label }}</h3>
          <span class="badge badge-sm">{{ cardsByStatus(col.key).length }}</span>
        </div>

        <div class="space-y-2">
          <KanbanCard
            v-for="card in cardsByStatus(col.key)"
            :key="card.id"
            :card="card"
            @status-change="handleStatusChange"
            @delete="handleDelete"
          />
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
