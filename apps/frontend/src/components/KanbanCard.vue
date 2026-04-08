<script setup lang="ts">
import { ref } from 'vue'
import type { Card, CardStatus, Comment, CardPriority } from '../types'
import { COLUMNS, PRIORITY_COLORS } from '../types'
import { getComments, createComment } from '../api'
import CommentSection from './CommentSection.vue'

const props = defineProps<{ card: Card }>()
const emit = defineEmits<{
  'status-change': [cardId: string, status: CardStatus]
  delete: [cardId: string]
  'drag-start': [cardId: string]
  update: [cardId: string, patch: Partial<Card>]
}>()

const showDetails = ref(false)
const comments = ref<Comment[]>([])

const editTitle = ref('')
const editDescription = ref('')
const editAssignee = ref('')
const editLabels = ref('')
const editStatus = ref<CardStatus>('backlog')
const editPriority = ref<CardPriority>('medium')
const editDueDate = ref('')

function toDateTimeLocal(value?: string) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  const pad = (num: number) => String(num).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(
    date.getHours(),
  )}:${pad(date.getMinutes())}`
}

function syncEditFields() {
  editTitle.value = props.card.title
  editDescription.value = props.card.description
  editAssignee.value = props.card.assignee
  editLabels.value = (props.card.labels || []).join(', ')
  editStatus.value = props.card.status
  editPriority.value = props.card.priority
  editDueDate.value = toDateTimeLocal(props.card.due_date)
}

async function openDetails() {
  showDetails.value = true
  syncEditFields()
  comments.value = await getComments(props.card.id)
}

async function addComment(author: string, body: string) {
  const comment = await createComment(props.card.id, author, body)
  comments.value.push(comment)
}

function saveEdits() {
  if (!editTitle.value.trim()) return

  const labels = editLabels.value
    .split(',')
    .map((label) => label.trim())
    .filter(Boolean)

  const dueDateValue = editDueDate.value ? new Date(editDueDate.value) : null
  const dueDateISO =
    dueDateValue && !Number.isNaN(dueDateValue.getTime()) ? dueDateValue.toISOString() : undefined

  emit('update', props.card.id, {
    title: editTitle.value,
    description: editDescription.value,
    assignee: editAssignee.value,
    labels,
    status: editStatus.value,
    priority: editPriority.value,
    due_date: dueDateISO,
  })
  showDetails.value = false
}
</script>

<template>
  <div
    class="card bg-base-100 shadow-sm cursor-pointer hover:shadow-md transition-shadow"
    draggable="true"
    @dragstart="emit('drag-start', card.id)"
    @click="openDetails"
  >
    <div class="card-body p-4">
      <div class="flex items-start justify-between">
        <h4 class="font-medium text-sm">{{ card.title }}</h4>
        <div class="dropdown dropdown-end">
          <button class="btn btn-ghost btn-xs" @click.stop>⋮</button>
          <ul class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-32 z-10">
            <li v-for="col in COLUMNS" :key="col.key">
              <a @click="emit('status-change', card.id, col.key)">{{ col.label }}</a>
            </li>
            <li><a class="text-error" @click="emit('delete', card.id)">Delete</a></li>
          </ul>
        </div>
      </div>

      <p v-if="card.description" class="text-xs text-base-content/60 line-clamp-2">
        {{ card.description }}
      </p>

      <div class="flex flex-wrap gap-1 mt-1">
        <span v-for="label in card.labels" :key="label" class="badge badge-xs badge-outline">
          {{ label }}
        </span>
        <span :class="['badge badge-xs', PRIORITY_COLORS[card.priority]]">
          {{ card.priority }}
        </span>
      </div>

      <div v-if="card.due_date" class="flex items-center gap-1 mt-1 text-xs text-base-content/60">
        <span>📅</span>
        <span>{{ new Date(card.due_date).toLocaleDateString() }}</span>
      </div>

      <div v-if="card.assignee" class="flex items-center gap-1 mt-1">
        <div class="avatar placeholder">
          <div class="bg-neutral text-neutral-content rounded-full w-5">
            <span class="text-xs">{{ card.assignee[0].toUpperCase() }}</span>
          </div>
        </div>
        <span class="text-xs text-base-content/60">{{ card.assignee }}</span>
      </div>
    </div>
  </div>

  <dialog :class="{ modal: true, 'modal-open': showDetails }">
    <div class="modal-box">
      <h3 class="text-lg font-bold">Edit Card</h3>

      <div class="form-control mt-4">
        <label class="label"><span class="label-text">Title</span></label>
        <input v-model="editTitle" type="text" class="input input-bordered" />
      </div>

      <div class="form-control mt-2">
        <label class="label"><span class="label-text">Description</span></label>
        <textarea v-model="editDescription" class="textarea textarea-bordered"></textarea>
      </div>

      <div class="grid grid-cols-2 gap-4 mt-2">
        <div class="form-control">
          <label class="label"><span class="label-text">Status</span></label>
          <select v-model="editStatus" class="select select-bordered">
            <option v-for="col in COLUMNS" :key="col.key" :value="col.key">{{ col.label }}</option>
          </select>
        </div>

        <div class="form-control">
          <label class="label"><span class="label-text">Priority</span></label>
          <select v-model="editPriority" class="select select-bordered">
            <option value="low">Low</option>
            <option value="medium">Medium</option>
            <option value="high">High</option>
            <option value="critical">Critical</option>
          </select>
        </div>
      </div>

      <div class="form-control mt-2">
        <label class="label"><span class="label-text">Due Date</span></label>
        <input v-model="editDueDate" type="datetime-local" class="input input-bordered" />
      </div>

      <div class="form-control mt-2">
        <label class="label"><span class="label-text">Assignee</span></label>
        <input v-model="editAssignee" type="text" class="input input-bordered" />
      </div>

      <div class="form-control mt-2">
        <label class="label"><span class="label-text">Labels (comma-separated)</span></label>
        <input v-model="editLabels" type="text" class="input input-bordered" />
      </div>

      <div class="mt-4">
        <CommentSection :comments="comments" @add-comment="addComment" />
      </div>

      <div class="modal-action">
        <button class="btn btn-ghost btn-sm" @click="showDetails = false">Close</button>
        <button class="btn btn-primary btn-sm" @click="saveEdits">Save</button>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button @click="showDetails = false">close</button>
    </form>
  </dialog>
</template>
