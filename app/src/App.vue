<template lang="pug">
#app-root
  NavBar
  main.main-content
    router-view
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import NavBar from './components/NavBar.vue'
import { useAuthStore } from './stores/auth'

const authStore = useAuthStore()

onMounted(async () => {
  if (authStore.token) {
    authStore.fetchUser()
  } else {
    // Auto-login as dummy user when auth is disabled so all requests carry a bearer token
    await authStore.login('dummy@example.com', 'Dummy@example123')
  }
})
</script>

<style lang="less">
@import '@/assets/main.less';

#app-root {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-content {
  flex: 1;
}
</style>
