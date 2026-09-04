import React, { useState, useEffect } from 'react'
import { useSearchParams, Link } from 'react-router-dom'
import { Search, ShieldCheck, Package, Clock, CheckCircle2, ChevronRight } from 'lucide-react'
import api from '../api/client'
import toast from 'react-hot-toast'

export default function OrderTrackingPage() {
  const [searchParams] = useSearchParams()
  const initialOrder = searchParams.get('order') || ''
  const [orderQuery, setOrderQuery] = useState(initialOrder)
  const [orderData, setOrderData] = useState(null)
  const [loading, setLoading] = useState(false)

  const performTrack = async (num) => {
    if (!num.trim()) return
    setLoading(true)
    try {
      const data = await api.get(`/orders/track/${encodeURIComponent(num.trim())}`)
      if (data?.order) {
        setOrderData(data.order)
      }
    } catch (err) {
      toast.error(err.message || 'No order found with that tracking/order number')
      setOrderData(null)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (initialOrder) {
      performTrack(initialOrder)
    }
  }, [initialOrder])

  const handleTrackSubmit = (e) => {
    e.preventDefault()
    performTrack(orderQuery)
  }

  const steps = [
    { label: 'Order Placed', status: 'PENDING' },
    { label: 'Confirmed', status: 'CONFIRMED' },
    { label: 'Processing', status: 'PROCESSING' },
    { label: 'Shipped', status: 'SHIPPED' },
    { label: 'Delivered', status: 'DELIVERED' },
  ]

  const getStepIndex = (status) => {
    switch (status) {
      case 'CONFIRMED': return 1
      case 'PROCESSING': return 2
      case 'SHIPPED': return 3
      case 'DELIVERED': return 4
      default: return 0
    }
  }

  return (
    <div className="container" style={{ padding: '32px 16px', maxWidth: 800 }}>
      <div className="breadcrumbs">
        <Link to="/">Home</Link> <ChevronRight size={12} style={{ display: 'inline' }} /> <span>Track Your Order</span>
      </div>

      <div style={{ textAlign: 'center', marginBottom: 28 }}>
        <ShieldCheck size={48} color="#ef4a23" style={{ margin: '0 auto 12px' }} />
        <h1 style={{ fontSize: 26, fontWeight: 800, color: '#0f172a' }}>Real-Time Order Tracking</h1>
        <p style={{ color: '#64748b', fontSize: 14 }}>
          Enter your StarTech Order Number (e.g. ST-...) or courier tracking ID
        </p>
      </div>

      {/* Tracking Search Input */}
      <div style={{
        background: '#fff',
        borderRadius: 8,
        border: '1px solid #e2e8f0',
        padding: 24,
        boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
        marginBottom: 28,
      }}>
        <form onSubmit={handleTrackSubmit} style={{ display: 'flex', gap: 10 }}>
          <input
            type="text"
            required
            placeholder="e.g. ST-1772798122-1234"
            value={orderQuery}
            onChange={(e) => setOrderQuery(e.target.value)}
            style={{
              flex: 1,
              padding: '12px 16px',
              border: '1px solid #cbd5e1',
              borderRadius: 4,
              fontSize: 15,
            }}
          />
          <button type="submit" className="btn btn-primary btn-lg" disabled={loading}>
            {loading ? 'Searching...' : 'Track Status'} <Search size={16} />
          </button>
        </form>
      </div>

      {/* Tracking Result */}
      {orderData && (
        <div style={{
          background: '#fff',
          borderRadius: 8,
          border: '1px solid #e2e8f0',
          padding: 32,
          boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
        }}>
          <div style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            borderBottom: '1px solid #e2e8f0',
            paddingBottom: 16,
            marginBottom: 24,
          }}>
            <div>
              <h2 style={{ fontSize: 18, fontWeight: 700, color: '#0f172a' }}>
                Order #{orderData.order_number}
              </h2>
              <small style={{ color: '#64748b' }}>
                Placed on {new Date(orderData.created_at).toLocaleDateString()} · {orderData.payment_method}
              </small>
            </div>
            <div style={{ textAlign: 'right' }}>
              <span className={`badge ${orderData.status === 'DELIVERED' ? 'badge-stock' : 'badge-discount'}`} style={{ fontSize: 13, padding: '4px 10px' }}>
                {orderData.status}
              </span>
              <div style={{ fontSize: 18, fontWeight: 800, color: '#ef4a23', marginTop: 4 }}>
                ৳{orderData.total?.toLocaleString()}
              </div>
            </div>
          </div>

          {/* Progress Tracker */}
          <div style={{ display: 'flex', justifyContent: 'space-between', position: 'relative', margin: '36px 0 24px' }}>
            {steps.map((step, idx) => {
              const currentIndex = getStepIndex(orderData.status)
              const isDone = idx <= currentIndex
              return (
                <div key={step.label} style={{ textAlign: 'center', flex: 1, position: 'relative', zIndex: 2 }}>
                  <div style={{
                    width: 36,
                    height: 36,
                    borderRadius: '50%',
                    background: isDone ? '#10b981' : '#e2e8f0',
                    color: isDone ? '#fff' : '#64748b',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    margin: '0 auto 8px',
                    fontWeight: 700,
                    fontSize: 14,
                  }}>
                    {isDone ? <CheckCircle2 size={18} /> : idx + 1}
                  </div>
                  <small style={{ fontWeight: isDone ? 700 : 500, color: isDone ? '#0f172a' : '#94a3b8' }}>
                    {step.label}
                  </small>
                </div>
              )
            })}
          </div>

          {/* Details */}
          <div style={{
            background: '#f8fafc',
            borderRadius: 6,
            padding: 16,
            fontSize: 13,
            lineHeight: 1.6,
          }}>
            <p><b>Recipient:</b> {orderData.shipping_name} ({orderData.shipping_phone})</p>
            <p><b>Delivery Address:</b> {orderData.shipping_address}, {orderData.shipping_district}, {orderData.shipping_division}</p>
            {orderData.tracking_number && (
              <p><b>Courier Tracking Reference:</b> <span style={{ color: '#3749bb', fontWeight: 700 }}>{orderData.tracking_number}</span></p>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
