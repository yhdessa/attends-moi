<script setup lang="ts">
import { ref } from 'vue'
import type { Comment } from '../types'

defineProps<{ comments: Comment[] }>()
const emit = defineEmits<{ 'add-comment': [author: string, body: string] }>()

const author = ref('')
const body = ref('')

function submit() {
  if (!author.value.trim() || !body.value.trim()) return
  emit('add-comment', author.value, body.value)
  body.value = ''
}
</script>

<template>
  <div>
    <h4 class="font-semibold text-sm mb-2">Comments</h4>

    <div class="space-y-2 max-h-48 overflow-y-auto mb-3">
      <div
        v-for="comment in comments"
        :key="comment.id"
        class="chat chat-start"
      >
        <div class="chat-header">
          {{ comment.author }}
          <time class="text-xs opacity-50 ml-2">
            {{ new Date(comment.created_at).toLocaleString() }}
          </time>
        </div>
        <div class="chat-bubble chat-bubble-primary text-sm">{{ comment.body }}</div>
      </div>

      <p v-if="comments.length === 0" class="text-sm text-base-content/40 text-center py-4">
        No comments yet
      </p>
    </div>

    <div class="flex gap-2">
      <input
        v-model="author"
        type="text"
        placeholder="Your name"
        class="input input-bordered input-sm w-32"
      />
      <input
        v-model="body"
        type="text"
        placeholder="Write a comment..."
        class="input input-bordered input-sm flex-1"
        @keyup.enter="submit"
      />
      <button class="btn btn-primary btn-sm" @click="submit">Send</button>
    </div>
  </div>
</template>
