import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import { Phone, Mail, MapPin, Send, HelpCircle, ShieldCheck, CreditCard, ChevronRight } from 'lucide-react'
import toast from 'react-hot-toast'

export function AboutPage() {
  return (
    <div className="container" style={{ padding: '32px 16px', maxWidth: 900 }}>
      <div className="breadcrumbs">
        <Link to="/">Home</Link> <ChevronRight size={12} style={{ display: 'inline' }} /> <span>About Us</span>
      </div>
      <div className="checkout-card">
        <h1 style={{ fontSize: 26, fontWeight: 800, marginBottom: 16 }}>About Star Tech Clone</h1>
        <p style={{ fontSize: 15, lineHeight: 1.8, color: '#334155', marginBottom: 16 }}>
          Star Tech is one of the leading tech retailers in Bangladesh. Founded with a vision to make technology accessible, dependable, and affordable, we deliver computing systems, accessories, and professional tech gear to tech enthusiasts and businesses across the nation.
        </p>
        <h3 style={{ fontSize: 18, fontWeight: 700, margin: '20px 0 10px' }}>Our Mission</h3>
        <p style={{ fontSize: 14, lineHeight: 1.8, color: '#475569', marginBottom: 16 }}>
          To provide genuine technology hardware, backed by official warranties, honest pricing, and prompt customer support.
        </p>
        <h3 style={{ fontSize: 18, fontWeight: 700, margin: '20px 0 10px' }}>Nationwide Branches</h3>
        <p style={{ fontSize: 14, lineHeight: 1.8, color: '#475569' }}>
          We operate 18+ retail outlets across Dhaka, Chittagong, Rajshahi, Khulna, Sylhet, Rangpur, and Mymensingh, supported by an advanced online order fulfillment center.
        </p>
      </div>
    </div>
  )
}

