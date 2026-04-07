<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Board } from '../types'
import { getBoards } from '../api'
import KanbanBoard from '../components/KanbanBoard.vue'

const props = defineProps<{ id: string }>()
const board = ref<Board | null>(null)

onMounted(async () => {
  const boards = await getBoards()
  board.value = boards.find((b) => b.id === props.id) || null
})
</script>

<template>
  <div v-if="board">
    <KanbanBoard :board="board" />
  </div>
  <div v-else class="text-center py-10">Loading...</div>
</template>
