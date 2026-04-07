<script setup lang="ts">
import { ref } from 'vue'
import { createBoard } from './api'
import { useRouter } from 'vue-router'

const showCreateBoard = ref(false)
const newBoardTitle = ref('')
const newBoardDesc = ref('')
const router = useRouter()

async function handleCreateBoard() {
  if (!newBoardTitle.value.trim()) return
  const board = await createBoard(newBoardTitle.value, newBoardDesc.value)
  showCreateBoard.value = false
  newBoardTitle.value = ''
  newBoardDesc.value = ''
  router.push(`/boards/${board.id}`)
}
</script>

<template>
  <div class="min-h-screen bg-base-300">
    <nav class="navbar bg-base-100 shadow-lg px-6">
      <div class="flex-1">
        <router-link to="/" class="text-xl font-bold text-primary hover:underline">attends-moi</router-link>
      </div>
      <div class="flex-none gap-2">
        <button class="btn btn-primary btn-sm" @click="showCreateBoard = true">
          + New Board
        </button>
      </div>
    </nav>

    <div class="container mx-auto p-6">
      <router-view />
    </div>

    <dialog :class="{ modal: true, 'modal-open': showCreateBoard }">
      <div class="modal-box">
        <h3 class="text-lg font-bold">Create New Board</h3>
        <div class="form-control mt-4">
          <label class="label"><span class="label-text">Title</span></label>
          <input v-model="newBoardTitle" type="text" placeholder="e.g. Backend API" class="input input-bordered" />
        </div>
        <div class="form-control mt-2">
          <label class="label"><span class="label-text">Description</span></label>
          <textarea v-model="newBoardDesc" placeholder="Optional description..." class="textarea textarea-bordered"></textarea>
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
