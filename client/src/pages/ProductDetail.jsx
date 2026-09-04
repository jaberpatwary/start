import React, { useState, useEffect } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { Heart, Scale, ShoppingCart, Truck, ShieldCheck, Star, ChevronRight, Plus, Minus, CheckCircle, MessageSquare } from 'lucide-react'
import api from '../api/client'
import { useShop } from '../context/ShopContext'
import { useAuth } from '../context/AuthContext'
import ProductCard from '../components/product/ProductCard'
import toast from 'react-hot-toast'

export default function ProductDetail() {
  const { slug } = useParams()
  const navigate = useNavigate()
  const { addToCart, toggleWishlist, isInWishlist, addToCompare } = useShop()
  const { isAuthenticated, user } = useAuth()

  const [product, setProduct] = useState(null)
  const [relatedProducts, setRelatedProducts] = useState([])
  const [selectedImage, setSelectedImage] = useState(0)
  const [quantity, setQuantity] = useState(1)
  const [activeTab, setActiveTab] = useState('specs') // 'specs' | 'desc' | 'reviews'
  const [loading, setLoading] = useState(true)

  // Review submission state
  const [reviewRating, setReviewRating] = useState(5)
  const [reviewComment, setReviewComment] = useState('')
  const [submittingReview, setSubmittingReview] = useState(false)

  useEffect(() => {
    const fetchProduct = async () => {
      setLoading(true)
      try {
        const data = await api.get(`/products/slug/${slug}`)
        if (data?.product) {
          setProduct(data.product)
          setSelectedImage(0)
          setQuantity(1)

          // Fetch related
          if (data.product.category_id) {
            const relData = await api.get(`/products?category_id=${data.product.category_id}&limit=4`)
            setRelatedProducts((relData?.results || []).filter((p) => p.id !== data.product.id))
          }
        }
      } catch (err) {
        console.error('Failed to load product:', err)
        toast.error('Product not found')
      } finally {
        setLoading(false)
      }
    }
    fetchProduct()
  }, [slug])

  if (loading) {
    return (
      <div className="container" style={{ padding: '80px 0', textAlign: 'center', color: '#64748b' }}>
        <h2>Loading product details...</h2>
      </div>
    )
  }

  if (!product) {
    return (
      <div className="container" style={{ padding: '80px 0', textAlign: 'center' }}>
        <h2>Product not found</h2>
        <Link to="/" className="btn btn-primary" style={{ marginTop: 16 }}>
          Return to Home
        </Link>
      </div>
    )
  }

  const isDiscounted = product.discount_price && product.discount_price < product.price
  const displayPrice = product.discount_price || product.price
  const inWishlist = isInWishlist(product.id)

  const images = product.images && product.images.length > 0
    ? product.images.map((img) => img.url)
    : ['https://placehold.co/600x600/1a1a2e/ef4a23?text=' + encodeURIComponent(product.name)]

  // Parse specs JSON
  let specsObj = {}
  try {
    specsObj = typeof product.specs === 'string' ? JSON.parse(product.specs) : product.specs || {}
  } catch {
    specsObj = {}
  }

  const handleBuyNow = () => {
    addToCart(product, quantity)
    navigate('/checkout')
  }

  const handleAddToCart = () => {
    addToCart(product, quantity)
  }

  const handleReviewSubmit = async (e) => {
    e.preventDefault()
    if (!isAuthenticated) {
      toast.error('Please login to write a review')
      navigate('/login')
      return
    }
    if (!reviewComment.trim()) {
      toast.error('Please enter your review comments')
      return
    }

    try {
      setSubmittingReview(true)
      await api.post(`/products/${product.id}/reviews`, {
        rating: reviewRating,
        comment: reviewComment,
      })
      toast.success('Thank you! Your review has been submitted for moderation.')
      setReviewComment('')
    } catch (err) {
      toast.error(err.message || 'Failed to submit review')
    } finally {
      setSubmittingReview(false)
    }
  }

  return (
    <div className="container" style={{ padding: '20px 16px' }}>
      {/* Breadcrumbs */}
      <div className="breadcrumbs">
        <Link to="/">Home</Link> <ChevronRight size={12} style={{ display: 'inline' }} />{' '}
        <Link to={`/category/${product.category?.slug || 'desktop'}`}>{product.category?.name || 'Category'}</Link>{' '}
        <ChevronRight size={12} style={{ display: 'inline' }} /> <span>{product.name}</span>
      </div>

      {/* Product Detail Main Block */}
      <div className="detail-layout">
        {/* Left: Gallery */}
        <div>
          <div className="gallery-main">
            <img src={images[selectedImage]} alt={product.name} />
          </div>
          {images.length > 1 && (
            <div className="gallery-thumbs">
              {images.map((img, idx) => (
                <button
                  key={idx}
                  className={`thumb-btn ${selectedImage === idx ? 'active' : ''}`}
                  onClick={() => setSelectedImage(idx)}
                >
                  <img src={img} alt={`${product.name} thumb ${idx}`} />
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Right: Info */}
        <div className="detail-info">
          <h1>{product.name}</h1>

          <div className="detail-meta-row">
            <span>Brand: <b style={{ color: '#3749bb' }}>{product.brand?.name || 'StarTech'}</b></span>
            <span>SKU: <b>{product.sku}</b></span>
            <span>
              Status:{' '}
              <b style={{ color: product.stock > 0 ? '#10b981' : '#ef4444' }}>
                {product.stock > 0 ? `In Stock (${product.stock})` : 'Out of Stock'}
              </b>
            </span>
          </div>

          {/* Pricing */}
          <div className="detail-price-box">
            <span className="price">৳{displayPrice.toLocaleString()}</span>
            {isDiscounted && (
              <>
                <span className="old">৳{product.price.toLocaleString()}</span>
                <span className="badge badge-discount">
                  Save ৳{(product.price - product.discount_price).toLocaleString()}
                </span>
              </>
            )}
          </div>

          {/* Key Features */}
          <div className="key-features">
            <h4>Key Features:</h4>
            <ul>
              {Object.entries(specsObj).slice(0, 5).map(([k, v]) => (
                <li key={k}>
                  <b>{k}:</b> {v}
                </li>
              ))}
              <li><b>Warranty:</b> 1 Year Official Warranty</li>
            </ul>
          </div>

          {/* Quantity & Buy Actions */}
          <div className="detail-actions-row">
            <div className="qty-control">
              <button
                onClick={() => setQuantity(Math.max(1, quantity - 1))}
                disabled={quantity <= 1}
              >
                <Minus size={14} />
              </button>
              <span>{quantity}</span>
              <button
                onClick={() => setQuantity(Math.min(product.stock || 10, quantity + 1))}
                disabled={quantity >= (product.stock || 10)}
              >
                <Plus size={14} />
              </button>
            </div>

            <button
              className="btn btn-primary btn-lg"
              onClick={handleBuyNow}
              disabled={product.stock <= 0}
            >
              <ShoppingCart size={18} /> Buy Now
            </button>

            <button
              className="btn btn-outline btn-lg"
              onClick={handleAddToCart}
              disabled={product.stock <= 0}
            >
              Add to Cart
            </button>

            <button
              className={`icon-action-btn ${inWishlist ? 'active' : ''}`}
              onClick={() => toggleWishlist(product)}
              style={{ width: 44, height: 44 }}
              title="Add to Wishlist"
            >
              <Heart size={20} fill={inWishlist ? '#ef4a23' : 'none'} />
            </button>

            <button
              className="icon-action-btn"
              onClick={() => addToCompare(product)}
              style={{ width: 44, height: 44 }}
              title="Add to Compare"
            >
              <Scale size={20} />
            </button>
          </div>

          {/* Trust badges */}
          <div style={{
            display: 'grid',
            gridTemplateColumns: '1fr 1fr',
            gap: 12,
            padding: 16,
            background: '#f8fafc',
            borderRadius: 8,
            border: '1px solid #e2e8f0',
            fontSize: 13,
          }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 8, color: '#334155' }}>
              <Truck size={18} color="#3749bb" />
              <span><b>Fast Delivery:</b> 2-3 Days nationwide</span>
            </div>
            <div style={{ display: 'flex', alignItems: 'center', gap: 8, color: '#334155' }}>
              <ShieldCheck size={18} color="#10b981" />
              <span><b>100% Genuine:</b> Official distributor warranty</span>
            </div>
          </div>
        </div>
      </div>

      {/* Tabs: Specification, Description, Reviews */}
      <section className="detail-tabs-section">
        <div className="tabs-nav">
          <button
            className={`tab-btn ${activeTab === 'specs' ? 'active' : ''}`}
            onClick={() => setActiveTab('specs')}
          >
            Full Specification
          </button>
          <button
            className={`tab-btn ${activeTab === 'desc' ? 'active' : ''}`}
            onClick={() => setActiveTab('desc')}
          >
            Description
          </button>
          <button
            className={`tab-btn ${activeTab === 'reviews' ? 'active' : ''}`}
            onClick={() => setActiveTab('reviews')}
          >
            Reviews ({product.reviews?.length || 0})
          </button>
        </div>

        <div className="tab-content">
          {/* Specs Table */}
          {activeTab === 'specs' && (
            <div>
              <h3 style={{ fontSize: 18, marginBottom: 16 }}>Technical Specifications</h3>
              <table className="specs-table">
                <tbody>
                  <tr>
                    <th>Product Model</th>
                    <td>{product.name}</td>
                  </tr>
                  <tr>
                    <th>SKU / Part Number</th>
                    <td>{product.sku}</td>
                  </tr>
                  <tr>
                    <th>Brand</th>
                    <td>{product.brand?.name || 'StarTech'}</td>
                  </tr>
                  <tr>
                    <th>Category</th>
                    <td>{product.category?.name || 'Electronics'}</td>
                  </tr>
                  {Object.entries(specsObj).map(([k, v]) => (
                    <tr key={k}>
                      <th>{k}</th>
                      <td>{v}</td>
                    </tr>
                  ))}
                  <tr>
                    <th>Warranty</th>
                    <td>1 Year Official Brand Warranty</td>
                  </tr>
                </tbody>
              </table>
            </div>
          )}

          {/* Description */}
          {activeTab === 'desc' && (
            <div style={{ lineHeight: 1.8, color: '#334155', fontSize: 15 }}>
              <h3 style={{ fontSize: 18, marginBottom: 12 }}>Product Description</h3>
              <p style={{ marginBottom: 16 }}>
                {product.description || product.short_description || 'High quality tech product backed by official brand warranty.'}
              </p>
              <p>
                Get the authentic <b>{product.name}</b> at the best price in Bangladesh from Star Tech Clone.
                Order online or visit your nearest Star Tech shop outlet.
              </p>
            </div>
          )}

          {/* Reviews */}
          {activeTab === 'reviews' && (
            <div>
              <div className="reviews-summary">
                <div className="reviews-score">
                  <b>{product.rating_avg ? product.rating_avg.toFixed(1) : '5.0'}</b>
                  <div style={{ color: '#f59e0b', fontSize: 18 }}>★★★★★</div>
                  <small style={{ color: '#64748b' }}>Based on {product.reviews?.length || 0} review(s)</small>
                </div>
              </div>

              {/* Existing approved reviews */}
              {product.reviews && product.reviews.length > 0 ? (
                product.reviews.map((r) => (
                  <div key={r.id} className="review-item">
                    <div className="review-header">
                      <div>
                        <b style={{ fontSize: 14 }}>{r.user?.name || 'Verified Buyer'}</b>
                        <span style={{ fontSize: 12, color: '#94a3b8', marginLeft: 8 }}>
                          {new Date(r.created_at).toLocaleDateString()}
                        </span>
                      </div>
                      <div style={{ color: '#f59e0b', fontSize: 14 }}>
                        {'★'.repeat(r.rating)}{'☆'.repeat(5 - r.rating)}
                      </div>
                    </div>
                    <p style={{ fontSize: 14, color: '#475569' }}>{r.comment}</p>
                  </div>
                ))
              ) : (
                <p style={{ color: '#64748b', marginBottom: 24 }}>No reviews yet. Be the first to review this product!</p>
              )}

              {/* Review Submission Form */}
              <div style={{
                background: '#f8fafc',
                borderRadius: 8,
                padding: 24,
                marginTop: 24,
                border: '1px solid #e2e8f0',
              }}>
                <h4 style={{ fontSize: 16, marginBottom: 16, display: 'flex', alignItems: 'center', gap: 8 }}>
                  <MessageSquare size={18} color="#ef4a23" /> Write a Review
                </h4>
                <form onSubmit={handleReviewSubmit}>
                  <div style={{ marginBottom: 14 }}>
                    <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 6 }}>
                      Your Rating:
                    </label>
                    <select
                      className="sort-select"
                      value={reviewRating}
                      onChange={(e) => setReviewRating(Number(e.target.value))}
                    >
                      <option value={5}>★★★★★ (5 - Excellent)</option>
                      <option value={4}>★★★★☆ (4 - Good)</option>
                      <option value={3}>★★★☆☆ (3 - Average)</option>
                      <option value={2}>★★☆☆☆ (2 - Poor)</option>
                      <option value={1}>★☆☆☆☆ (1 - Terrible)</option>
                    </select>
                  </div>

                  <div style={{ marginBottom: 14 }}>
                    <label style={{ display: 'block', fontSize: 13, fontWeight: 600, marginBottom: 6 }}>
                      Your Review Feedback:
                    </label>
                    <textarea
                      rows={4}
                      style={{
                        width: '100%',
                        padding: 12,
                        borderRadius: 4,
                        border: '1px solid #cbd5e1',
                        fontSize: 14,
                        outline: 'none',
                      }}
                      placeholder="Share your experience with this product..."
                      value={reviewComment}
                      onChange={(e) => setReviewComment(e.target.value)}
                    />
                  </div>

                  <button
                    type="submit"
                    className="btn btn-primary"
                    disabled={submittingReview}
                  >
                    {submittingReview ? 'Submitting...' : 'Submit Review'}
                  </button>
                </form>
              </div>
            </div>
          )}
        </div>
      </section>

      {/* Related Products */}
      {relatedProducts.length > 0 && (
        <section>
          <div className="section-header">
            <h2 className="section-title">Related Products</h2>
          </div>
          <div className="products-grid">
            {relatedProducts.map((p) => (
              <ProductCard key={p.id} product={p} />
            ))}
          </div>
        </section>
      )}
    </div>
  )
}
