import axios from 'axios'

const API_URL = 'http://localhost:8000'

export default {
  register(user) {
    return axios.post(`${API_URL}/register`, user)
  },

  login(user) {
    return axios.post(`${API_URL}/login`, user)
      .then(response => {
        if (response.data.token) {
          localStorage.setItem('token', response.data.token)
        }
        return response.data
      })
  },

  logout() {
    localStorage.removeItem('token')
  },

  getProtectedContent() {
    return axios.get(`${API_URL}/protected`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
  },

  isAuthenticated() {
    return localStorage.getItem('token') !== null
  }
}
