// 建材通 - 客户端
Http.init(APP_CONFIG.API_BASE);
Store.init();

// ========== 首页 ==========
Router.register('/', async () => {
  const city = Store.get('city');
  let products = [], categories = [], merchants = [];
  try {
    const [pData, cats, mData] = await Promise.all([
      Http.get('/client/products', { city_id: city?.id, page_size: 12 }),
      Http.get('/common/categories'),
      Http.get('/client/merchants', { city_id: city?.id, page_size: 5 }),
    ]);
    products = pData.list || [];
    categories = cats || [];
    merchants = mData.list || [];
  } catch (e) { /* empty */ }

  const cityName = city ? city.name : '选择城市';
  return `
    <div class="page">
      ${UI.renderHomeHeader(cityName)}
      ${UI.renderBanner()}
      ${UI.renderCategoryGrid(categories)}
      <div class="section">
        <div class="section-head"><h2>猜你喜欢</h2><span class="more" onclick="Router.navigate('/category')">查看全部 ›</span></div>
      </div>
      ${UI.renderProductGrid(products)}
      ${merchants.length ? `
        <div class="section">
          <div class="section-head"><h2>附近商家</h2></div>
          ${merchants.map(m => UI.renderMerchantCard(m)).join('')}
        </div>` : ''}
    </div>`;
});

// ========== 城市选择 ==========
Router.register('/city', async () => {
  const cities = await Http.get('/common/cities');
  return `
    ${UI.renderNavbar('选择城市', true)}
    <div class="page">
      ${cities.map(c => `
        <div class="card" style="cursor:pointer" onclick="selectCity(${c.id},'${c.name}')">
          <strong>${c.name}</strong>
          <span style="color:var(--text-secondary);margin-left:8px">${c.province}</span>
        </div>
      `).join('')}
    </div>`;
});

function selectCity(id, name) {
  Store.set('city', { id, name });
  UI.toast('已切换到 ' + name);
  Router.navigate('/');
}

// ========== 搜索 ==========
Router.register('/search', async (params) => {
  const city = Store.get('city');
  const data = await Http.get('/client/products', {
    keyword: params.keyword, city_id: city?.id, page_size: 30,
  });
  const list = data.list || [];
  return `
    ${UI.renderNavbar('搜索结果', true)}
    ${list.length ? list.map(UI.renderProductCard).join('') :
      '<div class="empty"><div class="icon">🔍</div><p>未找到相关商品</p></div>'}
  `;
});

// ========== 分类 ==========
Router.register('/category', async () => {
  const cats = await Http.get('/common/categories');
  return `
    ${UI.renderNavbar('全部分类')}
    <div class="page">
      <div class="category-grid" style="margin:12px">
        ${(cats || []).map(c => `
          <div class="category-grid-item" onclick="Router.navigate('/category-list?id=${c.id}&name=${encodeURIComponent(c.name)}')">
            <div class="cg-icon">${UI.CAT_ICONS[c.name] || '📦'}</div>
            <span>${c.name}</span>
          </div>`).join('')}
      </div>
    </div>`;
});

Router.register('/category-list', async (params) => {
  const city = Store.get('city');
  const data = await Http.get('/client/products', {
    category_id: params.id, city_id: city?.id, page_size: 30,
  });
  const list = data.list || [];
  return `
    ${UI.renderNavbar(decodeURIComponent(params.name || '分类'), true)}
    <div class="page">${UI.renderProductGrid(list)}</div>
  `;
});

