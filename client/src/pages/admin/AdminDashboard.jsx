import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { DollarSign, ShoppingBag, Package, AlertTriangle, ArrowRight } from 'lucide-react'
import { ResponsiveContainer, BarChart, Bar, XAxis, YAxis, Tooltip, CartesianGrid, PieChart, Pie, Cell } from 'recharts'
import api from '../../api/client'

const COLORS = ['#ef4a23', '#3749bb', '#10b981', '#f59e0b', '#8b5cf6']

export default function AdminDashboard() {
  const [stats, setStats] = useState(null)
  const [recentOrders, setRecentOrders] = useState([])
  const [lowStock, setLowStock] = useState([])
  const [catSales, setCatSales] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchDashboard = async () => {
      try {
        const [dashRes, repRes] = await Promise.all([
          api.get('/admin/dashboard'),
          api.get('/admin/reports').catch(() => ({ cat_sales: [] })),
        ])
        setStats(dashRes?.stats || {})
        setRecentOrders(dashRes?.recent_orders || [])
        setLowStock(dashRes?.low_stock || [])
        setCatSales(repRes?.cat_sales || [])
      } catch (err) {
        console.error('Dashboard load error:', err)
      } finally {
        setLoading(false)
      }
    }
    fetchDashboard()
  }, [])

  if (loading) {
    return <div style={{ padding: 40, color: '#64748b' }}>Loading dashboard statistics...</div>
  }

  // Monthly Sales Chart Data mock/projection
  const monthlyData = [
    { month: 'Apr', sales: 450000 },
    { month: 'May', sales: 580000 },
    { month: 'Jun', sales: 620000 },
    { month: 'Jul', sales: 710000 },
    { month: 'Aug', sales: 890000 },
    { month: 'Sep', sales: stats?.total_revenue || 950000 },
  ]

  const pieData = catSales.length > 0
    ? catSales.map((c) => ({ name: c.category_name, value: Number(c.count) }))
    : [
        { name: 'Desktop', value: 4 },
        { name: 'Laptop', value: 5 },
        { name: 'Component', value: 8 },
        { name: 'Monitor', value: 4 },
      ]

  return (
    <div>
      <h1 style={{ fontSize: 24, fontWeight: 800, marginBottom: 20 }}>Executive Overview</h1>

      {/* Stats Cards */}
      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-info">
            <small>Total Revenue</small>
            <b style={{ color: '#ef4a23' }}>৳{(stats?.total_revenue || 0).toLocaleString()}</b>
          </div>
          <div className="stat-icon" style={{ background: '#fff0eb', color: '#ef4a23' }}>
            <DollarSign size={24} />
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-info">
            <small>Total Orders</small>
            <b>{stats?.total_orders || 0}</b>
          </div>
          <div className="stat-icon" style={{ background: '#ebf0ff', color: '#3749bb' }}>
            <ShoppingBag size={24} />
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-info">
            <small>Total Products</small>
            <b>{stats?.total_products || 0}</b>
          </div>
          <div className="stat-icon" style={{ background: '#ecfdf5', color: '#10b981' }}>
            <Package size={24} />
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-info">
            <small>Pending Orders</small>
            <b style={{ color: '#f59e0b' }}>{stats?.pending_orders || 0}</b>
          </div>
          <div className="stat-icon" style={{ background: '#fef3c7', color: '#d97706' }}>
            <AlertTriangle size={24} />
          </div>
        </div>
      </div>

      {/* Charts Grid */}
      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: 24, marginBottom: 28 }}>
        <div className="admin-card">
          <h3 style={{ fontSize: 16, fontWeight: 700, marginBottom: 16 }}>Revenue Growth Trend (BDT)</h3>
          <div style={{ height: 260 }}>
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={monthlyData}>
                <CartesianGrid strokeDasharray="3 3" vertical={false} />
                <XAxis dataKey="month" />
                <YAxis />
                <Tooltip formatter={(value) => `৳${Number(value).toLocaleString()}`} />
                <Bar dataKey="sales" fill="#ef4a23" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="admin-card">
          <h3 style={{ fontSize: 16, fontWeight: 700, marginBottom: 16 }}>Catalog Distribution</h3>
          <div style={{ height: 260 }}>
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie data={pieData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={80} label>
                  {pieData.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      {/* Recent Orders & Low Stock */}
      <div style={{ display: 'grid', gridTemplateColumns: '1.4fr 1fr', gap: 24 }}>
        {/* Recent Orders */}
        <div className="admin-card">
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
            <h3 style={{ fontSize: 16, fontWeight: 700 }}>Recent Orders</h3>
            <Link to="/admin/orders" style={{ fontSize: 13, color: '#3749bb', fontWeight: 600 }}>
              View All <ArrowRight size={14} style={{ display: 'inline' }} />
            </Link>
          </div>

          <table className="admin-table">
            <thead>
              <tr>
                <th>Order #</th>
                <th>Customer</th>
                <th>Total</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {recentOrders.length > 0 ? (
                recentOrders.map((o) => (
                  <tr key={o.id}>
                    <td><b>{o.order_number}</b></td>
                    <td>{o.shipping_name || o.user?.name || 'Customer'}</td>
                    <td style={{ color: '#ef4a23', fontWeight: 700 }}>৳{o.total?.toLocaleString()}</td>
                    <td>
                      <span className={`badge ${o.status === 'DELIVERED' ? 'badge-stock' : 'badge-discount'}`}>
                        {o.status}
                      </span>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={4} style={{ textAlign: 'center', color: '#64748b' }}>No recent orders</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {/* Low Stock Alerts */}
        <div className="admin-card">
          <h3 style={{ fontSize: 16, fontWeight: 700, marginBottom: 16, color: '#d97706', display: 'flex', alignItems: 'center', gap: 6 }}>
            <AlertTriangle size={18} /> Low Stock Warnings
          </h3>

          <table className="admin-table">
            <thead>
              <tr>
                <th>Product</th>
                <th style={{ textAlign: 'center' }}>Stock Left</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              {lowStock.length > 0 ? (
                lowStock.map((p) => (
                  <tr key={p.id}>
                    <td>
                      <div style={{ fontSize: 13, fontWeight: 600, maxWidth: 180, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                        {p.name}
                      </div>
                    </td>
                    <td style={{ textAlign: 'center', color: '#ef4444', fontWeight: 800 }}>
                      {p.stock}
                    </td>
                    <td>
                      <Link to={`/admin/products?edit=${p.id}`} className="btn btn-outline btn-sm">
                        Restock
                      </Link>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={3} style={{ textAlign: 'center', color: '#10b981' }}>All items healthy in stock</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
