<script setup lang="ts">
import { onMounted } from 'vue'
import { useUserStore } from './stores/user'

const userStore = useUserStore()

onMounted(() => {
  if (userStore.token) {
    userStore.fetchProfile()
  }
})
</script>

<template>
  <router-view />
</template>
