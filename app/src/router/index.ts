import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import SignupView from '../views/SignupView.vue'
import BooksView from '../views/BooksView.vue'
import BookDetailView from '../views/BookDetailView.vue'
import ChapterView from '../views/ChapterView.vue'
import FavoritesView from '../views/FavoritesView.vue'
import SourcesView from '../views/SourcesView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/books' },
    { path: '/login', component: LoginView },
    { path: '/signup', component: SignupView },
    { path: '/books', component: BooksView },
    { path: '/books/:id', component: BookDetailView },
    { path: '/chapters/:id', component: ChapterView },
    { path: '/favorites', component: FavoritesView },
    { path: '/sources', component: SourcesView }
  ],
  scrollBehavior() {
    return { top: 0 }
  }
})

export default router
