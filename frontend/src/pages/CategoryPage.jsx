import React, { useState, useEffect } from 'react'
import { useParams, Link, useSearchParams } from 'react-router-dom'
import { Filter, SlidersHorizontal, ChevronRight } from 'lucide-react'
import api from '../api/client'
import ProductCard from '../components/product/ProductCard'

export default function CategoryPage({ isSearch = false }) {
  const { slug } = useParams()
  const [searchParams] = useSearchParams()
  const queryParam = searchParams.get('q') || ''

  const [products, setProducts] = useState([])
  const [brands, setBrands] = useState([])
  const [totalResults, setTotalResults] = useState(0)
  const [totalPages, setTotalPages] = useState(1)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)

  // Filters state
  const [selectedBrands, setSelectedBrands] = useState([])
  const [priceRange, setPriceRange] = useState({ min: '', max: '' })
  const [inStockOnly, setInStockOnly] = useState(false)
  const [sortBy, setSortBy] = useState('newest')

  // Fetch available brands
  useEffect(() => {
    api.get('/brands').then((res) => setBrands(res?.brands || [])).catch(() => {})
  }, [])

  // Reset page when category or search changes
  useEffect(() => {
    setPage(1)
    setSelectedBrands([])
    setPriceRange({ min: '', max: '' })
    setInStockOnly(false)
  }, [slug, queryParam])

  // Fetch products
  useEffect(() => {
    const fetchFilteredProducts = async () => {
      setLoading(true)
      try {
        let endpoint = `/products?page=${page}&limit=12`

        if (isSearch && queryParam) {
          endpoint += `&search=${encodeURIComponent(queryParam)}`
        } else if (slug) {
          endpoint += `&category=${encodeURIComponent(slug)}`
        }

        if (selectedBrands.length === 1) {
          endpoint += `&brand_id=${selectedBrands[0]}`
        }
        if (priceRange.min) {
          endpoint += `&min_price=${priceRange.min}`
        }
        if (priceRange.max) {
          endpoint += `&max_price=${priceRange.max}`
        }
        if (inStockOnly) {
          endpoint += `&in_stock=true`
        }

        if (sortBy === 'price_asc') endpoint += '&sort_by=price_asc'
        else if (sortBy === 'price_desc') endpoint += '&sort_by=price_desc'
        else if (sortBy === 'popular') endpoint += '&sort_by=popular'
        else if (sortBy === 'rating') endpoint += '&sort_by=rating'
        // default: newest (no sort_by param, backend defaults to created_at DESC)

        const data = await api.get(endpoint)
        setProducts(data?.results || [])
        setTotalResults(data?.total_results || 0)
        setTotalPages(data?.total_pages || 1)
      } catch (err) {
        console.error('Filter error:', err)
        setProducts([])
      } finally {
        setLoading(false)
      }
    }
    fetchFilteredProducts()
  }, [slug, queryParam, isSearch, page, selectedBrands, inStockOnly, sortBy, priceRange.min, priceRange.max])

  const toggleBrand = (brandId) => {
    setSelectedBrands((prev) =>
      prev.includes(brandId) ? prev.filter((id) => id !== brandId) : [...prev, brandId]
    )
    setPage(1)
  }

  const categoryTitle = isSearch
    ? `Search Results for "${queryParam}"`
    : slug
    ? slug.charAt(0).toUpperCase() + slug.slice(1)
    : 'All Products'

  return (
    <div className="container" style={{ padding: '20px 16px' }}>
      {/* Breadcrumbs */}
      <div className="breadcrumbs">
        <Link to="/">Home</Link> <ChevronRight size={12} style={{ display: 'inline' }} />{' '}
        {isSearch ? (
          <span>Search</span>
        ) : (
          <span>{categoryTitle}</span>
        )}
      </div>

      <div className="listing-layout">
        {/* Filters Sidebar */}
        <aside className="filters-sidebar">
          <h3>
            <span><SlidersHorizontal size={18} style={{ display: 'inline', marginRight: 6 }} /> Filters</span>
            {(selectedBrands.length > 0 || priceRange.min || priceRange.max || inStockOnly) && (
              <button
                onClick={() => {
                  setSelectedBrands([])
                  setPriceRange({ min: '', max: '' })
                  setInStockOnly(false)
                }}
                style={{ fontSize: 12, color: '#ef4a23', fontWeight: 600 }}
              >
                Reset
              </button>
            )}
          </h3>

          {/* Price Range Filter */}
          <div className="filter-group">
            <h4 className="filter-title">Price Range (৳)</h4>
            <div className="price-inputs">
              <input
                type="number"
                placeholder="Min ৳"
                value={priceRange.min}
                onChange={(e) => setPriceRange({ ...priceRange, min: e.target.value })}
              />
              <input
                type="number"
                placeholder="Max ৳"
                value={priceRange.max}
                onChange={(e) => setPriceRange({ ...priceRange, max: e.target.value })}
              />
            </div>
          </div>

          {/* Availability */}
          <div className="filter-group">
            <h4 className="filter-title">Availability</h4>
            <label className="filter-label">
              <input
                type="checkbox"
                checked={inStockOnly}
                onChange={(e) => setInStockOnly(e.target.checked)}
              />
              In Stock Only
            </label>
          </div>

          {/* Brand Filter */}
          <div className="filter-group">
            <h4 className="filter-title">Brands</h4>
            <div className="filter-checkbox-list">
              {brands.map((b) => (
                <label key={b.id} className="filter-label">
                  <input
                    type="checkbox"
                    checked={selectedBrands.includes(b.id)}
                    onChange={() => toggleBrand(b.id)}
                  />
                  {b.name}
                </label>
              ))}
            </div>
          </div>
        </aside>

        {/* Main Products Content */}
        <main>
          {/* Header */}
          <div className="listing-header">
            <div>
              <h1 style={{ fontSize: 20, fontWeight: 800, color: '#0f172a' }}>{categoryTitle}</h1>
              <small style={{ color: '#64748b' }}>Showing {products.length} of {totalResults} items</small>
            </div>

            <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
              <span style={{ fontSize: 13, color: '#64748b' }}>Sort By:</span>
              <select
                className="sort-select"
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value)}
              >
                <option value="newest">Newest First</option>
                <option value="price_asc">Price: Low to High</option>
                <option value="price_desc">Price: High to Low</option>
                <option value="popular">Most Popular</option>
                <option value="rating">Top Rated</option>
              </select>
            </div>
          </div>

          {/* Products Grid */}
          {loading ? (
            <div style={{ textAlign: 'center', padding: '60px 0', color: '#64748b' }}>
              Loading products...
            </div>
          ) : products.length > 0 ? (
            <div className="products-grid" style={{ gridTemplateColumns: 'repeat(3, 1fr)' }}>
              {products.map((p) => (
                <ProductCard key={p.id} product={p} />
              ))}
            </div>
          ) : (
            <div style={{
              background: '#fff',
              borderRadius: 8,
              padding: '60px 20px',
              textAlign: 'center',
              border: '1px solid #e2e8f0',
            }}>
              <h3 style={{ fontSize: 18, color: '#334155', marginBottom: 8 }}>No products found</h3>
              <p style={{ fontSize: 14, color: '#64748b', marginBottom: 16 }}>
                Try adjusting your filters or searching for different keywords.
              </p>
              <button
                className="btn btn-primary"
                onClick={() => {
                  setSelectedBrands([])
                  setPriceRange({ min: '', max: '' })
                  setInStockOnly(false)
                }}
              >
                Clear All Filters
              </button>
            </div>
          )}

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="pagination">
              {Array.from({ length: totalPages }, (_, i) => i + 1).map((pNum) => (
                <button
                  key={pNum}
                  className={`page-btn ${page === pNum ? 'active' : ''}`}
                  onClick={() => {
                    setPage(pNum)
                    window.scrollTo({ top: 0, behavior: 'smooth' })
                  }}
                >
                  {pNum}
                </button>
              ))}
            </div>
          )}
        </main>
      </div>
    </div>
  )
}
