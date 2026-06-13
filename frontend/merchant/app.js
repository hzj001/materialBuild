// 建材通 - 商家端
Http.init(APP_CONFIG.API_BASE);
Store.init();

function requireAuth() {
  if (!Store.get('user')) {
    Router.navigate('/login');
    return false;
  }
  return true;
}

// ========== 工作台 ==========
Router.register('/', async () => {
  if (!requireAuth()) return '';
  let orders = { list: [] };
  try { orders = await Http.get('/merchant/orders', { page_size: 5 }); } catch (e) { /* */ }

  const pending = (orders.list || []).filter(o => o.status <= 1).length;
  return `
    <div class="page">
      <div class="merchant-header">
        <h1>商家工作台</h1>
        <p>建材通 · 高效经营每一单</p>
      </div>
      <div class="stat-grid">
        <div class="stat-item"><div class="num">${pending}</div><div class="label">待处理</div></div>
        <div class="stat-item"><div class="num">${(orders.list || []).length}</div><div class="label">近期订单</div></div>
        <div class="stat-item"><div class="num">★</div><div class="label">店铺运营</div></div>
      </div>
      <div class="card menu-list">
        <div class="menu-item" onclick="Router.navigate('/product-add')">
          <span class="icon">➕</span>上架商品<span class="arrow">›</span>
        </div>
        <div class="menu-item" onclick="Router.navigate('/orders')">
          <span class="icon">📋</span>订单管理<span class="arrow">›</span>
        </div>
        <div class="menu-item" onclick="Router.navigate('/after-sales')">
          <span class="icon">🔧</span>售后处理<span class="arrow">›</span>
        </div>
        <div class="menu-item" onclick="Router.navigate('/delivery')">
          <span class="icon">🚚</span>配送设置<span class="arrow">›</span>
        </div>
        <div class="menu-item" onclick="Router.navigate('/promotions')">
          <span class="icon">🏷️</span>促销活动<span class="arrow">›</span>
        </div>
      </div>
    </div>
  `;
});

// ========== 商品管理 ==========
Router.register('/products', async () => {
  if (!requireAuth()) return '';
  const data = await Http.get('/merchant/products');
  const list = data.list || [];
  return `
    ${UI.renderNavbar('商品管理')}
    <div style="padding:12px">
      <button class="btn btn-primary btn-block" onclick="Router.navigate('/product-add')">+ 上架新商品</button>
    </div>
    ${list.length ? list.map(p => `
      <div class="product-card" onclick="Router.navigate('/product-edit?id=${p.id}')">
        <img src="${p.cover_image || 'https://via.placeholder.com/90'}" alt="">
        <div class="info">
          <div class="name">${p.name}</div>
          <div class="price">¥${p.sale_price} · 库存 ${p.stock}</div>
          <div class="merchant">${p.status === 1 ? '已上架' : '已下架'}</div>
        </div>
      </div>
    `).join('') : '<div class="empty"><div class="icon">📦</div><p>暂无商品</p></div>'}
  `;
});

Router.register('/product-add', () => {
  if (!requireAuth()) return '';
  return renderProductForm();
});

Router.register('/product-edit', async (params) => {
  if (!requireAuth()) return '';
  const data = await Http.get('/merchant/products');
  const p = (data.list || []).find(i => String(i.id) === params.id);
  return renderProductForm(p);
});

