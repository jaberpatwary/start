import React, { useState, useEffect } from 'react'
import { Settings, Save, Store, Truck, DollarSign, Bell } from 'lucide-react'
import api from '../../api/client'
import toast from 'react-hot-toast'

export default function AdminSettings() {
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  const [settings, setSettings] = useState({
    store_name: 'MI-Tech Computer Hardware Store',
    contact_email: 'support@mitech.bd',
    contact_phone: '+880 1700-000000',
    store_address: '608 & 609, Multiplan Center, Elephant Road, Dhaka-1205',
    currency_symbol: '৳',
    delivery_fee_dhaka: '60',
    delivery_fee_outside: '120',
    free_shipping_threshold: '50000',
    announcement_bar: '⚡ Special Offer: Get 5% Cashback on BKash/Nagad Payment!',
    maintenance_mode: 'false',
  })

  useEffect(() => {
    const fetchSettings = async () => {
      try {
        const res = await api.get('/admin/settings')
        if (res?.settings) {
          setSettings((prev) => ({ ...prev, ...res.settings }))
        }
      } catch (err) {
        console.warn('Could not load remote settings, using default state')
      } finally {
        setLoading(false)
      }
    }
    fetchSettings()
  }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    setSaving(true)
    try {
      await api.put('/admin/settings', settings)
      toast.success('Store settings saved successfully')
    } catch (err) {
      toast.error(err.message || 'Failed to update settings')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return <div style={{ padding: 40, color: '#64748b' }}>Loading store settings...</div>
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
        <div>
          <h1 style={{ fontSize: 24, fontWeight: 800 }}>Store Configuration & Settings</h1>
          <small style={{ color: '#64748b' }}>Configure store details, delivery charges, contact info, and site banners</small>
        </div>
      </div>

      <form onSubmit={handleSubmit} style={{ maxWidth: 800 }}>
        {/* Store General Info */}
        <div className="admin-card" style={{ padding: 24, marginBottom: 24 }}>
          <h3 style={{ fontSize: 16, fontWeight: 700, marginBottom: 16, display: 'flex', alignItems: 'center', gap: 8, color: '#3749bb' }}>
            <Store size={18} /> Store General Identity
          </h3>
          <div className="form-grid">
            <div className="form-group form-full">
              <label>Store Name Title *</label>
              <input
                type="text"
                required
                value={settings.store_name}
                onChange={(e) => setSettings({ ...settings, store_name: e.target.value })}
              />
            </div>
            <div className="form-group">
              <label>Support Email *</label>
              <input
                type="email"
                required
                value={settings.contact_email}
                onChange={(e) => setSettings({ ...settings, contact_email: e.target.value })}
              />
            </div>
            <div className="form-group">
              <label>Hotline Phone Number *</label>
              <input
                type="text"
                required
                value={settings.contact_phone}
                onChange={(e) => setSettings({ ...settings, contact_phone: e.target.value })}
              />
            </div>
            <div className="form-group form-full">
              <label>Physical Address</label>
              <input
                type="text"
                value={settings.store_address}
                onChange={(e) => setSettings({ ...settings, store_address: e.target.value })}
              />
            </div>
          </div>
        </div>

        {/* Shipping & Currency Settings */}
        <div className="admin-card" style={{ padding: 24, marginBottom: 24 }}>
          <h3 style={{ fontSize: 16, fontWeight: 700, marginBottom: 16, display: 'flex', alignItems: 'center', gap: 8, color: '#ef4a23' }}>
            <Truck size={18} /> Delivery Charges & Shipping Rates
          </h3>
          <div className="form-grid">
            <div className="form-group">
              <label>Delivery Fee Inside Dhaka (৳)</label>
              <input
                type="number"
                value={settings.delivery_fee_dhaka}
                onChange={(e) => setSettings({ ...settings, delivery_fee_dhaka: e.target.value })}
              />
            </div>
            <div className="form-group">
              <label>Delivery Fee Outside Dhaka (৳)</label>
              <input
                type="number"
                value={settings.delivery_fee_outside}
                onChange={(e) => setSettings({ ...settings, delivery_fee_outside: e.target.value })}
              />
            </div>
            <div className="form-group">
              <label>Free Shipping Threshold (৳)</label>
              <input
                type="number"
                value={settings.free_shipping_threshold}
                onChange={(e) => setSettings({ ...settings, free_shipping_threshold: e.target.value })}
              />
            </div>
            <div className="form-group">
              <label>Currency Symbol</label>
              <input
                type="text"
                value={settings.currency_symbol}
                onChange={(e) => setSettings({ ...settings, currency_symbol: e.target.value })}
              />
            </div>
          </div>
        </div>

        {/* Announcement Bar */}
        <div className="admin-card" style={{ padding: 24, marginBottom: 24 }}>
          <h3 style={{ fontSize: 16, fontWeight: 700, marginBottom: 16, display: 'flex', alignItems: 'center', gap: 8, color: '#10b981' }}>
            <Bell size={18} /> Top Announcement Bar Banner
          </h3>
          <div className="form-group form-full">
            <label>Header Announcement Text</label>
            <input
              type="text"
              value={settings.announcement_bar}
              onChange={(e) => setSettings({ ...settings, announcement_bar: e.target.value })}
            />
          </div>
        </div>

        <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: 16 }}>
          <button type="submit" className="btn btn-primary" disabled={saving}>
            <Save size={16} /> {saving ? 'Saving Settings...' : 'Save Settings'}
          </button>
        </div>
      </form>
    </div>
  )
}
