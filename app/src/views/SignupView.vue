<template lang="pug">
.page.auth-page
  .auth-card
    h1.auth-card__title Create Account
    p.auth-card__note(v-if="!authEnabled") Registration requires auth to be enabled on the server.
    form.auth-card__form(@submit.prevent="handleSignup")
      .form-group
        label(for="name") Display Name
        input#name(
          type="text"
          v-model="name"
          placeholder="Your Name"
          autocomplete="name"
          required
        )
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
          placeholder="Min 8 characters"
          autocomplete="new-password"
          required
          minlength="8"
        )
      p.error-msg(v-if="error") {{ error }}
      p.success-msg(v-if="success") {{ success }}
      button.btn.btn--primary(type="submit" :disabled="loading" style="width:100%")
        span(v-if="loading") Creating account…
        span(v-else) Sign Up
    p.auth-card__footer
      | Already have an account?&nbsp;
      router-link(to="/login") Sign In
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '../api/auth'

const router = useRouter()

// Auth is disabled by default per SERVICE_AUTH_ENABLE=false
const authEnabled = false

const name = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const success = ref('')

async function handleSignup() {
  loading.value = true
  error.value = ''
  success.value = ''
  try {
    await authApi.register(name.value, email.value, password.value)
    success.value = 'Account created! Redirecting to login…'
    setTimeout(() => router.push('/login'), 1500)
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Registration failed'
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
    margin-bottom: 1.5rem;
    color: @text-primary;
  }

  &__note {
    font-size: @fs-sm;
    color: @warning;
    text-align: center;
    margin-bottom: 1rem;
    padding: 0.5rem;
    background-color: rgba(255, 152, 0, 0.1);
    border-radius: @radius-sm;
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

.success-msg {
  color: @success;
  font-size: @fs-sm;
  text-align: center;
  padding: 0.5rem 0;
}
</style>