function renderProductForm(product) {
  const p = product || {};
  return `
    ${UI.renderNavbar(product ? '编辑商品' : '上架商品', true)}
    <div class="card">
      <div class="form-group">
        <label class="form-label">商品名称</label>
        <input class="form-input" id="pName" value="${p.name || ''}" placeholder="如：海螺水泥 P.O 42.5">
      </div>
      <div class="form-group">
        <label class="form-label">原价 (¥)</label>
        <input class="form-input" id="pPrice" type="number" value="${p.price || ''}" placeholder="0.00">
      </div>
      <div class="form-group">
        <label class="form-label">售价 (¥)</label>
        <input class="form-input" id="pSalePrice" type="number" value="${p.sale_price || ''}" placeholder="0.00">
      </div>
      <div class="form-group">
        <label class="form-label">单位</label>
        <input class="form-input" id="pUnit" value="${p.unit || '袋'}" placeholder="袋/吨/件">
      </div>
      <div class="form-group">
        <label class="form-label">库存</label>
        <input class="form-input" id="pStock" type="number" value="${p.stock || 0}">
      </div>
      <div class="form-group">
        <label class="form-label">封面图 URL</label>
        <input class="form-input" id="pCover" value="${p.cover_image || ''}" placeholder="图片地址">
      </div>
      <div class="form-group">
        <label class="form-label">商品描述</label>
        <textarea class="form-input" id="pDesc" rows="3">${p.description || ''}</textarea>
      </div>
      <button class="btn btn-primary btn-block" onclick="saveProduct(${p.id || 'null'})">保存</button>
    </div>
  `;
}

async function saveProduct(id) {
  const body = {
    name: document.getElementById('pName').value,
    price: parseFloat(document.getElementById('pPrice').value),
    sale_price: parseFloat(document.getElementById('pSalePrice').value),
    unit: document.getElementById('pUnit').value,
    stock: parseInt(document.getElementById('pStock').value),
    cover_image: document.getElementById('pCover').value,
    description: document.getElementById('pDesc').value,
  };
  try {
    if (id) {
      await Http.put('/merchant/products/' + id, body);
    } else {
      await Http.post('/merchant/products', body);
    }
    UI.toast('保存成功');
    Router.navigate('/products');
  } catch (e) { UI.toast(e.message); }
}

// ========== 订单管理 ==========
Router.register('/orders', async () => {
  if (!requireAuth()) return '';
  const data = await Http.get('/merchant/orders');
  const list = data.list || [];
  return `
    ${UI.renderNavbar('订单管理')}
    ${list.length ? list.map(o => `
      <div class="card">
        <div style="display:flex;justify-content:space-between;margin-bottom:8px">
          <span style="font-size:12px">${o.order_no}</span>
          ${UI.orderStatusTag(o.status)}
        </div>
        ${(o.items || []).map(i => `<div>${i.product_name} × ${i.quantity}</div>`).join('')}
        <div style="margin-top:8px;font-weight:600">¥${o.pay_amount}</div>
        ${o.delivery ? `<div style="font-size:12px;color:var(--text-secondary);margin-top:4px">${o.delivery.receiver_name} ${o.delivery.receiver_phone}<br>${o.delivery.address}</div>` : ''}
        <div style="margin-top:8px;display:flex;gap:8px">
          ${o.status === 1 ? `<button class="btn btn-sm btn-primary" onclick="updateOrderStatus(${o.id},2)">开始配送</button>` : ''}
          ${o.status === 2 ? `<button class="btn btn-sm btn-primary" onclick="updateOrderStatus(${o.id},3)">确认送达</button>` : ''}
        </div>
      </div>
    `).join('') : '<div class="empty"><div class="icon">📋</div><p>暂无订单</p></div>'}
  `;
});

async function updateOrderStatus(id, status) {
  try {
    await Http.put('/merchant/orders/' + id + '/status', { status });
    UI.toast('状态已更新');
    Router.navigate('/orders');
  } catch (e) { UI.toast(e.message); }
}

// ========== 售后 ==========
Router.register('/after-sales', async () => {
  if (!requireAuth()) return '';
  const data = await Http.get('/merchant/after-sales');
  const list = data.list || [];
  return `
    ${UI.renderNavbar('售后处理', true)}
    ${list.length ? list.map(a => `
      <div class="card">
        <div>类型：${a.type} · 订单 #${a.order_id}</div>
        <div style="margin:8px 0">${a.reason}</div>
        <div style="display:flex;gap:8px;margin-top:8px">
          <button class="btn btn-sm btn-primary" onclick="handleAfterSale(${a.id},2,'已同意')">同意</button>
          <button class="btn btn-sm btn-outline" onclick="handleAfterSale(${a.id},3,'已拒绝')">拒绝</button>
        </div>
      </div>
    `).join('') : '<div class="empty"><div class="icon">✅</div><p>暂无售后工单</p></div>'}
  `;
});

