import React from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Heart, Scale, ShoppingCart, Star } from 'lucide-react'
import { useShop } from '../../context/ShopContext'

export default function ProductCard({ product }) {
  const { addToCart, toggleWishlist, isInWishlist, addToCompare } = useShop()
  const navigate = useNavigate()

  if (!product) return null

  const isDiscounted = product.discount_price && product.discount_price < product.price
  const discountPercent = isDiscounted
    ? Math.round(((product.price - product.discount_price) / product.price) * 100)
    : 0

  const displayPrice = product.discount_price || product.price
  const inWishlist = isInWishlist(product.id)

  // Parse specs if string/json
  let specsObj = {}
  try {
    specsObj = typeof product.specs === 'string' ? JSON.parse(product.specs) : product.specs || {}
  } catch {
    specsObj = {}
  }
  const specEntries = Object.entries(specsObj).slice(0, 3)

  const handleBuyNow = (e) => {
    e.preventDefault()
    addToCart(product, 1)
    navigate('/checkout')
  }

  const handleAddToCart = (e) => {
    e.preventDefault()
    addToCart(product, 1)
  }

  const mainImage = product.images?.[0]?.url || 'https://placehold.co/400x400/1a1a2e/ef4a23?text=' + encodeURIComponent(product.name)

  return (
    <div className="product-card">
      {/* Badges */}
      <div className="product-badges">
        {isDiscounted && (
          <span className="badge badge-discount">-{discountPercent}% OFF</span>
        )}
        {product.stock > 0 ? (
          <span className="badge badge-stock">In Stock</span>
        ) : (
          <span className="badge badge-out">Out of Stock</span>
        )}
      </div>

      {/* Quick Action Icons */}
      <div className="product-quick-actions">
        <button
          className={`icon-action-btn ${inWishlist ? 'active' : ''}`}
          onClick={(e) => {
            e.preventDefault()
            toggleWishlist(product)
          }}
          title={inWishlist ? 'Remove from Wishlist' : 'Add to Wishlist'}
        >
          <Heart size={16} fill={inWishlist ? '#ef4a23' : 'none'} />
        </button>
        <button
          className="icon-action-btn"
          onClick={(e) => {
            e.preventDefault()
            addToCompare(product)
          }}
          title="Add to Compare"
        >
          <Scale size={16} />
        </button>
      </div>

      {/* Image */}
      <Link to={`/product/${product.slug}`} className="product-image-container">
        <img src={mainImage} alt={product.name} loading="lazy" />
      </Link>

      {/* Details */}
      <div className="product-details">
        <div className="product-brand-cat">
          {product.brand?.name || 'MI-Tech'} · {product.category?.name || 'Tech'}
        </div>

        <Link to={`/product/${product.slug}`}>
          <h3 className="product-title" title={product.name}>
            {product.name}
          </h3>
        </Link>

        {/* Specs snippet */}
        <ul className="product-specs-list">
          {specEntries.length > 0 ? (
            specEntries.map(([k, v]) => (
              <li key={k}>
                <b>{k}:</b> {v}
              </li>
            ))
          ) : (
            <li>{product.short_description || 'Official warranty guaranteed.'}</li>
          )}
        </ul>

        {/* Pricing */}
        <div className="product-pricing">
          <span className="current-price">৳{displayPrice.toLocaleString()}</span>
          {isDiscounted && (
            <span className="old-price">৳{product.price.toLocaleString()}</span>
          )}
        </div>

        {/* Actions */}
        <div className="product-card-actions">
          <button
            className="buy-now-btn"
            onClick={handleBuyNow}
            disabled={product.stock <= 0}
          >
            <ShoppingCart size={15} /> Buy Now
          </button>
          <button
            className="cart-add-btn"
            onClick={handleAddToCart}
            disabled={product.stock <= 0}
            title="Add to Cart"
          >
            <ShoppingCart size={16} />
          </button>
        </div>
      </div>
    </div>
  )
}
