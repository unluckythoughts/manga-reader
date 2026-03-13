import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import './assets/main.less'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)

async function bootstrap() {
	// Bootstrap auth before mounting so every component has a valid token from the start.
	const authStore = useAuthStore()
	await (authStore.token ? authStore.fetchUser() : authStore.login('dummy@example.com', 'Dummy@example123'))

	app.use(router)
	app.mount('#app')
}

void bootstrap()
