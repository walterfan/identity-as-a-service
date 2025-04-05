<template>
  <MainLayout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="text-center">
        <h1 class="text-3xl font-extrabold text-gray-900 dark:text-white sm:text-4xl">
          Blog
        </h1>
        <p class="mt-3 max-w-2xl mx-auto text-xl text-gray-500 dark:text-gray-300 sm:mt-4">
          Thoughts, ideas, and tutorials
        </p>
      </div>

      <div class="mt-12">
        <div v-if="loading" class="flex justify-center">
          <svg class="animate-spin h-10 w-10 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>

        <div v-else-if="posts.length > 0" class="grid gap-10 lg:grid-cols-2">
          <div v-for="post in posts" :key="post.id" class="flex flex-col rounded-lg shadow-lg overflow-hidden bg-white dark:bg-gray-800">
            <div class="flex-shrink-0">
              <img v-if="post.imageURL" :src="post.imageURL" alt="" class="h-48 w-full object-cover">
              <div v-else class="h-48 w-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
                <svg class="h-12 w-12 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
                </svg>
              </div>
            </div>
            <div class="flex-1 p-6 flex flex-col justify-between">
              <div class="flex-1">
                <p class="text-sm font-medium text-blue-600 dark:text-blue-400">
                  {{ post.tags ? post.tags.split(',')[0] : 'Blog' }}
                </p>
                <router-link :to="`/blog/${post.slug}`" class="block mt-2">
                  <p class="text-xl font-semibold text-gray-900 dark:text-white">{{ post.title }}</p>
                  <p class="mt-3 text-base text-gray-500 dark:text-gray-300">{{ post.excerpt }}</p>
                </router-link>
              </div>
              <div class="mt-6 flex items-center">
                <div class="flex-shrink-0">
                  <span class="sr-only">{{ post.user?.name || 'Author' }}</span>
                  <div class="h-10 w-10 rounded-full bg-blue-500 flex items-center justify-center text-white">
                    {{ post.user?.name?.charAt(0) || 'A' }}
                  </div>
                </div>
                <div class="ml-3">
                  <p class="text-sm font-medium text-gray-900 dark:text-white">
                    {{ post.user?.name || 'Author' }}
                  </p>
                  <div class="flex space-x-1 text-sm text-gray-500 dark:text-gray-400">
                    <time :datetime="post.publishedAt">
                      {{ formatDate(post.publishedAt) }}
                    </time>
                    <span aria-hidden="true">&middot;</span>
                    <span>{{ post.readTime || '5 min read' }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="text-center text-gray-500 dark:text-gray-400">
          <p>No blog posts to display yet.</p>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="mt-12 flex justify-center">
          <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
            <button 
              @click="changePage(currentPage - 1)" 
              :disabled="currentPage === 1"
              class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm font-medium text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span class="sr-only">Previous</span>
              <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
            </button>
            <button 
              v-for="page in paginationItems" 
              :key="page"
              @click="changePage(page)"
              :class="[
                page === currentPage ? 'z-10 bg-blue-50 dark:bg-blue-900/20 border-blue-500 text-blue-600 dark:text-blue-400' : 'bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700',
                'relative inline-flex items-center px-4 py-2 border text-sm font-medium'
              ]"
            >
              {{ page }}
            </button>
            <button 
              @click="changePage(currentPage + 1)" 
              :disabled="currentPage === totalPages"
              class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm font-medium text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span class="sr-only">Next</span>
              <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
              </svg>
            </button>
          </nav>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import MainLayout from '@/layouts/MainLayout.vue';
import api from '@/services/api';

const route = useRoute();
const router = useRouter();
const posts = ref([]);
const loading = ref(true);
const currentPage = ref(parseInt(route.query.page) || 1);
const totalPosts = ref(0);
const postsPerPage = 10;

const totalPages = computed(() => Math.ceil(totalPosts.value / postsPerPage));

const paginationItems = computed(() => {
  const items = [];
  const maxVisiblePages = 5;
  
  if (totalPages.value <= maxVisiblePages) {
    for (let i = 1; i <= totalPages.value; i++) {
      items.push(i);
    }
  } else {
    // Always show first page
    items.push(1);
    
    // Calculate start and end of visible pages
    let start = Math.max(2, currentPage.value - 1);
    let end = Math.min(totalPages.value - 1, start + 2);
    
    // Adjust if we're near the end
    if (end === totalPages.value - 1) {
      start = Math.max(2, end - 2);
    }
    
    // Add ellipsis if needed
    if (start > 2) {
      items.push('...');
    }
    
    // Add visible pages
    for (let i = start; i <= end; i++) {
      items.push(i);
    }
    
    // Add ellipsis if needed
    if (end < totalPages.value - 1) {
      items.push('...');
    }
    
    // Always show last page
    items.push(totalPages.value);
  }
  
  return items;
});

const fetchPosts = async () => {
  loading.value = true;
  try {
    const response = await api.getPosts({
      page: currentPage.value,
      limit: postsPerPage
    });
    posts.value = response.data.posts;
    totalPosts.value = response.data.total;
  } catch (error) {
    console.error('Error fetching posts:', error);
  } finally {
    loading.value = false;
  }
};

const changePage = (page) => {
  if (typeof page === 'number' && page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
    router.push({ query: { ...route.query, page } });
  }
};

const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
};

onMounted(fetchPosts);

watch(() => route.query.page, (newPage) => {
  currentPage.value = parseInt(newPage) || 1;
  fetchPosts();
});
</script>