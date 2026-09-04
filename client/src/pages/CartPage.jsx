import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useShop } from '../context/ShopContext'
import { Trash2, Plus, Minus, ShoppingCart, ArrowRight, Tag, ChevronRight } from 'lucide-react'

export default function CartPage() {
  const {
    cart,
    updateCartQty,
    removeFromCart,
    clearCart,
    subtotal,
    discount,
    deliveryFee,
    grandTotal,
    appliedCoupon,
    applyCoupon,
    removeCoupon,
  } = useShop()

  const [couponInput, setCouponInput] = useState('')
  const [validating, setValidating] = useState(false)
  const navigate = useNavigate()

  const handleApplyCoupon = async (e) => {
    e.preventDefault()
    if (!couponInput.trim()) return
    setValidating(true)
    await applyCoupon(couponInput.trim().toUpperCase())
    setValidating(false)
  }

  if (cart.length === 0) {
    return (
      <div className="container" style={{ padding: '80px 16px', textAlign: 'center' }}>
        <div style={{
          background: '#fff',
          borderRadius: 8,
          padding: 60,
          maxWidth: 480,
          margin: '0 auto',
          border: '1px solid #e2e8f0',
        }}>
          <ShoppingCart size={54} color="#94a3b8" style={{ margin: '0 auto 16px' }} />
          <h2 style={{ fontSize: 22, fontWeight: 700, marginBottom: 8 }}>Your Shopping Cart is Empty</h2>
          <p style={{ color: '#64748b', marginBottom: 24, fontSize: 14 }}>
            Looks like you haven't added any tech items to your cart yet.
          </p>
          <Link to="/" className="btn btn-primary btn-lg">
            Start Shopping <ArrowRight size={16} />
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="container" style={{ padding: '24px 16px' }}>
      <div className="breadcrumbs">
        <Link to="/">Home</Link> <ChevronRight size={12} style={{ display: 'inline' }} /> <span>Shopping Cart</span>
      </div>

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
        <h1 style={{ fontSize: 24, fontWeight: 800 }}>Shopping Cart ({cart.length} items)</h1>
        <button className="btn btn-outline btn-sm" onClick={clearCart}>
          <Trash2 size={14} /> Clear Cart
        </button>
      </div>

      <div className="cart-layout">
        {/* Left: Items Table */}
        <div className="cart-items-card">
          <table className="cart-table">
            <thead>
              <tr>
                <th>Product Description</th>
                <th style={{ textAlign: 'center' }}>Unit Price</th>
                <th style={{ textAlign: 'center' }}>Quantity</th>
                <th style={{ textAlign: 'right' }}>Total</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {cart.map(({ product, quantity }) => {
                const unitPrice = product.discount_price || product.price
                const itemTotal = unitPrice * quantity
                return (
                  <tr key={product.id}>
                    <td>
                      <div className="cart-item-info">
                        <img
                          src={product.images?.[0]?.url || 'https://placehold.co/100x100'}
                          alt={product.name}
                        />
                        <div>
                          <Link
                            to={`/product/${product.slug}`}
                            style={{ fontWeight: 600, fontSize: 14, color: '#0f172a' }}
                          >
                            {product.name}
                          </Link>
                          <div style={{ fontSize: 12, color: '#64748b', marginTop: 2 }}>
                            Brand: {product.brand?.name || 'StarTech'}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td style={{ textAlign: 'center', fontWeight: 600 }}>
                      ৳{unitPrice.toLocaleString()}
                    </td>
                    <td>
                      <div className="qty-control" style={{ margin: '0 auto', width: 'fit-content' }}>
                        <button onClick={() => updateCartQty(product.id, quantity - 1)}>
                          <Minus size={12} />
                        </button>
                        <span>{quantity}</span>
                        <button onClick={() => updateCartQty(product.id, quantity + 1)}>
                          <Plus size={12} />
                        </button>
                      </div>
                    </td>
                    <td style={{ textAlign: 'right', fontWeight: 700, color: '#ef4a23', fontSize: 15 }}>
                      ৳{itemTotal.toLocaleString()}
                    </td>
                    <td style={{ textAlign: 'right' }}>
                      <button
                        onClick={() => removeFromCart(product.id)}
                        style={{ color: '#ef4444', padding: 6 }}
                        title="Remove item"
                      >
                        <Trash2 size={16} />
                      </button>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>

        {/* Right: Order Summary & Coupon */}
        <div>
          <div className="order-summary-card">
            <h3>Order Summary</h3>

            <div className="summary-row">
              <span style={{ color: '#64748b' }}>Subtotal:</span>
              <b>৳{subtotal.toLocaleString()}</b>
            </div>

            {discount > 0 && (
              <div className="summary-row" style={{ color: '#10b981' }}>
                <span>Discount ({appliedCoupon?.code}):</span>
                <b>-৳{discount.toLocaleString()}</b>
              </div>
            )}

            <div className="summary-row">
              <span style={{ color: '#64748b' }}>Estimated Delivery:</span>
              <b>{deliveryFee > 0 ? `৳${deliveryFee}` : 'Free'}</b>
            </div>

            <div className="summary-row summary-total">
              <span>Grand Total:</span>
              <span>৳{grandTotal.toLocaleString()}</span>
            </div>

            {/* Coupon input */}
            <div style={{ marginTop: 20, paddingTop: 16, borderTop: '1px solid #f1f5f9' }}>
              <label style={{ fontSize: 12, fontWeight: 700, color: '#475569', display: 'flex', alignItems: 'center', gap: 6 }}>
                <Tag size={14} color="#ef4a23" /> Have a Coupon Code?
              </label>
              {appliedCoupon ? (
                <div style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  background: '#ecfdf5',
                  padding: '8px 12px',
                  borderRadius: 4,
                  marginTop: 8,
                  fontSize: 13,
                  color: '#065f46',
                }}>
                  <span><b>{appliedCoupon.code}</b> Applied</span>
                  <button onClick={removeCoupon} style={{ color: '#ef4444', fontWeight: 600 }}>
                    Remove
                  </button>
                </div>
              ) : (
                <form onSubmit={handleApplyCoupon} className="coupon-box">
                  <input
                    type="text"
                    placeholder="Enter Coupon (e.g. WELCOME500)"
                    value={couponInput}
                    onChange={(e) => setCouponInput(e.target.value)}
                  />
                  <button type="submit" className="btn btn-secondary btn-sm" disabled={validating}>
                    {validating ? '...' : 'Apply'}
                  </button>
                </form>
              )}
            </div>

            <button
              className="btn btn-primary btn-lg"
              style={{ width: '100%', marginTop: 20 }}
              onClick={() => navigate('/checkout')}
            >
              Proceed to Checkout <ArrowRight size={18} />
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
