import axios from 'axios'
import { useAuthStore } from '../stores/AuthStore';

const API_URL = 'http://localhost:8000'

export default {
  register(user) {
    return axios.post(`${API_URL}/register`, user)
  },

  login(user_pass_pair) {
    return axios.post(`${API_URL}/login`, user_pass_pair)
      .then(response => {
        if (response.data.token) {
          localStorage.setItem('token', response.data.token)
          const authStore = useAuthStore();
          authStore.setAuthenticated(true);
        }
        return response.data
      })
  },

  logout() {
    localStorage.removeItem('token')
    const authStore = useAuthStore();
    authStore.setAuthenticated(false);
  },

  getProtectedContent() {
    return axios.get(`${API_URL}/protected`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
  },

  isAuthenticated() {
    console.log("isAuthenticated",localStorage.getItem('token') !== null)
    return localStorage.getItem('token') !== null
  }
}
