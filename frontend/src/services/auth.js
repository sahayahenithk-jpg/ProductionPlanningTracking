import api from './api'

export const AUTH_TOKEN_KEY = 'token'
export const AUTH_ROLE_KEY = 'userRole'
export const AUTH_NAME_KEY = 'userName'

export const getToken = () => localStorage.getItem(AUTH_TOKEN_KEY)
export const getUserRole = () => localStorage.getItem(AUTH_ROLE_KEY) || ''
export const getUserName = () => localStorage.getItem(AUTH_NAME_KEY) || ''

export const setAuth = async (token) => {
  localStorage.setItem(AUTH_TOKEN_KEY, token)
  await refreshUserProfile()
}

export const refreshUserProfile = async () => {
  const response = await api.get('/user')
  localStorage.setItem(AUTH_ROLE_KEY, response.data.role)
  localStorage.setItem(AUTH_NAME_KEY, response.data.name)
  return response.data
}

export const clearAuth = () => {
  localStorage.removeItem(AUTH_TOKEN_KEY)
  localStorage.removeItem(AUTH_ROLE_KEY)
  localStorage.removeItem(AUTH_NAME_KEY)
}
