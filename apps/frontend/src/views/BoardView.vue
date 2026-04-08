<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Board } from '../types'
import { getBoard, updateBoard } from '../api'
import KanbanBoard from '../components/KanbanBoard.vue'

const props = defineProps<{ id: string }>()
const board = ref<Board | null>(null)
const loading = ref(true)

const showEdit = ref(false)
const editTitle = ref('')
const editDescription = ref('')

onMounted(async () => {
  board.value = await getBoard(props.id)
  loading.value = false
})

function openEdit() {
  if (!board.value) return
  editTitle.value = board.value.title
  editDescription.value = board.value.description
  showEdit.value = true
}

async function saveEdit() {
  if (!board.value) return
  if (!editTitle.value.trim()) return

  const updated = await updateBoard(board.value.id, {
    title: editTitle.value,
    description: editDescription.value,
  })
  board.value = updated
  showEdit.value = false
}
</script>

<template>
  <div v-if="loading" class="text-center py-10">Loading...</div>
  <div v-else-if="board">
    <KanbanBoard :board="board" @edit-board="openEdit" />
  </div>
  <div v-else class="text-center py-10 text-base-content/60">Board not found</div>

  <dialog :class="{ modal: true, 'modal-open': showEdit }">
    <div class="modal-box">
      <h3 class="text-lg font-bold">Edit Board</h3>
      <div class="form-control mt-4">
        <label class="label"><span class="label-text">Title</span></label>
        <input v-model="editTitle" type="text" class="input input-bordered" />
      </div>
      <div class="form-control mt-2">
        <label class="label"><span class="label-text">Description</span></label>
        <textarea v-model="editDescription" class="textarea textarea-bordered"></textarea>
      </div>
      <div class="modal-action">
        <button class="btn btn-ghost" @click="showEdit = false">Cancel</button>
        <button class="btn btn-primary" @click="saveEdit">Save</button>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button @click="showEdit = false">close</button>
    </form>
  </dialog>
</template>