export function ContactPage() {
  const [formData, setFormData] = useState({ name: '', email: '', phone: '', subject: '', message: '' })
  const [sent, setSent] = useState(false)

  const handleSubmit = (e) => {
    e.preventDefault()
    setSent(true)
    toast.success('Thank you! Your inquiry has been sent to our customer care team.')
    setFormData({ name: '', email: '', phone: '', subject: '', message: '' })
  }

  return (
    <div className="container" style={{ padding: '32px 16px', maxWidth: 1000 }}>
      <div className="breadcrumbs">
        <Link to="/">Home</Link> <ChevronRight size={12} style={{ display: 'inline' }} /> <span>Contact Us</span>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1.2fr', gap: 28 }}>
        <div className="checkout-card">
          <h2 style={{ fontSize: 22, fontWeight: 800, marginBottom: 16 }}>Get in Touch</h2>
          <p style={{ color: '#64748b', fontSize: 14, marginBottom: 24 }}>
            Our customer care specialists are available 7 days a week from 9 AM to 8 PM.
          </p>

          <div style={{ display: 'flex', flexDirection: 'column', gap: 16, fontSize: 14 }}>
            <div style={{ display: 'flex', gap: 12 }}>
              <Phone color="#ef4a23" size={20} />
              <div>
                <b>Helpline:</b>
                <p>16793 (Local calls)</p>
                <p>09612316793 (Overseas / Mobile)</p>
              </div>
            </div>

            <div style={{ display: 'flex', gap: 12 }}>
              <Mail color="#ef4a23" size={20} />
              <div>
                <b>Email Support:</b>
                <p>support@mitech.local</p>
                <p>corporate@mitech.local</p>
              </div>
            </div>

            <div style={{ display: 'flex', gap: 12 }}>
              <MapPin color="#ef4a23" size={20} />
              <div>
                <b>Headquarters:</b>
                <p>28 Kazi Nazrul Islam Avenue, Shahbagh, Dhaka 1000</p>
              </div>
            </div>
          </div>
        </div>

        <div className="checkout-card">
          <h2 style={{ fontSize: 20, fontWeight: 700, marginBottom: 16 }}>Send Us a Message</h2>
          <form onSubmit={handleSubmit} className="form-grid">
            <div className="form-group">
              <label>Your Name *</label>
              <input
                type="text"
                required
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              />
            </div>
            <div className="form-group">
              <label>Your Email *</label>
              <input
                type="email"
                required
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              />
            </div>
            <div className="form-group form-full">
              <label>Subject</label>
              <input
                type="text"
                required
                value={formData.subject}
                onChange={(e) => setFormData({ ...formData, subject: e.target.value })}
              />
            </div>
            <div className="form-group form-full">
              <label>Message *</label>
              <textarea
                rows={4}
                required
                value={formData.message}
                onChange={(e) => setFormData({ ...formData, message: e.target.value })}
              />
            </div>
            <div className="form-full">
              <button type="submit" className="btn btn-primary btn-lg">
                <Send size={16} /> Send Message
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}

export function TermsPage() {
  return (
    <div className="container" style={{ padding: '32px 16px', maxWidth: 900 }}>
      <div className="checkout-card">
        <h1 style={{ fontSize: 24, fontWeight: 800, marginBottom: 16 }}>Terms & Conditions</h1>
        <p style={{ fontSize: 14, lineHeight: 1.8, color: '#475569', marginBottom: 12 }}>
          1. <b>Pricing and Availability:</b> All prices quoted on MI-Tech are in Bangladeshi Taka (BDT) including applicable VAT. Product availability is updated regularly.
        </p>
        <p style={{ fontSize: 14, lineHeight: 1.8, color: '#475569', marginBottom: 12 }}>
          2. <b>Order Confirmation:</b> Orders placed online will receive an automated order reference number. MI-Tech customer care verifies orders before shipping.
        </p>
        <p style={{ fontSize: 14, lineHeight: 1.8, color: '#475569' }}>
          3. <b>Warranty Claims:</b> All branded warranty services are provided in accordance with the official distributor warranty policy.
        </p>
      </div>
    </div>
  )
}

export function WarrantyPolicyPage() {
  return (
    <div className="container" style={{ padding: '32px 16px', maxWidth: 900 }}>
      <div className="checkout-card">
        <h1 style={{ fontSize: 24, fontWeight: 800, marginBottom: 16 }}>Official Warranty Policy</h1>
        <div style={{ display: 'flex', gap: 12, background: '#f8fafc', padding: 16, borderRadius: 6, marginBottom: 20 }}>
          <ShieldCheck color="#10b981" size={24} />
          <div>
            <b>100% Official Brand Warranty</b>
            <p style={{ fontSize: 13, color: '#64748b' }}>Every product sold at MI-Tech comes with full manufacturer warranty coverage.</p>
          </div>
        </div>
        <p style={{ fontSize: 14, lineHeight: 1.8, color: '#475569', marginBottom: 12 }}>
          • <b>Laptops & Desktops:</b> 1 to 3 Years comprehensive warranty depending on manufacturer terms.
        </p>
        <p style={{ fontSize: 14, lineHeight: 1.8, color: '#475569', marginBottom: 12 }}>
          • <b>Components:</b> Processors (3 Years), Motherboards (3 Years), GPUs (2-3 Years), RAM (Lifetime / 10 Years).
        </p>
        <p style={{ fontSize: 14, lineHeight: 1.8, color: '#475569' }}>
          • <b>Monitors:</b> 3 Years warranty (Panel + Parts + Service).
        </p>
      </div>
    </div>
  )
}

export function EMIPage() {
  return (
    <div className="container" style={{ padding: '32px 16px', maxWidth: 900 }}>
      <div className="checkout-card">
        <h1 style={{ fontSize: 24, fontWeight: 800, marginBottom: 16 }}>0% Interest EMI Facilities</h1>
        <p style={{ fontSize: 14, lineHeight: 1.8, color: '#475569', marginBottom: 20 }}>
          Enjoy 0% interest EMI options on purchases exceeding ৳10,000 for up to 36 months using credit cards from 20+ partner banks.
        </p>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 12 }}>
          {['City Bank (Amex)', 'BRAC Bank', 'Standard Chartered', 'Eastern Bank (EBL)', 'Dutch-Bangla Bank', 'Mutual Trust Bank', 'Prime Bank', 'UCB', 'Dhaka Bank'].map((b) => (
            <div key={b} style={{ background: '#f8fafc', padding: 12, borderRadius: 6, border: '1px solid #e2e8f0', fontSize: 13, fontWeight: 600 }}>
              💳 {b}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

export function FAQPage() {
  const faqs = [
    { q: 'How can I track my order?', a: 'Visit our Order Tracking page and enter your Order Number (ST-...) or courier tracking ID.' },
    { q: 'What payment options do you support?', a: 'We accept Cash on Delivery (COD), bKash, Nagad, Rocket, VISA, Mastercard, and Bank Wire.' },
    { q: 'How long does delivery take?', a: 'Inside Dhaka: 24 to 48 hours. Outside Dhaka: 48 to 72 hours via express courier.' },
    { q: 'Are all products authentic?', a: 'Yes, 100% of our inventory is sourced from authorized global brand distributors.' },
  ]
  return (
    <div className="container" style={{ padding: '32px 16px', maxWidth: 900 }}>
      <div className="checkout-card">
        <h1 style={{ fontSize: 24, fontWeight: 800, marginBottom: 20 }}>Frequently Asked Questions</h1>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
          {faqs.map((f, i) => (
            <div key={i} style={{ borderBottom: '1px solid #f1f5f9', paddingBottom: 16 }}>
              <h3 style={{ fontSize: 16, fontWeight: 700, color: '#0f172a', marginBottom: 6 }}>
                ❓ {f.q}
              </h3>
              <p style={{ fontSize: 14, color: '#475569', lineHeight: 1.6 }}>{f.a}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