// ========== 商品详情 ==========
Router.register('/product', async (params) => {
  const p = await Http.get('/client/products/' + params.id);
  const media = (p.media || []).map(m =>
    m.type === 'video'
      ? `<video src="${m.url}" controls></video>`
      : `<img src="${m.url}" alt="">`
  ).join('');

  return `
    ${UI.renderNavbar('商品详情', true)}
    <div class="page">
      <div class="detail-hero">
        <img class="detail-img" src="${p.cover_image || 'https://via.placeholder.com/400'}" alt="">
      </div>
      <div class="detail-panel">
        <div class="detail-price">¥${Number(p.sale_price).toFixed(2)}<small>/${p.unit}</small>
          ${p.price > p.sale_price ? `<del style="font-size:14px;color:#bbb;font-weight:400;margin-left:8px">¥${p.price}</del>` : ''}
        </div>
        <div class="detail-name">${p.name}</div>
        ${p.merchant ? `<div style="margin-top:8px;font-size:13px;color:var(--text-secondary)">🏪 ${p.merchant.name} · ★ ${Number(p.merchant.rating||5).toFixed(1)}</div>` : ''}
      </div>
      ${p.merchant ? `
        <div class="contact-bar">
          <button class="btn btn-ghost btn-sm" onclick="location.href='tel:${p.merchant.contact_phone}'">📞 电话商家</button>
          <button class="btn btn-outline btn-sm" onclick="Router.navigate('/chat?merchant_id=${p.merchant_id}')">💬 在线咨询</button>
        </div>` : ''}
      ${media ? `<div class="media-gallery">${media}</div>` : ''}
      <div class="detail-desc">${p.description || '优质建材，同城配送，欢迎选购。'}</div>
      <div style="height:80px"></div>
      <div class="bottom-bar">
        <button class="btn btn-outline btn-cart" onclick="addCart(${JSON.stringify(p).replace(/"/g, '&quot;')})">加购物车</button>
        <button class="btn btn-primary" onclick="buyNow(${p.id}, ${p.merchant_id})">立即购买</button>
      </div>
    </div>
  `;
});

function addCart(product) {
  Store.addToCart(typeof product === 'object' ? product : JSON.parse(product));
  UI.updateCartBadge();
  UI.toast('已加入购物车');
}

function buyNow(productId, merchantId) {
  Store.addToCart({ id: productId, merchant_id: merchantId }, 1);
  Router.navigate('/checkout?merchant_id=' + merchantId);
}

// ========== 购物车 ==========
Router.register('/cart', async () => {
  const cart = Store.get('cart') || [];
  if (!cart.length) {
    return `${UI.renderNavbar('购物车')}
      <div class="empty"><div class="icon">🛒</div><p>购物车是空的</p>
      <button class="btn btn-primary" style="margin-top:16px" onclick="Router.navigate('/')">去逛逛</button></div>`;
  }

  const items = cart.map((item, i) => `
    <div class="product-card">
      <img src="${item.image || 'https://via.placeholder.com/90'}" alt="">
      <div class="info">
        <div class="name">${item.name}</div>
        <div class="price">¥${item.price} × ${item.quantity}</div>
      </div>
      <button class="btn btn-sm btn-outline" onclick="removeCart(${i})">删除</button>
    </div>
  `).join('');

  return `
    ${UI.renderNavbar('购物车')}
    ${items}
    <div class="card" style="text-align:right">
      <span>合计：</span><span class="price" style="font-size:20px">${UI.formatPrice(Store.cartTotal())}</span>
    </div>
    <div style="padding:12px">
      <button class="btn btn-primary btn-block" onclick="Router.navigate('/checkout?merchant_id=${cart[0].merchant_id}')">去结算</button>
    </div>
  `;
});

function removeCart(index) {
  const cart = [...Store.get('cart')];
  cart.splice(index, 1);
  Store.set('cart', cart);
  Router.navigate('/cart');
}

// ========== 结算 ==========
Router.register('/checkout', async (params) => {
  const user = Store.get('user');
  if (!user) { Router.navigate('/login'); return ''; }

  const cart = Store.get('cart').filter(i => String(i.merchant_id) === params.merchant_id);
  return `
    ${UI.renderNavbar('确认订单', true)}
    <div class="card">
      <div class="form-group">
        <label class="form-label">收货人</label>
        <input class="form-input" id="receiverName" placeholder="姓名">
      </div>
      <div class="form-group">
        <label class="form-label">联系电话</label>
        <input class="form-input" id="receiverPhone" placeholder="手机号">
      </div>
      <div class="form-group">
        <label class="form-label">收货地址</label>
        <input class="form-input" id="address" placeholder="详细地址">
      </div>
      <div class="form-group">
        <label class="form-label">备注</label>
        <input class="form-input" id="remark" placeholder="选填">
      </div>
    </div>
    <div class="card">
      ${cart.map(i => `<div style="padding:4px 0">${i.name} × ${i.quantity} = ¥${(i.price * i.quantity).toFixed(2)}</div>`).join('')}
      <div style="text-align:right;margin-top:8px;font-weight:600">合计：${UI.formatPrice(Store.cartTotal())}</div>
    </div>
    <div style="padding:12px">
      <button class="btn btn-primary btn-block" onclick="submitOrder(${params.merchant_id})">提交订单</button>
    </div>
  `;
});

