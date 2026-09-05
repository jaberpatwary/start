import React, { useState } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import toast from 'react-hot-toast'
import { Lock, Mail, User, Phone, ArrowRight, ShieldCheck } from 'lucide-react'

export function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const { login } = useAuth()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const redirect = searchParams.get('redirect') || '/'

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    try {
      const user = await login(email, password)
      if (user?.role === 'ADMIN') {
        navigate('/admin')
      } else {
        navigate(redirect)
      }
    } catch (err) {
      toast.error(err.message || 'Invalid email or password')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container" style={{ padding: '60px 16px', maxWidth: 460 }}>
      <div style={{
        background: '#fff',
        borderRadius: 8,
        border: '1px solid #e2e8f0',
        padding: 32,
        boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
      }}>
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <img src="/logo.png" alt="MI-Tech" style={{ height: 65, width: 'auto', marginBottom: 12, objectFit: 'contain' }} />
          <h1 style={{ fontSize: 22, fontWeight: 800, color: '#0f172a' }}>Account Login</h1>
          <p style={{ fontSize: 13, color: '#64748b' }}>Please login to your MI-Tech account</p>
        </div>

        <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
          <div className="form-group">
            <label>Email Address or Username</label>
            <div style={{ position: 'relative' }}>
              <input
                type="text"
                required
                placeholder="admin or email address"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
          </div>

          <div className="form-group">
            <label>Password</label>
            <input
              type="password"
              required
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>

          {/* Demo Credentials Quick Fill */}
          <div style={{
            background: '#f8fafc',
            border: '1px dashed #cbd5e1',
            padding: 12,
            borderRadius: 6,
            fontSize: 12,
            color: '#475569',
          }}>
            <b style={{ color: '#0f172a', display: 'block', marginBottom: 4 }}>Demo Admin Credentials:</b>
            <span>admin / admin123</span>
            <button
              type="button"
              onClick={() => {
                setEmail('admin')
                setPassword('admin123')
              }}
              style={{
                display: 'block',
                marginTop: 6,
                color: '#ef4a23',
                fontWeight: 700,
                fontSize: 11,
                textDecoration: 'underline',
              }}
            >
              Fill Admin Credentials
            </button>
          </div>

          <button
            type="submit"
            className="btn btn-primary btn-lg"
            disabled={loading}
            style={{ width: '100%', marginTop: 8 }}
          >
            {loading ? 'Signing in...' : 'Sign In'} <ArrowRight size={16} />
          </button>
        </form>

        <div style={{ textAlign: 'center', marginTop: 24, fontSize: 13, color: '#64748b' }}>
          Don't have an account?{' '}
          <Link to={`/register?redirect=${redirect}`} style={{ color: '#ef4a23', fontWeight: 700 }}>
            Register Now
          </Link>
        </div>
      </div>
    </div>
  )
}

export function RegisterPage() {
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    phone: '',
    password: '',
    confirmPassword: '',
  })
  const [loading, setLoading] = useState(false)
  const { register } = useAuth()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const redirect = searchParams.get('redirect') || '/'

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (formData.password !== formData.confirmPassword) {
      toast.error('Passwords do not match')
      return
    }
    if (formData.password.length < 6) {
      toast.error('Password must be at least 6 characters')
      return
    }

    setLoading(true)
    try {
      await register({
        name: formData.name,
        email: formData.email,
        phone: formData.phone,
        password: formData.password,
      })
      navigate(redirect)
    } catch (err) {
      toast.error(err.message || 'Registration failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container" style={{ padding: '60px 16px', maxWidth: 480 }}>
      <div style={{
        background: '#fff',
        borderRadius: 8,
        border: '1px solid #e2e8f0',
        padding: 32,
        boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
      }}>
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <img src="/logo.png" alt="MI-Tech" style={{ height: 65, width: 'auto', marginBottom: 12, objectFit: 'contain' }} />
          <h1 style={{ fontSize: 22, fontWeight: 800, color: '#0f172a' }}>Register Account</h1>
          <p style={{ fontSize: 13, color: '#64748b' }}>Create a new MI-Tech account for fast orders and tracking</p>
        </div>

        <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
          <div className="form-group">
            <label>Full Name *</label>
            <input
              type="text"
              required
              placeholder="e.g. Asif Mahmud"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            />
          </div>

          <div className="form-group">
            <label>Email Address *</label>
            <input
              type="email"
              required
              placeholder="name@example.com"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
            />
          </div>

          <div className="form-group">
            <label>Phone Number *</label>
            <input
              type="tel"
              required
              placeholder="01700000000"
              value={formData.phone}
              onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
            />
          </div>

          <div className="form-group">
            <label>Password (min 6 chars) *</label>
            <input
              type="password"
              required
              placeholder="••••••••"
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
            />
          </div>

          <div className="form-group">
            <label>Confirm Password *</label>
            <input
              type="password"
              required
              placeholder="••••••••"
              value={formData.confirmPassword}
              onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
            />
          </div>

          <button
            type="submit"
            className="btn btn-primary btn-lg"
            disabled={loading}
            style={{ width: '100%', marginTop: 8 }}
          >
            {loading ? 'Creating account...' : 'Create Account'} <ArrowRight size={16} />
          </button>
        </form>

        <div style={{ textAlign: 'center', marginTop: 24, fontSize: 13, color: '#64748b' }}>
          Already registered?{' '}
          <Link to={`/login?redirect=${redirect}`} style={{ color: '#ef4a23', fontWeight: 700 }}>
            Login here
          </Link>
        </div>
      </div>
    </div>
  )
}
