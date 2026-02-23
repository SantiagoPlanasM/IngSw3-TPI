import { useState } from 'react';
import ProductList from './components/ProductList';
import CreateOrder from './components/CreateOrder';
import OrderHistory from './components/OrderHistory';

function App() {
  const [activeTab, setActiveTab] = useState('products');
  const [cart, setCart] = useState([]);
  const [refreshOrders, setRefreshOrders] = useState(0);

  const addToCart = (product) => {
    const existingItem = cart.find(item => item.id === product.id);
    
    if (existingItem) {
      setCart(cart.map(item =>
        item.id === product.id
          ? { ...item, quantity: item.quantity + 1 }
          : item
      ));
    } else {
      setCart([...cart, { ...product, quantity: 1 }]);
    }
    
    // Switch to cart tab
    setActiveTab('cart');
  };

  const clearCart = () => {
    setCart([]);
  };

  const handleOrderCreated = () => {
    setRefreshOrders(prev => prev + 1);
    setActiveTab('orders');
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-purple-50">
      {/* Header */}
      <header className="bg-white shadow-md">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center justify-between">
            <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
              🛍️ Sistema de Gestión de Pedidos v1.0
            </h1>
            {cart.length > 0 && (
              <div className="bg-blue-500 text-white px-4 py-2 rounded-full font-semibold">
                🛒 {cart.length} {cart.length === 1 ? 'producto' : 'productos'}
              </div>
            )}
          </div>
        </div>
      </header>

      {/* Navigation Tabs */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-6">
        <div className="flex gap-2 border-b border-gray-200">
          <button
            onClick={() => setActiveTab('products')}
            className={`px-6 py-3 font-semibold transition-colors ${
              activeTab === 'products'
                ? 'border-b-2 border-blue-500 text-blue-600'
                : 'text-gray-600 hover:text-gray-800'
            }`}
          >
            📦 Productos
          </button>
          <button
            onClick={() => setActiveTab('cart')}
            className={`px-6 py-3 font-semibold transition-colors relative ${
              activeTab === 'cart'
                ? 'border-b-2 border-blue-500 text-blue-600'
                : 'text-gray-600 hover:text-gray-800'
            }`}
          >
            🛒 Carrito
            {cart.length > 0 && (
              <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center">
                {cart.length}
              </span>
            )}
          </button>
          <button
            onClick={() => setActiveTab('orders')}
            className={`px-6 py-3 font-semibold transition-colors ${
              activeTab === 'orders'
                ? 'border-b-2 border-blue-500 text-blue-600'
                : 'text-gray-600 hover:text-gray-800'
            }`}
          >
            📋 Historial de Pedidos
          </button>
        </div>
      </div>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {activeTab === 'products' && (
          <div>
            <h2 className="text-2xl font-bold mb-6 text-gray-800">Catálogo de Productos</h2>
            <ProductList onAddToCart={addToCart} />
          </div>
        )}

        {activeTab === 'cart' && (
          <div>
            <h2 className="text-2xl font-bold mb-6 text-gray-800">Carrito de Compras</h2>
            <CreateOrder 
              cart={cart} 
              onClearCart={clearCart}
              onOrderCreated={handleOrderCreated}
            />
          </div>
        )}

        {activeTab === 'orders' && (
          <div>
            <h2 className="text-2xl font-bold mb-6 text-gray-800">Historial de Pedidos</h2>
            <OrderHistory refreshTrigger={refreshOrders} />
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-gray-200 mt-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <p className="text-center text-gray-600">
            Sistema de Gestión de Pedidos - DevOps CI/CD Demo
          </p>
        </div>
      </footer>
    </div>
  );
}

export default App;