async function submitOrder(merchantId) {
  const cart = Store.get('cart').filter(i => i.merchant_id === merchantId);
  try {
    const order = await Http.post('/client/orders', {
      merchant_id: merchantId,
      items: cart.map(i => ({ product_id: i.product_id, quantity: i.quantity })),
      receiver_name: document.getElementById('receiverName').value,
      receiver_phone: document.getElementById('receiverPhone').value,
      address: document.getElementById('address').value,
      remark: document.getElementById('remark').value,
    });
    Store.set('cart', Store.get('cart').filter(i => i.merchant_id !== merchantId));
    UI.toast('下单成功');
    Router.navigate('/order-detail?id=' + order.id);
  } catch (e) {
    UI.toast(e.message);
  }
}

// ========== 订单 ==========
Router.register('/orders', async () => {
  const data = await Http.get('/client/orders');
  const list = data.list || [];
  return `
    ${UI.renderNavbar('我的订单', true)}
    ${list.length ? list.map(o => `
      <div class="card" style="cursor:pointer" onclick="Router.navigate('/order-detail?id=${o.id}')">
        <div style="display:flex;justify-content:space-between;margin-bottom:8px">
          <span style="font-size:12px;color:var(--text-secondary)">${o.order_no}</span>
          ${UI.orderStatusTag(o.status)}
        </div>
        ${(o.items || []).map(i => `<div style="padding:4px 0">${i.product_name} × ${i.quantity}</div>`).join('')}
        <div style="text-align:right;margin-top:8px;font-weight:600">¥${o.pay_amount}</div>
      </div>
    `).join('') : '<div class="empty"><div class="icon">📋</div><p>暂无订单</p></div>'}
  `;
});

Router.register('/order-detail', async (params) => {
  const data = await Http.get('/client/orders');
  const order = (data.list || []).find(o => String(o.id) === params.id);
  if (!order) return '<div class="empty"><p>订单不存在</p></div>';

  return `
    ${UI.renderNavbar('订单详情', true)}
    <div class="card">
      <div>订单号：${order.order_no}</div>
      <div>状态：${UI.orderStatusText(order.status)}</div>
      <div>金额：¥${order.pay_amount}</div>
      ${order.delivery ? `<div>地址：${order.delivery.address}</div>` : ''}
    </div>
    ${order.status === 0 ? `<div style="padding:12px"><button class="btn btn-primary btn-block" onclick="payOrder(${order.id})">立即支付</button></div>` : ''}
    ${order.status >= 1 && order.status <= 3 ? `<div style="padding:12px"><button class="btn btn-outline btn-block" onclick="Router.navigate('/after-sale?order_id=${order.id}')">申请售后</button></div>` : ''}
  `;
});

async function payOrder(id) {
  try {
    await Http.post('/client/orders/' + id + '/pay');
    UI.toast('支付成功');
    Router.navigate('/order-detail?id=' + id);
  } catch (e) { UI.toast(e.message); }
}

// ========== 售后 ==========
Router.register('/after-sale', async (params) => {
  return `
    ${UI.renderNavbar('申请售后', true)}
    <div class="card">
      <div class="form-group">
        <label class="form-label">售后类型</label>
        <select class="form-input" id="asType">
          <option value="refund">退款</option>
          <option value="return">退货</option>
          <option value="exchange">换货</option>
          <option value="complaint">投诉</option>
        </select>
      </div>
      <div class="form-group">
        <label class="form-label">原因说明</label>
        <textarea class="form-input" id="asReason" rows="4" placeholder="请描述问题"></textarea>
      </div>
      <button class="btn btn-primary btn-block" onclick="submitAfterSale(${params.order_id})">提交</button>
    </div>
  `;
});

