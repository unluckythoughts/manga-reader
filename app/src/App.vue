<template lang="pug">
#app-root
  NavBar
  main.main-content
    router-view
</template>

<script setup lang="ts">
import NavBar from './components/NavBar.vue'
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
