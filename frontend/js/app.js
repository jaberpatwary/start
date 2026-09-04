const products=[
  {name:'Asus ROG Strix G16 Gaming Laptop',category:'Laptop',price:189900},
  {name:'HP Pavilion 15 Core i5 Laptop',category:'Laptop',price:82500},
  {name:'Corsair Vengeance Gaming PC',category:'Desktop',price:162000},
  {name:'Samsung Odyssey G5 Monitor',category:'Monitor',price:38500},
  {name:'Gigabyte RTX 4060 Windforce',category:'Component',price:49500},
  {name:'MSI B760M Mortar WiFi',category:'Component',price:24500},
  {name:'Dell Inspiron 15 3530',category:'Laptop',price:74500},
  {name:'Samsung 990 EVO 1TB SSD',category:'Component',price:10500}
];
let category='All',cart=0;
function renderCategories(){const el=document.getElementById('categories');el.innerHTML=['All','Desktop','Laptop','Component','Monitor'].map(c=>`<button class="category-btn ${c===category?'active':''}" onclick="chooseCategory('${c}')">${c}</button>`).join('')}
function chooseCategory(value){category=value;renderCategories();renderProducts()}
function renderProducts(){const query=document.getElementById('searchInput').value.toLowerCase();const list=products.filter(p=>(category==='All'||p.category===category)&&p.name.toLowerCase().includes(query));document.getElementById('products').innerHTML=list.map(p=>`<article class="card"><img src="https://placehold.co/600x600/f8fafc/e63946?text=${encodeURIComponent(p.name)}" alt="${p.name}"><small>${p.category}</small><h3>${p.name}</h3><span class="price">৳${p.price.toLocaleString()}</span><button class="add" onclick="addToCart()">Add to Cart</button></article>`).join('')||'<p>No products found.</p>'}
function addToCart(){cart++;document.getElementById('cartCount').textContent=cart}
document.getElementById('searchInput').addEventListener('input',renderProducts);renderCategories();renderProducts();
