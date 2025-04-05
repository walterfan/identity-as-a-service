<template>
  <MainLayout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Hero section -->
      <div class="py-12 md:py-20">
        <div class="text-center">
          <h1 class="text-4xl tracking-tight font-extrabold text-gray-900 dark:text-white sm:text-5xl md:text-6xl">
            <span class="block">Hi, I'm</span>
            <span class="block text-blue-600 dark:text-blue-400">Your Name</span>
          </h1>
          <p class="mt-3 max-w-md mx-auto text-base text-gray-500 dark:text-gray-300 sm:text-lg md:mt-5 md:text-xl md:max-w-3xl">
            I'm a full-stack developer passionate about building web applications with modern technologies.
          </p>
          <div class="mt-5 max-w-md mx-auto sm:flex sm:justify-center md:mt-8">
            <div class="rounded-md shadow">
              <router-link to="/projects" class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 md:py-4 md:text-lg md:px-10">
                View My Work
              </router-link>
            </div>
            <div class="mt-3 rounded-md shadow sm:mt-0 sm:ml-3">
              <router-link to="/contact" class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-blue-600 bg-white hover:bg-gray-50 dark:text-blue-400 dark:bg-gray-800 dark:hover:bg-gray-700 md:py-4 md:text-lg md:px-10">
                Contact Me
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- Featured projects section -->
      <div class="py-12">
        <div class="text-center">
          <h2 class="text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white sm:text-4xl">
            Featured Projects
          </h2>
          <p class="mt-3 max-w-2xl mx-auto text-xl text-gray-500 dark:text-gray-300 sm:mt-4">
            Check out some of my recent work
          </p>
        </div>

        <div class="mt-10" v-if="projects.length > 0">
          <div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
            <div v-for="project in projects" :key="project.id" class="flex flex-col rounded-lg shadow-lg overflow-hidden bg-white dark:bg-gray-800">
              <div class="flex-shrink-0">
                <img v-if="project.imageURL" :src="project.imageURL" alt="" class="h-48 w-full object-cover">
                <div v-else class="h-48 w-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
                  <svg class="h-12 w-12 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                </div>
              </div>
              <div class="flex-1 p-6 flex flex-col justify-between">
                <div class="flex-1">
                  <p class="text-sm font-medium text-blue-600 dark:text-blue-400">
                    {{ project.tags ? project.tags.split(',')[0] : 'Project' }}
                  </p>
                  <router-link :to="`/projects/${project.id}`" class="block mt-2">
                    <p class="text-xl font-semibold text-gray-900 dark:text-white">{{ project.title }}</p>
                    <p class="mt-3 text-base text-gray-500 dark:text-gray-300">{{ project.description }}</p>
                  </router-link>
                </div>
                <div class="mt-6 flex items-center">
                  <div class="flex space-x-2">
                    <a v-if="project.projectURL" :href="project.projectURL" target="_blank" class="text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white">
                      <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                      </svg>
                    </a>
                    <a v-if="project.githubURL" :href="project.githubURL" target="_blank" class="text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white">
                      <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                        <path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd" />
                      </svg>
                    </a>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="mt-10 text-center">
            <router-link to="/projects" class="inline-flex items-center px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700">
              View All Projects
            </router-link>
          </div>
        </div>
        <div v-else class="mt-10 text-center text-gray-500 dark:text-gray-400">
          <p>No projects to display yet.</p>
        </div>
      </div>

      <!-- Latest blog posts section -->
      <div class="py-12">
        <div class="text-center">
          <h2 class="text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white sm:text-4xl">
            Latest Blog Posts
          </h2>
          <p class="mt-3 max-w-2xl mx-auto text-xl text-gray-500 dark:text-gray-300 sm:mt-4">
            Thoughts, ideas, and tutorials
          </p>
        </div>

        <div class="mt-10" v-if="posts.length > 0">
          <div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
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
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="mt-10 text-center">
            <router-link to="/blog" class="inline-flex items-center px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700">
              View All Posts
            </router-link>
          </div>
        </div>
        <div v-else class="mt-10 text-center text-gray-500 dark:text-gray-400">
          <p>No blog posts to display yet.</p>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import MainLayout from '@/layouts/MainLayout.vue';
import api from '@/services/api';

const projects = ref([]);
const posts = ref([]);

onMounted(async () => {
  try {
    // Fetch featured projects
    const projectsResponse = await api.getProjects({ featured: true });
    projects.value = projectsResponse.data.projects.slice(0, 3);
    
    // Fetch latest blog posts
    const postsResponse = await api.getPosts({ limit: 3 });
    posts.value = postsResponse.data.posts;
  } catch (error) {
    console.error('Error fetching data:', error);
  }
});

const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
};
</script>