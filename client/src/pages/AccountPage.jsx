import React, { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { useShop } from '../context/ShopContext'
import api from '../api/client'
import toast from 'react-hot-toast'
import { User, Package, MapPin, Heart, ShieldCheck, LogOut, ChevronRight, Plus, Trash2, Search } from 'lucide-react'

export default function AccountPage() {
  const { user, isAuthenticated, logout, updateProfile, refreshUser } = useAuth()
  const { wishlist, removeFromCart, addToCart, toggleWishlist } = useShop()
  const navigate = useNavigate()

  const [activeTab, setActiveTab] = useState('orders') // 'profile' | 'orders' | 'tracking' | 'addresses' | 'wishlist'
  const [orders, setOrders] = useState([])
  const [loadingOrders, setLoadingOrders] = useState(false)

  // Profile Edit State
  const [name, setName] = useState('')
  const [phone, setPhone] = useState('')

  // New Address State
  const [newAddress, setNewAddress] = useState({
    full_name: '',
    phone: '',
    division: 'Dhaka',
    district: 'Dhaka',
    thana: '',
    address_line: '',
    postal_code: '',
  })
  const [showAddressModal, setShowAddressModal] = useState(false)

  // Tracking query state
  const [trackingNumber, setTrackingNumber] = useState('')
  const [trackedOrder, setTrackedOrder] = useState(null)
  const [trackingLoading, setTrackingLoading] = useState(false)

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login?redirect=/account')
    } else if (user) {
      setName(user.name || '')
      setPhone(user.phone || '')
    }
  }, [isAuthenticated, user, navigate])

  // Fetch orders
  useEffect(() => {
    if (isAuthenticated) {
      setLoadingOrders(true)
      api.get('/orders')
        .then((res) => setOrders(res?.orders || []))
        .catch(() => {})
        .finally(() => setLoadingOrders(false))
    }
  }, [isAuthenticated])

  const handleProfileSave = async (e) => {
    e.preventDefault()
    try {
      await updateProfile({ name, phone })
    } catch (err) {
      toast.error(err.message || 'Failed to update profile')
    }
  }

  const handleAddAddress = async (e) => {
    e.preventDefault()
    try {
      await api.post('/users/addresses', newAddress)
      toast.success('Address added successfully')
      setShowAddressModal(false)
      refreshUser()
      setNewAddress({
        full_name: '',
        phone: '',
        division: 'Dhaka',
        district: 'Dhaka',
        thana: '',
        address_line: '',
        postal_code: '',
      })
    } catch (err) {
      toast.error(err.message || 'Failed to add address')
    }
  }

  const handleDeleteAddress = async (id) => {
    try {
      await api.delete(`/users/addresses/${id}`)
      toast.success('Address removed')
      refreshUser()
    } catch (err) {
      toast.error(err.message || 'Failed to remove address')
    }
  }

  const handleTrackSearch = async (e) => {
    e.preventDefault()
    if (!trackingNumber.trim()) return
    setTrackingLoading(true)
    try {
      const data = await api.get(`/orders/track/${encodeURIComponent(trackingNumber.trim())}`)
      if (data?.order) {
        setTrackedOrder(data.order)
      }
    } catch (err) {
      toast.error(err.message || 'No order found with this tracking number')
      setTrackedOrder(null)
    } finally {
      setTrackingLoading(false)
    }
  }

  if (!user) return null

  return (
    <div className="container" style={{ padding: '24px 16px' }}>
      <div className="breadcrumbs">
        <Link to="/">Home</Link> <ChevronRight size={12} style={{ display: 'inline' }} /> <span>My Account</span>
      </div>

      <div className="account-layout">
        {/* Sidebar Nav */}
        <aside className="account-nav">
          <div className="account-user-banner">
            <div style={{
              width: 60,
              height: 60,
              borderRadius: '50%',
              background: '#ef4a23',
              color: '#fff',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              margin: '0 auto 8px',
              fontSize: 24,
              fontWeight: 800,
            }}>
              {user.name.charAt(0).toUpperCase()}
            </div>
            <h3>{user.name}</h3>
            <small style={{ color: '#94a3b8' }}>{user.email}</small>
          </div>

          <div>
            <button
              className={`account-nav-btn ${activeTab === 'orders' ? 'active' : ''}`}
              onClick={() => setActiveTab('orders')}
            >
              <Package size={18} /> My Orders ({orders.length})
            </button>
            <button
              className={`account-nav-btn ${activeTab === 'profile' ? 'active' : ''}`}
              onClick={() => setActiveTab('profile')}
            >
              <User size={18} /> Profile Details
            </button>
            <button
              className={`account-nav-btn ${activeTab === 'tracking' ? 'active' : ''}`}
              onClick={() => setActiveTab('tracking')}
            >
              <ShieldCheck size={18} /> Order Tracking
            </button>
            <button
              className={`account-nav-btn ${activeTab === 'addresses' ? 'active' : ''}`}
              onClick={() => setActiveTab('addresses')}
            >
              <MapPin size={18} /> Saved Addresses
            </button>
            <button
              className={`account-nav-btn ${activeTab === 'wishlist' ? 'active' : ''}`}
              onClick={() => setActiveTab('wishlist')}
            >
              <Heart size={18} /> Wishlist ({wishlist.length})
            </button>
            <button
              className="account-nav-btn"
              onClick={logout}
              style={{ color: '#ef4444' }}
            >
              <LogOut size={18} /> Logout
            </button>
          </div>
        </aside>

        {/* Main Content Area */}
        <main>
          {/* Tab: Orders */}
          {activeTab === 'orders' && (
            <div className="checkout-card">
              <h2 style={{ fontSize: 20, fontWeight: 700, marginBottom: 16 }}>My Order History</h2>
              {loadingOrders ? (
                <p style={{ color: '#64748b' }}>Loading orders...</p>
              ) : orders.length > 0 ? (
                <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
                  {orders.map((o) => (
                    <div
                      key={o.id}
                      style={{
                        border: '1px solid #e2e8f0',
                        borderRadius: 8,
                        padding: 16,
                        background: '#f8fafc',
                      }}
                    >
                      <div style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        borderBottom: '1px solid #e2e8f0',
                        paddingBottom: 12,
                        marginBottom: 12,
                      }}>
                        <div>
                          <b style={{ fontSize: 15, color: '#0f172a' }}>Order #: {o.order_number}</b>
                          <div style={{ fontSize: 12, color: '#64748b', marginTop: 2 }}>
                            Placed on {new Date(o.created_at).toLocaleDateString()} · {o.payment_method}
                          </div>
                        </div>
                        <div style={{ textAlign: 'right' }}>
                          <span className={`badge ${o.status === 'DELIVERED' ? 'badge-stock' : 'badge-discount'}`}>
                            {o.status}
                          </span>
                          <div style={{ fontSize: 16, fontWeight: 800, color: '#ef4a23', marginTop: 4 }}>
                            ৳{o.total?.toLocaleString()}
                          </div>
                        </div>
                      </div>

                      {/* Items */}
                      <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                        {o.items?.map((item) => (
                          <div key={item.id} style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', fontSize: 13 }}>
                            <span>{item.name} × <b>{item.quantity}</b></span>
                            <b>৳{(item.price * item.quantity).toLocaleString()}</b>
                          </div>
                        ))}
                      </div>

                      <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: 12 }}>
                        <Link to={`/order-success/${o.id}`} state={{ order: o }} className="btn btn-outline btn-sm">
                          View Invoice
                        </Link>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div style={{ textAlign: 'center', padding: 40, color: '#64748b' }}>
                  <Package size={40} style={{ margin: '0 auto 12px' }} />
                  <p>You haven't placed any orders yet.</p>
                  <Link to="/" className="btn btn-primary btn-sm" style={{ marginTop: 12 }}>
                    Start Shopping
                  </Link>
                </div>
              )}
            </div>
          )}

          {/* Tab: Profile */}
          {activeTab === 'profile' && (
            <div className="checkout-card">
              <h2 style={{ fontSize: 20, fontWeight: 700, marginBottom: 16 }}>Personal Profile</h2>
              <form onSubmit={handleProfileSave} style={{ maxWidth: 500, display: 'flex', flexDirection: 'column', gap: 16 }}>
                <div className="form-group">
                  <label>Full Name</label>
                  <input
                    type="text"
                    required
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                  />
                </div>

                <div className="form-group">
                  <label>Email Address</label>
                  <input
                    type="email"
                    disabled
                    value={user.email}
                    style={{ background: '#f1f5f9', cursor: 'not-allowed' }}
                  />
                </div>

                <div className="form-group">
                  <label>Phone Number</label>
                  <input
                    type="tel"
                    placeholder="01700000000"
                    value={phone}
                    onChange={(e) => setPhone(e.target.value)}
                  />
                </div>

                <button type="submit" className="btn btn-primary" style={{ width: 'fit-content' }}>
                  Save Profile Changes
                </button>
              </form>
            </div>
          )}

          {/* Tab: Tracking */}
          {activeTab === 'tracking' && (
            <div className="checkout-card">
              <h2 style={{ fontSize: 20, fontWeight: 700, marginBottom: 16 }}>Track Your Order</h2>
              <p style={{ fontSize: 13, color: '#64748b', marginBottom: 16 }}>
                Enter your Order Number (e.g. ST-...) or courier Tracking ID to check status.
              </p>

              <form onSubmit={handleTrackSearch} style={{ display: 'flex', gap: 8, maxWidth: 500, marginBottom: 24 }}>
                <input
                  type="text"
                  required
                  placeholder="Enter Order or Tracking #"
                  style={{ flex: 1, padding: '10px 14px', border: '1px solid #cbd5e1', borderRadius: 4 }}
                  value={trackingNumber}
                  onChange={(e) => setTrackingNumber(e.target.value)}
                />
                <button type="submit" className="btn btn-primary" disabled={trackingLoading}>
                  <Search size={16} /> Track
                </button>
              </form>

              {trackedOrder && (
                <div style={{
                  background: '#f8fafc',
                  border: '1px solid #e2e8f0',
                  borderRadius: 8,
                  padding: 20,
                }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 12 }}>
                    <div>
                      <h3 style={{ fontSize: 16, fontWeight: 700 }}>Order: {trackedOrder.order_number}</h3>
                      <small style={{ color: '#64748b' }}>Status: <b>{trackedOrder.status}</b></small>
                    </div>
                    <div style={{ fontSize: 18, fontWeight: 800, color: '#ef4a23' }}>
                      ৳{trackedOrder.total?.toLocaleString()}
                    </div>
                  </div>
                  <p style={{ fontSize: 13 }}><b>Destination:</b> {trackedOrder.shipping_address}, {trackedOrder.shipping_district}</p>
                  {trackedOrder.tracking_number && (
                    <p style={{ fontSize: 13, marginTop: 4 }}><b>Courier Tracking #:</b> {trackedOrder.tracking_number}</p>
                  )}
                </div>
              )}
            </div>
          )}

          {/* Tab: Addresses */}
          {activeTab === 'addresses' && (
            <div className="checkout-card">
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                <h2 style={{ fontSize: 20, fontWeight: 700 }}>Saved Shipping Addresses</h2>
                <button className="btn btn-secondary btn-sm" onClick={() => setShowAddressModal(!showAddressModal)}>
                  <Plus size={14} /> Add New Address
                </button>
              </div>

              {showAddressModal && (
                <form onSubmit={handleAddAddress} style={{
                  background: '#f8fafc',
                  padding: 20,
                  borderRadius: 8,
                  border: '1px solid #cbd5e1',
                  marginBottom: 20,
                }} className="form-grid">
                  <div className="form-group">
                    <label>Recipient Name</label>
                    <input
                      type="text"
                      required
                      value={newAddress.full_name}
                      onChange={(e) => setNewAddress({ ...newAddress, full_name: e.target.value })}
                    />
                  </div>
                  <div className="form-group">
                    <label>Phone</label>
                    <input
                      type="tel"
                      required
                      value={newAddress.phone}
                      onChange={(e) => setNewAddress({ ...newAddress, phone: e.target.value })}
                    />
                  </div>
                  <div className="form-group">
                    <label>District</label>
                    <input
                      type="text"
                      required
                      value={newAddress.district}
                      onChange={(e) => setNewAddress({ ...newAddress, district: e.target.value })}
                    />
                  </div>
                  <div className="form-group">
                    <label>Thana</label>
                    <input
                      type="text"
                      value={newAddress.thana}
                      onChange={(e) => setNewAddress({ ...newAddress, thana: e.target.value })}
                    />
                  </div>
                  <div className="form-group form-full">
                    <label>Street Address</label>
                    <input
                      type="text"
                      required
                      value={newAddress.address_line}
                      onChange={(e) => setNewAddress({ ...newAddress, address_line: e.target.value })}
                    />
                  </div>
                  <div className="form-full" style={{ display: 'flex', gap: 8 }}>
                    <button type="submit" className="btn btn-primary btn-sm">Save Address</button>
                    <button type="button" className="btn btn-outline btn-sm" onClick={() => setShowAddressModal(false)}>Cancel</button>
                  </div>
                </form>
              )}

              {user.addresses && user.addresses.length > 0 ? (
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
                  {user.addresses.map((addr) => (
                    <div
                      key={addr.id}
                      style={{
                        border: '1px solid #e2e8f0',
                        borderRadius: 6,
                        padding: 16,
                        background: '#fff',
                        position: 'relative',
                      }}
                    >
                      <button
                        onClick={() => handleDeleteAddress(addr.id)}
                        style={{ position: 'absolute', top: 12, right: 12, color: '#ef4444' }}
                        title="Delete address"
                      >
                        <Trash2 size={16} />
                      </button>
                      <b style={{ fontSize: 14, display: 'block', marginBottom: 4 }}>{addr.full_name}</b>
                      <p style={{ fontSize: 13, color: '#64748b' }}>{addr.phone}</p>
                      <p style={{ fontSize: 13, marginTop: 4 }}>{addr.address_line}, {addr.district}, {addr.division}</p>
                    </div>
                  ))}
                </div>
              ) : (
                <p style={{ color: '#64748b' }}>No saved addresses. Add one to speed up checkout.</p>
              )}
            </div>
          )}

          {/* Tab: Wishlist */}
          {activeTab === 'wishlist' && (
            <div className="checkout-card">
              <h2 style={{ fontSize: 20, fontWeight: 700, marginBottom: 16 }}>My Wishlist ({wishlist.length})</h2>
              {wishlist.length > 0 ? (
                <div className="products-grid" style={{ gridTemplateColumns: 'repeat(3, 1fr)' }}>
                  {wishlist.map((product) => (
                    <div key={product.id} className="product-card">
                      <div className="product-image-container">
                        <img
                          src={product.images?.[0]?.url || 'https://placehold.co/200x200'}
                          alt={product.name}
                        />
                      </div>
                      <div className="product-details">
                        <Link to={`/product/${product.slug}`} className="product-title">
                          {product.name}
                        </Link>
                        <div className="product-pricing">
                          <span className="current-price">৳{(product.discount_price || product.price).toLocaleString()}</span>
                        </div>
                        <div style={{ display: 'flex', gap: 8, marginTop: 'auto' }}>
                          <button
                            className="btn btn-primary btn-sm"
                            style={{ flex: 1 }}
                            onClick={() => {
                              addToCart(product, 1)
                              toggleWishlist(product)
                            }}
                          >
                            Move to Cart
                          </button>
                          <button
                            className="btn btn-outline btn-sm"
                            onClick={() => toggleWishlist(product)}
                            style={{ color: '#ef4444' }}
                          >
                            <Trash2 size={14} />
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <p style={{ color: '#64748b' }}>Your wishlist is empty.</p>
              )}
            </div>
          )}
        </main>
      </div>
    </div>
  )
}
