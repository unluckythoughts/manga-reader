<template lang="pug">
nav.navbar
  .navbar__brand
    router-link(to="/books") 📖 Book Reader
  .navbar__links
    router-link(to="/books") Books
    router-link(to="/favorites") Favorites
    router-link(to="/sources") Sources
  .navbar__auth
    template(v-if="authStore.isLoggedIn")
      span.navbar__user {{ authStore.user?.name || authStore.user?.email }}
      button.btn.btn--ghost(@click="handleLogout") Logout
    template(v-else)
      router-link.btn.btn--ghost(to="/login") Login
      router-link.btn.btn--primary(to="/signup") Sign Up
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style lang="less" scoped>
.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2rem;
  height: 60px;
  background-color: @bg-secondary;
  border-bottom: 1px solid @border-color;
  position: sticky;
  top: 0;
  z-index: 100;
  gap: 1rem;

  &__brand {
    a {
      font-size: @fs-xl;
      font-weight: 700;
      color: @text-primary;
      text-decoration: none;
      white-space: nowrap;

      &:hover {
        color: @accent;
        text-decoration: none;
      }
    }
  }

  &__links {
    display: flex;
    gap: 0.25rem;
    flex: 1;
    justify-content: center;

    a {
      color: @text-secondary;
      font-weight: 500;
      padding: 0.35rem 0.85rem;
      border-radius: @radius-md;
      text-decoration: none;
      transition: color @transition-fast, background-color @transition-fast;

      &:hover {
        color: @text-primary;
        background-color: @bg-card;
        text-decoration: none;
      }

      &.router-link-active {
        color: @accent;
        background-color: @accent-muted;
      }
    }
  }

  &__auth {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    white-space: nowrap;
  }

  &__user {
    font-size: @fs-sm;
    color: @text-secondary;
    max-width: 140px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
</style>
