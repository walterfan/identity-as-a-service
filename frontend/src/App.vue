<template>
  <router-view />
</template>

<script setup>
import { onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();

onMounted(() => {
  // Check if user is logged in and fetch profile
  if (authStore.isAuthenticated) {
    authStore.fetchProfile();
  }
});
</script>