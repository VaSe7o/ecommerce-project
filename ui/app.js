window.addEventListener('DOMContentLoaded', () => {
  console.log("DOM loaded, attaching event listeners.");
  let currentRole = null;
  let currentCustomerName = "";
  let cartItems = [];
  let globalProductsCache = [];

  function showLoginPage() {
    document.getElementById('login-page').style.display = 'block';
    document.getElementById('owner-page').style.display = 'none';
    document.getElementById('customer-page').style.display = 'none';
  }
  function showOwnerPage() {
    document.getElementById('login-page').style.display = 'none';
    document.getElementById('owner-page').style.display = 'block';
    document.getElementById('customer-page').style.display = 'none';
  }
  function showCustomerPage() {
    document.getElementById('login-page').style.display = 'none';
    document.getElementById('owner-page').style.display = 'none';
    document.getElementById('customer-page').style.display = 'block';
    document.getElementById('customer-shop-page').style.display = 'block';
    document.getElementById('cart-page').style.display = 'none';
    document.getElementById('customer-welcome-name').textContent = currentCustomerName;
  }

  const roleSelect = document.getElementById('role-select');
  const ownerPasswordGroup = document.getElementById('owner-password-group');
  const customerNameGroup = document.getElementById('customer-name-group');
  const ownerPasswordInput = document.getElementById('owner-password');
  const customerNameInput = document.getElementById('customer-name');
  const loginBtn = document.getElementById('login-btn');

  const ownerLogoutBtn = document.getElementById('owner-logout-btn');
  const ownerAddProductForm = document.getElementById('owner-add-product-form');
  const ownerProductsTableBody = document.getElementById('owner-products-table').querySelector('tbody');
  const ownerOrdersTableBody = document.getElementById('owner-orders-table').querySelector('tbody');
  const ownerStatsTableBody = document.getElementById('owner-stats-table').querySelector('tbody');
  const sortNameBtn = document.getElementById('sort-name-btn');
  const sortQtyBtn = document.getElementById('sort-qty-btn');
  const sortPriceBtn = document.getElementById('sort-price-btn');
  const sortIdBtn = document.getElementById('sort-id-btn');

  const customerLogoutBtn = document.getElementById('customer-logout-btn');
  const customerListProductsBtn = document.getElementById('customer-list-products-btn');
  const customerProductsTableBody = document.getElementById('customer-products-table').querySelector('tbody');
  const goToCartBtn = document.getElementById('go-to-cart-btn');

  const cartTableBody = document.getElementById('cart-table').querySelector('tbody');
  const cartSubtotalEl = document.getElementById('cart-subtotal');
  const cartFinalTotalEl = document.getElementById('cart-final-total');
  const extraFeesMsgEl = document.getElementById('extra-fees-msg');
  const checkoutEmail = document.getElementById('checkout-email');
  const checkoutPhone = document.getElementById('checkout-phone');
  const checkoutAddress = document.getElementById('checkout-address');
  const checkoutPaymentMethod = document.getElementById('checkout-payment-method');
  const backToShopBtn = document.getElementById('back-to-shop-btn');
  const checkoutBtn = document.getElementById('checkout-btn');

  roleSelect.addEventListener('change', (e) => {
    console.log("Role changed to:", e.target.value);
    if (e.target.value === 'owner') {
      ownerPasswordGroup.style.display = 'block';
      customerNameGroup.style.display = 'none';
    } else {
      ownerPasswordGroup.style.display = 'none';
      customerNameGroup.style.display = 'block';
    }
  });

  loginBtn.addEventListener('click', () => {
    console.log("Login clicked");
    if (roleSelect.value === 'owner') {
      if (ownerPasswordInput.value === '12345') {
        currentRole = 'owner';
        console.log("Owner login success");
        showOwnerPage();
      } else {
        alert("Wrong owner password. Try 12345.");
      }
    } else {
      const nameVal = customerNameInput.value.trim();
      if (!nameVal) {
        alert("Please enter your name");
        return;
      }
      currentRole = 'customer';
      currentCustomerName = nameVal;
      console.log("Customer login success:", currentCustomerName);
      showCustomerPage();
      listProductsForCustomer();
    }
  });

  ownerLogoutBtn.addEventListener('click', () => {
    console.log("Owner logout clicked");
    currentRole = null;
    ownerPasswordInput.value = '';
    showLoginPage();
  });

  ownerAddProductForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    console.log("Owner add product submitted");
    const name = document.getElementById('owner-product-name').value.trim();
    const desc = document.getElementById('owner-product-desc').value.trim();
    const price = parseFloat(document.getElementById('owner-product-price').value || '0');
    const qty = parseInt(document.getElementById('owner-product-qty').value || '0', 10);
    try {
      const res = await fetch('/api/v1/products', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, description: desc, price, quantity: qty })
      });
      if (!res.ok) throw new Error('Failed to create product');
      const newProd = await res.json();
      alert('Product created:\n' + JSON.stringify(newProd, null, 2));
      ownerAddProductForm.reset();
    } catch (err) {
      alert('Error: ' + err.message);
    }
  });

  document.getElementById('owner-list-products-btn').addEventListener('click', async () => {
    console.log("Owner list products clicked");
    await listProductsForOwner();
  });

  if (sortNameBtn) sortNameBtn.addEventListener('click', () => { console.log("Sorting by name"); sortOwnerProducts('name'); });
  if (sortQtyBtn) sortQtyBtn.addEventListener('click', () => { console.log("Sorting by qty"); sortOwnerProducts('qty'); });
  if (sortPriceBtn) sortPriceBtn.addEventListener('click', () => { console.log("Sorting by price"); sortOwnerProducts('price'); });
  if (sortIdBtn) sortIdBtn.addEventListener('click', () => { console.log("Sorting by id"); sortOwnerProducts('id'); });

  async function listProductsForOwner() {
    ownerProductsTableBody.innerHTML = '';
    try {
      const res = await fetch('/api/v1/products');
      if (!res.ok) throw new Error('Cannot list products');
      const products = await res.json();
      globalProductsCache = products;
      renderOwnerProducts(globalProductsCache);
    } catch (err) {
      alert('Error listing products: ' + err.message);
    }
  }

  function renderOwnerProducts(products) {
    ownerProductsTableBody.innerHTML = '';
    products.forEach(p => {
      const tr = document.createElement('tr');
      tr.innerHTML = `
        <td>${p.id}</td>
        <td>${p.name}</td>
        <td>${p.description}</td>
        <td>${p.price}</td>
        <td>${p.quantity}</td>
        <td>
          <button class="btn btn-sm btn-danger delete-product-btn" data-id="${p.id}">
            Delete
          </button>
        </td>
      `;
      ownerProductsTableBody.appendChild(tr);
    });

    ownerProductsTableBody.querySelectorAll('.delete-product-btn').forEach(button => {
      button.addEventListener('click', async () => {
        const productId = button.getAttribute('data-id');
        if (confirm(`Are you sure you want to delete product ${productId}?`)) {
          try {
            const res = await fetch(`/api/v1/products/${productId}`, {
              method: 'DELETE'
            });
            if (!res.ok) {
              throw new Error('Failed to delete product');
            }
            alert(`Product ${productId} deleted successfully.`);
           
            globalProductsCache = globalProductsCache.filter(prod => prod.id.toString() !== productId);
            
            renderOwnerProducts(globalProductsCache);
          } catch (err) {
            alert('Error: ' + err.message);
          }
        }
      });
    });
  }

  function sortOwnerProducts(criteria) {
    const sorted = [...globalProductsCache];
    if (criteria === 'name') {
      sorted.sort((a, b) => a.name.localeCompare(b.name));
    } else if (criteria === 'qty') {
      sorted.sort((a, b) => a.quantity - b.quantity);
    } else if (criteria === 'price') {
      sorted.sort((a, b) => a.price - b.price);
    } else if (criteria === 'id') {
      sorted.sort((a, b) => a.id.toString().localeCompare(b.id.toString()));
    }
    renderOwnerProducts(sorted);
  }

  document.getElementById('owner-list-orders-btn').addEventListener('click', async () => {
    console.log("Owner list orders clicked");
    ownerOrdersTableBody.innerHTML = '';
    try {
      const res = await fetch('/api/v1/orders');
      if (!res.ok) throw new Error('Cannot list orders');
      const orders = await res.json();
      orders.forEach(o => {
        const itemsStr = o.items.map(it => `(${it.product_id} x ${it.quantity} @ ${it.price})`).join(', ');
        const tr = document.createElement('tr');
        tr.innerHTML = `
          <td>${o.id}</td>
          <td>${o.customer_name}</td>
          <td>${o.email}</td>
          <td>${o.phone}</td>
          <td>${o.address}</td>
          <td>${o.payment_method}</td>
          <td>${itemsStr}</td>
          <td>$${(o.final_total || 0).toFixed(2)}</td>
        `;
        ownerOrdersTableBody.appendChild(tr);
      });
    } catch (err) {
      alert('Error listing orders: ' + err.message);
    }
  });

  document.getElementById('owner-stats-btn').addEventListener('click', async () => {
    console.log("Owner stats clicked");
    ownerStatsTableBody.innerHTML = '';
    try {
      const res = await fetch('/api/v1/orders');
      if (!res.ok) throw new Error('Cannot fetch orders for stats');
      const orders = await res.json();
      const productTotals = {};
      orders.forEach(o => {
        o.items.forEach(it => {
          if (!productTotals[it.product_id]) productTotals[it.product_id] = 0;
          productTotals[it.product_id] += it.quantity;
        });
      });
      const sorted = Object.entries(productTotals).sort((a, b) => b[1] - a[1]);
      sorted.forEach(([pid, tot]) => {
        const tr = document.createElement('tr');
        tr.innerHTML = `<td>${pid}</td><td>${tot}</td>`;
        ownerStatsTableBody.appendChild(tr);
      });
    } catch (err) {
      alert('Error fetching stats: ' + err.message);
    }
  });

  customerLogoutBtn.addEventListener('click', () => {
    console.log("Customer logout clicked");
    currentRole = null;
    currentCustomerName = "";
    cartItems = [];
    showLoginPage();
  });

  customerListProductsBtn.addEventListener('click', () => {
    console.log("Customer list products clicked");
    listProductsForCustomer();
  });

  goToCartBtn.addEventListener('click', () => {
    console.log("Customer go to cart clicked");
    showCartPage();
  });

  async function listProductsForCustomer() {
    customerProductsTableBody.innerHTML = '';
    try {
      const res = await fetch('/api/v1/products');
      if (!res.ok) throw new Error('Cannot fetch products');
      const products = await res.json();
      products.forEach(p => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
          <td>${p.name}</td>
          <td>${p.description}</td>
          <td>$${p.price}</td>
          <td>${p.quantity}</td>
          <td><button class="btn btn-sm btn-primary" data-id="${p.id}">Add</button></td>
        `;
        customerProductsTableBody.appendChild(tr);
      });
      customerProductsTableBody.querySelectorAll('button').forEach(btn => {
        btn.addEventListener('click', () => {
          const productId = btn.getAttribute('data-id');
          addToCart(productId, products);
        });
      });
    } catch (err) {
      alert('Error listing products: ' + err.message);
    }
  }

  function addToCart(productId, productList) {
    console.log("Adding product to cart:", productId);
    const prod = productList.find(p => p.id.toString() === productId);
    if (!prod) return;
    const existing = cartItems.find(ci => ci.productId.toString() === productId);
    if (existing) {
      existing.quantity += 1;
    } else {
      cartItems.push({
        productId: prod.id,
        productName: prod.name,
        price: prod.price,
        quantity: 1
      });
    }
    alert(`Added ${prod.name} to cart.`);
    renderCart();
  }

  function showCartPage() {
    console.log("Showing cart page");
    document.getElementById('customer-shop-page').style.display = 'none';
    document.getElementById('cart-page').style.display = 'block';
    renderCart();
  }

  backToShopBtn.addEventListener('click', () => {
    console.log("Back to shop clicked");
    document.getElementById('cart-page').style.display = 'none';
    document.getElementById('customer-shop-page').style.display = 'block';
  });

  function renderCart() {
    console.log("Rendering cart");
    cartTableBody.innerHTML = '';
    let subtotal = 0;
    cartItems.forEach(item => {
      const st = item.price * item.quantity;
      subtotal += st;
      const tr = document.createElement('tr');
      tr.innerHTML = `
        <td>${item.productName}</td>
        <td>${item.quantity}</td>
        <td>$${item.price.toFixed(2)}</td>
        <td>$${st.toFixed(2)}</td>
      `;
      cartTableBody.appendChild(tr);
    });
    cartSubtotalEl.textContent = subtotal.toFixed(2);
    const payMethod = checkoutPaymentMethod.value;
    let extraFees = "";
    let finalTotal = subtotal;
    if (subtotal < 50) {
      finalTotal += 5;
      extraFees += `Order < $50 => +$5 delivery\n`;
    }
    if (payMethod === "cash") {
      const surcharge = finalTotal * 0.05;
      finalTotal += surcharge;
      extraFees += `Paying cash => +5%\n`;
    }
    extraFeesMsgEl.textContent = extraFees;
    cartFinalTotalEl.textContent = finalTotal.toFixed(2);
  }

  checkoutPaymentMethod.addEventListener('change', () => {
    console.log("Payment method changed:", checkoutPaymentMethod.value);
    renderCart();
  });

  checkoutBtn.addEventListener('click', async () => {
    console.log("Checkout button clicked");
    if (cartItems.length === 0) {
      alert('Cart is empty!');
      return;
    }
    let subtotal = 0;
    cartItems.forEach(ci => { subtotal += ci.price * ci.quantity });
    let finalTotal = subtotal;
    if (subtotal < 50) finalTotal += 5;
    if (checkoutPaymentMethod.value === 'cash') {
      finalTotal += finalTotal * 0.05;
    }
    const email = checkoutEmail.value.trim();
    const phone = checkoutPhone.value.trim();
    const address = checkoutAddress.value.trim();
    const paymentMethod = checkoutPaymentMethod.value;
    if (!email || !phone || !address) {
      alert('Please fill email, phone, address.');
      return;
    }
    const items = cartItems.map(ci => ({
      product_id: ci.productId,
      quantity: ci.quantity,
      price: ci.price
    }));
    console.log("Creating order with payment method:", paymentMethod);
    try {
      const res = await fetch('/api/v1/orders', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          customer_name: currentCustomerName,
          phone,
          email,
          address,
          payment_method: paymentMethod,
          final_total: parseFloat(finalTotal.toFixed(2)),
          items
        })
      });
      if (!res.ok) throw new Error('Failed to create order');
      const orderResp = await res.json();
      alert('Order created!\n' + JSON.stringify(orderResp, null, 2));
      if (paymentMethod === 'card') {
        console.log("Opening card payment window");
        window.open('./card.html', '_blank');
      }
      cartItems = [];
      checkoutEmail.value = '';
      checkoutPhone.value = '';
      checkoutAddress.value = '';
      checkoutPaymentMethod.value = 'card';
      renderCart();
      document.getElementById('cart-page').style.display = 'none';
      document.getElementById('customer-shop-page').style.display = 'block';
    } catch (err) {
      alert('Checkout error: ' + err.message);
    }
  });

  showLoginPage();
  console.log("All event listeners attached.");
});
