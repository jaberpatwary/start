import React, { useState, useEffect } from 'react'
import { Boxes, AlertTriangle, CheckCircle, RefreshCw, Save, Search } from 'lucide-react'
import api from '../../api/client'
import toast from 'react-hot-toast'

export default function AdminInventory() {
  const [products, setProducts] = useState([])
  const [stockFilter, setStockFilter] = useState('')
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)

  // Edit stock map: { [productId]: newStockValue }
  const [stockEdits, setStockEdits] = useState({})
  const [savingId, setSavingId] = useState(null)

  const fetchInventory = async () => {
    setLoading(true)
    try {
      let endpoint = '/products?limit=100'
      if (search) endpoint += `&search=${encodeURIComponent(search)}`
      const res = await api.get(endpoint)
      setProducts(res?.results || [])
    } catch (err) {
      toast.error('Failed to load inventory')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchInventory()
  }, [search])

  const handleStockChange = (id, val) => {
    setStockEdits((prev) => ({ ...prev, [id]: val }))
  }

  const handleSaveStock = async (p) => {
    const newStock = stockEdits[p.id]
    if (newStock === undefined || newStock === '' || isNaN(newStock)) return

    setSavingId(p.id)
    try {
      await api.put(`/admin/products/${p.id}`, {
        name: p.name,
        sku: p.sku,
        category_id: p.category_id,
        brand_id: p.brand_id,
        price: p.price,
        discount_price: p.discount_price,
        stock: Number(newStock),
        short_description: p.short_description,
        description: p.description,
        featured: p.featured,
        status: p.status,
      })
      toast.success(`Stock updated for "${p.name}" to ${newStock}`)
      // Update local state
      setProducts((prev) =>
        prev.map((item) => (item.id === p.id ? { ...item, stock: Number(newStock) } : item))
      )
      setStockEdits((prev) => {
        const copy = { ...prev }
        delete copy[p.id]
        return copy
      })
    } catch (err) {
      toast.error(err.message || 'Failed to update stock')
    } finally {
      setSavingId(null)
    }
  }

  const outOfStockCount = products.filter((p) => p.stock <= 0).length
  const lowStockCount = products.filter((p) => p.stock > 0 && p.stock <= 5).length
  const totalStockCount = products.reduce((acc, p) => acc + (p.stock || 0), 0)

  const filteredProducts = products.filter((p) => {
    if (stockFilter === 'out') return p.stock <= 0
    if (stockFilter === 'low') return p.stock > 0 && p.stock <= 5
    if (stockFilter === 'healthy') return p.stock > 5
    return true
  })

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
        <div>
          <h1 style={{ fontSize: 24, fontWeight: 800 }}>Inventory & Stock Management</h1>
          <small style={{ color: '#64748b' }}>Monitor live inventory, reorder stock, and update physical quantities instantly</small>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="stats-grid" style={{ marginBottom: 24 }}>
        <div className="stat-card">
          <div className="stat-info">
            <small>Total Units in Stock</small>
            <b style={{ color: '#10b981' }}>{totalStockCount.toLocaleString()} units</b>
          </div>
          <div className="stat-icon" style={{ background: '#ecfdf5', color: '#10b981' }}>
            <Boxes size={24} />
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-info">
            <small>Low Stock Warnings (≤ 5)</small>
            <b style={{ color: '#d97706' }}>{lowStockCount} items</b>
          </div>
          <div className="stat-icon" style={{ background: '#fef3c7', color: '#d97706' }}>
            <AlertTriangle size={24} />
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-info">
            <small>Out of Stock Items</small>
            <b style={{ color: '#ef4444' }}>{outOfStockCount} items</b>
          </div>
          <div className="stat-icon" style={{ background: '#fee2e2', color: '#ef4444' }}>
            <AlertTriangle size={24} />
          </div>
        </div>
      </div>

      {/* Search & Stock Filter Bar */}
      <div className="admin-card" style={{ padding: 16, marginBottom: 20 }}>
        <div style={{ display: 'flex', gap: 12, alignItems: 'center', flexWrap: 'wrap' }}>
          <input
            type="text"
            placeholder="Search by product title or SKU..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            style={{ flex: 1, minWidth: 260, padding: '10px 14px', borderRadius: 4, border: '1px solid #cbd5e1', fontSize: 14 }}
          />

          <div style={{ display: 'flex', gap: 6 }}>
            {['', 'out', 'low', 'healthy'].map((st) => (
              <button
                key={st}
                onClick={() => setStockFilter(st)}
                className={`btn btn-sm ${stockFilter === st ? 'btn-primary' : 'btn-outline'}`}
              >
                {st === '' ? 'All Items' : st === 'out' ? 'Out of Stock' : st === 'low' ? 'Low Stock (≤5)' : 'Healthy (>5)'}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Inventory Table */}
      <div className="admin-card">
        <table className="admin-table">
          <thead>
            <tr>
              <th>SKU</th>
              <th>Product Name</th>
              <th>Category</th>
              <th>Current Stock</th>
              <th>Status</th>
              <th>Quick Restock / Adjust</th>
              <th style={{ textAlign: 'right' }}>Save</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr><td colSpan={7} style={{ textAlign: 'center', padding: 32 }}>Loading inventory...</td></tr>
            ) : filteredProducts.length > 0 ? (
              filteredProducts.map((p) => {
                const currentEditVal = stockEdits[p.id] !== undefined ? stockEdits[p.id] : p.stock
                const isModified = stockEdits[p.id] !== undefined && Number(stockEdits[p.id]) !== p.stock

                return (
                  <tr key={p.id}>
                    <td>
                      <span style={{ fontFamily: 'monospace', fontSize: 12, color: '#3749bb' }}>{p.sku}</span>
                    </td>
                    <td>
                      <b style={{ fontSize: 14, color: '#0f172a' }}>{p.name}</b>
                      <small style={{ display: 'block', color: '#64748b' }}>Brand: {p.brand?.name || '—'}</small>
                    </td>
                    <td>{p.category?.name || '—'}</td>
                    <td>
                      <b style={{
                        fontSize: 16,
                        color: p.stock <= 0 ? '#ef4444' : p.stock <= 5 ? '#d97706' : '#10b981'
                      }}>
                        {p.stock}
                      </b>
                    </td>
                    <td>
                      <span style={{
                        padding: '4px 8px',
                        borderRadius: 4,
                        fontSize: 11,
                        fontWeight: 800,
                        background: p.stock <= 0 ? '#fee2e2' : p.stock <= 5 ? '#fef3c7' : '#dcfce7',
                        color: p.stock <= 0 ? '#ef4444' : p.stock <= 5 ? '#d97706' : '#15803d',
                      }}>
                        {p.stock <= 0 ? 'OUT OF STOCK' : p.stock <= 5 ? 'LOW STOCK' : 'IN STOCK'}
                      </span>
                    </td>
                    <td>
                      <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                        <input
                          type="number"
                          value={currentEditVal}
                          onChange={(e) => handleStockChange(p.id, e.target.value)}
                          style={{
                            width: 80,
                            padding: '6px 10px',
                            borderRadius: 4,
                            border: isModified ? '2px solid #ef4a23' : '1px solid #cbd5e1',
                            fontWeight: 700,
                            textAlign: 'center',
                          }}
                        />
                        <div style={{ display: 'flex', gap: 4 }}>
                          <button
                            type="button"
                            className="btn btn-outline btn-sm"
                            onClick={() => handleStockChange(p.id, (Number(currentEditVal) || 0) + 10)}
                            style={{ fontSize: 11, padding: '4px 6px' }}
                          >
                            +10
                          </button>
                          <button
                            type="button"
                            className="btn btn-outline btn-sm"
                            onClick={() => handleStockChange(p.id, (Number(currentEditVal) || 0) + 50)}
                            style={{ fontSize: 11, padding: '4px 6px' }}
                          >
                            +50
                          </button>
                        </div>
                      </div>
                    </td>
                    <td style={{ textAlign: 'right' }}>
                      <button
                        className="btn btn-primary btn-sm"
                        disabled={!isModified || savingId === p.id}
                        onClick={() => handleSaveStock(p)}
                      >
                        <Save size={14} /> {savingId === p.id ? 'Saving...' : 'Update'}
                      </button>
                    </td>
                  </tr>
                )
              })
            ) : (
              <tr><td colSpan={7} style={{ textAlign: 'center', padding: 32 }}>No products match this inventory criteria</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
