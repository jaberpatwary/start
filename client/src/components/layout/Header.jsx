import React, { useState, useEffect } from 'react'
import { Link, NavLink, useNavigate } from 'react-router-dom'
import { Search, ShoppingCart, Heart, Scale, User, ShieldCheck, Phone, MapPin, Menu, X, ChevronDown } from 'lucide-react'
import { useAuth } from '../../context/AuthContext'
import { useShop } from '../../context/ShopContext'
import api from '../../api/client'

const MENU_CATEGORIES = [
  { name: 'Desktop', slug: 'desktop' },
  { name: 'Laptop', slug: 'laptop' },
  { name: 'Component', slug: 'component' },
  { name: 'Monitor', slug: 'monitor' },
]

export default function Header() {
  const { user, isAuthenticated, isAdmin, logout } = useAuth()
  const { cartCount, wishlistCount, compareCount, grandTotal } = useShop()
  const [searchQuery, setSearchQuery] = useState('')
  const [suggestions, setSuggestions] = useState([])
  const [isSearching, setIsSearching] = useState(false)
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const [userDropdownOpen, setUserDropdownOpen] = useState(false)
  const navigate = useNavigate()

  // Live search debounce
  useEffect(() => {
    if (!searchQuery.trim()) {
      setSuggestions([])
      return
    }
    const timer = setTimeout(async () => {
      try {
        setIsSearching(true)
        const data = await api.get(`/products?search=${encodeURIComponent(searchQuery)}&limit=5`)
        setSuggestions(data?.results || [])
      } catch {
        setSuggestions([])
      } finally {
        setIsSearching(false)
      }
    }, 250)
    return () => clearTimeout(timer)
  }, [searchQuery])

  const handleSearchSubmit = (e) => {
    e.preventDefault()
    if (searchQuery.trim()) {
      setSuggestions([])
      navigate(`/search?q=${encodeURIComponent(searchQuery.trim())}`)
    }
  }

  return (
    <header>
      {/* Top Utility Bar */}
      <div className="top-bar">
        <div className="container">
          <div style={{ display: 'flex', alignItems: 'center', gap: '20px' }}>
            <span><Phone size={12} style={{ display: 'inline', marginRight: 4 }} /> Hotline: <b>16793</b> (9 AM - 8 PM)</span>
            <span><MapPin size={12} style={{ display: 'inline', marginRight: 4 }} /> Store Locator: 18+ Outlets Across Bangladesh</span>
          </div>
          <div className="top-bar-right">
            <Link to="/order-tracking"><ShieldCheck size={12} style={{ display: 'inline', marginRight: 4 }} /> Order Tracking</Link>
            <Link to="/emi-info">EMI Facility</Link>
            <Link to="/warranty-policy">Warranty Policy</Link>
            {isAdmin && (
              <Link to="/admin" style={{ color: '#ef4a23', fontWeight: 700 }}>
                ⚙️ Admin Panel
              </Link>
            )}
          </div>
        </div>
      </div>

      {/* Main Header */}
      <div className="header-main">
        <div className="container">
          {/* Logo */}
          <Link to="/" className="logo" style={{ display: 'flex', alignItems: 'center', gap: 10, textDecoration: 'none' }}>
            <img src="/icon-dark.png" alt="MI-Tech" style={{ height: 38, width: 'auto', objectFit: 'contain' }} />
            <span style={{ fontSize: 24, fontWeight: 900, letterSpacing: '-0.5px', color: '#fff', display: 'flex', alignItems: 'center' }}>
              MI-<span style={{ color: '#00C2FF' }}>Tech</span>
            </span>
          </Link>

          {/* Search Box with Autocomplete */}
          <div className="search-box-wrapper">
            <form className="search-form" onSubmit={handleSearchSubmit}>
              <input
                type="text"
                className="search-input"
                placeholder="Search for Desktop, Laptop, Graphics Card, Monitor..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              <button type="submit" className="search-btn">
                <Search size={18} />
              </button>
            </form>

            {/* Suggestions dropdown */}
            {suggestions.length > 0 && (
              <div style={{
                position: 'absolute',
                top: '100%',
                left: 0,
                right: 0,
                background: '#fff',
                borderRadius: '0 0 8px 8px',
                boxShadow: '0 10px 25px rgba(0,0,0,0.2)',
                zIndex: 200,
                color: '#1f2937',
                marginTop: 2,
                overflow: 'hidden',
                border: '1px solid #e2e8f0',
              }}>
                {suggestions.map((p) => (
                  <Link
                    key={p.id}
                    to={`/product/${p.slug}`}
                    onClick={() => {
                      setSuggestions([])
                      setSearchQuery('')
                    }}
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      gap: 12,
                      padding: '10px 16px',
                      borderBottom: '1px solid #f1f5f9',
                      transition: 'background 0.2s',
                    }}
                    onMouseEnter={(e) => (e.currentTarget.style.background = '#f8fafc')}
                    onMouseLeave={(e) => (e.currentTarget.style.background = '#fff')}
                  >
                    <img
                      src={p.images?.[0]?.url || 'https://placehold.co/60x60'}
                      alt={p.name}
                      style={{ width: 40, height: 40, objectFit: 'contain' }}
                    />
                    <div style={{ flex: 1 }}>
                      <div style={{ fontSize: 13, fontWeight: 600, color: '#0f172a' }}>{p.name}</div>
                      <div style={{ fontSize: 11, color: '#64748b' }}>{p.brand?.name} · {p.category?.name}</div>
                    </div>
                    <div style={{ fontSize: 14, fontWeight: 700, color: '#ef4a23' }}>
                      ৳{(p.discount_price || p.price).toLocaleString()}
                    </div>
                  </Link>
                ))}
                <button
                  onClick={handleSearchSubmit}
                  style={{
                    width: '100%',
                    padding: '10px',
                    textAlign: 'center',
                    background: '#f8fafc',
                    color: '#3749bb',
                    fontWeight: 600,
                    fontSize: 13,
                  }}
                >
                  View all results for "{searchQuery}"
                </button>
              </div>
            )}
          </div>

          {/* Action Buttons */}
          <div className="header-actions">
            {/* Compare */}
            <Link to="/compare" className="action-item" title="Compare Products">
              <div className="action-icon-wrapper">
                <Scale size={20} />
                {compareCount > 0 && <span className="badge-count">{compareCount}</span>}
              </div>
              <div className="action-text">
                <small>Compare</small>
                <strong>({compareCount})</strong>
              </div>
            </Link>

            {/* Wishlist */}
            <Link to="/wishlist" className="action-item" title="My Wishlist">
              <div className="action-icon-wrapper">
                <Heart size={20} />
                {wishlistCount > 0 && <span className="badge-count">{wishlistCount}</span>}
              </div>
              <div className="action-text">
                <small>Wishlist</small>
                <strong>({wishlistCount})</strong>
              </div>
            </Link>

            {/* Cart */}
            <Link to="/cart" className="action-item" title="Shopping Cart">
              <div className="action-icon-wrapper">
                <ShoppingCart size={20} />
                {cartCount > 0 && <span className="badge-count">{cartCount}</span>}
              </div>
              <div className="action-text">
                <small>Cart</small>
                <strong>৳{grandTotal.toLocaleString()}</strong>
              </div>
            </Link>

            {/* Account dropdown */}
            <div style={{ position: 'relative' }}>
              <button
                className="action-item"
                onClick={() => setUserDropdownOpen(!userDropdownOpen)}
                style={{ cursor: 'pointer' }}
              >
                <div className="action-icon-wrapper">
                  <User size={20} />
                </div>
                <div className="action-text">
                  <small>{isAuthenticated ? 'Hello,' : 'Account'}</small>
                  <strong>{isAuthenticated ? user.name.split(' ')[0] : 'Login'}</strong>
                </div>
                <ChevronDown size={14} />
              </button>

              {userDropdownOpen && (
                <div
                  style={{
                    position: 'absolute',
                    top: '100%',
                    right: 0,
                    background: '#fff',
                    borderRadius: 8,
                    boxShadow: '0 10px 25px rgba(0,0,0,0.15)',
                    padding: 8,
                    minWidth: 180,
                    zIndex: 200,
                    border: '1px solid #e2e8f0',
                  }}
                  onClick={() => setUserDropdownOpen(false)}
                >
                  {isAuthenticated ? (
                    <>
                      <Link to="/account" style={{ display: 'block', padding: '8px 12px', fontSize: 13, fontWeight: 600, color: '#1e293b', borderRadius: 4 }}>
                        👤 My Profile & Orders
                      </Link>
                      {isAdmin && (
                        <Link to="/admin" style={{ display: 'block', padding: '8px 12px', fontSize: 13, fontWeight: 600, color: '#ef4a23', borderRadius: 4 }}>
                          ⚙️ Admin Dashboard
                        </Link>
                      )}
                      <button
                        onClick={logout}
                        style={{ width: '100%', textAlign: 'left', padding: '8px 12px', fontSize: 13, fontWeight: 600, color: '#ef4444', borderRadius: 4 }}
                      >
                        🚪 Logout
                      </button>
                    </>
                  ) : (
                    <>
                      <Link to="/login" style={{ display: 'block', padding: '8px 12px', fontSize: 13, fontWeight: 600, color: '#1e293b', borderRadius: 4 }}>
                        🔑 Login
                      </Link>
                      <Link to="/register" style={{ display: 'block', padding: '8px 12px', fontSize: 13, fontWeight: 600, color: '#1e293b', borderRadius: 4 }}>
                        📝 Register
                      </Link>
                    </>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Navigation Menu Bar (Exactly 4 Categories) */}
      <nav className="nav-bar">
        <div className="container">
          <div className="menu-categories">
            {MENU_CATEGORIES.map((cat) => (
              <NavLink
                key={cat.slug}
                to={`/category/${cat.slug}`}
                className={({ isActive }) => `menu-item ${isActive ? 'active' : ''}`}
              >
                {cat.name}
              </NavLink>
            ))}
          </div>

          <div className="nav-right-links">
            <Link to="/compare" style={{ color: '#475569' }}>
              Compare ({compareCount})
            </Link>
            <Link to="/category/component" className="pc-builder-btn">
              ⚡ PC Builder
            </Link>
          </div>
        </div>
      </nav>
    </header>
  )
}
