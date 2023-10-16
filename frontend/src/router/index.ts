import { createRouter, createWebHashHistory } from 'vue-router'

/**
 * 基础路由
 * @type { *[] }
 */
const router: any[] = [
  {
    path: '/',
    name: 'keyboard',
    component: () => import('../views/keyboard/Keyboard.vue')
  }
]

const Router = createRouter({
  history: createWebHashHistory(),
  routes: router,
})

export default Router
