<template lang="pug">
.page.auth-page
  .auth-card
    h1.auth-card__title Sign In
    form.auth-card__form(@submit.prevent="handleLogin")
      .form-group
        label(for="email") Email
        input#email(
          type="email"
          v-model="email"
          placeholder="you@example.com"
          autocomplete="email"
          required
        )
      .form-group
        label(for="password") Password
        input#password(
          type="password"
          v-model="password"
          placeholder="••••••••"
          autocomplete="current-password"
          required
        )
      p.error-msg(v-if="error") {{ error }}
      button.btn.btn--primary(type="submit" :disabled="loading" style="width:100%")
        span(v-if="loading") Signing in…
        span(v-else) Sign In
    p.auth-card__footer
      | Don't have an account?&nbsp;
      router-link(to="/signup") Sign Up
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  loading.value = true
  error.value = ''
  try {
    await authStore.login(email.value, password.value)
    router.push('/books')
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>

<style lang="less" scoped>
.auth-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 60px);
}

.auth-card {
  background-color: @bg-card;
  border: 1px solid @border-color;
  border-radius: @radius-lg;
  padding: 2.5rem;
  width: 100%;
  max-width: 420px;
  box-shadow: @shadow-lg;

  &__title {
    font-size: @fs-3xl;
    font-weight: 700;
    text-align: center;
    margin-bottom: 2rem;
    color: @text-primary;
  }

  &__form {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  &__footer {
    margin-top: 1.5rem;
    text-align: center;
    font-size: @fs-sm;
    color: @text-secondary;
  }
}
</style>
