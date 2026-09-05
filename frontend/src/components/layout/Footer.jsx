import React from 'react'
import { Link } from 'react-router-dom'
import { Phone, Mail, MapPin, ShieldCheck, Truck, Clock, CreditCard } from 'lucide-react'

export default function Footer() {
  return (
    <footer>
      <div className="container">
        <div className="footer-grid">
          {/* Col 1: Store info */}
          <div className="footer-col">
            <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 12 }}>
              <img src="/icon-dark.png" alt="MI-Tech" style={{ height: 32, width: 'auto', objectFit: 'contain' }} />
              <h3 style={{ fontSize: 20, color: '#fff', margin: 0 }}>
                MI-<span style={{ color: '#00C2FF' }}>Tech</span>
              </h3>
            </div>
            <p style={{ fontSize: 13, lineHeight: 1.6, marginBottom: 16 }}>
              MI-Tech & Engineering Ltd. is Bangladesh’s premier computer and tech retail store.
              Providing 100% genuine products, official warranties, and fast nationwide delivery.
            </p>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 8, fontSize: 13 }}>
              <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <Phone size={14} color="#00C2FF" /> Hotline: <b>16793</b> / <b>09612316793</b>
              </span>
              <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <Mail size={14} color="#00C2FF" /> support@mitech.local
              </span>
              <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <MapPin size={14} color="#00C2FF" /> Head Office: 28 Kazi Nazrul Islam Ave, Dhaka 1000
              </span>
            </div>
          </div>

          {/* Col 2: Customer Care */}
          <div className="footer-col">
            <h3>Customer Care</h3>
            <ul className="footer-links">
              <li><Link to="/contact">Contact Us</Link></li>
              <li><Link to="/order-tracking">Order Tracking</Link></li>
              <li><Link to="/warranty-policy">Warranty Policy</Link></li>
              <li><Link to="/return-policy">Return & Refund</Link></li>
              <li><Link to="/faq">Frequently Asked Questions</Link></li>
              <li><Link to="/emi-info">EMI Information</Link></li>
            </ul>
          </div>

          {/* Col 3: About & Policies */}
          <div className="footer-col">
            <h3>About & Terms</h3>
            <ul className="footer-links">
              <li><Link to="/about">About MI-Tech</Link></li>
              <li><Link to="/terms">Terms & Conditions</Link></li>
              <li><Link to="/privacy">Privacy Policy</Link></li>
              <li><Link to="/compare">Compare Products</Link></li>
              <li><Link to="/cart">View Shopping Cart</Link></li>
              <li><Link to="/account">My Account</Link></li>
            </ul>
          </div>

          {/* Col 4: Payments & Security */}
          <div className="footer-col">
            <h3>Payment & Security</h3>
            <p style={{ fontSize: 13, marginBottom: 12 }}>
              We support all major payment channels across Bangladesh:
            </p>
            <div style={{
              display: 'flex',
              flexWrap: 'wrap',
              gap: 8,
              background: '#0f1f2e',
              padding: 12,
              borderRadius: 6,
              marginBottom: 16,
              fontSize: 12,
              fontWeight: 700,
              color: '#fff',
            }}>
              <span style={{ background: '#e2136e', padding: '4px 8px', borderRadius: 4 }}>bKash</span>
              <span style={{ background: '#f7941d', padding: '4px 8px', borderRadius: 4 }}>Nagad</span>
              <span style={{ background: '#8b1d88', padding: '4px 8px', borderRadius: 4 }}>Rocket</span>
              <span style={{ background: '#1a1f71', padding: '4px 8px', borderRadius: 4 }}>VISA</span>
              <span style={{ background: '#eb001b', padding: '4px 8px', borderRadius: 4 }}>Mastercard</span>
              <span style={{ background: '#2563eb', padding: '4px 8px', borderRadius: 4 }}>COD</span>
            </div>
            <div style={{ display: 'flex', gap: 16, fontSize: 12 }}>
              <span style={{ display: 'flex', alignItems: 'center', gap: 4 }}><ShieldCheck size={14} color="#10b981" /> 100% Secure</span>
              <span style={{ display: 'flex', alignItems: 'center', gap: 4 }}><Truck size={14} color="#38bdf8" /> Express Delivery</span>
            </div>
          </div>
        </div>

        {/* Bottom */}
        <div className="footer-bottom">
          <p>© 2026 MI-Tech Ltd. | All Rights Reserved. Bangladesh's Leading Computer Hardware & Technology Store.</p>
          <p>Powered by Golang Echo + GORM + PostgreSQL + React Vite</p>
        </div>
      </div>
    </footer>
  )
}
