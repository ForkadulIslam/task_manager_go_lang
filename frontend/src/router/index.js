import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth';
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import AppLayout from '../layouts/AppLayout.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/',
      component: AppLayout,
      meta: { requiresAuth: true }, // This meta field marks the route as protected
      children: [
        {
          path: '',
          redirect: '/dashboard' // Redirect root of layout to dashboard
        },
        {
          path: 'dashboard',
          name: 'dashboard',
          component: DashboardView
        },
        {
          path: 'tasks',
          name: 'tasks',
          // route level code-splitting
          // this generates a separate chunk (About.[hash].js) for this route
          // which is lazy-loaded when the route is visited.
          component: () => import('../views/TasksListView.vue')
        }
        // Other authenticated routes like /tasks will be added here
      ]
    }
  ]
})

// Navigation Guard
router.beforeEach((to, from, next) => {
  // We need to initialize the store here because we are outside a component
  const authStore = useAuthStore();
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);

  if (requiresAuth && !authStore.token) {
    // If the route requires auth and user is not authenticated,
    // redirect to the login page.
    next('/login');
  } else {
    // Otherwise, allow navigation.
    next();
  }
});

export default router
