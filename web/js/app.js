// Home Manager PWA
(() => {
  const API_URL = '';
  const main = document.getElementById('main');
  const nav = document.getElementById('nav');
  const toast = document.getElementById('toast');

  let token = localStorage.getItem('hm_token') || '';
  let currentUser = null;
  let currentView = localStorage.getItem('hm_view') || 'lists';

  // ─── Router ───────────────────────────────────────────────────────

  async function render() {
    if (!token) {
      nav.hidden = true;
      renderLogin();
      return;
    }

    if (!currentUser) {
      try {
        const data = await api('GET', '/api/v1/me');
        currentUser = data;
      } catch (err) {
        showToast('Sesión inválida');
        logout();
        return;
      }
    }

    nav.hidden = false;
    updateNav();

    switch (currentView) {
      case 'products':
        renderProducts();
        break;
      case 'expenses':
        renderExpenses();
        break;
      case 'incomes':
        renderIncomes();
        break;
      case 'dashboard':
        renderDashboard();
        break;
      default:
        renderLists();
    }
  }

  function navigate(view) {
    currentView = view;
    localStorage.setItem('hm_view', view);
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

  function logout() {
    token = '';
    currentUser = null;
    localStorage.removeItem('hm_token');
    render();
  }

  // ─── Auth views ───────────────────────────────────────────────────

  function renderLogin() {
    main.innerHTML = `
      <div class="card" style="margin-top: 2rem;">
        <h2 class="card-title">Iniciar sesión</h2>
        <form id="login-form">
          <div class="form-group">
            <label for="email">Email</label>
            <input type="email" id="email" class="form-control" required>
          </div>
          <div class="form-group">
            <label for="password">Contraseña</label>
            <input type="password" id="password" class="form-control" required>
          </div>
          <button type="submit" class="btn btn-primary">Entrar</button>
        </form>
        <p class="text-center mt-1">
          <a href="#" id="show-register">Crear cuenta</a>
        </p>
      </div>
    `;

    document.getElementById('login-form').addEventListener('submit', async (e) => {
      e.preventDefault();
      const email = document.getElementById('email').value;
      const password = document.getElementById('password').value;

      try {
        const data = await api('POST', '/api/v1/auth/login', { email, password });
        token = data.token;
        currentUser = data.user;
        localStorage.setItem('hm_token', token);
        render();
      } catch (err) {
        showToast(err.message);
      }
    });

    document.getElementById('show-register').addEventListener('click', (e) => {
      e.preventDefault();
      renderRegister();
    });
  }

  function renderRegister() {
    main.innerHTML = `
      <div class="card" style="margin-top: 2rem;">
        <h2 class="card-title">Crear cuenta</h2>
        <form id="register-form">
          <div class="form-group">
            <label for="email">Email</label>
            <input type="email" id="email" class="form-control" required>
          </div>
          <div class="form-group">
            <label for="password">Contraseña</label>
            <input type="password" id="password" class="form-control" minlength="8" required>
          </div>
          <button type="submit" class="btn btn-primary">Registrarme</button>
        </form>
        <p class="text-center mt-1">
          <a href="#" id="show-login">Ya tengo cuenta</a>
        </p>
      </div>
    `;

    document.getElementById('register-form').addEventListener('submit', async (e) => {
      e.preventDefault();
      const email = document.getElementById('email').value;
      const password = document.getElementById('password').value;

      try {
        const data = await api('POST', '/api/v1/auth/register', { email, password });
        token = data.token;
        currentUser = data.user;
        localStorage.setItem('hm_token', token);
        showToast('Cuenta creada');
        render();
      } catch (err) {
        showToast(err.message);
      }
    });

    document.getElementById('show-login').addEventListener('click', (e) => {
      e.preventDefault();
      renderLogin();
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
    const storeName = prompt('Nombre de la tienda:', 'Mercadona');
    if (!storeName) return;
    const amount = parseFloat(prompt('Precio:', '0'));
    if (!amount || amount <= 0) return;

    try {
      await api('POST', `/api/v1/products/${productId}/prices`, { store: storeName, amount });
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

  // ─── Expenses view ────────────────────────────────────────────────

  async function renderExpenses() {
    const today = new Date().toISOString().split('T')[0];
    main.innerHTML = `
      <div class="card">
        <h2 class="card-title">Nuevo gasto</h2>
        <form id="expense-form">
          <div class="form-group">
            <input type="text" id="exp-description" class="form-control" placeholder="Descripción" required>
          </div>
          <div class="form-group">
            <input type="number" step="0.01" id="exp-amount" class="form-control" placeholder="Importe (€)" required>
          </div>
          <div class="form-group">
            <input type="text" id="exp-category" class="form-control" placeholder="Categoría (supermercado, luz...)">
          </div>
          <div class="form-group">
            <label for="exp-visibility">Visibilidad</label>
            <select id="exp-visibility" class="form-control">
              <option value="private">Privado (solo yo)</option>
              <option value="shared">Compartido</option>
            </select>
          </div>
          <div class="form-group" id="split-group" hidden>
            <label for="exp-split">Tu porcentaje (%)</label>
            <input type="number" step="0.01" min="0" max="100" id="exp-split" class="form-control" value="50">
          </div>
          <div class="form-group">
            <label for="exp-date">Fecha</label>
            <input type="date" id="exp-date" class="form-control" value="${today}" required>
          </div>
          <div class="form-group">
            <label class="flex gap-1" style="align-items:center">
              <input type="checkbox" id="exp-recurring"> Recurrente mensual
            </label>
          </div>
          <button type="submit" class="btn btn-primary">Añadir gasto</button>
        </form>
      </div>
      <div class="card">
        <div class="flex gap-1">
          <select id="exp-filter" class="form-control">
            <option value="">Todos</option>
            <option value="shared">Compartidos</option>
            <option value="private">Privados</option>
          </select>
          <input type="month" id="exp-month" class="form-control" value="${today.slice(0, 7)}">
        </div>
      </div>
      <div id="expenses-list"></div>
    `;

    const visibilitySelect = document.getElementById('exp-visibility');
    const splitGroup = document.getElementById('split-group');
    visibilitySelect.addEventListener('change', () => {
      splitGroup.hidden = visibilitySelect.value !== 'shared';
    });

    document.getElementById('expense-form').addEventListener('submit', async (e) => {
      e.preventDefault();
      const body = {
        description: document.getElementById('exp-description').value,
        amount: parseFloat(document.getElementById('exp-amount').value),
        category: document.getElementById('exp-category').value,
        visibility: document.getElementById('exp-visibility').value,
        split_percentage: parseFloat(document.getElementById('exp-split').value) || 50,
        expense_date: document.getElementById('exp-date').value,
        is_recurring: document.getElementById('exp-recurring').checked,
      };

      try {
        await api('POST', '/api/v1/expenses', body);
        showToast('Gasto añadido');
        renderExpenses();
      } catch (err) {
        showToast(err.message);
      }
    });

    const filterSelect = document.getElementById('exp-filter');
    const monthInput = document.getElementById('exp-month');
    filterSelect.addEventListener('change', loadExpenses);
    monthInput.addEventListener('change', loadExpenses);

    await loadExpenses();
  }

  async function loadExpenses() {
    const filter = document.getElementById('exp-filter').value;
    const month = document.getElementById('exp-month').value;
    const { from, to } = getMonthBounds(month);

    try {
      const params = new URLSearchParams();
      params.set('from', from);
      params.set('to', to);
      if (filter) params.set('visibility', filter);

      const data = await api('GET', `/api/v1/expenses?${params.toString()}`);
      const container = document.getElementById('expenses-list');
      if (!data.data || data.data.length === 0) {
        container.innerHTML = '<div class="empty-state">No hay gastos este mes</div>';
        return;
      }

      const total = data.data.reduce((sum, e) => sum + e.amount, 0);
      const grouped = groupByDay(data.data, 'expense_date');

      container.innerHTML = `
        <div class="flex-between mb-1">
          <span class="list-item-subtitle">Total del mes: <strong>${total.toFixed(2)} €</strong></span>
        </div>
        ${grouped.map(([day, items]) => {
          const dayTotal = items.reduce((sum, e) => sum + e.amount, 0);
          return `
            <div class="day-group">
              <div class="day-header">
                <span class="day-label">${formatDayLabel(day)}</span>
                <span class="day-total">${dayTotal.toFixed(2)} €</span>
              </div>
              ${items.map(e => `
                <div class="list-item day-item">
                  <div class="list-item-info">
                    <p class="list-item-title">${escapeHtml(e.description)} ${visibilityBadge(e.visibility)}</p>
                    <p class="list-item-subtitle">${e.category ? `${escapeHtml(e.category)} · ` : ''}${e.is_recurring ? '🔄 · ' : ''}${e.visibility === 'shared' ? `${e.split_percentage}% tuyo` : ''}</p>
                  </div>
                  <div class="list-item-actions">
                    <span class="total-badge">${e.amount.toFixed(2)} €</span>
                    ${e.user_id === (currentUser?.id || currentUser?.user_id) ? `<button class="btn-icon delete-expense" data-id="${e.id}">🗑️</button>` : ''}
                  </div>
                </div>
              `).join('')}
            </div>
          `;
        }).join('')}
      `;

      container.querySelectorAll('.delete-expense').forEach(btn => {
        btn.addEventListener('click', async () => {
          if (!confirm('¿Eliminar este gasto?')) return;
          try {
            await api('DELETE', `/api/v1/expenses/${btn.dataset.id}`);
            showToast('Gasto eliminado');
            loadExpenses();
          } catch (err) {
            showToast(err.message);
          }
        });
      });
    } catch (err) {
      showToast(err.message);
    }
  }

  // ─── Incomes view ─────────────────────────────────────────────────

  async function renderIncomes() {
    const today = new Date().toISOString().split('T')[0];
    main.innerHTML = `
      <div class="card">
        <h2 class="card-title">Nuevo ingreso</h2>
        <form id="income-form">
          <div class="form-group">
            <input type="text" id="inc-description" class="form-control" placeholder="Descripción" required>
          </div>
          <div class="form-group">
            <input type="number" step="0.01" id="inc-amount" class="form-control" placeholder="Importe (€)" required>
          </div>
          <div class="form-group">
            <input type="text" id="inc-category" class="form-control" placeholder="Categoría (nómina, extra...)">
          </div>
          <div class="form-group">
            <label for="inc-visibility">Visibilidad</label>
            <select id="inc-visibility" class="form-control">
              <option value="private">Privado (solo yo)</option>
              <option value="shared">Compartido</option>
            </select>
          </div>
          <div class="form-group">
            <label for="inc-date">Fecha</label>
            <input type="date" id="inc-date" class="form-control" value="${today}" required>
          </div>
          <div class="form-group">
            <label class="flex gap-1" style="align-items:center">
              <input type="checkbox" id="inc-recurring"> Recurrente mensual
            </label>
          </div>
          <button type="submit" class="btn btn-primary">Añadir ingreso</button>
        </form>
      </div>
      <div class="card">
        <div class="flex gap-1">
          <select id="inc-filter" class="form-control">
            <option value="">Todos</option>
            <option value="shared">Compartidos</option>
            <option value="private">Privados</option>
          </select>
          <input type="month" id="inc-month" class="form-control" value="${today.slice(0, 7)}">
        </div>
      </div>
      <div id="incomes-list"></div>
    `;

    document.getElementById('income-form').addEventListener('submit', async (e) => {
      e.preventDefault();
      const body = {
        description: document.getElementById('inc-description').value,
        amount: parseFloat(document.getElementById('inc-amount').value),
        category: document.getElementById('inc-category').value,
        visibility: document.getElementById('inc-visibility').value,
        income_date: document.getElementById('inc-date').value,
        is_recurring: document.getElementById('inc-recurring').checked,
      };

      try {
        await api('POST', '/api/v1/incomes', body);
        showToast('Ingreso añadido');
        renderIncomes();
      } catch (err) {
        showToast(err.message);
      }
    });

    const filterSelect = document.getElementById('inc-filter');
    const monthInput = document.getElementById('inc-month');
    filterSelect.addEventListener('change', loadIncomes);
    monthInput.addEventListener('change', loadIncomes);

    await loadIncomes();
  }

  async function loadIncomes() {
    const filter = document.getElementById('inc-filter').value;
    const month = document.getElementById('inc-month').value;
    const { from, to } = getMonthBounds(month);

    try {
      const params = new URLSearchParams();
      params.set('from', from);
      params.set('to', to);
      if (filter) params.set('visibility', filter);

      const data = await api('GET', `/api/v1/incomes?${params.toString()}`);
      const container = document.getElementById('incomes-list');
      if (!data.data || data.data.length === 0) {
        container.innerHTML = '<div class="empty-state">No hay ingresos este mes</div>';
        return;
      }

      const total = data.data.reduce((sum, i) => sum + i.amount, 0);
      const grouped = groupByDay(data.data, 'income_date');

      container.innerHTML = `
        <div class="flex-between mb-1">
          <span class="list-item-subtitle">Total del mes: <strong>${total.toFixed(2)} €</strong></span>
        </div>
        ${grouped.map(([day, items]) => {
          const dayTotal = items.reduce((sum, i) => sum + i.amount, 0);
          return `
            <div class="day-group">
              <div class="day-header">
                <span class="day-label">${formatDayLabel(day)}</span>
                <span class="day-total">${dayTotal.toFixed(2)} €</span>
              </div>
              ${items.map(i => `
                <div class="list-item day-item">
                  <div class="list-item-info">
                    <p class="list-item-title">${escapeHtml(i.description)} ${visibilityBadge(i.visibility)}</p>
                    <p class="list-item-subtitle">${i.category ? `${escapeHtml(i.category)} · ` : ''}${i.is_recurring ? '🔄 · ' : ''}</p>
                  </div>
                  <div class="list-item-actions">
                    <span class="total-badge">${i.amount.toFixed(2)} €</span>
                    ${i.user_id === (currentUser?.id || currentUser?.user_id) ? `<button class="btn-icon delete-income" data-id="${i.id}">🗑️</button>` : ''}
                  </div>
                </div>
              `).join('')}
            </div>
          `;
        }).join('')}
      `;

      container.querySelectorAll('.delete-income').forEach(btn => {
        btn.addEventListener('click', async () => {
          if (!confirm('¿Eliminar este ingreso?')) return;
          try {
            await api('DELETE', `/api/v1/incomes/${btn.dataset.id}`);
            showToast('Ingreso eliminado');
            loadIncomes();
          } catch (err) {
            showToast(err.message);
          }
        });
      });
    } catch (err) {
      showToast(err.message);
    }
  }

  // ─── Dashboard view ───────────────────────────────────────────────

  async function renderDashboard() {
    const currentYear = new Date().getFullYear();
    main.innerHTML = `
      <div class="card">
        <h2 class="card-title">Resumen mensual</h2>
        <div class="form-group">
          <label for="dash-month">Mes</label>
          <input type="month" id="dash-month" class="form-control" value="${new Date().toISOString().slice(0, 7)}">
        </div>
        <div id="summary-cards" class="grid-3"></div>
      </div>
      <div class="card">
        <h2 class="card-title">Evolución anual</h2>
        <div class="form-group">
          <select id="dash-year" class="form-control">
            <option value="${currentYear}" selected>${currentYear}</option>
            <option value="${currentYear - 1}">${currentYear - 1}</option>
          </select>
        </div>
        <div class="chart-container">
          <canvas id="bar-chart"></canvas>
        </div>
        <div class="chart-container mt-2">
          <canvas id="line-chart"></canvas>
        </div>
      </div>
    `;

    document.getElementById('dash-month').addEventListener('change', loadSummary);
    document.getElementById('dash-year').addEventListener('change', loadMonthly);

    await loadSummary();
    await loadMonthly();
  }

  async function loadSummary() {
    const month = document.getElementById('dash-month').value;
    try {
      const data = await api('GET', `/api/v1/dashboard/summary?month=${month}`);
      const summary = data.data;
      const container = document.getElementById('summary-cards');
      container.innerHTML = `
        <div class="summary-card income">
          <p class="summary-label">Tus ingresos</p>
          <p class="summary-value">${summary.income_total.toFixed(2)} €</p>
        </div>
        <div class="summary-card expense">
          <p class="summary-label">Tus gastos</p>
          <p class="summary-value">${summary.expense_total.toFixed(2)} €</p>
          <p class="summary-detail" style="font-size:0.75rem;color:#6B7280;margin-top:0.25rem;">
            Compartidos: ${summary.shared_expense_total.toFixed(2)} € · Privados: ${summary.private_expense_total.toFixed(2)} €
          </p>
        </div>
        <div class="summary-card ${summary.balance >= 0 ? 'positive' : 'negative'}">
          <p class="summary-label">Balance</p>
          <p class="summary-value">${summary.balance.toFixed(2)} €</p>
        </div>
      `;
    } catch (err) {
      showToast(err.message);
    }
  }

  async function loadMonthly() {
    const year = document.getElementById('dash-year').value;
    try {
      const data = await api('GET', `/api/v1/dashboard/monthly?year=${year}`);
      const months = data.data;
      drawBarChart(months);
      drawLineChart(months);
    } catch (err) {
      showToast(err.message);
    }
  }

  function drawBarChart(months) {
    const canvas = document.getElementById('bar-chart');
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    const labels = months.map(m => m.month.slice(5));
    const income = months.map(m => m.income_total);
    const expense = months.map(m => m.expense_total);
    drawGroupedBarChart(canvas, ctx, labels, income, expense, 'Ingresos vs Gastos (€)');
  }

  function drawLineChart(months) {
    const canvas = document.getElementById('line-chart');
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    const labels = months.map(m => m.month.slice(5));
    const balance = months.map(m => m.balance);
    drawLineChartCanvas(canvas, ctx, labels, balance, 'Balance mensual (€)');
  }

  // ─── Canvas charts ────────────────────────────────────────────────

  function setupCanvas(canvas) {
    const dpr = window.devicePixelRatio || 1;
    const rect = canvas.getBoundingClientRect();
    canvas.width = rect.width * dpr;
    canvas.height = rect.height * dpr;
    const ctx = canvas.getContext('2d');
    ctx.scale(dpr, dpr);
    return { width: rect.width, height: rect.height };
  }

  function drawGroupedBarChart(canvas, ctx, labels, income, expense, title) {
    const { width, height } = setupCanvas(canvas);
    ctx.clearRect(0, 0, width, height);

    const padding = { top: 40, right: 20, bottom: 40, left: 50 };
    const chartWidth = width - padding.left - padding.right;
    const chartHeight = height - padding.top - padding.bottom;

    const maxValue = Math.max(...income, ...expense, 1);

    // Title
    ctx.fillStyle = '#1F2937';
    ctx.font = 'bold 14px sans-serif';
    ctx.textAlign = 'center';
    ctx.fillText(title, width / 2, 20);

    // Axes
    ctx.strokeStyle = '#D9EDC8';
    ctx.lineWidth = 1;
    ctx.beginPath();
    ctx.moveTo(padding.left, padding.top);
    ctx.lineTo(padding.left, height - padding.bottom);
    ctx.lineTo(width - padding.right, height - padding.bottom);
    ctx.stroke();

    // Y-axis labels
    ctx.fillStyle = '#6B7280';
    ctx.font = '11px sans-serif';
    ctx.textAlign = 'right';
    for (let i = 0; i <= 5; i++) {
      const value = (maxValue / 5) * i;
      const y = height - padding.bottom - (chartHeight / 5) * i;
      ctx.fillText(value.toFixed(0), padding.left - 8, y + 4);
      if (i > 0) {
        ctx.strokeStyle = '#F3F4F6';
        ctx.beginPath();
        ctx.moveTo(padding.left, y);
        ctx.lineTo(width - padding.right, y);
        ctx.stroke();
      }
    }

    // Bars
    const groupWidth = chartWidth / labels.length;
    const barWidth = groupWidth * 0.35;
    const gap = groupWidth * 0.1;

    labels.forEach((label, i) => {
      const x = padding.left + i * groupWidth + gap;
      const incomeHeight = (income[i] / maxValue) * chartHeight;
      const expenseHeight = (expense[i] / maxValue) * chartHeight;

      // Income bar
      ctx.fillStyle = '#22C55E';
      ctx.fillRect(x, height - padding.bottom - incomeHeight, barWidth, incomeHeight);

      // Expense bar
      ctx.fillStyle = '#EF4444';
      ctx.fillRect(x + barWidth + 2, height - padding.bottom - expenseHeight, barWidth, expenseHeight);

      // Label
      ctx.fillStyle = '#6B7280';
      ctx.font = '10px sans-serif';
      ctx.textAlign = 'center';
      ctx.fillText(label, x + barWidth, height - padding.bottom + 15);
    });

    // Legend
    const legendY = 18;
    ctx.fillStyle = '#22C55E';
    ctx.fillRect(width - 110, legendY - 8, 10, 10);
    ctx.fillStyle = '#1F2937';
    ctx.font = '11px sans-serif';
    ctx.textAlign = 'left';
    ctx.fillText('Ingresos', width - 95, legendY);

    ctx.fillStyle = '#EF4444';
    ctx.fillRect(width - 50, legendY - 8, 10, 10);
    ctx.fillStyle = '#1F2937';
    ctx.fillText('Gastos', width - 35, legendY);
  }

  function drawLineChartCanvas(canvas, ctx, labels, values, title) {
    const { width, height } = setupCanvas(canvas);
    ctx.clearRect(0, 0, width, height);

    const padding = { top: 40, right: 20, bottom: 40, left: 50 };
    const chartWidth = width - padding.left - padding.right;
    const chartHeight = height - padding.top - padding.bottom;

    const minValue = Math.min(...values, 0);
    const maxValue = Math.max(...values, 1);
    const range = maxValue - minValue || 1;

    // Title
    ctx.fillStyle = '#1F2937';
    ctx.font = 'bold 14px sans-serif';
    ctx.textAlign = 'center';
    ctx.fillText(title, width / 2, 20);

    // Axes
    ctx.strokeStyle = '#D9EDC8';
    ctx.lineWidth = 1;
    ctx.beginPath();
    ctx.moveTo(padding.left, padding.top);
    ctx.lineTo(padding.left, height - padding.bottom);
    ctx.lineTo(width - padding.right, height - padding.bottom);
    ctx.stroke();

    // Y-axis labels
    ctx.fillStyle = '#6B7280';
    ctx.font = '11px sans-serif';
    ctx.textAlign = 'right';
    for (let i = 0; i <= 5; i++) {
      const value = minValue + (range / 5) * i;
      const y = height - padding.bottom - (chartHeight / 5) * i;
      ctx.fillText(value.toFixed(0), padding.left - 8, y + 4);
    }

    // Line
    const stepX = chartWidth / (labels.length - 1 || 1);
    ctx.strokeStyle = '#8FB87A';
    ctx.lineWidth = 3;
    ctx.beginPath();
    values.forEach((value, i) => {
      const x = padding.left + i * stepX;
      const y = height - padding.bottom - ((value - minValue) / range) * chartHeight;
      if (i === 0) ctx.moveTo(x, y);
      else ctx.lineTo(x, y);
    });
    ctx.stroke();

    // Points and labels
    values.forEach((value, i) => {
      const x = padding.left + i * stepX;
      const y = height - padding.bottom - ((value - minValue) / range) * chartHeight;

      ctx.fillStyle = '#B4D89E';
      ctx.beginPath();
      ctx.arc(x, y, 4, 0, Math.PI * 2);
      ctx.fill();

      ctx.fillStyle = '#6B7280';
      ctx.font = '10px sans-serif';
      ctx.textAlign = 'center';
      ctx.fillText(labels[i], x, height - padding.bottom + 15);
    });
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

  function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('es-ES', { day: '2-digit', month: '2-digit', year: 'numeric' });
  }

  function formatDayLabel(dateString) {
    const date = new Date(dateString);
    const today = new Date();
    const yesterday = new Date(today);
    yesterday.setDate(yesterday.getDate() - 1);

    const isSameDay = (a, b) =>
      a.getFullYear() === b.getFullYear() &&
      a.getMonth() === b.getMonth() &&
      a.getDate() === b.getDate();

    if (isSameDay(date, today)) return 'Hoy';
    if (isSameDay(date, yesterday)) return 'Ayer';

    return date.toLocaleDateString('es-ES', {
      weekday: 'long',
      day: 'numeric',
      month: 'long',
      year: 'numeric',
    });
  }

  function groupByDay(items, dateField) {
    const groups = new Map();
    items.forEach(item => {
      const day = item[dateField].slice(0, 10);
      if (!groups.has(day)) groups.set(day, []);
      groups.get(day).push(item);
    });
    return Array.from(groups.entries());
  }

  function getMonthBounds(month) {
    const [year, monthNum] = month.split('-').map(Number);
    const from = `${month}-01`;
    const toDate = new Date(year, monthNum, 0);
    const to = `${year}-${String(monthNum).padStart(2, '0')}-${String(toDate.getDate()).padStart(2, '0')}`;
    return { from, to };
  }

  function visibilityBadge(visibility) {
    if (visibility === 'shared') {
      return '<span class="badge badge-shared">Compartido</span>';
    }
    return '<span class="badge badge-private">Privado</span>';
  }

  // ─── Init ─────────────────────────────────────────────────────────

  render();
})();
