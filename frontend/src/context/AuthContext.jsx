import React, { createContext, useContext, useEffect, useState } from 'react'
import api from '../api/client'
import toast from 'react-hot-toast'

const AuthContext = createContext()

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [token, setToken] = useState(localStorage.getItem('startech_token') || '')
  const [loading, setLoading] = useState(true)

  const fetchCurrentUser = async () => {
    try {
      if (!token) {
        setLoading(false)
        return
      }
      const data = await api.get('/auth/me')
      if (data?.user) {
        setUser(data.user)
      }
    } catch (err) {
      console.error('Session validation error:', err.message)
      // Invalid or expired token
      localStorage.removeItem('startech_token')
      setToken('')
      setUser(null)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchCurrentUser()
  }, [token])

  const login = async (email, password) => {
    const data = await api.post('/auth/login', { email, password })
    if (data?.token) {
      localStorage.setItem('startech_token', data.token)
      setToken(data.token)
      setUser(data.user)
      toast.success(`Welcome back, ${data.user.name}!`)
      return data.user
    }
  }

  const register = async (userData) => {
    const data = await api.post('/auth/register', userData)
    if (data?.token) {
      localStorage.setItem('startech_token', data.token)
      setToken(data.token)
      setUser(data.user)
      toast.success('Account created successfully!')
      return data.user
    }
  }

  const logout = async () => {
    try {
      await api.post('/auth/logout')
    } catch {
      // ignore
    }
    localStorage.removeItem('startech_token')
    setToken('')
    setUser(null)
    toast.success('Logged out successfully')
  }

  const updateProfile = async (profileData) => {
    const data = await api.put('/users/profile', profileData)
    if (data?.user) {
      setUser(data.user)
      toast.success('Profile updated successfully')
    }
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        loading,
        isAuthenticated: !!user,
        isAdmin: user?.role === 'ADMIN',
        login,
        register,
        logout,
        updateProfile,
        refreshUser: fetchCurrentUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => useContext(AuthContext)
