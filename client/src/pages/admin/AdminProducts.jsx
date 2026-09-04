import React, { useState, useEffect } from 'react'
import { Plus, Edit2, Trash2, Search, Upload, X, Check } from 'lucide-react'
import api from '../../api/client'
import toast from 'react-hot-toast'

export default function AdminProducts() {
  const [products, setProducts] = useState([])
  const [categories, setCategories] = useState([])
  const [brands, setBrands] = useState([])
  const [totalResults, setTotalResults] = useState(0)
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(false)

  // Modal State
  const [modalOpen, setModalOpen] = useState(false)
  const [editingProduct, setEditingProduct] = useState(null)
  const [uploadingImage, setUploadingImage] = useState(false)

  // Form State
  const [formData, setFormData] = useState({
    name: '',
    sku: '',
    category_id: '',
    brand_id: '',
    price: '',
    discount_price: '',
    stock: '',
    short_description: '',
    description: '',
    featured: false,
    status: 'ACTIVE',
    images: [],
    specs: {},
  })

  // Specs helper state (Key-Value)
  const [specKey, setSpecKey] = useState('')
  const [specVal, setSpecVal] = useState('')

  const fetchCatalogMeta = async () => {
    try {
      const [catRes, brandRes] = await Promise.all([
        api.get('/categories'),
        api.get('/brands'),
      ])
      setCategories(catRes?.categories || [])
      setBrands(brandRes?.brands || [])
    } catch {}
  }

  const fetchProducts = async () => {
    setLoading(true)
    try {
      let endpoint = `/products?page=${page}&limit=10`
      if (search.trim()) endpoint += `&search=${encodeURIComponent(search.trim())}`
      const res = await api.get(endpoint)
      setProducts(res?.results || [])
      setTotalResults(res?.total_results || 0)
    } catch (err) {
      toast.error('Failed to load products')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchCatalogMeta()
  }, [])

  useEffect(() => {
    fetchProducts()
  }, [page, search])

  const handleOpenCreate = () => {
    setEditingProduct(null)
    setFormData({
      name: '',
      sku: 'ST-' + Math.floor(100000 + Math.random() * 900000),
      category_id: categories[0]?.id || '',
      brand_id: brands[0]?.id || '',
      price: '',
      discount_price: '',
      stock: '10',
      short_description: '',
      description: '',
      featured: false,
      status: 'ACTIVE',
      images: [],
      specs: {},
    })
    setModalOpen(true)
  }

  const handleOpenEdit = (p) => {
    setEditingProduct(p)
    let parsedSpecs = {}
    try {
      parsedSpecs = typeof p.specs === 'string' ? JSON.parse(p.specs) : p.specs || {}
    } catch {
      parsedSpecs = {}
    }

    setFormData({
      name: p.name,
      sku: p.sku,
      category_id: p.category_id,
      brand_id: p.brand_id,
      price: p.price,
      discount_price: p.discount_price || '',
      stock: p.stock,
      short_description: p.short_description || '',
      description: p.description || '',
      featured: p.featured,
      status: p.status,
      images: p.images ? p.images.map((i) => i.url) : [],
      specs: parsedSpecs,
    })
    setModalOpen(true)
  }

  const handleImageUpload = async (e) => {
    const file = e.target.files[0]
    if (!file) return
    const uploadData = new FormData()
    uploadData.append('image', file)

    setUploadingImage(true)
    try {
      const res = await api.post('/admin/upload', uploadData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      if (res?.url) {
        setFormData((prev) => ({
          ...prev,
          images: [...prev.images, res.url],
        }))
        toast.success('Image uploaded successfully')
      }
    } catch (err) {
      toast.error(err.message || 'Image upload failed')
    } finally {
      setUploadingImage(false)
    }
  }

  const handleAddSpec = () => {
    if (!specKey.trim() || !specVal.trim()) return
    setFormData((prev) => ({
      ...prev,
      specs: { ...prev.specs, [specKey.trim()]: specVal.trim() },
    }))
    setSpecKey('')
    setSpecVal('')
  }

  const handleRemoveSpec = (key) => {
    setFormData((prev) => {
      const copy = { ...prev.specs }
      delete copy[key]
      return { ...prev, specs: copy }
    })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!formData.name || !formData.price || !formData.category_id || !formData.brand_id) {
      toast.error('Please fill in required fields')
      return
    }

    const payload = {
      name: formData.name,
      sku: formData.sku,
      category_id: formData.category_id,
      brand_id: formData.brand_id,
      price: Number(formData.price),
      discount_price: formData.discount_price ? Number(formData.discount_price) : null,
      stock: Number(formData.stock || 0),
      short_description: formData.short_description,
      description: formData.description,
      featured: Boolean(formData.featured),
      status: formData.status,
      images: formData.images,
    }

    try {
      if (editingProduct) {
        await api.put(`/admin/products/${editingProduct.id}`, payload)
        toast.success('Product updated successfully')
      } else {
        await api.post('/admin/products', payload)
        toast.success('Product created successfully')
      }
      setModalOpen(false)
      fetchProducts()
    } catch (err) {
      toast.error(err.message || 'Operation failed')
    }
  }

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this product?')) return
    try {
      await api.delete(`/admin/products/${id}`)
      toast.success('Product deleted')
      fetchProducts()
    } catch (err) {
      toast.error(err.message || 'Delete failed')
    }
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
        <div>
          <h1 style={{ fontSize: 24, fontWeight: 800 }}>Products Management</h1>
          <small style={{ color: '#64748b' }}>Total products: {totalResults}</small>
        </div>
        <button className="btn btn-primary" onClick={handleOpenCreate}>
          <Plus size={16} /> Add New Product
        </button>
      </div>

      {/* Search Filter */}
      <div className="admin-card" style={{ padding: 16, marginBottom: 20 }}>
        <div style={{ display: 'flex', gap: 12 }}>
          <div style={{ position: 'relative', flex: 1 }}>
            <input
              type="text"
              placeholder="Search products by title, SKU, brand..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              style={{
                width: '100%',
                padding: '10px 14px',
                borderRadius: 4,
                border: '1px solid #cbd5e1',
                fontSize: 14,
              }}
            />
          </div>
        </div>
      </div>

      {/* Products Table */}
      <div className="admin-card">
        <table className="admin-table">
          <thead>
            <tr>
              <th>Image</th>
              <th>Product Title</th>
              <th>Category</th>
              <th>Brand</th>
              <th>Price</th>
              <th>Stock</th>
              <th>Status</th>
              <th style={{ textAlign: 'right' }}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr><td colSpan={8} style={{ textAlign: 'center', padding: 32 }}>Loading products...</td></tr>
            ) : products.length > 0 ? (
              products.map((p) => (
                <tr key={p.id}>
                  <td>
                    <img
                      src={p.images?.[0]?.url || 'https://placehold.co/50x50'}
                      alt=""
                      style={{ width: 44, height: 44, objectFit: 'contain', borderRadius: 4, border: '1px solid #e2e8f0' }}
                    />
                  </td>
                  <td>
                    <b style={{ fontSize: 14, color: '#0f172a', display: 'block' }}>{p.name}</b>
                    <small style={{ color: '#64748b' }}>SKU: {p.sku}</small>
                  </td>
                  <td>{p.category?.name || '—'}</td>
                  <td>{p.brand?.name || '—'}</td>
                  <td>
                    <b>৳{(p.discount_price || p.price).toLocaleString()}</b>
                    {p.discount_price && (
                      <small style={{ display: 'block', color: '#94a3b8', textDecoration: 'line-through' }}>
                        ৳{p.price.toLocaleString()}
                      </small>
                    )}
                  </td>
                  <td>
                    <b style={{ color: p.stock > 0 ? '#10b981' : '#ef4444' }}>{p.stock}</b>
                  </td>
                  <td>
                    <span className={`badge ${p.status === 'ACTIVE' ? 'badge-stock' : 'badge-out'}`}>
                      {p.status}
                    </span>
                  </td>
                  <td style={{ textAlign: 'right' }}>
                    <div style={{ display: 'flex', gap: 6, justifyContent: 'flex-end' }}>
                      <button className="btn btn-outline btn-sm" onClick={() => handleOpenEdit(p)}>
                        <Edit2 size={14} />
                      </button>
                      <button
                        className="btn btn-outline btn-sm"
                        onClick={() => handleDelete(p.id)}
                        style={{ color: '#ef4444' }}
                      >
                        <Trash2 size={14} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))
            ) : (
              <tr><td colSpan={8} style={{ textAlign: 'center', padding: 32 }}>No products found</td></tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Modal for Create/Edit */}
      {modalOpen && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          background: 'rgba(0,0,0,0.5)',
          zIndex: 1000,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          padding: 20,
        }}>
          <div style={{
            background: '#fff',
            borderRadius: 8,
            maxWidth: 750,
            width: '100%',
            maxHeight: '90vh',
            overflowY: 'auto',
            padding: 28,
            boxShadow: '0 20px 25px -5px rgba(0,0,0,0.2)',
          }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20, borderBottom: '1px solid #e2e8f0', paddingBottom: 12 }}>
              <h2 style={{ fontSize: 20, fontWeight: 700 }}>
                {editingProduct ? 'Edit Product' : 'Add New Product'}
              </h2>
              <button onClick={() => setModalOpen(false)}>
                <X size={20} />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="form-grid">
              <div className="form-group form-full">
                <label>Product Title *</label>
                <input
                  type="text"
                  required
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                />
              </div>

              <div className="form-group">
                <label>SKU Code *</label>
                <input
                  type="text"
                  required
                  value={formData.sku}
                  onChange={(e) => setFormData({ ...formData, sku: e.target.value })}
                />
              </div>

              <div className="form-group">
                <label>Category *</label>
                <select
                  value={formData.category_id}
                  onChange={(e) => setFormData({ ...formData, category_id: e.target.value })}
                >
                  {categories.map((c) => (
                    <option key={c.id} value={c.id}>{c.name}</option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label>Brand *</label>
                <select
                  value={formData.brand_id}
                  onChange={(e) => setFormData({ ...formData, brand_id: e.target.value })}
                >
                  {brands.map((b) => (
                    <option key={b.id} value={b.id}>{b.name}</option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label>Regular Price (৳) *</label>
                <input
                  type="number"
                  required
                  value={formData.price}
                  onChange={(e) => setFormData({ ...formData, price: e.target.value })}
                />
              </div>

              <div className="form-group">
                <label>Discounted / Special Price (৳)</label>
                <input
                  type="number"
                  placeholder="Optional"
                  value={formData.discount_price}
                  onChange={(e) => setFormData({ ...formData, discount_price: e.target.value })}
                />
              </div>

              <div className="form-group">
                <label>Stock Quantity</label>
                <input
                  type="number"
                  value={formData.stock}
                  onChange={(e) => setFormData({ ...formData, stock: e.target.value })}
                />
              </div>

              <div className="form-group form-full">
                <label>Short Key Description (Features)</label>
                <input
                  type="text"
                  placeholder="e.g. Intel Core i7, 16GB DDR5, RTX 4060"
                  value={formData.short_description}
                  onChange={(e) => setFormData({ ...formData, short_description: e.target.value })}
                />
              </div>

              <div className="form-group form-full">
                <label>Detailed Description</label>
                <textarea
                  rows={3}
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                />
              </div>

              {/* Image Upload */}
              <div className="form-group form-full">
                <label>Product Images</label>
                <div style={{ display: 'flex', gap: 10, flexWrap: 'wrap', alignItems: 'center' }}>
                  {formData.images.map((img, idx) => (
                    <div key={idx} style={{ position: 'relative' }}>
                      <img src={img} alt="" style={{ width: 60, height: 60, objectFit: 'contain', border: '1px solid #e2e8f0', borderRadius: 4 }} />
                      <button
                        type="button"
                        onClick={() => setFormData({ ...formData, images: formData.images.filter((_, i) => i !== idx) })}
                        style={{ position: 'absolute', top: -6, right: -6, background: '#ef4444', color: '#fff', borderRadius: '50%', width: 18, height: 18, fontSize: 10 }}
                      >
                        ×
                      </button>
                    </div>
                  ))}
                  <label style={{
                    border: '2px dashed #cbd5e1',
                    borderRadius: 4,
                    padding: '12px 18px',
                    cursor: 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    gap: 6,
                    fontSize: 13,
                    color: '#64748b',
                  }}>
                    <Upload size={16} /> {uploadingImage ? 'Uploading...' : 'Upload Image'}
                    <input type="file" accept="image/*" onChange={handleImageUpload} style={{ display: 'none' }} />
                  </label>
                </div>
              </div>

              {/* Specs Builder */}
              <div className="form-group form-full">
                <label>Technical Specifications</label>
                <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
                  <input
                    type="text"
                    placeholder="Key (e.g. Processor)"
                    value={specKey}
                    onChange={(e) => setSpecKey(e.target.value)}
                  />
                  <input
                    type="text"
                    placeholder="Value (e.g. Intel Core i9)"
                    value={specVal}
                    onChange={(e) => setSpecVal(e.target.value)}
                  />
                  <button type="button" className="btn btn-secondary btn-sm" onClick={handleAddSpec}>
                    Add Spec
                  </button>
                </div>
                <div style={{ display: 'flex', flexWrap: 'wrap', gap: 6 }}>
                  {Object.entries(formData.specs).map(([k, v]) => (
                    <span key={k} style={{ background: '#f1f5f9', padding: '4px 10px', borderRadius: 4, fontSize: 12, display: 'flex', alignItems: 'center', gap: 6 }}>
                      <b>{k}:</b> {v}
                      <button type="button" onClick={() => handleRemoveSpec(k)} style={{ color: '#ef4444' }}>×</button>
                    </span>
                  ))}
                </div>
              </div>

              <div className="form-full" style={{ display: 'flex', justifyContent: 'flex-end', gap: 10, marginTop: 16 }}>
                <button type="button" className="btn btn-outline" onClick={() => setModalOpen(false)}>
                  Cancel
                </button>
                <button type="submit" className="btn btn-primary">
                  Save Product
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
