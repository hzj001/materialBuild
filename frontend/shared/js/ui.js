/**
 * 建材通 UI 组件库
 */
const UI = {
  CAT_ICONS: {
    '水泥砂石': '🏗️', '瓷砖地板': '🧱', '油漆涂料': '🎨',
    '管材管件': '🔧', '五金工具': '🔩', '门窗定制': '🚪'
  },

  toast(msg, duration = 2000) {
    const el = document.createElement('div');
    el.className = 'toast';
    el.textContent = msg;
    document.body.appendChild(el);
    setTimeout(() => el.remove(), duration);
  },

  formatPrice(price) {
    return '¥' + Number(price).toFixed(2);
  },

  orderStatusText(status) {
    const map = { 0: '待支付', 1: '已支付', 2: '配送中', 3: '已完成', 4: '已取消', 5: '售后中' };
    return map[status] || '未知';
  },

  orderStatusTag(status) {
    const cls = { 0: 'tag-pending', 1: 'tag-success', 2: 'tag-pending', 3: 'tag-success', 4: 'tag-danger', 5: 'tag-pending' };
    return `<span class="tag ${cls[status] || ''}">${this.orderStatusText(status)}</span>`;
  },

  merchantStatusText(status) {
    const map = { 0: '待审核', 1: '营业中', 2: '休息中', 3: '已禁用' };
    return map[status] || '未知';
  },

  updateCartBadge() {
    const el = document.getElementById('cartBadge');
    if (!el) return;
    const n = Store.cartCount();
    if (n > 0) {
      el.textContent = n > 99 ? '99+' : n;
      el.style.display = 'flex';
    } else {
      el.style.display = 'none';
    }
  },

  renderNavbar(title, showBack = false, transparent = false) {
    return `<div class="navbar${transparent ? ' transparent' : ''}">
      ${showBack ? '<button class="back" onclick="Router.back()">‹</button>' : ''}
      ${title}
    </div>`;
  },

  renderHomeHeader(cityName) {
    const name = cityName || '选择城市';
    return `
      <div class="home-header">
        <div class="brand-row">
          <div class="logo">建材通 <span class="logo-badge">本地购</span></div>
          <div class="location" onclick="Router.navigate('/city')">📍 ${name} ›</div>
        </div>
        <div class="home-search">
          <span>🔍</span>
          <input type="text" id="searchInput" placeholder="搜索水泥、瓷砖、管材、五金..."
            onkeydown="if(event.key==='Enter') Router.navigate('/search?keyword='+encodeURIComponent(this.value))">
          <button class="search-btn" onclick="var v=document.getElementById('searchInput').value;Router.navigate('/search?keyword='+encodeURIComponent(v))">搜索</button>
        </div>
      </div>
      <div class="promise-bar">
        <div class="promise-item"><span class="pi-icon">✓</span><span>正品保障</span></div>
        <div class="promise-item"><span class="pi-icon">🚚</span><span>同城配送</span></div>
        <div class="promise-item"><span class="pi-icon">💰</span><span>价格透明</span></div>
        <div class="promise-item"><span class="pi-icon">🛡️</span><span>售后无忧</span></div>
      </div>`;
  },

  renderBanner() {
    return `
      <div class="banner-wrap">
        <div class="banner-swiper">
          <div class="banner-slide"><h3>本地建材 · 一站购齐</h3><p>水泥瓷砖管材 · 30分钟响应 · 送货上门</p></div>
          <div class="banner-slide"><h3>新用户专享</h3><p>精选商家 · 限时优惠 · 满减活动进行中</p></div>
          <div class="banner-dots"><span></span><span></span></div>
        </div>
      </div>`;
  },

  renderCategoryGrid(categories) {
    if (!categories || !categories.length) return '';
    const items = categories.slice(0, 8).map(c => `
      <div class="category-grid-item" onclick="Router.navigate('/category-list?id=${c.id}&name=${encodeURIComponent(c.name)}')">
        <div class="cg-icon">${this.CAT_ICONS[c.name] || '📦'}</div>
        <span>${c.name}</span>
      </div>`).join('');
    return `
      <div class="section">
        <div class="section-head"><h2>全部分类</h2><span class="more" onclick="Router.navigate('/category')">更多 ›</span></div>
        <div class="category-grid">${items}</div>
      </div>`;
  },

  renderProductGrid(products) {
    if (!products || !products.length) {
      return '<div class="empty"><div class="icon">📦</div><p>暂无商品<br>请先选择城市或稍后再来</p></div>';
    }
    return `<div class="product-grid">${products.map(p => this.renderProductGridCard(p)).join('')}</div>`;
  },

  renderProductGridCard(p) {
    const discount = p.price > p.sale_price;
    return `
      <div class="product-grid-card" onclick="Router.navigate('/product?id=${p.id}')">
        <img src="${p.cover_image || 'https://via.placeholder.com/300'}" alt="" loading="lazy">
        <div class="pgc-body">
          <div class="pgc-name">${p.name}</div>
          <div class="pgc-meta">${p.merchant ? p.merchant.name : '本地商家'}</div>
          <div class="pgc-price"><small>¥</small>${Number(p.sale_price).toFixed(2)}<small>/${p.unit || '件'}</small></div>
          ${discount ? '<span class="pgc-tag">特惠</span>' : ''}
        </div>
      </div>`;
  },

  renderProductCard(p) {
    const discount = p.price > p.sale_price;
    return `
      <div class="product-card" onclick="Router.navigate('/product?id=${p.id}')">
        <img src="${p.cover_image || 'https://via.placeholder.com/96'}" alt="" loading="lazy">
        <div class="info">
          <div class="name">${p.name}</div>
          <div class="merchant">${p.merchant ? p.merchant.name : ''} · 同城配送</div>
          <div class="price-row">
            <div class="price"><small>¥</small>${Number(p.sale_price).toFixed(2)}<small>/${p.unit || '件'}</small></div>
            ${discount ? `<span class="price-old">¥${Number(p.price).toFixed(2)}</span>` : ''}
          </div>
        </div>
      </div>`;
  },

  renderMerchantCard(m) {
    return `
      <div class="merchant-card" onclick="Router.navigate('/search?merchant=${m.id}')">
        <div class="mc-logo">🏪</div>
        <div class="mc-info">
          <div class="mc-name">${m.name}</div>
          <div class="mc-sub">${m.address || '本地建材商家'}</div>
          <div class="mc-rating">★ ${Number(m.rating || 5).toFixed(1)} · 营业中</div>
        </div>
      </div>`;
  },

  renderMineHeader(user) {
    if (!user) {
      return `
        <div class="mine-header">
          <div class="mine-user">
            <div class="mine-avatar">👤</div>
            <div><div class="mine-name">登录享更多服务</div><div class="mine-phone">订单 · 售后 · 专属优惠</div></div>
          </div>
        </div>`;
    }
    return `
      <div class="mine-header">
        <div class="mine-user">
          <div class="mine-avatar">👤</div>
          <div>
            <div class="mine-name">${user.nickname || '建材通用户'}</div>
            <div class="mine-phone">${user.phone || ''}</div>
          </div>
        </div>
      </div>`;
  },

  renderLoginPage(roleLabel, onLoginFn, onRegisterFn) {
    return `
      <div class="login-page">
        <div class="login-brand">
          <h1>建材通</h1>
          <p>${roleLabel}</p>
        </div>
        <div class="login-form-wrap">
          <div class="form-group">
            <label class="form-label">手机号</label>
            <input class="form-input" id="phone" type="tel" placeholder="请输入手机号">
          </div>
          <div class="form-group">
            <label class="form-label">密码</label>
            <input class="form-input" id="password" type="password" placeholder="请输入密码">
          </div>
          <button class="btn btn-primary btn-block" onclick="${onLoginFn}()">登录</button>
          <button class="btn btn-outline btn-block" style="margin-top:12px" onclick="${onRegisterFn}()">注册账号</button>
        </div>
      </div>`;
  },
};
