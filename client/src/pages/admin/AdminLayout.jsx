import React from 'react'
import { Link, NavLink, Outlet, useNavigate } from 'react-router-dom'
import { LayoutDashboard, Package, Tags, ShoppingBag, Users, Star, Boxes, BarChart3, Settings, LogOut, ArrowLeft } from 'lucide-react'
import { useAuth } from '../../context/AuthContext'

export default function AdminLayout() {
  const { user, logout, isAdmin } = useAuth()
  const navigate = useNavigate()

  if (!isAdmin) {
    return (
      <div className="container" style={{ padding: '80px 16px', textAlign: 'center' }}>
        <h2>Access Denied: Admin privileges required</h2>
        <p style={{ color: '#64748b', margin: '12px 0 24px' }}>
          Please login with an administrator account (e.g. admin@startech.local / Admin@123)
        </p>
        <Link to="/login?redirect=/admin" className="btn btn-primary">
          Login as Admin
        </Link>
      </div>
    )
  }

  return (
    <div className="admin-layout">
      {/* Sidebar */}
      <aside className="admin-sidebar">
        <div className="admin-brand">
          Star<span>Tech</span> Admin
        </div>

        <nav className="admin-menu">
          <NavLink end to="/admin" className={({ isActive }) => `admin-nav-link ${isActive ? 'active' : ''}`}>
            <LayoutDashboard size={18} /> Dashboard
          </NavLink>
          <NavLink to="/admin/products" className={({ isActive }) => `admin-nav-link ${isActive ? 'active' : ''}`}>
            <Package size={18} /> Products CRUD
          </NavLink>
          <NavLink to="/admin/catalog" className={({ isActive }) => `admin-nav-link ${isActive ? 'active' : ''}`}>
            <Tags size={18} /> Categories & Brands
          </NavLink>
          <NavLink to="/admin/orders" className={({ isActive }) => `admin-nav-link ${isActive ? 'active' : ''}`}>
            <ShoppingBag size={18} /> Orders Management
          </NavLink>
          <NavLink to="/admin/users-reviews" className={({ isActive }) => `admin-nav-link ${isActive ? 'active' : ''}`}>
            <Users size={18} /> Users & Reviews
          </NavLink>
          <NavLink to="/admin/inventory" className={({ isActive }) => `admin-nav-link ${isActive ? 'active' : ''}`}>
            <Boxes size={18} /> Inventory & Stock
          </NavLink>
          <NavLink to="/admin/reports" className={({ isActive }) => `admin-nav-link ${isActive ? 'active' : ''}`}>
            <BarChart3 size={18} /> Sales Reports
          </NavLink>
          <NavLink to="/admin/settings" className={({ isActive }) => `admin-nav-link ${isActive ? 'active' : ''}`}>
            <Settings size={18} /> Store Settings
          </NavLink>
        </nav>

        <div style={{ padding: 16, borderTop: '1px solid #1e293b' }}>
          <Link to="/" style={{ display: 'flex', alignItems: 'center', gap: 8, color: '#94a3b8', fontSize: 13, marginBottom: 12 }}>
            <ArrowLeft size={16} /> Back to Customer Store
          </Link>
          <button
            onClick={() => {
              logout()
              navigate('/login')
            }}
            style={{ display: 'flex', alignItems: 'center', gap: 8, color: '#ef4444', fontSize: 13, fontWeight: 600 }}
          >
            <LogOut size={16} /> Logout Admin
          </button>
        </div>
      </aside>

      {/* Main Content Area */}
      <div style={{ display: 'flex', flexDirection: 'column', flex: 1, minHeight: '100vh', overflowX: 'hidden' }}>
        <header className="admin-topbar">
          <div>
            <span style={{ fontSize: 14, color: '#64748b' }}>Store Administration Console</span>
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <span style={{ fontSize: 13, fontWeight: 600 }}>{user?.name} (Admin)</span>
            <div style={{
              width: 36,
              height: 36,
              borderRadius: '50%',
              background: '#ef4a23',
              color: '#fff',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontWeight: 700,
            }}>
              A
            </div>
          </div>
        </header>

        <main className="admin-main">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
