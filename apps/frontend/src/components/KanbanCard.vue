<script setup lang="ts">
import { ref } from 'vue'
import type { Card, CardStatus, Comment } from '../types'
import { COLUMNS, PRIORITY_COLORS } from '../types'
import { getComments, createComment } from '../api'
import CommentSection from './CommentSection.vue'

const props = defineProps<{ card: Card }>()
const emit = defineEmits<{
  'status-change': [cardId: string, status: CardStatus]
  delete: [cardId: string]
}>()

const showDetails = ref(false)
const comments = ref<Comment[]>([])

async function openDetails() {
  showDetails.value = true
  comments.value = await getComments(props.card.id)
}

async function addComment(author: string, body: string) {
  const comment = await createComment(props.card.id, author, body)
  comments.value.push(comment)
}
</script>

<template>
  <div
    class="card bg-base-100 shadow-sm cursor-pointer hover:shadow-md transition-shadow"
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
        <span
          v-for="label in card.labels"
          :key="label"
          class="badge badge-xs badge-outline"
        >
          {{ label }}
        </span>
        <span :class="['badge badge-xs', PRIORITY_COLORS[card.priority]]">
          {{ card.priority }}
        </span>
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
      <h3 class="text-lg font-bold">{{ card.title }}</h3>
      <p class="text-sm text-base-content/60 mt-2">{{ card.description }}</p>

      <div class="flex flex-wrap gap-2 mt-3">
        <span :class="['badge', PRIORITY_COLORS[card.priority]]">{{ card.priority }}</span>
        <span v-for="label in card.labels" :key="label" class="badge badge-outline">
          {{ label }}
        </span>
      </div>

      <div class="mt-4">
        <CommentSection
          :comments="comments"
          @add-comment="addComment"
        />
      </div>

      <div class="modal-action">
        <button class="btn btn-sm" @click="showDetails = false">Close</button>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button @click="showDetails = false">close</button>
    </form>
  </dialog>
</template>
