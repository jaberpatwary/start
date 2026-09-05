import React, { createContext, useContext, useEffect, useState } from 'react'
import api from '../api/client'
import { useAuth } from './AuthContext'
import toast from 'react-hot-toast'

const ShopContext = createContext()

export function ShopProvider({ children }) {
  const { isAuthenticated } = useAuth()

  // Cart state
  const [cart, setCart] = useState(() => {
    try {
      const saved = localStorage.getItem('startech_cart')
      return saved ? JSON.parse(saved) : []
    } catch {
      return []
    }
  })

  // Wishlist state
  const [wishlist, setWishlist] = useState(() => {
    try {
      const saved = localStorage.getItem('startech_wishlist')
      return saved ? JSON.parse(saved) : []
    } catch {
      return []
    }
  })

  // Compare state (up to 4 products)
  const [compare, setCompare] = useState(() => {
    try {
      const saved = localStorage.getItem('startech_compare')
      return saved ? JSON.parse(saved) : []
    } catch {
      return []
    }
  })

  // Coupon state
  const [appliedCoupon, setAppliedCoupon] = useState(null)

  // Sync to local storage
  useEffect(() => {
    localStorage.setItem('startech_cart', JSON.stringify(cart))
  }, [cart])

  useEffect(() => {
    localStorage.setItem('startech_wishlist', JSON.stringify(wishlist))
  }, [wishlist])

  useEffect(() => {
    localStorage.setItem('startech_compare', JSON.stringify(compare))
  }, [compare])

  // Cart operations
  const addToCart = (product, quantity = 1) => {
    setCart((prev) => {
      const existing = prev.find((item) => item.product.id === product.id)
      if (existing) {
        toast.success(`Updated quantity for ${product.name}`)
        return prev.map((item) =>
          item.product.id === product.id
            ? { ...item, quantity: item.quantity + quantity }
            : item
        )
      } else {
        toast.success(`Added ${product.name} to cart`)
        return [...prev, { product, quantity }]
      }
    })
  }

  const updateCartQty = (productId, quantity) => {
    if (quantity < 1) {
      removeFromCart(productId)
      return
    }
    setCart((prev) =>
      prev.map((item) =>
        item.product.id === productId ? { ...item, quantity } : item
      )
    )
  }

  const removeFromCart = (productId) => {
    setCart((prev) => prev.filter((item) => item.product.id !== productId))
    toast.success('Item removed from cart')
  }

  const clearCart = () => {
    setCart([])
    setAppliedCoupon(null)
    localStorage.removeItem('startech_cart')
  }

  // Calculate totals
  const subtotal = cart.reduce((sum, item) => {
    const price = item.product.discount_price || item.product.price
    return sum + price * item.quantity
  }, 0)

  const discount = appliedCoupon
    ? appliedCoupon.type === 'PERCENT'
      ? Math.round((subtotal * appliedCoupon.value) / 100)
      : appliedCoupon.value
    : 0

  const deliveryFee = subtotal >= 50000 ? 0 : subtotal > 0 ? 120 : 0
  const grandTotal = Math.max(0, subtotal - discount + deliveryFee)

  const applyCoupon = async (code) => {
    try {
      const data = await api.get(`/coupons/validate?code=${encodeURIComponent(code)}`)
      if (data?.coupon) {
        if (subtotal < data.coupon.min_order) {
          toast.error(`Minimum order amount for this coupon is ৳${data.coupon.min_order.toLocaleString()}`)
          return false
        }
        setAppliedCoupon(data.coupon)
        toast.success(`Coupon "${data.coupon.code}" applied successfully!`)
        return true
      }
    } catch (err) {
      toast.error(err.message || 'Invalid or expired coupon code')
      return false
    }
  }

  const removeCoupon = () => {
    setAppliedCoupon(null)
    toast.success('Coupon removed')
  }

  // Wishlist operations
  const toggleWishlist = (product) => {
    const exists = wishlist.some((item) => item.id === product.id)
    if (exists) {
      setWishlist((prev) => prev.filter((item) => item.id !== product.id))
      toast.success('Removed from wishlist')
    } else {
      setWishlist((prev) => [...prev, product])
      toast.success('Added to wishlist')
    }
  }

  const isInWishlist = (productId) => {
    return wishlist.some((item) => item.id === productId)
  }

  // Compare operations (max 4 products)
  const addToCompare = (product) => {
    if (compare.some((item) => item.id === product.id)) {
      toast('Product is already in compare list')
      return
    }
    if (compare.length >= 4) {
      toast.error('You can compare a maximum of 4 products at a time')
      return
    }
    setCompare((prev) => [...prev, product])
    toast.success(`Added ${product.name} to comparison`)
  }

  const removeFromCompare = (productId) => {
    setCompare((prev) => prev.filter((item) => item.id !== productId))
    toast.success('Product removed from comparison')
  }

  const clearCompare = () => {
    setCompare([])
  }

  return (
    <ShopContext.Provider
      value={{
        cart,
        cartCount: cart.reduce((n, i) => n + i.quantity, 0),
        addToCart,
        updateCartQty,
        removeFromCart,
        clearCart,
        subtotal,
        discount,
        deliveryFee,
        grandTotal,
        appliedCoupon,
        applyCoupon,
        removeCoupon,
        wishlist,
        wishlistCount: wishlist.length,
        toggleWishlist,
        isInWishlist,
        compare,
        compareCount: compare.length,
        addToCompare,
        removeFromCompare,
        clearCompare,
      }}
    >
      {children}
    </ShopContext.Provider>
  )
}

export const useShop = () => useContext(ShopContext)
