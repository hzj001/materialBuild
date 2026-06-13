/**
 * 轻量状态管理
 */
const Store = {
  _state: {
    user: null,
    city: null,
    cart: [],
  },

  init() {
    try {
      const user = localStorage.getItem('user');
      const city = localStorage.getItem('city');
      if (user) this._state.user = JSON.parse(user);
      if (city) this._state.city = JSON.parse(city);
      const cart = localStorage.getItem('cart');
      if (cart) this._state.cart = JSON.parse(cart);
    } catch (e) { /* ignore */ }
  },

  get(key) { return this._state[key]; },

  set(key, value) {
    this._state[key] = value;
    if (['user', 'city', 'cart'].includes(key)) {
      localStorage.setItem(key, JSON.stringify(value));
    }
  },

  setUser(user) {
    this.set('user', user);
    if (user) localStorage.setItem('token', user.token || localStorage.getItem('token'));
  },

  logout() {
    this._state.user = null;
    localStorage.removeItem('user');
    localStorage.removeItem('token');
  },

  addToCart(product, quantity = 1) {
    const cart = [...this._state.cart];
    const idx = cart.findIndex(i => i.product_id === product.id);
    if (idx >= 0) {
      cart[idx].quantity += quantity;
    } else {
      cart.push({
        product_id: product.id,
        merchant_id: product.merchant_id,
        name: product.name,
        image: product.cover_image,
        price: product.sale_price,
        unit: product.unit,
        quantity,
      });
    }
    this.set('cart', cart);
  },

  cartCount() {
    return this._state.cart.reduce((s, i) => s + i.quantity, 0);
  },

  cartTotal() {
    return this._state.cart.reduce((s, i) => s + i.price * i.quantity, 0);
  },
};