async function submitAfterSale(orderId) {
  try {
    await Http.post('/client/after-sales', {
      order_id: orderId,
      type: document.getElementById('asType').value,
      reason: document.getElementById('asReason').value,
    });
    UI.toast('售后申请已提交');
    Router.navigate('/orders');
  } catch (e) { UI.toast(e.message); }
}

// ========== 客服 ==========
Router.register('/chat', async (params) => {
  return `
    ${UI.renderNavbar('在线客服', true)}
    <div class="card">
      <p style="color:var(--text-secondary);text-align:center;padding:40px 0">
        客服功能开发中，请先拨打商家电话或联系官方客服<br>
        <strong style="color:var(--primary)">400-888-6688</strong>
      </p>
    </div>
  `;
});

// ========== 我的 ==========
Router.register('/mine', async () => {
  const user = Store.get('user');
  if (!user) {
    return `
      ${UI.renderMineHeader(null)}
      <div class="mine-menu">
        <div class="mine-menu-item" onclick="Router.navigate('/login')">
          <span class="mi-icon">🔐</span><span class="mi-label">登录 / 注册</span><span class="mi-arrow">›</span>
        </div>
      </div>
      <div class="card" style="margin-top:12px;text-align:center">
        <p style="color:var(--text-secondary);font-size:13px">登录后可查看订单、申请售后</p>
        <button class="btn btn-primary" style="margin-top:16px" onclick="Router.navigate('/login')">立即登录</button>
      </div>`;
  }
  return `
    ${UI.renderMineHeader(user)}
    <div class="mine-menu">
      <div class="mine-menu-item" onclick="Router.navigate('/orders')">
        <span class="mi-icon">📋</span><span class="mi-label">我的订单</span><span class="mi-arrow">›</span>
      </div>
      <div class="mine-menu-item" onclick="Router.navigate('/cart')">
        <span class="mi-icon">🛒</span><span class="mi-label">购物车</span><span class="mi-arrow">›</span>
      </div>
      <div class="mine-menu-item" onclick="Router.navigate('/chat')">
        <span class="mi-icon">💬</span><span class="mi-label">官方客服</span><span class="mi-arrow">›</span>
      </div>
      <div class="mine-menu-item" onclick="Router.navigate('/city')">
        <span class="mi-icon">📍</span><span class="mi-label">切换城市</span><span class="mi-arrow">›</span>
      </div>
      <div class="mine-menu-item" onclick="logout()">
        <span class="mi-icon">🚪</span><span class="mi-label">退出登录</span><span class="mi-arrow">›</span>
      </div>
    </div>`;
});

// ========== 登录 ==========
Router.register('/login', () => {
  document.getElementById('tabbar').style.display = 'none';
  return UI.renderLoginPage('本地建材 · 一站购齐', 'doLogin', 'doRegister');
});

async function doLogin() {
  try {
    const data = await Http.post('/auth/login', {
      phone: document.getElementById('phone').value,
      password: document.getElementById('password').value,
      role: APP_CONFIG.ROLE,
    });
    localStorage.setItem('token', data.token);
    Store.setUser(data.user);
    document.getElementById('tabbar').style.display = 'flex';
    UI.toast('登录成功');
    Router.navigate('/mine');
  } catch (e) { UI.toast(e.message); }
}

async function doRegister() {
  const city = Store.get('city');
  try {
    const data = await Http.post('/auth/register', {
      phone: document.getElementById('phone').value,
      password: document.getElementById('password').value,
      nickname: '用户' + document.getElementById('phone').value.slice(-4),
      role: APP_CONFIG.ROLE,
      city_id: city?.id,
    });
    localStorage.setItem('token', data.token);
    Store.setUser(data.user);
    document.getElementById('tabbar').style.display = 'flex';
    UI.toast('注册成功');
    Router.navigate('/mine');
  } catch (e) { UI.toast(e.message); }
}

function logout() {
  Store.logout();
  UI.toast('已退出');
  Router.navigate('/mine');
}

// 启动
Router.start('/');
UI.updateCartBadge();
