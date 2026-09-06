import React, { useState, useEffect } from 'react'
import { useParams, useLocation, Link } from 'react-router-dom'
import { CheckCircle2, Printer, ArrowRight, ShieldCheck } from 'lucide-react'
import api from '../api/client'

export default function OrderSuccessPage() {
  const { id } = useParams()
  const location = useLocation()
  const [order, setOrder] = useState(location.state?.order || null)
  const [loading, setLoading] = useState(!order)

  useEffect(() => {
    if (!order && id) {
      api.get(`/orders/${id}`)
        .then((res) => setOrder(res?.order))
        .catch(() => {})
        .finally(() => setLoading(false))
    }
  }, [id, order])

  const handlePrint = () => {
    window.print()
  }

  if (loading) {
    return (
      <div className="container" style={{ padding: '80px 0', textAlign: 'center', color: '#64748b' }}>
        <h2>Loading invoice...</h2>
      </div>
    )
  }

  return (
    <div className="container" style={{ padding: '32px 16px' }}>
      <div style={{ textAlign: 'center', marginBottom: 24 }}>
        <CheckCircle2 size={56} color="#10b981" style={{ margin: '0 auto 12px' }} />
        <h1 style={{ fontSize: 26, fontWeight: 800, color: '#0f172a' }}>Thank You For Your Order!</h1>
        <p style={{ color: '#64748b', fontSize: 15 }}>
          Your order has been placed successfully and is currently being processed by MI-Tech.
        </p>
      </div>

      {/* Invoice Card */}
      <div className="invoice-card">
        {/* Invoice Header */}
        <div className="invoice-header">
          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <img src="/logo.png" alt="MI-Tech" style={{ height: 48, width: 'auto', objectFit: 'contain' }} />
            <div>
              <h2 style={{ fontSize: 22, fontWeight: 800, margin: 0 }}>MI-<span style={{ color: '#ef4a23' }}>Tech</span></h2>
              <p style={{ fontSize: 12, color: '#64748b' }}>MI-Tech & Engineering Ltd.</p>
              <p style={{ fontSize: 12, color: '#64748b' }}>Hotline: 16793 | support@mitech.local</p>
            </div>
          </div>
          <div style={{ textAlign: 'right' }}>
            <h3 style={{ fontSize: 18, color: '#ef4a23', fontWeight: 800 }}>INVOICE</h3>
            <p style={{ fontSize: 13 }}>Order #: <b>{order?.order_number || id}</b></p>
            <p style={{ fontSize: 12, color: '#64748b' }}>
              Date: {order?.created_at ? new Date(order.created_at).toLocaleDateString() : new Date().toLocaleDateString()}
            </p>
            <p style={{ fontSize: 12, color: '#64748b' }}>Payment: <b>{order?.payment_method || 'COD'}</b></p>
          </div>
        </div>

        {/* Customer / Shipping Info */}
        <div style={{
          display: 'grid',
          gridTemplateColumns: '1fr 1fr',
          gap: 24,
          background: '#f8fafc',
          padding: 16,
          borderRadius: 6,
          marginBottom: 24,
          fontSize: 13,
        }}>
          <div>
            <b style={{ color: '#334155', display: 'block', marginBottom: 4 }}>Bill & Ship To:</b>
            <p><b>Name:</b> {order?.shipping_name || 'Customer'}</p>
            <p><b>Phone:</b> {order?.shipping_phone || 'N/A'}</p>
            <p><b>Address:</b> {order?.shipping_address || 'N/A'}, {order?.shipping_district}</p>
          </div>
          <div>
            <b style={{ color: '#334155', display: 'block', marginBottom: 4 }}>Order Details:</b>
            <p><b>Status:</b> <span style={{ color: '#ef4a23', fontWeight: 700 }}>{order?.status || 'PENDING'}</span></p>
            <p><b>Payment Status:</b> {order?.payment_status || 'UNPAID'}</p>
            {order?.tracking_number && <p><b>Tracking #:</b> {order.tracking_number}</p>}
          </div>
        </div>

        {/* Itemized Table */}
        <table className="specs-table" style={{ marginBottom: 24 }}>
          <thead>
            <tr style={{ background: '#0f172a', color: '#fff' }}>
              <th style={{ color: '#fff' }}>Item Description</th>
              <th style={{ color: '#fff', textAlign: 'center' }}>Price</th>
              <th style={{ color: '#fff', textAlign: 'center' }}>Qty</th>
              <th style={{ color: '#fff', textAlign: 'right' }}>Total</th>
            </tr>
          </thead>
          <tbody>
            {order?.items?.map((item) => (
              <tr key={item.id}>
                <td>{item.name}</td>
                <td style={{ textAlign: 'center' }}>৳{item.price?.toLocaleString()}</td>
                <td style={{ textAlign: 'center' }}>{item.quantity}</td>
                <td style={{ textAlign: 'right', fontWeight: 700 }}>
                  ৳{(item.price * item.quantity).toLocaleString()}
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {/* Total Summary */}
        <div style={{ width: 280, marginLeft: 'auto', fontSize: 14 }}>
          <div className="summary-row">
            <span>Subtotal:</span>
            <b>৳{order?.subtotal?.toLocaleString()}</b>
          </div>
          {order?.discount > 0 && (
            <div className="summary-row" style={{ color: '#10b981' }}>
              <span>Discount ({order.coupon_code}):</span>
              <b>-৳{order?.discount?.toLocaleString()}</b>
            </div>
          )}
          <div className="summary-row">
            <span>Shipping Fee:</span>
            <b>৳{order?.shipping_fee?.toLocaleString()}</b>
          </div>
          <div className="summary-row summary-total">
            <span>Total Payable:</span>
            <span>৳{order?.total?.toLocaleString()}</span>
          </div>
        </div>

        {/* Print / Track Actions */}
        <div style={{
          display: 'flex',
          justifyContent: 'space-between',
          marginTop: 32,
          paddingTop: 20,
          borderTop: '1px solid #e2e8f0',
        }}>
          <button className="btn btn-outline" onClick={handlePrint}>
            <Printer size={16} /> Print Invoice
          </button>
          <div style={{ display: 'flex', gap: 12 }}>
            <Link to={`/order-tracking?order=${order?.order_number || ''}`} className="btn btn-secondary">
              Track Order <ArrowRight size={16} />
            </Link>
            <Link to="/" className="btn btn-primary">
              Continue Shopping
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}
