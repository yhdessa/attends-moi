<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Board } from './types'
import { getBoards, createBoard } from './api'
import KanbanBoard from './components/KanbanBoard.vue'

const boards = ref<Board[]>([])
const selectedBoard = ref<Board | null>(null)
const showCreateBoard = ref(false)
const newBoardTitle = ref('')
const newBoardDesc = ref('')

onMounted(async () => {
  boards.value = await getBoards()
})

async function handleCreateBoard() {
  if (!newBoardTitle.value.trim()) return
  const board = await createBoard(newBoardTitle.value, newBoardDesc.value)
  boards.value.push(board)
  selectedBoard.value = board
  showCreateBoard.value = false
  newBoardTitle.value = ''
  newBoardDesc.value = ''
}
</script>

<template>
  <div class="min-h-screen bg-base-300">
    <nav class="navbar bg-base-100 shadow-lg px-6">
      <div class="flex-1">
        <h1 class="text-xl font-bold text-primary">attends-moi</h1>
      </div>
      <div class="flex-none gap-2">
        <button class="btn btn-primary btn-sm" @click="showCreateBoard = true">
          + New Board
        </button>
      </div>
    </nav>

    <div class="container mx-auto p-6">
      <div v-if="!selectedBoard" class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div
          v-for="board in boards"
          :key="board.id"
          class="card bg-base-100 shadow-xl cursor-pointer hover:shadow-2xl transition-shadow"
          @click="selectedBoard = board"
        >
          <div class="card-body">
            <h2 class="card-title">{{ board.title }}</h2>
            <p v-if="board.description" class="text-sm text-base-content/60">
              {{ board.description }}
            </p>
            <div class="card-actions justify-end">
              <div class="badge badge-outline">{{ board.updated_at.slice(0, 10) }}</div>
            </div>
          </div>
        </div>

        <div
          v-if="boards.length === 0"
          class="col-span-full text-center py-20 text-base-content/40"
        >
          <p class="text-lg">No boards yet. Create your first one!</p>
        </div>
      </div>

      <KanbanBoard
        v-else
        :board="selectedBoard"
        @back="selectedBoard = null"
      />
    </div>

    <dialog :class="{ modal: true, 'modal-open': showCreateBoard }">
      <div class="modal-box">
        <h3 class="text-lg font-bold">Create New Board</h3>
        <div class="form-control mt-4">
          <label class="label"><span class="label-text">Title</span></label>
          <input
            v-model="newBoardTitle"
            type="text"
            placeholder="e.g. Backend API"
            class="input input-bordered"
          />
        </div>
        <div class="form-control mt-2">
          <label class="label"><span class="label-text">Description</span></label>
          <textarea
            v-model="newBoardDesc"
            placeholder="Optional description..."
            class="textarea textarea-bordered"
          ></textarea>
        </div>
        <div class="modal-action">
          <button class="btn btn-ghost" @click="showCreateBoard = false">Cancel</button>
          <button class="btn btn-primary" @click="handleCreateBoard">Create</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="showCreateBoard = false">close</button>
      </form>
    </dialog>
  </div>
</template>
