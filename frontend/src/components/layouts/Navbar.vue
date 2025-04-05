<template>
  <nav class="bg-white dark:bg-gray-800 shadow">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex">
          <div class="flex-shrink-0 flex items-center">
            <router-link to="/" class="text-xl font-bold text-gray-800 dark:text-white">
              My Portfolio
            </router-link>
          </div>
          <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
            <router-link 
              to="/" 
              class="inline-flex items-center px-1 pt-1 border-b-2" 
              :class="[$route.path === '/' ? 'border-blue-500 text-gray-900 dark:text-white' : 'border-transparent text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200']"
            >
              Home
            </router-link>
            <router-link 
              to="/projects" 
              class="inline-flex items-center px-1 pt-1 border-b-2" 
              :class="[$route.path.startsWith('/projects') ? 'border-blue-500 text-gray-900 dark:text-white' : 'border-transparent text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200']"
            >
              Projects
            </router-link>
            <router-link 
              to="/blog" 
              class="inline-flex items-center px-1 pt-1 border-b-2" 
              :class="[$route.path.startsWith('/blog') ? 'border-blue-500 text-gray-900 dark:text-white' : 'border-transparent text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200']"
            >
              Blog
            </router-link>
            <router-link 
              to="/contact" 
              class="inline-flex items-center px-1 pt-1 border-b-2" 
              :class="[$route.path === '/contact' ? 'border-blue-500 text-gray-900 dark:text-white' : 'border-transparent text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200']"
            >
              Contact
            </router-link>
          </div>
        </div>
        <div class="hidden sm:ml-6 sm:flex sm:items-center">
          <!-- Dark mode toggle -->
          <button 
            @click="toggleDarkMode" 
            class="p-1 rounded-full text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200 focus:outline-none"
          >
            <span v-if="isDarkMode" class="sr-only">Switch to light mode</span>
            <span v-else class="sr-only">Switch to dark mode</span>
            <svg v-if="isDarkMode" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
          </button>

          <!-- Profile dropdown -->
          <div v-if="authStore.isAuthenticated" class="ml-3 relative">
            <div>
              <button 
                @click="isProfileOpen = !isProfileOpen" 
                class="flex text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                <span class="sr-only">Open user menu</span>
                <div v-if="authStore.user?.avatar" class="h-8 w-8 rounded-full overflow-hidden">
                  <img :src="authStore.user.avatar" alt="User avatar" class="h-full w-full object-cover" />
                </div>
                <div v-else class="h-8 w-8 rounded-full bg-blue-500 flex items-center justify-center text-white">
                  {{ authStore.user?.name?.charAt(0) || authStore.user?.username?.charAt(0) || 'U' }}
                </div>
              </button>
            </div>
            
            <div 
              v-if="isProfileOpen" 
              class="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg py-1 bg-white dark:bg-gray-700 ring-1 ring-black ring-opacity-5 focus:outline-none z-10"
            >
              <router-link 
                v-if="authStore.isAdmin" 
                to="/admin" 
                class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-600"
                @click="isProfileOpen = false"
              >
                Dashboard
              </router-link>
              <router-link 
                to="/admin/profile" 
                class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-600"
                @click="isProfileOpen = false"
              >
                Profile
              </router-link>
              <button 
                @click="logout" 
                class="block w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-600"
              >
                Sign out
              </button>
            </div>
          </div>
          
          <div v-else class="flex items-center space-x-4">
            <router-link 
              to="/login" 
              class="text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200"
            >
              Login
            </router-link>
            <router-link 
              to="/register" 
              class="btn-primary"
            >
              Register
            </router-link>
          </div>
        </div>
        
        <!-- Mobile menu button -->
        <div class="flex items-center sm:hidden">
          <button 
            @click="toggleDarkMode" 
            class="p-1 mr-2 rounded-full text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200 focus:outline-none"
          >
            <svg v-if="isDarkMode" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
          </button>
          <button 
            @click="isMobileMenuOpen = !isMobileMenuOpen" 
            class="inline-flex items-center justify-center p-2 rounded-md text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500"
          >
            <span class="sr-only">Open main menu</span>
            <svg 
              v-if="!isMobileMenuOpen" 
              class="block h-6 w-6" 
              xmlns="http://www.w3.org/2000/svg" 
              fill="none" 
              viewBox="0 0 24 24" 
              stroke="currentColor" 
              aria-hidden="true"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
            <svg 
              v-else 
              class="block h-6 w-6" 
              xmlns="http://www.w3.org/2000/svg" 
              fill="none" 
              viewBox="0 0 24 24" 
              stroke="currentColor" 
              aria-hidden="true"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>
    
    <!-- Mobile menu -->
    <div v-if="isMobileMenuOpen" class="sm:hidden">
      <div class="pt-2 pb-3 space-y-1">
        <router-link 
          to="/" 
          class="block pl-3 pr-4 py-2 border-l-4"
          :class="[$route.path === '/' ? 'border-blue-500 text-blue-700 dark:text-blue-300 bg-blue-50 dark:bg-blue-900/20' : 'border-transparent text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700']"
          @click="isMobileMenuOpen = false"
        >
          Home
        </router-link>
        <router-link 
          to="/projects" 
          class="block pl-3 pr-4 py-2 border-l-4"
          :class="[$route.path.startsWith('/projects') ? 'border-blue-500 text-blue-700 dark:text-blue-300 bg-blue-50 dark:bg-blue-900/20' : 'border-transparent text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700']"
          @click="isMobileMenuOpen = false"
        >
          Projects
        </router-link>
        <router-link 
          to="/blog" 
          class="block pl-3 pr-4 py-2 border-l-4"
          :class="[$route.path.startsWith('/blog') ? 'border-blue-500 text-blue-700 dark:text-blue-300 bg-blue-50 dark:bg-blue-900/20' : 'border-transparent text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700']"
          @click="isMobileMenuOpen = false"
        >
          Blog
        </router-link>
        <router-link 
          to="/contact" 
          class="block pl-3 pr-4 py-2 border-l-4"
          :class="[$route.path === '/contact' ? 'border-blue-500 text-blue-700 dark:text-blue-300 bg-blue-50 dark:bg-blue-900/20' : 'border-transparent text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700']"
          @click="isMobileMenuOpen = false"
        >
          Contact
        </router-link>
      </div>
      
      <div v-if="authStore.isAuthenticated" class="pt-4 pb-3 border-t border-gray-200 dark:border-gray-600">
        <div class="flex items-center px-4">
          <div class="flex-shrink-0">
            <div v-if="authStore.user?.avatar" class="h-10 w-10 rounded-full overflow-hidden">
              <img :src="authStore.user.avatar" alt="User avatar" class="h-full w-full object-cover" />
            </div>
            <div v-else class="h-10 w-10 rounded-full bg-blue-500 flex items-center justify-center text-white">
              {{ authStore.user?.name?.charAt(0) || authStore.user?.username?.charAt(0) || 'U' }}
            </div>
          </div>
          <div class="ml-3">
            <div class="text-base font-medium text-gray-800 dark:text-white">{{ authStore.user?.name }}</div>
            <div class="text-sm font-medium text-gray-500 dark:text-gray-300">{{ authStore.user?.email }}</div>
          </div>
        </div>
        <div class="mt-3 space-y-1">
          <router-link 
            v-if="authStore.isAdmin" 
            to="/admin" 
            class="block px-4 py-2 text-base font-medium text-gray-500 dark:text-gray-300 hover:text-gray-800 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700"
            @click="isMobileMenuOpen = false"
          >
            Dashboard
          </router-link>
          <router-link 
            to="/admin/profile" 
            class="block px-4 py-2 text-base font-medium text-gray-500 dark:text-gray-300 hover:text-gray-800 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700"
            @click="isMobileMenuOpen = false"
          >
            Profile
          </router-link>
          <button 
            @click="logout" 
            class="block w-full text-left px-4 py-2 text-base font-medium text-gray-500 dark:text-gray-300 hover:text-gray-800 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700"
          >
            Sign out
          </button>
        </div>
      </div>
      
      <div v-else class="pt-4 pb-3 border-t border-gray-200 dark:border-gray-600">
        <div class="space-y-1">
          <router-link 
            to="/login" 
            class="block px-4 py-2 text-base font-medium text-gray-500 dark:text-gray-300 hover:text-gray-800 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700"
            @click="isMobileMenuOpen = false"
          >
            Login
          </router-link>
          <router-link 
            to="/register" 
            class="block px-4 py-2 text-base font-medium text-gray-500 dark:text-gray-300 hover:text-gray-800 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700"
            @click="isMobileMenuOpen = false"
          >
            Register
          </router-link>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const authStore = useAuthStore();

const isProfileOpen = ref(false);
const isMobileMenuOpen = ref(false);
const isDarkMode = ref(false);

// Close dropdown when clicking outside
const handleClickOutside = (event) => {
  if (isProfileOpen.value && !event.target.closest('.relative')) {
    isProfileOpen.value = false;
  }
};

onMounted(() => {
  // Check if dark mode is enabled in localStorage
  const darkMode = localStorage.getItem('darkMode') === 'true';
  isDarkMode.value = darkMode;
  
  if (darkMode) {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
  
  document.addEventListener('click', handleClickOutside);
});

const toggleDarkMode = () => {
  isDarkMode.value = !isDarkMode.value;
  localStorage.setItem('darkMode', isDarkMode.value);
  
  if (isDarkMode.value) {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
};

const logout = () => {
  authStore.logout();
  isProfileOpen.value = false;
  isMobileMenuOpen.value = false;
  router.push('/');
};
</script>