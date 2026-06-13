/**
 * 建材通 HTTP 请求封装
 */
const Http = {
  baseURL: '',

  init(baseURL) {
    this.baseURL = baseURL;
  },

  getToken() {
    return localStorage.getItem('token') || '';
  },

  async request(method, path, data) {
    const headers = { 'Content-Type': 'application/json' };
    const token = this.getToken();
    if (token) headers['Authorization'] = `Bearer ${token}`;

    const opts = { method, headers };
    if (data && method !== 'GET') {
      opts.body = JSON.stringify(data);
    }

    let url = this.baseURL + path;
    if (data && method === 'GET') {
      const params = new URLSearchParams();
      Object.entries(data).forEach(([k, v]) => {
        if (v !== undefined && v !== null && v !== '') params.append(k, v);
      });
      const qs = params.toString();
      if (qs) url += '?' + qs;
    }

    const res = await fetch(url, opts);
    const json = await res.json();

    if (json.code !== 0) {
      if (json.code === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        Router.navigate('/login');
      }
      throw new Error(json.message || '请求失败');
    }
    return json.data;
  },

  get(path, params) { return this.request('GET', path, params); },
  post(path, data) { return this.request('POST', path, data); },
  put(path, data) { return this.request('PUT', path, data); },
  delete(path) { return this.request('DELETE', path); },
};
