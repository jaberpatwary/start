import React, { useState, useEffect } from 'react'
import { Users, Star, Lock, Unlock, Check, X, ShieldAlert, Search } from 'lucide-react'
import api from '../../api/client'
import toast from 'react-hot-toast'

export default function AdminUsersReviews() {
  const [activeTab, setActiveTab] = useState('users')
  const [users, setUsers] = useState([])
  const [reviews, setReviews] = useState([])
  const [searchUser, setSearchUser] = useState('')
  const [reviewFilter, setReviewFilter] = useState('')
  const [loading, setLoading] = useState(false)

  const fetchUsers = async () => {
    setLoading(true)
    try {
      const res = await api.get('/admin/users')
      setUsers(res?.results || [])
    } catch (err) {
      toast.error('Failed to load users')
    } finally {
      setLoading(false)
    }
  }

  const fetchReviews = async () => {
    setLoading(true)
    try {
      let endpoint = '/admin/reviews'
      if (reviewFilter) endpoint += `?status=${reviewFilter}`
      const res = await api.get(endpoint)
      setReviews(res?.results || [])
    } catch (err) {
      toast.error('Failed to load reviews')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (activeTab === 'users') fetchUsers()
    else fetchReviews()
  }, [activeTab, reviewFilter])

  const handleToggleBlock = async (userId, isBlocked) => {
    try {
      const res = await api.patch(`/admin/users/${userId}/toggle-block`)
      toast.success(res.message || 'User status updated')
      fetchUsers()
    } catch (err) {
      toast.error(err.message || 'Action failed')
    }
  }

  const handleUpdateReviewStatus = async (reviewId, newStatus) => {
    try {
      await api.patch(`/admin/reviews/${reviewId}/status`, { status: newStatus })
      toast.success(`Review ${newStatus.toLowerCase()}`)
      fetchReviews()
    } catch (err) {
      toast.error(err.message || 'Failed to update review status')
    }
  }

  const filteredUsers = users.filter(
    (u) => u.name?.toLowerCase().includes(searchUser.toLowerCase()) || u.email?.toLowerCase().includes(searchUser.toLowerCase())
  )

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
        <div>
          <h1 style={{ fontSize: 24, fontWeight: 800 }}>Users & Product Reviews</h1>
          <small style={{ color: '#64748b' }}>Manage customer accounts, permissions, and moderate user product reviews</small>
        </div>
      </div>

      {/* Tabs */}
      <div style={{ display: 'flex', gap: 12, borderBottom: '2px solid #e2e8f0', marginBottom: 24 }}>
        <button
          onClick={() => setActiveTab('users')}
          style={{
            padding: '10px 18px',
            fontSize: 14,
            fontWeight: 700,
            color: activeTab === 'users' ? '#ef4a23' : '#64748b',
            borderBottom: activeTab === 'users' ? '3px solid #ef4a23' : 'none',
            background: 'none',
            display: 'flex',
            alignItems: 'center',
            gap: 8,
          }}
        >
          <Users size={16} /> User Accounts ({users.length})
        </button>
        <button
          onClick={() => setActiveTab('reviews')}
          style={{
            padding: '10px 18px',
            fontSize: 14,
            fontWeight: 700,
            color: activeTab === 'reviews' ? '#ef4a23' : '#64748b',
            borderBottom: activeTab === 'reviews' ? '3px solid #ef4a23' : 'none',
            background: 'none',
            display: 'flex',
            alignItems: 'center',
            gap: 8,
          }}
        >
          <Star size={16} /> Product Reviews Moderation
        </button>
      </div>

      {/* Users Tab */}
      {activeTab === 'users' && (
        <div>
          <div className="admin-card" style={{ padding: 16, marginBottom: 20 }}>
            <input
              type="text"
              placeholder="Search user by name or email..."
              value={searchUser}
              onChange={(e) => setSearchUser(e.target.value)}
              style={{ width: '100%', padding: '10px 14px', borderRadius: 4, border: '1px solid #cbd5e1', fontSize: 14 }}
            />
          </div>

          <div className="admin-card">
            <table className="admin-table">
              <thead>
                <tr>
                  <th>User</th>
                  <th>Email</th>
                  <th>Role</th>
                  <th>Account Status</th>
                  <th style={{ textAlign: 'right' }}>Actions</th>
                </tr>
              </thead>
              <tbody>
                {loading ? (
                  <tr><td colSpan={5} style={{ textAlign: 'center', padding: 32 }}>Loading user accounts...</td></tr>
                ) : filteredUsers.length > 0 ? (
                  filteredUsers.map((u) => (
                    <tr key={u.id}>
                      <td>
                        <b style={{ fontSize: 14 }}>{u.name}</b>
                      </td>
                      <td style={{ color: '#3749bb' }}>{u.email}</td>
                      <td>
                        <span className={`badge ${u.role === 'ADMIN' ? 'badge-discount' : 'badge-stock'}`}>
                          {u.role}
                        </span>
                      </td>
                      <td>
                        <span style={{
                          padding: '4px 8px',
                          borderRadius: 4,
                          fontSize: 12,
                          fontWeight: 700,
                          background: u.is_blocked ? '#fee2e2' : '#dcfce7',
                          color: u.is_blocked ? '#ef4444' : '#15803d',
                        }}>
                          {u.is_blocked ? 'BLOCKED' : 'ACTIVE'}
                        </span>
                      </td>
                      <td style={{ textAlign: 'right' }}>
                        {u.role !== 'ADMIN' && (
                          <button
                            className={`btn btn-sm ${u.is_blocked ? 'btn-primary' : 'btn-outline'}`}
                            onClick={() => handleToggleBlock(u.id, u.is_blocked)}
                            style={{ color: u.is_blocked ? '#fff' : '#ef4444' }}
                          >
                            {u.is_blocked ? <Unlock size={14} /> : <Lock size={14} />}
                            {u.is_blocked ? 'Unblock' : 'Block User'}
                          </button>
                        )}
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr><td colSpan={5} style={{ textAlign: 'center', padding: 32 }}>No users found</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Reviews Tab */}
      {activeTab === 'reviews' && (
        <div>
          <div className="admin-card" style={{ padding: 12, marginBottom: 20 }}>
            <div style={{ display: 'flex', gap: 8 }}>
              {['', 'PENDING', 'APPROVED', 'REJECTED'].map((st) => (
                <button
                  key={st}
                  onClick={() => setReviewFilter(st)}
                  className={`btn btn-sm ${reviewFilter === st ? 'btn-primary' : 'btn-outline'}`}
                >
                  {st || 'All Reviews'}
                </button>
              ))}
            </div>
          </div>

          <div className="admin-card">
            <table className="admin-table">
              <thead>
                <tr>
                  <th>Product</th>
                  <th>Customer</th>
                  <th>Rating</th>
                  <th>Review Comment</th>
                  <th>Status</th>
                  <th style={{ textAlign: 'right' }}>Actions</th>
                </tr>
              </thead>
              <tbody>
                {loading ? (
                  <tr><td colSpan={6} style={{ textAlign: 'center', padding: 32 }}>Loading reviews...</td></tr>
                ) : reviews.length > 0 ? (
                  reviews.map((r) => (
                    <tr key={r.id}>
                      <td><b>{r.product?.name || 'Product'}</b></td>
                      <td>{r.user?.name || 'Customer'}</td>
                      <td>
                        <div style={{ display: 'flex', color: '#f59e0b' }}>
                          {[...Array(5)].map((_, i) => (
                            <Star key={i} size={14} fill={i < r.rating ? '#f59e0b' : 'none'} />
                          ))}
                        </div>
                      </td>
                      <td style={{ maxWidth: 280, fontSize: 13 }}>{r.comment}</td>
                      <td>
                        <span className={`badge ${r.status === 'APPROVED' ? 'badge-stock' : r.status === 'PENDING' ? 'badge-discount' : 'badge-out'}`}>
                          {r.status}
                        </span>
                      </td>
                      <td style={{ textAlign: 'right' }}>
                        <div style={{ display: 'flex', gap: 6, justifyContent: 'flex-end' }}>
                          <button
                            className="btn btn-outline btn-sm"
                            onClick={() => handleUpdateReviewStatus(r.id, 'APPROVED')}
                            style={{ color: '#10b981' }}
                          >
                            <Check size={14} /> Approve
                          </button>
                          <button
                            className="btn btn-outline btn-sm"
                            onClick={() => handleUpdateReviewStatus(r.id, 'REJECTED')}
                            style={{ color: '#ef4444' }}
                          >
                            <X size={14} /> Reject
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr><td colSpan={6} style={{ textAlign: 'center', padding: 32 }}>No customer reviews found</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  )
}
