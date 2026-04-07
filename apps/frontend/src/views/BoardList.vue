<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import type { Board } from '../types'
import { getBoards } from '../api'

const boards = ref<Board[]>([])
const router = useRouter()

onMounted(async () => {
  boards.value = await getBoards()
})
</script>

<template>
  <div>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div
        v-for="board in boards"
        :key="board.id"
        class="card bg-base-100 shadow-xl cursor-pointer hover:shadow-2xl transition-shadow"
        @click="router.push(`/boards/${board.id}`)"
      >
        <div class="card-body">
          <h2 class="card-title">{{ board.title }}</h2>
          <p v-if="board.description" class="text-sm text-base-content/60">{{ board.description }}</p>
          <div class="card-actions justify-end">
            <div class="badge badge-outline">{{ board.updated_at.slice(0, 10) }}</div>
          </div>
        </div>
      </div>

      <div v-if="boards.length === 0" class="col-span-full text-center py-20 text-base-content/40">
        <p class="text-lg">No boards yet. Create your first one!</p>
      </div>
    </div>
  </div>
</template>
