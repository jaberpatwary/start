import React, { useState, useEffect } from 'react'
import { Plus, Edit2, Trash2, Tags, Award, Image as ImageIcon, Ticket, X, Upload } from 'lucide-react'
import api from '../../api/client'
import toast from 'react-hot-toast'

export default function AdminCatalog() {
  const [activeTab, setActiveTab] = useState('categories')
  const [categories, setCategories] = useState([])
  const [brands, setBrands] = useState([])
  const [banners, setBanners] = useState([])
  const [loading, setLoading] = useState(false)

  // Modals
  const [catModalOpen, setCatModalOpen] = useState(false)
  const [editingCat, setEditingCat] = useState(null)
  const [catForm, setCatForm] = useState({ name: '', icon: '', sort_order: 0 })

  const [brandModalOpen, setBrandModalOpen] = useState(false)
  const [editingBrand, setEditingBrand] = useState(null)
  const [brandForm, setBrandForm] = useState({ name: '', logo: '' })

  const [bannerModalOpen, setBannerModalOpen] = useState(false)
  const [bannerForm, setBannerForm] = useState({ title: '', subtitle: '', image_url: '', link_url: '', button_text: 'Shop Now' })

  const [couponModalOpen, setCouponModalOpen] = useState(false)
  const [couponForm, setCouponForm] = useState({ code: '', discount_type: 'PERCENTAGE', discount_value: 10, min_order_amount: 1000, max_discount: 500, expires_at: '' })

  const fetchData = async () => {
    setLoading(true)
    try {
      const [catRes, brandRes, banRes] = await Promise.all([
        api.get('/categories'),
        api.get('/brands'),
        api.get('/banners').catch(() => ({ banners: [] })),
      ])
      setCategories(catRes?.categories || [])
      setBrands(brandRes?.brands || [])
      setBanners(banRes?.banners || [])
    } catch (err) {
      toast.error('Failed to load catalog data')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  // --- Category Handlers ---
  const handleOpenCatModal = (cat = null) => {
    if (cat) {
      setEditingCat(cat)
      setCatForm({ name: cat.name, icon: cat.icon || '', sort_order: cat.sort_order || 0 })
    } else {
      setEditingCat(null)
      setCatForm({ name: '', icon: '', sort_order: categories.length + 1 })
    }
    setCatModalOpen(true)
  }

  const handleSaveCat = async (e) => {
    e.preventDefault()
    if (!catForm.name) return toast.error('Category name is required')
    try {
      if (editingCat) {
        await api.put(`/admin/categories/${editingCat.id}`, catForm)
        toast.success('Category updated')
      } else {
        await api.post('/admin/categories', catForm)
        toast.success('Category created')
      }
      setCatModalOpen(false)
      fetchData()
    } catch (err) {
      toast.error(err.message || 'Operation failed')
    }
  }

  const handleDeleteCat = async (id) => {
    if (!window.confirm('Delete this category?')) return
    try {
      await api.delete(`/admin/categories/${id}`)
      toast.success('Category deleted')
      fetchData()
    } catch (err) {
      toast.error(err.message || 'Delete failed')
    }
  }

  // --- Brand Handlers ---
  const handleOpenBrandModal = (b = null) => {
    if (b) {
      setEditingBrand(b)
      setBrandForm({ name: b.name, logo: b.logo || '' })
    } else {
      setEditingBrand(null)
      setBrandForm({ name: '', logo: '' })
    }
    setBrandModalOpen(true)
  }

  const handleSaveBrand = async (e) => {
    e.preventDefault()
    if (!brandForm.name) return toast.error('Brand name is required')
    try {
      if (editingBrand) {
        await api.put(`/admin/brands/${editingBrand.id}`, brandForm)
        toast.success('Brand updated')
      } else {
        await api.post('/admin/brands', brandForm)
        toast.success('Brand created')
      }
      setBrandModalOpen(false)
      fetchData()
    } catch (err) {
      toast.error(err.message || 'Operation failed')
    }
  }

  const handleDeleteBrand = async (id) => {
    if (!window.confirm('Delete this brand?')) return
    try {
      await api.delete(`/admin/brands/${id}`)
      toast.success('Brand deleted')
      fetchData()
    } catch (err) {
      toast.error(err.message || 'Delete failed')
    }
  }

  // --- Banner Handlers ---
  const handleSaveBanner = async (e) => {
    e.preventDefault()
    if (!bannerForm.image_url) return toast.error('Banner image URL is required')
    try {
      await api.post('/admin/banners', bannerForm)
      toast.success('Banner added successfully')
      setBannerModalOpen(false)
      fetchData()
    } catch (err) {
      toast.error(err.message || 'Failed to add banner')
    }
  }

  const handleDeleteBanner = async (id) => {
    if (!window.confirm('Delete this banner?')) return
    try {
      await api.delete(`/admin/banners/${id}`)
      toast.success('Banner deleted')
      fetchData()
    } catch (err) {
      toast.error(err.message || 'Delete failed')
    }
  }

  // --- Coupon Handlers ---
  const handleSaveCoupon = async (e) => {
    e.preventDefault()
    if (!couponForm.code) return toast.error('Coupon code is required')
    try {
      await api.post('/admin/coupons', {
        ...couponForm,
        code: couponForm.code.toUpperCase(),
        discount_value: Number(couponForm.discount_value),
        min_order_amount: Number(couponForm.min_order_amount),
        max_discount: Number(couponForm.max_discount),
      })
      toast.success('Coupon created successfully')
      setCouponModalOpen(false)
    } catch (err) {
      toast.error(err.message || 'Failed to create coupon')
    }
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
        <div>
          <h1 style={{ fontSize: 24, fontWeight: 800 }}>Categories & Catalog Management</h1>
          <small style={{ color: '#64748b' }}>Manage store categories, brand logos, promotional banners, and discount coupons</small>
        </div>
      </div>

      {/* Tabs */}
      <div style={{ display: 'flex', gap: 12, borderBottom: '2px solid #e2e8f0', marginBottom: 24 }}>
        <button
          onClick={() => setActiveTab('categories')}
          style={{
            padding: '10px 18px',
            fontSize: 14,
            fontWeight: 700,
            color: activeTab === 'categories' ? '#ef4a23' : '#64748b',
            borderBottom: activeTab === 'categories' ? '3px solid #ef4a23' : 'none',
            background: 'none',
            display: 'flex',
            alignItems: 'center',
            gap: 8,
          }}
        >
          <Tags size={16} /> Categories ({categories.length})
        </button>
        <button
          onClick={() => setActiveTab('brands')}
          style={{
            padding: '10px 18px',
            fontSize: 14,
            fontWeight: 700,
            color: activeTab === 'brands' ? '#ef4a23' : '#64748b',
            borderBottom: activeTab === 'brands' ? '3px solid #ef4a23' : 'none',
            background: 'none',
            display: 'flex',
            alignItems: 'center',
            gap: 8,
          }}
        >
          <Award size={16} /> Brands ({brands.length})
        </button>
        <button
          onClick={() => setActiveTab('banners')}
          style={{
            padding: '10px 18px',
            fontSize: 14,
            fontWeight: 700,
            color: activeTab === 'banners' ? '#ef4a23' : '#64748b',
            borderBottom: activeTab === 'banners' ? '3px solid #ef4a23' : 'none',
            background: 'none',
            display: 'flex',
            alignItems: 'center',
            gap: 8,
          }}
        >
          <ImageIcon size={16} /> Promo Banners ({banners.length})
        </button>
        <button
          onClick={() => setActiveTab('coupons')}
          style={{
            padding: '10px 18px',
            fontSize: 14,
            fontWeight: 700,
            color: activeTab === 'coupons' ? '#ef4a23' : '#64748b',
            borderBottom: activeTab === 'coupons' ? '3px solid #ef4a23' : 'none',
            background: 'none',
            display: 'flex',
            alignItems: 'center',
            gap: 8,
          }}
        >
          <Ticket size={16} /> Coupons & Discounts
        </button>
      </div>

      {/* Categories Content */}
      {activeTab === 'categories' && (
        <div className="admin-card">
          <div style={{ display: 'flex', justifyContent: 'space-between', padding: 16, borderBottom: '1px solid #e2e8f0' }}>
            <h3 style={{ fontSize: 16, fontWeight: 700 }}>Store Categories</h3>
            <button className="btn btn-primary btn-sm" onClick={() => handleOpenCatModal()}>
              <Plus size={14} /> Add Category
            </button>
          </div>
          <table className="admin-table">
            <thead>
              <tr>
                <th>Icon / Class</th>
                <th>Category Name</th>
                <th>Slug</th>
                <th>Sort Order</th>
                <th style={{ textAlign: 'right' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr><td colSpan={5} style={{ textAlign: 'center', padding: 32 }}>Loading categories...</td></tr>
              ) : categories.length > 0 ? (
                categories.map((cat) => (
                  <tr key={cat.id}>
                    <td>
                      <span style={{ padding: '4px 8px', background: '#f1f5f9', borderRadius: 4, fontFamily: 'monospace', fontSize: 12 }}>
                        {cat.icon || 'folder'}
                      </span>
                    </td>
                    <td><b>{cat.name}</b></td>
                    <td style={{ color: '#64748b', fontSize: 13 }}>{cat.slug}</td>
                    <td>{cat.sort_order}</td>
                    <td style={{ textAlign: 'right' }}>
                      <div style={{ display: 'flex', gap: 6, justifyContent: 'flex-end' }}>
                        <button className="btn btn-outline btn-sm" onClick={() => handleOpenCatModal(cat)}>
                          <Edit2 size={14} />
                        </button>
                        <button className="btn btn-outline btn-sm" style={{ color: '#ef4444' }} onClick={() => handleDeleteCat(cat.id)}>
                          <Trash2 size={14} />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr><td colSpan={5} style={{ textAlign: 'center', padding: 32 }}>No categories found</td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}

      {/* Brands Content */}
      {activeTab === 'brands' && (
        <div className="admin-card">
          <div style={{ display: 'flex', justifyContent: 'space-between', padding: 16, borderBottom: '1px solid #e2e8f0' }}>
            <h3 style={{ fontSize: 16, fontWeight: 700 }}>Manufacturer & Brands</h3>
            <button className="btn btn-primary btn-sm" onClick={() => handleOpenBrandModal()}>
              <Plus size={14} /> Add Brand
            </button>
          </div>
          <table className="admin-table">
            <thead>
              <tr>
                <th>Logo</th>
                <th>Brand Name</th>
                <th>Slug</th>
                <th style={{ textAlign: 'right' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr><td colSpan={4} style={{ textAlign: 'center', padding: 32 }}>Loading brands...</td></tr>
              ) : brands.length > 0 ? (
                brands.map((b) => (
                  <tr key={b.id}>
                    <td>
                      {b.logo ? (
                        <img src={b.logo} alt={b.name} style={{ height: 28, maxWidth: 80, objectFit: 'contain' }} />
                      ) : (
                        <span style={{ fontSize: 12, color: '#94a3b8' }}>No Logo</span>
                      )}
                    </td>
                    <td><b>{b.name}</b></td>
                    <td style={{ color: '#64748b', fontSize: 13 }}>{b.slug}</td>
                    <td style={{ textAlign: 'right' }}>
                      <div style={{ display: 'flex', gap: 6, justifyContent: 'flex-end' }}>
                        <button className="btn btn-outline btn-sm" onClick={() => handleOpenBrandModal(b)}>
                          <Edit2 size={14} />
                        </button>
                        <button className="btn btn-outline btn-sm" style={{ color: '#ef4444' }} onClick={() => handleDeleteBrand(b.id)}>
                          <Trash2 size={14} />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr><td colSpan={4} style={{ textAlign: 'center', padding: 32 }}>No brands found</td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}

      {/* Banners Content */}
      {activeTab === 'banners' && (
        <div className="admin-card">
          <div style={{ display: 'flex', justifyContent: 'space-between', padding: 16, borderBottom: '1px solid #e2e8f0' }}>
            <h3 style={{ fontSize: 16, fontWeight: 700 }}>Homepage Banners</h3>
            <button className="btn btn-primary btn-sm" onClick={() => setBannerModalOpen(true)}>
              <Plus size={14} /> Add Banner
            </button>
          </div>
          <table className="admin-table">
            <thead>
              <tr>
                <th>Preview</th>
                <th>Title / Subtitle</th>
                <th>Link Target</th>
                <th style={{ textAlign: 'right' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {banners.length > 0 ? (
                banners.map((ban) => (
                  <tr key={ban.id}>
                    <td>
                      <img src={ban.image_url} alt="" style={{ width: 120, height: 50, objectFit: 'cover', borderRadius: 4 }} />
                    </td>
                    <td>
                      <b>{ban.title}</b>
                      <small style={{ display: 'block', color: '#64748b' }}>{ban.subtitle}</small>
                    </td>
                    <td style={{ color: '#3749bb', fontSize: 13 }}>{ban.link_url || '/'}</td>
                    <td style={{ textAlign: 'right' }}>
                      <button className="btn btn-outline btn-sm" style={{ color: '#ef4444' }} onClick={() => handleDeleteBanner(ban.id)}>
                        <Trash2 size={14} />
                      </button>
                    </td>
                  </tr>
                ))
              ) : (
                <tr><td colSpan={4} style={{ textAlign: 'center', padding: 32 }}>No banners created yet</td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}

      {/* Coupons Content */}
      {activeTab === 'coupons' && (
        <div className="admin-card">
          <div style={{ display: 'flex', justifyContent: 'space-between', padding: 16, borderBottom: '1px solid #e2e8f0' }}>
            <h3 style={{ fontSize: 16, fontWeight: 700 }}>Discount Coupons</h3>
            <button className="btn btn-primary btn-sm" onClick={() => setCouponModalOpen(true)}>
              <Plus size={14} /> Create Coupon
            </button>
          </div>
          <div style={{ padding: 20 }}>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: 16 }}>
              <div style={{ border: '2px dashed #3749bb', borderRadius: 8, padding: 16, background: '#f8fafc' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <span className="badge badge-stock">ACTIVE</span>
                  <b style={{ color: '#ef4a23', fontSize: 18 }}>STARTECH10</b>
                </div>
                <p style={{ fontSize: 14, margin: '8px 0', fontWeight: 600 }}>10% OFF on all purchases</p>
                <small style={{ color: '#64748b', display: 'block' }}>Min Order: ৳2,000 | Max Disc: ৳1,000</small>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Category Modal */}
      {catModalOpen && (
        <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.5)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <div style={{ background: '#fff', padding: 24, borderRadius: 8, width: 420 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
              <h3>{editingCat ? 'Edit Category' : 'Add Category'}</h3>
              <button onClick={() => setCatModalOpen(false)}><X size={18} /></button>
            </div>
            <form onSubmit={handleSaveCat}>
              <div style={{ marginBottom: 12 }}>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 4 }}>Category Name *</label>
                <input type="text" required value={catForm.name} onChange={(e) => setCatForm({ ...catForm, name: e.target.value })} style={{ width: '100%', padding: 8, border: '1px solid #cbd5e1', borderRadius: 4 }} />
              </div>
              <div style={{ marginBottom: 12 }}>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 4 }}>Icon Identifier (e.g. laptop, monitor, cpu)</label>
                <input type="text" value={catForm.icon} onChange={(e) => setCatForm({ ...catForm, icon: e.target.value })} style={{ width: '100%', padding: 8, border: '1px solid #cbd5e1', borderRadius: 4 }} />
              </div>
              <div style={{ marginBottom: 16 }}>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 4 }}>Sort Order</label>
                <input type="number" value={catForm.sort_order} onChange={(e) => setCatForm({ ...catForm, sort_order: Number(e.target.value) })} style={{ width: '100%', padding: 8, border: '1px solid #cbd5e1', borderRadius: 4 }} />
              </div>
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
                <button type="button" className="btn btn-outline btn-sm" onClick={() => setCatModalOpen(false)}>Cancel</button>
                <button type="submit" className="btn btn-primary btn-sm">Save Category</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Brand Modal */}
      {brandModalOpen && (
        <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.5)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <div style={{ background: '#fff', padding: 24, borderRadius: 8, width: 420 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
              <h3>{editingBrand ? 'Edit Brand' : 'Add Brand'}</h3>
              <button onClick={() => setBrandModalOpen(false)}><X size={18} /></button>
            </div>
            <form onSubmit={handleSaveBrand}>
              <div style={{ marginBottom: 12 }}>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 4 }}>Brand Name *</label>
                <input type="text" required value={brandForm.name} onChange={(e) => setBrandForm({ ...brandForm, name: e.target.value })} style={{ width: '100%', padding: 8, border: '1px solid #cbd5e1', borderRadius: 4 }} />
              </div>
              <div style={{ marginBottom: 16 }}>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 4 }}>Logo Image URL</label>
                <input type="text" placeholder="https://example.com/logo.png" value={brandForm.logo} onChange={(e) => setBrandForm({ ...brandForm, logo: e.target.value })} style={{ width: '100%', padding: 8, border: '1px solid #cbd5e1', borderRadius: 4 }} />
              </div>
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
                <button type="button" className="btn btn-outline btn-sm" onClick={() => setBrandModalOpen(false)}>Cancel</button>
                <button type="submit" className="btn btn-primary btn-sm">Save Brand</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Banner Modal */}
      {bannerModalOpen && (
        <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.5)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <div style={{ background: '#fff', padding: 24, borderRadius: 8, width: 480 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
              <h3>Add Homepage Promo Banner</h3>
              <button onClick={() => setBannerModalOpen(false)}><X size={18} /></button>
            </div>
            <form onSubmit={handleSaveBanner}>
              <div style={{ marginBottom: 12 }}>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 4 }}>Banner Title</label>
                <input type="text" value={bannerForm.title} onChange={(e) => setBannerForm({ ...bannerForm, title: e.target.value })} style={{ width: '100%', padding: 8, border: '1px solid #cbd5e1', borderRadius: 4 }} />
              </div>
              <div style={{ marginBottom: 12 }}>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 4 }}>Subtitle</label>
                <input type="text" value={bannerForm.subtitle} onChange={(e) => setBannerForm({ ...bannerForm, subtitle: e.target.value })} style={{ width: '100%', padding: 8, border: '1px solid #cbd5e1', borderRadius: 4 }} />
              </div>
              <div style={{ marginBottom: 12 }}>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 4 }}>Banner Image URL *</label>
                <input type="text" required value={bannerForm.image_url} onChange={(e) => setBannerForm({ ...bannerForm, image_url: e.target.value })} style={{ width: '100%', padding: 8, border: '1px solid #cbd5e1', borderRadius: 4 }} />
              </div>
              <div style={{ marginBottom: 16 }}>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 4 }}>Click Target URL</label>
                <input type="text" value={bannerForm.link_url} onChange={(e) => setBannerForm({ ...bannerForm, link_url: e.target.value })} style={{ width: '100%', padding: 8, border: '1px solid #cbd5e1', borderRadius: 4 }} />
              </div>
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
                <button type="button" className="btn btn-outline btn-sm" onClick={() => setBannerModalOpen(false)}>Cancel</button>
                <button type="submit" className="btn btn-primary btn-sm">Save Banner</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Coupon Modal */}
      {couponModalOpen && (
        <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.5)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <div style={{ background: '#fff', padding: 24, borderRadius: 8, width: 440 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
              <h3>Create New Coupon</h3>
              <button onClick={() => setCouponModalOpen(false)}><X size={18} /></button>
            </div>
            <form onSubmit={handleSaveCoupon}>
              <div style={{ marginBottom: 12 }}>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 4 }}>Coupon Code *</label>
                <input type="text" required placeholder="e.g. FESTIVE20" value={couponForm.code} onChange={(e) => setCouponForm({ ...couponForm, code: e.target.value })} style={{ width: '100%', padding: 8, border: '1px solid #cbd5e1', borderRadius: 4, textTransform: 'uppercase' }} />
              </div>
              <div style={{ marginBottom: 12, display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8 }}>
                <div>
                  <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 4 }}>Type</label>
                  <select value={couponForm.discount_type} onChange={(e) => setCouponForm({ ...couponForm, discount_type: e.target.value })} style={{ width: '100%', padding: 8, border: '1px solid #cbd5e1', borderRadius: 4 }}>
                    <option value="PERCENTAGE">Percentage (%)</option>
                    <option value="FIXED">Fixed Amount (৳)</option>
                  </select>
                </div>
                <div>
                  <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 4 }}>Discount Value</label>
                  <input type="number" required value={couponForm.discount_value} onChange={(e) => setCouponForm({ ...couponForm, discount_value: e.target.value })} style={{ width: '100%', padding: 8, border: '1px solid #cbd5e1', borderRadius: 4 }} />
                </div>
              </div>
              <div style={{ marginBottom: 12, display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8 }}>
                <div>
                  <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 4 }}>Min Order (৳)</label>
                  <input type="number" value={couponForm.min_order_amount} onChange={(e) => setCouponForm({ ...couponForm, min_order_amount: e.target.value })} style={{ width: '100%', padding: 8, border: '1px solid #cbd5e1', borderRadius: 4 }} />
                </div>
                <div>
                  <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 4 }}>Max Disc (৳)</label>
                  <input type="number" value={couponForm.max_discount} onChange={(e) => setCouponForm({ ...couponForm, max_discount: e.target.value })} style={{ width: '100%', padding: 8, border: '1px solid #cbd5e1', borderRadius: 4 }} />
                </div>
              </div>
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8, marginTop: 16 }}>
                <button type="button" className="btn btn-outline btn-sm" onClick={() => setCouponModalOpen(false)}>Cancel</button>
                <button type="submit" className="btn btn-primary btn-sm">Create Coupon</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
