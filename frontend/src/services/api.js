import axios from 'axios';

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add a request interceptor to include the auth token in requests
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add a response interceptor to handle common errors
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response && error.response.status === 401) {
      // If unauthorized, clear token and redirect to login
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default {
  // Auth
  login(credentials) {
    return apiClient.post('/auth/login', credentials);
  },
  register(userData) {
    return apiClient.post('/auth/register', userData);
  },
  getProfile() {
    return apiClient.get('/profile');
  },
  updateProfile(profileData) {
    return apiClient.put('/profile', profileData);
  },
  
  // Projects
  getProjects(params = {}) {
    return apiClient.get('/projects', { params });
  },
  getProject(id) {
    return apiClient.get(`/projects/${id}`);
  },
  createProject(projectData) {
    return apiClient.post('/projects', projectData);
  },
  updateProject(id, projectData) {
    return apiClient.put(`/projects/${id}`, projectData);
  },
  deleteProject(id) {
    return apiClient.delete(`/projects/${id}`);
  },
  
  // Blog Posts
  getPosts(params = {}) {
    return apiClient.get('/posts', { params });
  },
  getPost(slug) {
    return apiClient.get(`/posts/${slug}`);
  },
  createPost(postData) {
    return apiClient.post('/posts', postData);
  },
  updatePost(id, postData) {
    return apiClient.put(`/posts/${id}`, postData);
  },
  deletePost(id) {
    return apiClient.delete(`/posts/${id}`);
  },
  
  // Contact
  submitContact(contactData) {
    return apiClient.post('/contact', contactData);
  },
  getContacts(params = {}) {
    return apiClient.get('/admin/contacts', { params });
  },
  getContact(id) {
    return apiClient.get(`/admin/contacts/${id}`);
  },
  markContactAsRead(id) {
    return apiClient.put(`/admin/contacts/${id}/read`);
  },
  deleteContact(id) {
    return apiClient.delete(`/admin/contacts/${id}`);
  },
};