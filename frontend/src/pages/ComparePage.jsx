import React from 'react'
import { Link } from 'react-router-dom'
import { useShop } from '../context/ShopContext'
import { Trash2, ShoppingCart, Scale, ChevronRight } from 'lucide-react'

export default function ComparePage() {
  const { compare, removeFromCompare, clearCompare, addToCart } = useShop()

  if (compare.length === 0) {
    return (
      <div className="container" style={{ padding: '80px 16px', textAlign: 'center' }}>
        <div style={{
          background: '#fff',
          borderRadius: 8,
          padding: 60,
          maxWidth: 500,
          margin: '0 auto',
          border: '1px solid #e2e8f0',
        }}>
          <Scale size={48} color="#94a3b8" style={{ margin: '0 auto 16px' }} />
          <h2 style={{ fontSize: 20, marginBottom: 8 }}>Your comparison list is empty</h2>
          <p style={{ color: '#64748b', marginBottom: 20 }}>
            Browse our catalog and click the compare icon on any product to compare specifications side by side.
          </p>
          <Link to="/" className="btn btn-primary">
            Explore Products
          </Link>
        </div>
      </div>
    )
  }

  // Collect all unique spec keys from all compare products
  const allSpecKeys = Array.from(
    new Set(
      compare.flatMap((p) => {
        try {
          const s = typeof p.specs === 'string' ? JSON.parse(p.specs) : p.specs || {}
          return Object.keys(s)
        } catch {
          return []
        }
      })
    )
  )

  return (
    <div className="container" style={{ padding: '24px 16px' }}>
      <div className="breadcrumbs">
        <Link to="/">Home</Link> <ChevronRight size={12} style={{ display: 'inline' }} /> <span>Product Comparison</span>
      </div>

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
        <h1 style={{ fontSize: 24, fontWeight: 800 }}>Product Comparison ({compare.length}/4)</h1>
        <button className="btn btn-outline btn-sm" onClick={clearCompare}>
          <Trash2 size={14} /> Clear All
        </button>
      </div>

      <div className="compare-container">
        <table className="compare-table">
          <thead>
            <tr>
              <th style={{ width: 200 }}>Product</th>
              {compare.map((p) => (
                <td key={p.id} style={{ verticalAlign: 'top', minWidth: 220 }}>
                  <button
                    onClick={() => removeFromCompare(p.id)}
                    style={{
                      float: 'right',
                      color: '#ef4444',
                      padding: 4,
                      borderRadius: 4,
                    }}
                    title="Remove from comparison"
                  >
                    <Trash2 size={16} />
                  </button>
                  <img
                    src={p.images?.[0]?.url || 'https://placehold.co/200x200'}
                    alt={p.name}
                    style={{ width: 140, height: 140, objectFit: 'contain', margin: '0 auto 12px' }}
                  />
                  <Link
                    to={`/product/${p.slug}`}
                    style={{ fontWeight: 700, fontSize: 14, color: '#0f172a', display: 'block', marginBottom: 8 }}
                  >
                    {p.name}
                  </Link>
                  <div style={{ fontSize: 18, fontWeight: 800, color: '#ef4a23', marginBottom: 12 }}>
                    ৳{(p.discount_price || p.price).toLocaleString()}
                  </div>
                  <button
                    className="btn btn-primary btn-sm"
                    style={{ width: '100%' }}
                    onClick={() => addToCart(p, 1)}
                  >
                    <ShoppingCart size={14} /> Add to Cart
                  </button>
                </td>
              ))}
            </tr>
          </thead>
          <tbody>
            <tr>
              <th>Brand</th>
              {compare.map((p) => (
                <td key={p.id}><b>{p.brand?.name || 'MI-Tech'}</b></td>
              ))}
            </tr>
            <tr>
              <th>Category</th>
              {compare.map((p) => (
                <td key={p.id}>{p.category?.name || 'Tech'}</td>
              ))}
            </tr>
            <tr>
              <th>Availability</th>
              {compare.map((p) => (
                <td key={p.id} style={{ color: p.stock > 0 ? '#10b981' : '#ef4444', fontWeight: 600 }}>
                  {p.stock > 0 ? 'In Stock' : 'Out of Stock'}
                </td>
              ))}
            </tr>
            <tr>
              <th>Rating</th>
              {compare.map((p) => (
                <td key={p.id} style={{ color: '#f59e0b', fontWeight: 700 }}>
                  ★ {p.rating_avg ? p.rating_avg.toFixed(1) : '5.0'} / 5.0
                </td>
              ))}
            </tr>
            {/* Dynamic Specs Rows */}
            {allSpecKeys.map((key) => (
              <tr key={key}>
                <th>{key}</th>
                {compare.map((p) => {
                  let val = '—'
                  try {
                    const s = typeof p.specs === 'string' ? JSON.parse(p.specs) : p.specs || {}
                    if (s[key]) val = s[key]
                  } catch {
                    val = '—'
                  }
                  return <td key={p.id}>{val}</td>
                })}
              </tr>
            ))}
            <tr>
              <th>Warranty</th>
              {compare.map((p) => (
                <td key={p.id}>1 Year Official Warranty</td>
              ))}
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  )
}
