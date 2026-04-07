<script setup lang="ts">
import { ref } from 'vue'
import type { CardStatus, CardPriority } from '../types'

defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: []; create: [card: Partial<import('../types').Card>] }>()

const title = ref('')
const description = ref('')
const assignee = ref('')
const status = ref<CardStatus>('backlog')
const priority = ref<CardPriority>('medium')
const labelsInput = ref('')
const dueDate = ref('')

function handleSubmit() {
  if (!title.value.trim()) return

  const labels = labelsInput.value
    .split(',')
    .map((l) => l.trim())
    .filter(Boolean)

  emit('create', {
    title: title.value,
    description: description.value,
    assignee: assignee.value,
    status: status.value,
    priority: priority.value,
    labels,
    due_date: dueDate.value || undefined,
  })

  title.value = ''
  description.value = ''
  assignee.value = ''
  labelsInput.value = ''
  dueDate.value = ''
}
</script>

<template>
  <dialog :class="{ modal: true, 'modal-open': open }">
    <div class="modal-box">
      <h3 class="text-lg font-bold">New Card</h3>

      <div class="form-control mt-4">
        <label class="label"><span class="label-text">Title *</span></label>
        <input v-model="title" type="text" class="input input-bordered" placeholder="Task title" />
      </div>

      <div class="form-control mt-2">
        <label class="label"><span class="label-text">Description</span></label>
        <textarea v-model="description" class="textarea textarea-bordered" placeholder="Details..."></textarea>
      </div>

      <div class="grid grid-cols-2 gap-4 mt-2">
        <div class="form-control">
          <label class="label"><span class="label-text">Status</span></label>
          <select v-model="status" class="select select-bordered">
            <option value="backlog">Backlog</option>
            <option value="todo">To Do</option>
            <option value="in_progress">In Progress</option>
            <option value="review">Review</option>
            <option value="done">Done</option>
          </select>
        </div>

        <div class="form-control">
          <label class="label"><span class="label-text">Priority</span></label>
          <select v-model="priority" class="select select-bordered">
            <option value="low">Low</option>
            <option value="medium">Medium</option>
            <option value="high">High</option>
            <option value="critical">Critical</option>
          </select>
        </div>
      </div>

      <div class="form-control mt-2">
        <label class="label"><span class="label-text">Due Date</span></label>
        <input v-model="dueDate" type="datetime-local" class="input input-bordered" />
      </div>

      <div class="form-control mt-2">
        <label class="label"><span class="label-text">Assignee</span></label>
        <input v-model="assignee" type="text" class="input input-bordered" placeholder="Name" />
      </div>

      <div class="form-control mt-2">
        <label class="label"><span class="label-text">Labels (comma-separated)</span></label>
        <input v-model="labelsInput" type="text" class="input input-bordered" placeholder="bug, feature, ops" />
      </div>

      <div class="modal-action">
        <button class="btn btn-ghost" @click="emit('close')">Cancel</button>
        <button class="btn btn-primary" @click="handleSubmit">Create</button>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button @click="emit('close')">close</button>
    </form>
  </dialog>
</template>
