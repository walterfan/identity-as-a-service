import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

// Public pages
import HomePage from '@/views/HomePage.vue';
import ProjectsPage from '@/views/ProjectsPage.vue';
import ProjectDetailPage from '@/views/ProjectDetailPage.vue';
import BlogPage from '@/views/BlogPage.vue';
import PostDetailPage from '@/views/PostDetailPage.vue';
import ContactPage from '@/views/ContactPage.vue';
import LoginPage from '@/views/LoginPage.vue';
import RegisterPage from '@/views/RegisterPage.vue';

// Admin pages
import AdminDashboard from '@/views/admin/DashboardPage.vue';
import AdminProjects from '@/views/admin/ProjectsPage.vue';
import AdminProjectEdit from '@/views/admin/ProjectEditPage.vue';
import AdminPosts from '@/views/admin/PostsPage.vue';
import AdminPostEdit from '@/views/admin/PostEditPage.vue';
import AdminContacts from '@/views/admin/ContactsPage.vue';
import AdminProfile from '@/views/admin/ProfilePage.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // Public routes
    {
      path: '/',
      name: 'home',
      component: HomePage,
      meta: { title: 'Home' }
    },
    {
      path: '/projects',
      name: 'projects',
      component: ProjectsPage,
      meta: { title: 'Projects' }
    },
    {
      path: '/projects/:id',
      name: 'project-detail',
      component: ProjectDetailPage,
      meta: { title: 'Project Detail' }
    },
    {
      path: '/blog',
      name: 'blog',
      component: BlogPage,
      meta: { title: 'Blog' }
    },
    {
      path: '/blog/:slug',
      name: 'post-detail',
      component: PostDetailPage,
      meta: { title: 'Blog Post' }
    },
    {
      path: '/contact',
      name: 'contact',
      component: ContactPage,
      meta: { title: 'Contact' }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginPage,
      meta: { title: 'Login', guest: true }
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterPage,
      meta: { title: 'Register', guest: true }
    },
    
    // Admin routes
    {
      path: '/admin',
      name: 'admin',
      component: AdminDashboard,
      meta: { requiresAuth: true, admin: true, title: 'Admin Dashboard' }
    },
    {
      path: '/admin/projects',
      name: 'admin-projects',
      component: AdminProjects,
      meta: { requiresAuth: true, admin: true, title: 'Manage Projects' }
    },
    {
      path: '/admin/projects/new',
      name: 'admin-project-new',
      component: AdminProjectEdit,
      meta: { requiresAuth: true, admin: true, title: 'New Project' }
    },
    {
      path: '/admin/projects/:id',
      name: 'admin-project-edit',
      component: AdminProjectEdit,
      meta: { requiresAuth: true, admin: true, title: 'Edit Project' }
    },
    {
      path: '/admin/posts',
      name: 'admin-posts',
      component: AdminPosts,
      meta: { requiresAuth: true, admin: true, title: 'Manage Posts' }
    },
    {
      path: '/admin/posts/new',
      name: 'admin-post-new',
      component: AdminPostEdit,
      meta: { requiresAuth: true, admin: true, title: 'New Post' }
    },
    {
      path: '/admin/posts/:id',
      name: 'admin-post-edit',
      component: AdminPostEdit,
      meta: { requiresAuth: true, admin: true, title: 'Edit Post' }
    },
    {
      path: '/admin/contacts',
      name: 'admin-contacts',
      component: AdminContacts,
      meta: { requiresAuth: true, admin: true, title: 'Manage Contacts' }
    },
    {
      path: '/admin/profile',
      name: 'admin-profile',
      component: AdminProfile,
      meta: { requiresAuth: true, title: 'Profile Settings' }
    },
    
    // 404 route
    {
      path: '/:pathMatch(.*)*',
      redirect: '/'
    }
  ]
});

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const requiresAdmin = to.matched.some(record => record.meta.admin);
  const isGuestOnly = to.matched.some(record => record.meta.guest);
  
  // Set page title
  document.title = to.meta.title ? `${to.meta.title} | Personal Website` : 'Personal Website';
  
  if (requiresAuth && !authStore.isAuthenticated) {
    next('/login');
  } else if (requiresAdmin && !authStore.isAdmin) {
    next('/');
  } else if (isGuestOnly && authStore.isAuthenticated) {
    next('/');
  } else {
    next();
  }
});

export default router;