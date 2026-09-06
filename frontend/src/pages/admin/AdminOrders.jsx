import React, { useState, useEffect } from 'react'
import { Eye, Printer, Filter, CheckCircle, Clock, Truck, XCircle, DollarSign, Package } from 'lucide-react'
import api from '../../api/client'
import toast from 'react-hot-toast'

const STATUS_BADGES = {
  PENDING: { label: 'PENDING', bg: '#fef3c7', text: '#d97706' },
  CONFIRMED: { label: 'CONFIRMED', bg: '#e0f2fe', text: '#0284c7' },
  PROCESSING: { label: 'PROCESSING', bg: '#fae8ff', text: '#c026d3' },
  SHIPPED: { label: 'SHIPPED', bg: '#e0e7ff', text: '#4338ca' },
  DELIVERED: { label: 'DELIVERED', bg: '#dcfce7', text: '#15803d' },
  CANCELLED: { label: 'CANCELLED', bg: '#fee2e2', text: '#b91c1c' },
}

export default function AdminOrders() {
  const [orders, setOrders] = useState([])
  const [statusFilter, setStatusFilter] = useState('')
  const [loading, setLoading] = useState(true)

  // Selected Order Modal
  const [selectedOrder, setSelectedOrder] = useState(null)
  const [trackingNo, setTrackingNo] = useState('')

  const fetchOrders = async () => {
    setLoading(true)
    try {
      let endpoint = '/admin/orders?limit=50'
      if (statusFilter) endpoint += `&status=${statusFilter}`
      const res = await api.get(endpoint)
      setOrders(res?.orders || [])
    } catch (err) {
      toast.error('Failed to fetch orders')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchOrders()
  }, [statusFilter])

  const handleUpdateStatus = async (orderId, newStatus) => {
    try {
      await api.patch(`/admin/orders/${orderId}/status`, {
        status: newStatus,
        tracking_number: trackingNo || undefined,
      })
      toast.success(`Order status updated to ${newStatus}`)
      if (selectedOrder && selectedOrder.id === orderId) {
        setSelectedOrder({ ...selectedOrder, status: newStatus })
      }
      fetchOrders()
    } catch (err) {
      toast.error(err.message || 'Status update failed')
    }
  }

  const handleUpdatePaymentStatus = async (orderId, newPaymentStatus) => {
    try {
      await api.patch(`/admin/orders/${orderId}/payment`, {
        payment_status: newPaymentStatus,
      })
      toast.success(`Payment status marked as ${newPaymentStatus}`)
      if (selectedOrder && selectedOrder.id === orderId) {
        setSelectedOrder({ ...selectedOrder, payment_status: newPaymentStatus })
      }
      fetchOrders()
    } catch (err) {
      toast.error(err.message || 'Payment update failed')
    }
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
        <div>
          <h1 style={{ fontSize: 24, fontWeight: 800 }}>Orders Management</h1>
          <small style={{ color: '#64748b' }}>Track, fulfill, print invoices, and update order statuses</small>
        </div>
      </div>

      {/* Filter Tabs */}
      <div className="admin-card" style={{ padding: 12, marginBottom: 20 }}>
        <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
          {['', 'PENDING', 'CONFIRMED', 'PROCESSING', 'SHIPPED', 'DELIVERED', 'CANCELLED'].map((st) => (
            <button
              key={st}
              onClick={() => setStatusFilter(st)}
              className={`btn btn-sm ${statusFilter === st ? 'btn-primary' : 'btn-outline'}`}
              style={{ textTransform: 'capitalize' }}
            >
              {st || 'All Orders'}
            </button>
          ))}
        </div>
      </div>

      {/* Orders Table */}
      <div className="admin-card">
        <table className="admin-table">
          <thead>
            <tr>
              <th>Order #</th>
              <th>Customer</th>
              <th>Items</th>
              <th>Total Amount</th>
              <th>Payment</th>
              <th>Status</th>
              <th>Date</th>
              <th style={{ textAlign: 'right' }}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr><td colSpan={8} style={{ textAlign: 'center', padding: 32 }}>Loading orders...</td></tr>
            ) : orders.length > 0 ? (
              orders.map((o) => {
                const badge = STATUS_BADGES[o.status] || STATUS_BADGES.PENDING
                return (
                  <tr key={o.id}>
                    <td>
                      <b style={{ color: '#3749bb' }}>{o.order_number}</b>
                    </td>
                    <td>
                      <b style={{ display: 'block', fontSize: 13 }}>{o.shipping_name || o.user?.name}</b>
                      <small style={{ color: '#64748b' }}>{o.shipping_phone}</small>
                    </td>
                    <td>{o.items?.length || 1} item(s)</td>
                    <td>
                      <b style={{ color: '#ef4a23' }}>৳{o.total?.toLocaleString()}</b>
                      <small style={{ display: 'block', color: '#94a3b8' }}>{o.payment_method}</small>
                    </td>
                    <td>
                      <select
                        value={o.payment_status}
                        onChange={(e) => handleUpdatePaymentStatus(o.id, e.target.value)}
                        style={{
                          padding: '3px 6px',
                          borderRadius: 4,
                          fontSize: 12,
                          fontWeight: 700,
                          border: '1px solid #cbd5e1',
                          background: o.payment_status === 'PAID' ? '#dcfce7' : '#fef3c7',
                          color: o.payment_status === 'PAID' ? '#15803d' : '#d97706',
                        }}
                      >
                        <option value="UNPAID">UNPAID</option>
                        <option value="PAID">PAID</option>
                        <option value="REFUNDED">REFUNDED</option>
                      </select>
                    </td>
                    <td>
                      <select
                        value={o.status}
                        onChange={(e) => handleUpdateStatus(o.id, e.target.value)}
                        style={{
                          padding: '4px 8px',
                          borderRadius: 4,
                          fontSize: 12,
                          fontWeight: 700,
                          border: `1px solid ${badge.text}`,
                          background: badge.bg,
                          color: badge.text,
                        }}
                      >
                        <option value="PENDING">PENDING</option>
                        <option value="CONFIRMED">CONFIRMED</option>
                        <option value="PROCESSING">PROCESSING</option>
                        <option value="SHIPPED">SHIPPED</option>
                        <option value="DELIVERED">DELIVERED</option>
                        <option value="CANCELLED">CANCELLED</option>
                      </select>
                    </td>
                    <td style={{ fontSize: 12, color: '#64748b' }}>
                      {new Date(o.created_at).toLocaleDateString('en-GB')}
                    </td>
                    <td style={{ textAlign: 'right' }}>
                      <button className="btn btn-outline btn-sm" onClick={() => setSelectedOrder(o)}>
                        <Eye size={14} /> View Invoice
                      </button>
                    </td>
                  </tr>
                )
              })
            ) : (
              <tr><td colSpan={8} style={{ textAlign: 'center', padding: 32 }}>No orders found for this status</td></tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Invoice & Order Detail Modal */}
      {selectedOrder && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          background: 'rgba(0,0,0,0.6)',
          zIndex: 1000,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          padding: 20,
        }}>
          <div style={{
            background: '#fff',
            borderRadius: 8,
            maxWidth: 700,
            width: '100%',
            maxHeight: '90vh',
            overflowY: 'auto',
            padding: 28,
            boxShadow: '0 20px 25px -5px rgba(0,0,0,0.3)',
          }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', borderBottom: '2px solid #ef4a23', paddingBottom: 12, marginBottom: 20 }}>
              <div>
                <h2 style={{ fontSize: 22, fontWeight: 800, color: '#0f172a' }}>INVOICE: #{selectedOrder.order_number}</h2>
                <small style={{ color: '#64748b' }}>Date: {new Date(selectedOrder.created_at).toLocaleString()}</small>
              </div>
              <div style={{ textAlign: 'right' }}>
                <button className="btn btn-primary btn-sm" onClick={() => window.print()} style={{ marginRight: 8 }}>
                  <Printer size={14} /> Print Invoice
                </button>
                <button className="btn btn-outline btn-sm" onClick={() => setSelectedOrder(null)}>Close</button>
              </div>
            </div>

            {/* Customer & Shipping Details */}
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16, marginBottom: 20, background: '#f8fafc', padding: 16, borderRadius: 6 }}>
              <div>
                <h4 style={{ fontSize: 13, textTransform: 'uppercase', color: '#64748b', marginBottom: 6 }}>Shipping Details</h4>
                <b style={{ display: 'block', fontSize: 15 }}>{selectedOrder.shipping_name}</b>
                <p style={{ fontSize: 13, margin: '4px 0', color: '#334155' }}>
                  {selectedOrder.shipping_address}, {selectedOrder.shipping_thana}, {selectedOrder.shipping_district}
                </p>
                <p style={{ fontSize: 13, margin: 0, color: '#334155' }}>Phone: <strong>{selectedOrder.shipping_phone}</strong></p>
              </div>
              <div>
                <h4 style={{ fontSize: 13, textTransform: 'uppercase', color: '#64748b', marginBottom: 6 }}>Order Meta</h4>
                <p style={{ fontSize: 13, margin: '2px 0' }}>Payment Method: <strong>{selectedOrder.payment_method}</strong></p>
                <p style={{ fontSize: 13, margin: '2px 0' }}>Payment Status: <strong>{selectedOrder.payment_status}</strong></p>
                <p style={{ fontSize: 13, margin: '2px 0' }}>Order Status: <strong>{selectedOrder.status}</strong></p>
              </div>
            </div>

            {/* Items Table */}
            <h4 style={{ fontSize: 15, fontWeight: 700, marginBottom: 10 }}>Order Items</h4>
            <table className="admin-table" style={{ marginBottom: 20 }}>
              <thead>
                <tr>
                  <th>Product</th>
                  <th>Price</th>
                  <th style={{ textAlign: 'center' }}>Qty</th>
                  <th style={{ textAlign: 'right' }}>Subtotal</th>
                </tr>
              </thead>
              <tbody>
                {selectedOrder.items?.map((item) => (
                  <tr key={item.id}>
                    <td>
                      <b>{item.product_name}</b>
                    </td>
                    <td>৳{item.price?.toLocaleString()}</td>
                    <td style={{ textAlign: 'center' }}>{item.quantity}</td>
                    <td style={{ textAlign: 'right', fontWeight: 700 }}>৳{(item.price * item.quantity).toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>

            {/* Summary */}
            <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
              <div style={{ width: 260, borderTop: '2px solid #e2e8f0', paddingTop: 10 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 6, fontSize: 14 }}>
                  <span>Subtotal:</span>
                  <span>৳{selectedOrder.subtotal?.toLocaleString()}</span>
                </div>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 6, fontSize: 14 }}>
                  <span>Delivery Charge:</span>
                  <span>৳{selectedOrder.shipping_cost?.toLocaleString() || 60}</span>
                </div>
                {selectedOrder.discount > 0 && (
                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 6, fontSize: 14, color: '#16a34a' }}>
                    <span>Discount:</span>
                    <span>-৳{selectedOrder.discount?.toLocaleString()}</span>
                  </div>
                )}
                <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: 18, fontWeight: 800, color: '#ef4a23', borderTop: '1px solid #cbd5e1', paddingTop: 8 }}>
                  <span>Grand Total:</span>
                  <span>৳{selectedOrder.total?.toLocaleString()}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
