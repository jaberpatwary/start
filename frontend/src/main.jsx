import React, { Suspense, lazy } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import './styles.css'

import { AuthProvider } from './context/AuthContext'
import { ShopProvider } from './context/ShopContext'

import Header from './components/layout/Header'
import Footer from './components/layout/Footer'

// Eagerly import auth pages (small + used often)
import { LoginPage, RegisterPage } from './pages/AuthPages'

// Lazy load heavier pages
const Home = lazy(() => import('./pages/Home'))
const CategoryPage = lazy(() => import('./pages/CategoryPage'))
const ProductDetail = lazy(() => import('./pages/ProductDetail'))
const CartPage = lazy(() => import('./pages/CartPage'))
const CheckoutPage = lazy(() => import('./pages/CheckoutPage'))
const OrderSuccessPage = lazy(() => import('./pages/OrderSuccessPage'))
const OrderTrackingPage = lazy(() => import('./pages/OrderTrackingPage'))
const AccountPage = lazy(() => import('./pages/AccountPage'))
const ComparePage = lazy(() => import('./pages/ComparePage'))
const StaticPages = lazy(() => import('./pages/StaticPages'))

// Admin pages
const AdminLayout = lazy(() => import('./pages/admin/AdminLayout'))
const AdminDashboard = lazy(() => import('./pages/admin/AdminDashboard'))
const AdminProducts = lazy(() => import('./pages/admin/AdminProducts'))
const AdminCatalog = lazy(() => import('./pages/admin/AdminCatalog'))
const AdminOrders = lazy(() => import('./pages/admin/AdminOrders'))
const AdminUsersReviews = lazy(() => import('./pages/admin/AdminUsersReviews'))
const AdminInventory = lazy(() => import('./pages/admin/AdminInventory'))
const AdminReports = lazy(() => import('./pages/admin/AdminReports'))
const AdminSettings = lazy(() => import('./pages/admin/AdminSettings'))

function LoadingSpinner() {
  return (
    <div style={{
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      minHeight: '60vh',
      flexDirection: 'column',
      gap: 16,
    }}>
      <div className="loading-spinner" />
      <span style={{ color: '#64748b', fontSize: 14 }}>Loading...</span>
    </div>
  )
}

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <ShopProvider>
          <Header />
          <main style={{ minHeight: '60vh' }}>
            <Suspense fallback={<LoadingSpinner />}>
              <Routes>
                {/* Public routes */}
                <Route path="/" element={<Home />} />
                <Route path="/category/:slug" element={<CategoryPage />} />
                <Route path="/search" element={<CategoryPage isSearch={true} />} />
                <Route path="/product/:slug" element={<ProductDetail />} />
                <Route path="/cart" element={<CartPage />} />
                <Route path="/checkout" element={<CheckoutPage />} />
                <Route path="/order-success" element={<OrderSuccessPage />} />
                <Route path="/order-tracking" element={<OrderTrackingPage />} />
                <Route path="/compare" element={<ComparePage />} />

                {/* Auth routes */}
                <Route path="/login" element={<LoginPage />} />
                <Route path="/register" element={<RegisterPage />} />

                {/* User account */}
                <Route path="/account/*" element={<AccountPage />} />

                {/* Admin routes */}
                <Route path="/admin" element={<AdminLayout />}>
                  <Route index element={<AdminDashboard />} />
                  <Route path="products" element={<AdminProducts />} />
                  <Route path="catalog" element={<AdminCatalog />} />
                  <Route path="orders" element={<AdminOrders />} />
                  <Route path="users-reviews" element={<AdminUsersReviews />} />
                  <Route path="inventory" element={<AdminInventory />} />
                  <Route path="reports" element={<AdminReports />} />
                  <Route path="settings" element={<AdminSettings />} />
                </Route>

                {/* Static pages */}
                <Route path="/about" element={<StaticPages page="about" />} />
                <Route path="/contact" element={<StaticPages page="contact" />} />
                <Route path="/privacy-policy" element={<StaticPages page="privacy" />} />
                <Route path="/terms-conditions" element={<StaticPages page="terms" />} />
                <Route path="/return-policy" element={<StaticPages page="return" />} />
                <Route path="/warranty-policy" element={<StaticPages page="warranty" />} />
                <Route path="/faq" element={<StaticPages page="faq" />} />
                <Route path="/emi-info" element={<StaticPages page="emi" />} />

                {/* Fallback */}
                <Route path="*" element={<Navigate to="/" replace />} />
              </Routes>
            </Suspense>
          </main>
          <Footer />
          <Toaster
            position="top-right"
            toastOptions={{
              duration: 3000,
              style: {
                background: '#1e293b',
                color: '#f1f5f9',
                borderRadius: 8,
                fontSize: 14,
              },
              success: { iconTheme: { primary: '#22c55e', secondary: '#fff' } },
              error: { iconTheme: { primary: '#ef4444', secondary: '#fff' } },
            }}
          />
        </ShopProvider>
      </AuthProvider>
    </BrowserRouter>
  )
}

createRoot(document.getElementById('root')).render(<App />)
