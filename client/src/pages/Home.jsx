import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Monitor, Laptop, Cpu, Tv, ChevronRight, ShieldCheck, Truck, Clock, RefreshCw, Flame, ArrowRight } from 'lucide-react'
import api from '../api/client'
import ProductCard from '../components/product/ProductCard'

export default function Home() {
  const [banners, setBanners] = useState([])
  const [featuredProducts, setFeaturedProducts] = useState([])
  const [newArrivals, setNewArrivals] = useState([])
  const [currentSlide, setCurrentSlide] = useState(0)
  const [dealTime, setDealTime] = useState({ hours: 14, minutes: 32, seconds: 45 })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [bannerRes, featRes, newRes] = await Promise.all([
          api.get('/banners').catch(() => ({ banners: [] })),
          api.get('/products?featured=true&limit=8').catch(() => ({ results: [] })),
          api.get('/products?limit=8').catch(() => ({ results: [] })),
        ])
        setBanners(bannerRes?.banners || [])
        setFeaturedProducts(featRes?.results || [])
        setNewArrivals(newRes?.results || [])
      } catch (err) {
        console.error('Home data load error:', err)
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [])

  // Auto banner slider
  useEffect(() => {
    if (banners.length <= 1) return
    const interval = setInterval(() => {
      setCurrentSlide((prev) => (prev + 1) % banners.length)
    }, 5000)
    return () => clearInterval(interval)
  }, [banners.length])

  // Countdown timer
  useEffect(() => {
    const timer = setInterval(() => {
      setDealTime((prev) => {
        if (prev.seconds > 0) return { ...prev, seconds: prev.seconds - 1 }
        if (prev.minutes > 0) return { ...prev, minutes: 59, seconds: 59 }
        if (prev.hours > 0) return { ...prev, hours: prev.hours - 1, minutes: 59, seconds: 59 }
        return { hours: 23, minutes: 59, seconds: 59 }
      })
    }, 1000)
    return () => clearInterval(timer)
  }, [])

  const defaultBanners = [
    {
      id: 'b1',
      title: 'Next-Gen Gaming Laptops',
      subtitle: 'Up to ৳15,000 Off on ROG & MSI Gaming Laptops with RTX 40 Series',
      image: 'https://placehold.co/1200x500/0f172a/ef4a23?text=Gaming+Laptops+Fest+2026',
      link: '/category/laptop',
    },
    {
      id: 'b2',
      title: 'Build Your Ultimate PC',
      subtitle: 'Intel 14th Gen & AMD Ryzen 7000 Series Processors with 3 Years Warranty',
      image: 'https://placehold.co/1200x500/1e1b4b/ef4a23?text=PC+Components+Special',
      link: '/category/component',
    },
    {
      id: 'b3',
      title: '4K & OLED Gaming Displays',
      subtitle: 'Samsung Odyssey & ASUS ROG Swift high refresh rate monitors in stock',
      image: 'https://placehold.co/1200x500/052e16/ef4a23?text=Monitors+Mega+Deal',
      link: '/category/monitor',
    },
  ]

  const activeBanners = banners.length > 0 ? banners : defaultBanners

  const categories = [
    { name: 'Desktop', slug: 'desktop', count: 'Prebuilt & Custom PCs', icon: Monitor },
    { name: 'Laptop', slug: 'laptop', count: 'Gaming & Ultrabooks', icon: Laptop },
    { name: 'Component', slug: 'component', count: 'CPU, GPU, RAM, SSD', icon: Cpu },
    { name: 'Monitor', slug: 'monitor', count: '4K, Curved & Esports', icon: Tv },
  ]

  const brands = ['ASUS', 'MSI', 'Gigabyte', 'Intel', 'AMD', 'Corsair', 'Samsung', 'Dell']

  return (
    <div>
      {/* Hero Section */}
      <section className="hero-section">
        <div className="container">
          <div className="hero-grid">
            {/* Main Slider */}
            <div className="hero-slider">
              {activeBanners.map((b, idx) => (
                <div
                  key={b.id || idx}
                  className="slider-slide"
                  style={{ display: idx === currentSlide ? 'block' : 'none' }}
                >
                  <img src={b.image} alt={b.title} />
                  <div className="slider-content">
                    <h2>{b.title}</h2>
                    <p>{b.subtitle}</p>
                    <Link to={b.link || '/category/laptop'} className="btn btn-primary">
                      Shop Now <ChevronRight size={16} />
                    </Link>
                  </div>
                </div>
              ))}
              <div className="slider-dots">
                {activeBanners.map((_, idx) => (
                  <button
                    key={idx}
                    className={`slider-dot ${idx === currentSlide ? 'active' : ''}`}
                    onClick={() => setCurrentSlide(idx)}
                  />
                ))}
              </div>
            </div>

            {/* Side Banners */}
            <div className="hero-banners-side">
              <div className="side-banner">
                <span className="badge badge-discount" style={{ width: 'fit-content', marginBottom: 8 }}>
                  ⚡ Hot Deal
                </span>
                <h3>Custom PC Builder</h3>
                <p>Select components with real-time compatibility checks.</p>
                <Link to="/category/component" className="btn btn-secondary btn-sm" style={{ width: 'fit-content' }}>
                  Start Building <ArrowRight size={14} />
                </Link>
              </div>
              <div className="side-banner" style={{ borderLeftColor: '#3749bb' }}>
                <span className="badge" style={{ background: '#3749bb', width: 'fit-content', marginBottom: 8 }}>
                  🌟 Official Warranty
                </span>
                <h3>Instant EMI at 0%</h3>
                <p>Available on up to 36 months for 20+ leading banks.</p>
                <Link to="/emi-info" className="btn btn-outline btn-sm" style={{ width: 'fit-content' }}>
                  View EMI Plans <ChevronRight size={14} />
                </Link>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Value Props */}
      <section>
        <div className="container">
          <div className="value-props">
            <div className="prop-card">
              <div className="prop-icon"><ShieldCheck size={22} /></div>
              <div className="prop-text">
                <h4>100% Genuine</h4>
                <p>Authentic products with official warranty</p>
              </div>
            </div>
            <div className="prop-card">
              <div className="prop-icon" style={{ background: '#ebf0ff', color: '#3749bb' }}><Truck size={22} /></div>
              <div className="prop-text">
                <h4>Fast Delivery</h4>
                <p>Nationwide home delivery in 24-72 hours</p>
              </div>
            </div>
            <div className="prop-card">
              <div className="prop-icon" style={{ background: '#ecfdf5', color: '#10b981' }}><RefreshCw size={22} /></div>
              <div className="prop-text">
                <h4>Easy Returns</h4>
                <p>7-day replacement policy for defects</p>
              </div>
            </div>
            <div className="prop-card">
              <div className="prop-icon" style={{ background: '#fef3c7', color: '#d97706' }}><Clock size={22} /></div>
              <div className="prop-text">
                <h4>24/7 Support</h4>
                <p>Dedicated tech helpline 16793</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Exactly 4 Categories Quick Row */}
      <section>
        <div className="container">
          <div className="section-header">
            <h2 className="section-title">Shop by Category</h2>
            <Link to="/category/desktop" style={{ fontSize: 14, fontWeight: 600, color: '#ef4a23', display: 'flex', alignItems: 'center' }}>
              Explore All <ChevronRight size={16} />
            </Link>
          </div>

          <div className="categories-grid">
            {categories.map((c) => {
              const Icon = c.icon
              return (
                <Link to={`/category/${c.slug}`} key={c.slug} className="category-card">
                  <div className="category-icon">
                    <Icon size={32} />
                  </div>
                  <h3>{c.name}</h3>
                  <small>{c.count}</small>
                </Link>
              )
            })}
          </div>
        </div>
      </section>

      {/* Featured Products */}
      <section>
        <div className="container">
          <div className="section-header">
            <h2 className="section-title">
              <Flame size={22} color="#ef4a23" /> Featured Products
            </h2>
            <Link to="/category/laptop" style={{ fontSize: 14, fontWeight: 600, color: '#3749bb', display: 'flex', alignItems: 'center' }}>
              View More <ChevronRight size={16} />
            </Link>
          </div>

          <div className="products-grid">
            {featuredProducts.slice(0, 8).map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        </div>
      </section>

      {/* Deal of the Day Banner */}
      <section>
        <div className="container">
          <div className="deal-section">
            <div className="deal-grid">
              <div className="deal-info">
                <span className="badge badge-discount" style={{ marginBottom: 12 }}>
                  🔥 Flash Deals
                </span>
                <h2>Deal of the Day</h2>
                <p>Special discounts on premium tech products. Offers expire soon!</p>
                <div className="countdown">
                  <div className="countdown-box">
                    <b>{String(dealTime.hours).padStart(2, '0')}</b>
                    <small>Hours</small>
                  </div>
                  <div className="countdown-box">
                    <b>{String(dealTime.minutes).padStart(2, '0')}</b>
                    <small>Mins</small>
                  </div>
                  <div className="countdown-box">
                    <b>{String(dealTime.seconds).padStart(2, '0')}</b>
                    <small>Secs</small>
                  </div>
                </div>
                <Link to="/category/component" className="btn btn-primary btn-lg">
                  Explore Flash Sales <ChevronRight size={18} />
                </Link>
              </div>

              {featuredProducts[0] && (
                <div style={{ maxWidth: 360, margin: '0 auto', width: '100%' }}>
                  <ProductCard product={featuredProducts[0]} />
                </div>
              )}
            </div>
          </div>
        </div>
      </section>

      {/* New Arrivals */}
      <section>
        <div className="container">
          <div className="section-header">
            <h2 className="section-title">New Arrivals</h2>
            <Link to="/category/component" style={{ fontSize: 14, fontWeight: 600, color: '#ef4a23', display: 'flex', alignItems: 'center' }}>
              View All <ChevronRight size={16} />
            </Link>
          </div>

          <div className="products-grid">
            {newArrivals.slice(0, 8).map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        </div>
      </section>

      {/* Brands Showcase */}
      <section>
        <div className="container">
          <div className="section-header">
            <h2 className="section-title">Top Brands</h2>
          </div>
          <div className="brands-row">
            {brands.map((b) => (
              <Link to={`/search?q=${encodeURIComponent(b)}`} key={b} className="brand-pill">
                {b}
              </Link>
            ))}
          </div>
        </div>
      </section>
    </div>
  )
}
