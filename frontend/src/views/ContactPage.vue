<template>
  <MainLayout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="max-w-3xl mx-auto">
        <h1 class="text-3xl font-extrabold text-gray-900 dark:text-white sm:text-4xl">
          Get in Touch
        </h1>
        <p class="mt-4 text-lg text-gray-500 dark:text-gray-300">
          Have a question or want to work together? Fill out the form below and I'll get back to you as soon as possible.
        </p>
        
        <div class="mt-12">
          <form @submit.prevent="submitForm" class="grid grid-cols-1 gap-y-6 sm:grid-cols-2 sm:gap-x-8">
            <div class="sm:col-span-2">
              <div v-if="formStatus.success" class="rounded-md bg-green-50 dark:bg-green-900/20 p-4 mb-6">
                <div class="flex">
                  <div class="flex-shrink-0">
                    <svg class="h-5 w-5 text-green-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                    </svg>
                  </div>
                  <div class="ml-3">
                    <p class="text-sm font-medium text-green-800 dark:text-green-200">
                      Your message has been sent successfully! I'll get back to you soon.
                    </p>
                  </div>
                </div>
              </div>
              
              <div v-if="formStatus.error" class="rounded-md bg-red-50 dark:bg-red-900/20 p-4 mb-6">
                <div class="flex">
                  <div class="flex-shrink-0">
                    <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                    </svg>
                  </div>
                  <div class="ml-3">
                    <p class="text-sm font-medium text-red-800 dark:text-red-200">
                      {{ formStatus.errorMessage }}
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <div>
              <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Name</label>
              <div class="mt-1">
                <input 
                  type="text" 
                  name="name" 
                  id="name" 
                  v-model="formData.name"
                  autocomplete="name" 
                  required
                  class="py-3 px-4 block w-full shadow-sm focus:ring-blue-500 focus:border-blue-500 border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-white"
                >
              </div>
            </div>
            
            <div>
              <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Email</label>
              <div class="mt-1">
                <input 
                  type="email" 
                  name="email" 
                  id="email" 
                  v-model="formData.email"
                  autocomplete="email" 
                  required
                  class="py-3 px-4 block w-full shadow-sm focus:ring-blue-500 focus:border-blue-500 border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-white"
                >
              </div>
            </div>
            
            <div class="sm:col-span-2">
              <label for="subject" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Subject</label>
              <div class="mt-1">
                <input 
                  type="text" 
                  name="subject" 
                  id="subject" 
                  v-model="formData.subject"
                  class="py-3 px-4 block w-full shadow-sm focus:ring-blue-500 focus:border-blue-500 border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-white"
                >
              </div>
            </div>
            
            <div class="sm:col-span-2">
              <label for="message" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Message</label>
              <div class="mt-1">
                <textarea 
                  id="message" 
                  name="message" 
                  rows="4" 
                  v-model="formData.message"
                  required
                  class="py-3 px-4 block w-full shadow-sm focus:ring-blue-500 focus:border-blue-500 border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-white"
                ></textarea>
              </div>
            </div>
            
            <div class="sm:col-span-2">
              <button 
                type="submit" 
                :disabled="formStatus.loading"
                class="w-full inline-flex items-center justify-center px-6 py-3 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <svg v-if="formStatus.loading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ formStatus.loading ? 'Sending...' : 'Send Message' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup>
import { ref, reactive } from 'vue';
import MainLayout from '@/layouts/MainLayout.vue';
import api from '@/services/api';

const formData = reactive({
  name: '',
  email: '',
  subject: '',
  message: ''
});

const formStatus = reactive({
  loading: false,
  success: false,
  error: false,
  errorMessage: ''
});

const resetForm = () => {
  formData.name = '';
  formData.email = '';
  formData.subject = '';
  formData.message = '';
};

const submitForm = async () => {
  // Reset status
  formStatus.loading = true;
  formStatus.success = false;
  formStatus.error = false;
  formStatus.errorMessage = '';
  
  try {
    await api.submitContact({
      name: formData.name,
      email: formData.email,
      subject: formData.subject,
      message: formData.message
    });
    
    formStatus.success = true;
    resetForm();
    
    // Hide success message after 5 seconds
    setTimeout(() => {
      formStatus.success = false;
    }, 5000);
    
  } catch (error) {
    formStatus.error = true;
    formStatus.errorMessage = error.response?.data?.error || 'Something went wrong. Please try again later.';
    
    // Hide error message after 5 seconds
    setTimeout(() => {
      formStatus.error = false;
    }, 5000);
  } finally {
    formStatus.loading = false;
  }
};
</script>