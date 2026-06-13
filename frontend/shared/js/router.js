/**
 * Hash 路由
 */
const Router = {
  routes: {},
  currentPath: '',

  register(path, handler) {
    this.routes[path] = handler;
  },

  navigate(path) {
    location.hash = path;
  },

  back() {
    history.back();
  },

  getParams() {
    const hash = location.hash.slice(1) || '/';
    const [path, query] = hash.split('?');
    const params = {};
    if (query) {
      new URLSearchParams(query).forEach((v, k) => { params[k] = v; });
    }
    return { path, params };
  },

  start(defaultPath = '/') {
    window.addEventListener('hashchange', () => this._render());
    if (!location.hash) location.hash = defaultPath;
    else this._render();
  },

  async _render() {
    const { path, params } = this.getParams();
    this.currentPath = path;

    const handler = this.routes[path];
    const app = document.getElementById('app');
    if (!app) return;

    if (handler) {
      app.innerHTML = '<div class="loading">加载中...</div>';
      try {
        app.innerHTML = await handler(params);
      } catch (e) {
        app.innerHTML = `<div class="empty"><div class="icon">⚠️</div><p>${e.message}</p></div>`;
      }
    } else {
      app.innerHTML = '<div class="empty"><div class="icon">🔍</div><p>页面不存在</p></div>';
    }

    this._updateTabbar();
    if (typeof UI !== 'undefined' && UI.updateCartBadge) UI.updateCartBadge();
  },

  _updateTabbar() {
    document.querySelectorAll('.tabbar-item').forEach(el => {
      const route = el.dataset.route;
      el.classList.toggle('active', this.currentPath === route);
    });
  },
};
