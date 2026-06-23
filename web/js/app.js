// Home Manager PWA
(() => {
  const API_URL = '';
  const main = document.getElementById('main');
  const nav = document.getElementById('nav');
  const toast = document.getElementById('toast');

  let token = localStorage.getItem('hm_token') || '';
  let currentView = 'lists';

  // ─── Router ───────────────────────────────────────────────────────

  function render() {
    if (!token) {
      nav.hidden = true;
      renderLogin();
      return;
    }

    nav.hidden = false;
    updateNav();

    if (currentView === 'products') {
      renderProducts();
    } else {
      renderLists();
    }
  }

  function navigate(view) {
    currentView = view;
    render();
  }

  function updateNav() {
    document.querySelectorAll('.nav-btn').forEach(btn => {
      btn.classList.toggle('active', btn.dataset.view === currentView);
    });
  }

  document.querySelectorAll('.nav-btn').forEach(btn => {
    btn.addEventListener('click', () => navigate(btn.dataset.view));
  });

  // ─── API helpers ──────────────────────────────────────────────────

  async function api(method, path, body) {
    const opts = {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
    };
    if (token) opts.headers['Authorization'] = `Bearer ${token}`;
    if (body) opts.body = JSON.stringify(body);

    const res = await fetch(`${API_URL}${path}`, opts);
    const data = await res.json().catch(() => ({}));

    if (!res.ok) {
      throw new Error(data.error || `Error ${res.status}`);
    }
    return data;
  }

  // ─── Login view ───────────────────────────────────────────────────

  function renderLogin() {
    main.innerHTML = `
      <div class="card" style="margin-top: 2rem;">
        <h2 class="card-title">Iniciar sesión</h2>
        <form id="login-form">
          <div class="form-group">
            <label for="email">Email</label>
            <input type="email" id="email" class="form-control" value="owner@home.local" required>
          </div>
          <div class="form-group">
            <label for="password">Contraseña</label>
            <input type="password" id="password" class="form-control" required>
          </div>
          <button type="submit" class="btn btn-primary">Entrar</button>
        </form>
      </div>
    `;

    document.getElementById('login-form').addEventListener('submit', async (e) => {
      e.preventDefault();
      const email = document.getElementById('email').value;
      const password = document.getElementById('password').value;

      try {
        const data = await api('POST', '/api/v1/auth/login', { email, password });
        token = data.token;
        localStorage.setItem('hm_token', token);
        render();
      } catch (err) {
        showToast(err.message);
      }
    });
  }

  // ─── Products view ────────────────────────────────────────────────

  async function renderProducts() {
    main.innerHTML = `
      <div class="card">
        <h2 class="card-title">Nuevo producto</h2>
        <form id="product-form">
          <div class="form-group">
            <input type="text" id="prod-name" class="form-control" placeholder="Nombre" required>
          </div>
          <div class="form-group">
            <input type="text" id="prod-unit" class="form-control" placeholder="Unidad (g, kg, ud...)" required>
          </div>
          <div class="form-group">
            <input type="text" id="prod-category" class="form-control" placeholder="Categoría">
          </div>
          <button type="submit" class="btn btn-primary">Añadir producto</button>
        </form>
      </div>
      <div id="products-list"></div>
    `;

    document.getElementById('product-form').addEventListener('submit', async (e) => {
      e.preventDefault();
      const name = document.getElementById('prod-name').value;
      const unit = document.getElementById('prod-unit').value;
      const category = document.getElementById('prod-category').value;

      try {
        await api('POST', '/api/v1/products', { name, unit, category });
        showToast('Producto añadido');
        renderProducts();
      } catch (err) {
        showToast(err.message);
      }
    });

    try {
      const data = await api('GET', '/api/v1/products?limit=100');
      const container = document.getElementById('products-list');
      if (!data.data || data.data.length === 0) {
        container.innerHTML = '<div class="empty-state">No hay productos aún</div>';
        return;
      }

      container.innerHTML = data.data.map(p => `
        <div class="list-item">
          <div class="list-item-info">
            <p class="list-item-title">${escapeHtml(p.name)}</p>
            <p class="list-item-subtitle">${escapeHtml(p.unit)}${p.category ? ` · ${escapeHtml(p.category)}` : ''}${p.latest_price ? ` · Último: ${p.latest_price.toFixed(2)} €` : ''}</p>
          </div>
          <button class="btn-icon" data-id="${p.id}" title="Añadir precio">💶</button>
        </div>
      `).join('');

      container.querySelectorAll('.btn-icon').forEach(btn => {
        btn.addEventListener('click', () => showPriceForm(btn.dataset.id));
      });
    } catch (err) {
      showToast(err.message);
    }
  }

  async function showPriceForm(productId) {
    const store = prompt('Nombre de la tienda:', 'Mercadona');
    if (!store) return;
    const amount = parseFloat(prompt('Precio:', '0'));
    if (!amount || amount <= 0) return;

    try {
      await api('POST', `/api/v1/products/${productId}/prices`, { store, amount });
      showToast('Precio añadido');
      renderProducts();
    } catch (err) {
      showToast(err.message);
    }
  }

  // ─── Lists view ───────────────────────────────────────────────────

  async function renderLists() {
    main.innerHTML = `
      <div class="card">
        <h2 class="card-title">Nueva lista</h2>
        <form id="list-form" class="flex gap-1">
          <input type="text" id="list-name" class="form-control" placeholder="Nombre de la lista" required style="flex:1">
          <button type="submit" class="btn btn-primary" style="width:auto">+</button>
        </form>
      </div>
      <div id="lists-list"></div>
    `;

    document.getElementById('list-form').addEventListener('submit', async (e) => {
      e.preventDefault();
      const name = document.getElementById('list-name').value;
      try {
        await api('POST', '/api/v1/lists', { name });
        showToast('Lista creada');
        renderLists();
      } catch (err) {
        showToast(err.message);
      }
    });

    try {
      const data = await api('GET', '/api/v1/lists?limit=100');
      const container = document.getElementById('lists-list');
      if (!data.data || data.data.length === 0) {
        container.innerHTML = '<div class="empty-state">No hay listas aún</div>';
        return;
      }

      container.innerHTML = data.data.map(l => `
        <div class="card">
          <div class="flex-between mb-1">
            <h3 class="card-title" style="margin:0">${escapeHtml(l.name)}</h3>
            <span class="total-badge">${l.estimated_total.toFixed(2)} €</span>
          </div>
          <p class="list-item-subtitle mb-1">${l.item_count} productos · ${l.status === 'completed' ? 'Completada' : 'Activa'}</p>
          <button class="btn btn-sm btn-primary" data-id="${l.id}">Ver / Añadir</button>
        </div>
      `).join('');

      container.querySelectorAll('.btn-primary').forEach(btn => {
        btn.addEventListener('click', () => renderListDetail(btn.dataset.id));
      });
    } catch (err) {
      showToast(err.message);
    }
  }

  async function renderListDetail(listId) {
    try {
      const data = await api('GET', `/api/v1/lists/${listId}`);
      const list = data.data.list;
      const items = data.data.items || [];

      main.innerHTML = `
        <div class="card">
          <div class="flex-between mb-1">
            <h2 class="card-title" style="margin:0">${escapeHtml(list.name)}</h2>
            <span class="total-badge">${list.estimated_total.toFixed(2)} €</span>
          </div>
          <button id="back-btn" class="btn btn-sm" style="background:var(--color-pistachio-light)">← Volver</button>
        </div>
        <div class="card">
          <h3 class="card-title">Añadir producto</h3>
          <form id="item-form">
            <div class="form-group">
              <select id="item-product" class="form-control" required>
                <option value="">Selecciona un producto</option>
              </select>
            </div>
            <div class="form-group">
              <input type="number" step="0.001" id="item-quantity" class="form-control" placeholder="Cantidad" required>
            </div>
            <div class="form-group">
              <input type="number" step="0.01" id="item-price" class="form-control" placeholder="Precio personalizado (opcional)">
            </div>
            <button type="submit" class="btn btn-primary">Añadir</button>
          </form>
        </div>
        <div id="items-list"></div>
      `;

      document.getElementById('back-btn').addEventListener('click', renderLists);

      const productsData = await api('GET', '/api/v1/products?limit=100');
      const select = document.getElementById('item-product');
      productsData.data.forEach(p => {
        const opt = document.createElement('option');
        opt.value = p.id;
        opt.textContent = `${p.name} (${p.unit})${p.latest_price ? ` - ${p.latest_price.toFixed(2)} €` : ''}`;
        select.appendChild(opt);
      });

      document.getElementById('item-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const product_id = document.getElementById('item-product').value;
        const quantity = parseFloat(document.getElementById('item-quantity').value);
        const priceInput = document.getElementById('item-price').value;
        const body = { product_id, quantity };
        if (priceInput) body.custom_price = parseFloat(priceInput);

        try {
          await api('POST', `/api/v1/lists/${listId}/items`, body);
          showToast('Producto añadido');
          renderListDetail(listId);
        } catch (err) {
          showToast(err.message);
        }
      });

      const itemsContainer = document.getElementById('items-list');
      if (items.length === 0) {
        itemsContainer.innerHTML = '<div class="empty-state">Lista vacía</div>';
      } else {
        itemsContainer.innerHTML = items.map(i => `
          <div class="list-item">
            <div class="list-item-info">
              <p class="list-item-title">${escapeHtml(i.product_name)} ${i.purchased ? '✅' : ''}</p>
              <p class="list-item-subtitle">${i.quantity} x ${i.unit_price.toFixed(2)} € = ${i.total.toFixed(2)} €</p>
            </div>
            <div class="list-item-actions">
              <button class="btn-icon toggle-purchase" data-id="${i.id}" data-purchased="${i.purchased}">${i.purchased ? '↩️' : '✅'}</button>
              <button class="btn-icon delete-item" data-id="${i.id}">🗑️</button>
            </div>
          </div>
        `).join('');

        itemsContainer.querySelectorAll('.toggle-purchase').forEach(btn => {
          btn.addEventListener('click', async () => {
            try {
              await api('PATCH', `/api/v1/lists/${listId}/items/${btn.dataset.id}`, { purchased: btn.dataset.purchased !== 'true' });
              renderListDetail(listId);
            } catch (err) {
              showToast(err.message);
            }
          });
        });

        itemsContainer.querySelectorAll('.delete-item').forEach(btn => {
          btn.addEventListener('click', async () => {
            if (!confirm('¿Eliminar este producto?')) return;
            try {
              await api('DELETE', `/api/v1/lists/${listId}/items/${btn.dataset.id}`);
              renderListDetail(listId);
            } catch (err) {
              showToast(err.message);
            }
          });
        });
      }
    } catch (err) {
      showToast(err.message);
    }
  }

  // ─── Utilities ────────────────────────────────────────────────────

  function showToast(message) {
    toast.textContent = message;
    toast.hidden = false;
    setTimeout(() => {
      toast.hidden = true;
    }, 3000);
  }

  function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  // ─── Init ─────────────────────────────────────────────────────────

  render();
})();
