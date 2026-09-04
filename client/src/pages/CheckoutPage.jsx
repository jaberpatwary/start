import React, { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useShop } from '../context/ShopContext'
import { useAuth } from '../context/AuthContext'
import api from '../api/client'
import toast from 'react-hot-toast'
import { ShieldCheck, ArrowRight, ArrowLeft, Check, ChevronRight } from 'lucide-react'

export default function CheckoutPage() {
  const { cart, subtotal, discount, deliveryFee, grandTotal, appliedCoupon, clearCart } = useShop()
  const { user, isAuthenticated } = useAuth()
  const navigate = useNavigate()

  const [currentStep, setCurrentStep] = useState(1) // 1: Shipping, 2: Payment, 3: Review
  const [submitting, setSubmitting] = useState(false)

  // Shipping Form State
  const [shippingData, setShippingData] = useState({
    shipping_name: user?.name || '',
    shipping_phone: user?.phone || '',
    email: user?.email || '',
    shipping_division: 'Dhaka',
    shipping_district: 'Dhaka',
    shipping_thana: '',
    shipping_address: '',
    shipping_postal: '',
    note: '',
  })

  // Payment Method State
  const [paymentMethod, setPaymentMethod] = useState('COD')

  const divisions = ['Dhaka', 'Chittagong', 'Rajshahi', 'Khulna', 'Sylhet', 'Barisal', 'Rangpur', 'Mymensingh']

  const handleShippingSubmit = (e) => {
    e.preventDefault()
    if (!shippingData.shipping_name || !shippingData.shipping_phone || !shippingData.shipping_address) {
      toast.error('Please fill in all required shipping fields')
      return
    }
    setCurrentStep(2)
  }

  const handlePaymentSubmit = (e) => {
    e.preventDefault()
    setCurrentStep(3)
  }

  const handlePlaceOrder = async () => {
    if (cart.length === 0) {
      toast.error('Your cart is empty')
      return
    }

    if (!isAuthenticated) {
      toast.error('Please login to finalize your order')
      navigate('/login?redirect=/checkout')
      return
    }

    setSubmitting(true)
    try {
      const orderPayload = {
        items: cart.map((item) => ({
          product_id: item.product.id,
          quantity: item.quantity,
        })),
        shipping_name: shippingData.shipping_name,
        shipping_phone: shippingData.shipping_phone,
        shipping_division: shippingData.shipping_division,
        shipping_district: shippingData.shipping_district,
        shipping_thana: shippingData.shipping_thana,
        shipping_address: shippingData.shipping_address,
        shipping_postal: shippingData.shipping_postal,
        payment_method: paymentMethod,
        coupon_code: appliedCoupon?.code || '',
        note: shippingData.note,
      }

      const res = await api.post('/orders', orderPayload)
      if (res?.order) {
        clearCart()
        toast.success('Order placed successfully!')
        navigate(`/order-success/${res.order.id}`, { state: { order: res.order } })
      }
    } catch (err) {
      toast.error(err.message || 'Failed to place order')
    } finally {
      setSubmitting(false)
    }
  }

  if (cart.length === 0) {
    return (
      <div className="container" style={{ padding: '80px 16px', textAlign: 'center' }}>
        <h2>No items to checkout</h2>
        <Link to="/" className="btn btn-primary" style={{ marginTop: 16 }}>
          Browse Products
        </Link>
      </div>
    )
  }

  return (
    <div className="container" style={{ padding: '24px 16px' }}>
      <div className="breadcrumbs">
        <Link to="/">Home</Link> <ChevronRight size={12} style={{ display: 'inline' }} />{' '}
        <Link to="/cart">Cart</Link> <ChevronRight size={12} style={{ display: 'inline' }} /> <span>Checkout</span>
      </div>

      <h1 style={{ fontSize: 26, fontWeight: 800, marginBottom: 20 }}>Checkout</h1>

      {/* 3 Step Indicator */}
      <div className="checkout-steps">
        <div className={`checkout-step ${currentStep === 1 ? 'active' : currentStep > 1 ? 'done' : ''}`}>
          {currentStep > 1 ? <Check size={16} style={{ display: 'inline', marginRight: 4 }} /> : '1.'} Shipping Address
        </div>
        <div className={`checkout-step ${currentStep === 2 ? 'active' : currentStep > 2 ? 'done' : ''}`}>
          {currentStep > 2 ? <Check size={16} style={{ display: 'inline', marginRight: 4 }} /> : '2.'} Payment Method
        </div>
        <div className={`checkout-step ${currentStep === 3 ? 'active' : ''}`}>
          3. Order Review & Confirm
        </div>
      </div>

      <div className="checkout-layout">
        {/* Step 1: Shipping Details */}
        {currentStep === 1 && (
          <div className="checkout-card">
            <h2 style={{ fontSize: 18, fontWeight: 700, marginBottom: 20 }}>1. Shipping Details</h2>
            <form onSubmit={handleShippingSubmit} className="form-grid">
              <div className="form-group">
                <label>Full Name *</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. Tanvir Ahmed"
                  value={shippingData.shipping_name}
                  onChange={(e) => setShippingData({ ...shippingData, shipping_name: e.target.value })}
                />
              </div>

              <div className="form-group">
                <label>Contact Phone Number *</label>
                <input
                  type="tel"
                  required
                  placeholder="e.g. 01700000000"
                  value={shippingData.shipping_phone}
                  onChange={(e) => setShippingData({ ...shippingData, shipping_phone: e.target.value })}
                />
              </div>

              <div className="form-group">
                <label>Division</label>
                <select
                  value={shippingData.shipping_division}
                  onChange={(e) => setShippingData({ ...shippingData, shipping_division: e.target.value })}
                >
                  {divisions.map((d) => (
                    <option key={d} value={d}>{d}</option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label>City / District *</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. Dhaka"
                  value={shippingData.shipping_district}
                  onChange={(e) => setShippingData({ ...shippingData, shipping_district: e.target.value })}
                />
              </div>

              <div className="form-group">
                <label>Thana / Upazila</label>
                <input
                  type="text"
                  placeholder="e.g. Dhanmondi"
                  value={shippingData.shipping_thana}
                  onChange={(e) => setShippingData({ ...shippingData, shipping_thana: e.target.value })}
                />
              </div>

              <div className="form-group">
                <label>Postal Code</label>
                <input
                  type="text"
                  placeholder="e.g. 1205"
                  value={shippingData.shipping_postal}
                  onChange={(e) => setShippingData({ ...shippingData, shipping_postal: e.target.value })}
                />
              </div>

              <div className="form-group form-full">
                <label>Full Street Address *</label>
                <textarea
                  rows={3}
                  required
                  placeholder="House #, Road #, Area details..."
                  value={shippingData.shipping_address}
                  onChange={(e) => setShippingData({ ...shippingData, shipping_address: e.target.value })}
                />
              </div>

              <div className="form-group form-full">
                <label>Order Special Notes (Optional)</label>
                <input
                  type="text"
                  placeholder="Special instructions for delivery..."
                  value={shippingData.note}
                  onChange={(e) => setShippingData({ ...shippingData, note: e.target.value })}
                />
              </div>

              <div className="form-full" style={{ display: 'flex', justifyContent: 'flex-end', marginTop: 12 }}>
                <button type="submit" className="btn btn-primary btn-lg">
                  Continue to Payment <ArrowRight size={18} />
                </button>
              </div>
            </form>
          </div>
        )}

        {/* Step 2: Payment Method */}
        {currentStep === 2 && (
          <div className="checkout-card">
            <h2 style={{ fontSize: 18, fontWeight: 700, marginBottom: 20 }}>2. Select Payment Method</h2>

            <div className="payment-methods-grid">
              {[
                { id: 'COD', label: 'Cash on Delivery', desc: 'Pay with cash upon receiving product' },
                { id: 'BKASH', label: 'bKash Online', desc: 'Instant mobile payment via bKash gateway' },
                { id: 'NAGAD', label: 'Nagad Online', desc: 'Pay directly using Nagad digital wallet' },
                { id: 'ROCKET', label: 'Rocket (DBBL)', desc: 'Dutch-Bangla Bank mobile payment' },
                { id: 'CARD', label: 'Credit / Debit Card', desc: 'VISA, Mastercard, AMEX with SSLCommerz' },
                { id: 'BANK', label: 'Bank Transfer / Wire', desc: 'Direct bank deposit into StarTech account' },
              ].map((m) => (
                <div
                  key={m.id}
                  className={`payment-card-opt ${paymentMethod === m.id ? 'selected' : ''}`}
                  onClick={() => setPaymentMethod(m.id)}
                >
                  <input
                    type="radio"
                    name="payment"
                    checked={paymentMethod === m.id}
                    onChange={() => setPaymentMethod(m.id)}
                  />
                  <div>
                    <b style={{ display: 'block', fontSize: 14 }}>{m.label}</b>
                    <small style={{ color: '#64748b' }}>{m.desc}</small>
                  </div>
                </div>
              ))}
            </div>

            <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 24 }}>
              <button
                type="button"
                className="btn btn-outline"
                onClick={() => setCurrentStep(1)}
              >
                <ArrowLeft size={16} /> Back to Shipping
              </button>
              <button
                type="button"
                className="btn btn-primary btn-lg"
                onClick={() => setCurrentStep(3)}
              >
                Review Order <ArrowRight size={18} />
              </button>
            </div>
          </div>
        )}

        {/* Step 3: Review & Confirmation */}
        {currentStep === 3 && (
          <div className="checkout-card">
            <h2 style={{ fontSize: 18, fontWeight: 700, marginBottom: 20 }}>3. Review & Confirm Order</h2>

            <div style={{
              background: '#f8fafc',
              padding: 16,
              borderRadius: 8,
              marginBottom: 20,
              border: '1px solid #e2e8f0',
              fontSize: 14,
            }}>
              <h4 style={{ fontSize: 15, marginBottom: 8, color: '#0f172a' }}>Delivery Address:</h4>
              <p><b>Recipient:</b> {shippingData.shipping_name} ({shippingData.shipping_phone})</p>
              <p><b>Address:</b> {shippingData.shipping_address}, {shippingData.shipping_district}, {shippingData.shipping_division}</p>
              <p style={{ marginTop: 8 }}><b>Payment Choice:</b> <span style={{ color: '#ef4a23', fontWeight: 700 }}>{paymentMethod}</span></p>
            </div>

            <div style={{ marginBottom: 20 }}>
              <h4 style={{ fontSize: 15, marginBottom: 12 }}>Items ({cart.length}):</h4>
              {cart.map(({ product, quantity }) => (
                <div
                  key={product.id}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    padding: '8px 0',
                    borderBottom: '1px solid #f1f5f9',
                    fontSize: 13,
                  }}
                >
                  <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
                    <img
                      src={product.images?.[0]?.url || 'https://placehold.co/50x50'}
                      alt=""
                      style={{ width: 40, height: 40, objectFit: 'contain' }}
                    />
                    <span>{product.name} × <b>{quantity}</b></span>
                  </div>
                  <b>৳{((product.discount_price || product.price) * quantity).toLocaleString()}</b>
                </div>
              ))}
            </div>

            {!isAuthenticated && (
              <div style={{
                background: '#fffbeb',
                border: '1px solid #fde68a',
                padding: 14,
                borderRadius: 6,
                marginBottom: 20,
                color: '#92400e',
                fontSize: 13,
              }}>
                ℹ️ You need to be logged in to place an order.{' '}
                <Link to="/login?redirect=/checkout" style={{ fontWeight: 700, textDecoration: 'underline' }}>
                  Click here to login or register.
                </Link>
              </div>
            )}

            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <button
                type="button"
                className="btn btn-outline"
                onClick={() => setCurrentStep(2)}
              >
                <ArrowLeft size={16} /> Back to Payment
              </button>
              <button
                type="button"
                className="btn btn-primary btn-lg"
                disabled={submitting}
                onClick={handlePlaceOrder}
              >
                {submitting ? 'Placing Order...' : 'Confirm & Place Order'} <ShieldCheck size={18} />
              </button>
            </div>
          </div>
        )}

        {/* Order Summary Sidebar */}
        <div>
          <div className="order-summary-card">
            <h3>Summary</h3>
            <div className="summary-row">
              <span>Items Total:</span>
              <b>৳{subtotal.toLocaleString()}</b>
            </div>
            {discount > 0 && (
              <div className="summary-row" style={{ color: '#10b981' }}>
                <span>Coupon ({appliedCoupon?.code}):</span>
                <b>-৳{discount.toLocaleString()}</b>
              </div>
            )}
            <div className="summary-row">
              <span>Delivery Fee:</span>
              <b>{deliveryFee > 0 ? `৳${deliveryFee}` : 'Free'}</b>
            </div>
            <div className="summary-row summary-total">
              <span>Payable:</span>
              <span>৳{grandTotal.toLocaleString()}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