async function handleAfterSale(id, status, reply) {
  try {
    await Http.put('/merchant/after-sales/' + id, { status, reply });
    UI.toast('已处理');
    Router.navigate('/after-sales');
  } catch (e) { UI.toast(e.message); }
}

// ========== 配送设置 ==========
Router.register('/delivery', async () => {
  if (!requireAuth()) return '';
  const profile = await Http.get('/merchant/profile');
  return `
    ${UI.renderNavbar('配送设置', true)}
    <div class="card">
      <div class="form-group">
        <label class="form-label">配送半径 (km)</label>
        <input class="form-input" value="${profile.delivery_radius}" disabled>
      </div>
      <div class="form-group">
        <label class="form-label">起送金额 (¥)</label>
        <input class="form-input" value="${profile.min_order_amount}" disabled>
      </div>
      <p style="font-size:12px;color:var(--text-secondary)">配送参数修改功能将在后续版本开放</p>
    </div>
  `;
});

// ========== 促销 ==========
Router.register('/promotions', () => {
  if (!requireAuth()) return '';
  return `
    ${UI.renderNavbar('促销活动', true)}
    <div class="empty">
      <div class="icon">🏷️</div>
      <p>促销管理功能开发中</p>
      <p style="font-size:12px;margin-top:8px">支持全店折扣、单品特价、限时活动</p>
    </div>
  `;
});

// ========== 店铺 ==========
Router.register('/mine', async () => {
  if (!requireAuth()) return '';
  let profile = null;
  try { profile = await Http.get('/merchant/profile'); } catch (e) { /* */ }

  return `
    ${UI.renderNavbar('我的店铺')}
    ${profile ? `
      <div class="card">
        <div style="font-size:18px;font-weight:600">${profile.name}</div>
        <div style="font-size:12px;color:var(--text-secondary);margin-top:4px">${profile.address}</div>
        <div style="margin-top:8px">📞 ${profile.contact_phone}</div>
        <div style="margin-top:4px">状态：${UI.merchantStatusText(profile.status)}</div>
      </div>
    ` : '<div class="card"><p>店铺信息待完善</p></div>'}
    <div class="card menu-list">
      <div class="menu-item" onclick="Router.navigate('/after-sales')">
        <span class="icon">🔧</span>售后管理<span class="arrow">›</span>
      </div>
      <div class="menu-item" onclick="logout()">
        <span class="icon">🚪</span>退出登录<span class="arrow">›</span>
      </div>
    </div>
  `;
});

// ========== 登录 ==========
Router.register('/login', () => {
  document.getElementById('tabbar').style.display = 'none';
  return UI.renderLoginPage('建材商家 · 入驻经营', 'doLogin', 'doRegister');
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
    Router.navigate('/');
  } catch (e) { UI.toast(e.message); }
}

async function doRegister() {
  try {
    const data = await Http.post('/auth/register', {
      phone: document.getElementById('phone').value,
      password: document.getElementById('password').value,
      nickname: '商家' + document.getElementById('phone').value.slice(-4),
      role: APP_CONFIG.ROLE,
    });
    localStorage.setItem('token', data.token);
    Store.setUser(data.user);
    document.getElementById('tabbar').style.display = 'flex';
    UI.toast('注册成功，等待管理员审核');
    Router.navigate('/');
  } catch (e) { UI.toast(e.message); }
}

function logout() {
  Store.logout();
  UI.toast('已退出');
  Router.navigate('/login');
}

Router.start('/');
