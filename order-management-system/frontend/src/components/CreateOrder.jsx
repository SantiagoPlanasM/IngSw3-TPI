import { useState, useEffect } from 'react';
import { orderService, userService } from '../services/api';

export default function CreateOrder({ cart, onClearCart, onOrderCreated }) {
  const [users, setUsers] = useState([]);
  const [selectedUser, setSelectedUser] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadUsers();
  }, []);

  const loadUsers = async () => {
    try {
      const response = await userService.getAll();
      setUsers(response.data);
      if (response.data.length > 0) {
        setSelectedUser(response.data[0].id.toString());
      }
    } catch (err) {
      console.error('Error loading users:', err);
    }
  };

  const getTotalPrice = () => {
    return cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (cart.length === 0) {
      setError('El carrito está vacío');
      return;
    }

    if (!selectedUser) {
      setError('Selecciona un usuario');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const orderData = {
        user_id: parseInt(selectedUser),
        items: cart.map(item => ({
          product_id: item.id,
          quantity: item.quantity
        }))
      };

      const response = await orderService.create(orderData);
      onClearCart();
      onOrderCreated(response.data);
      alert('¡Pedido creado exitosamente!');
    } catch (err) {
      setError(err.response?.data?.error || 'Error al crear el pedido');
    } finally {
      setLoading(false);
    }
  };

  const updateQuantity = (productId, delta) => {
    const item = cart.find(i => i.id === productId);
    if (item) {
      const newQuantity = item.quantity + delta;
      if (newQuantity > 0) {
        item.quantity = newQuantity;
      }
    }
  };

  if (cart.length === 0) {
    return (
      <div className="bg-gray-50 rounded-lg p-8 text-center">
        <p className="text-gray-500 text-lg">
          🛒 Tu carrito está vacío
        </p>
        <p className="text-gray-400 mt-2">
          Agrega productos para crear un pedido
        </p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-bold mb-4 text-gray-800">Resumen del Pedido</h2>
      
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Usuario
        </label>
        <select
          value={selectedUser}
          onChange={(e) => setSelectedUser(e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          {users.map(user => (
            <option key={user.id} value={user.id}>
              {user.name} ({user.email})
            </option>
          ))}
        </select>
      </div>

      <div className="border-t border-b border-gray-200 py-4 mb-4">
        <h3 className="font-semibold mb-3">Productos:</h3>
        {cart.map(item => (
          <div key={item.id} className="flex justify-between items-center mb-3 pb-3 border-b border-gray-100 last:border-0">
            <div className="flex-1">
              <p className="font-medium text-gray-800">{item.name}</p>
              <p className="text-sm text-gray-500">${item.price.toFixed(2)} c/u</p>
            </div>
            <div className="flex items-center gap-2">
              <button
                onClick={() => updateQuantity(item.id, -1)}
                className="w-8 h-8 rounded bg-gray-200 hover:bg-gray-300 flex items-center justify-center"
              >
                -
              </button>
              <span className="w-8 text-center font-semibold">{item.quantity}</span>
              <button
                onClick={() => updateQuantity(item.id, 1)}
                className="w-8 h-8 rounded bg-gray-200 hover:bg-gray-300 flex items-center justify-center"
              >
                +
              </button>
            </div>
            <div className="ml-4 font-semibold text-blue-600">
              ${(item.price * item.quantity).toFixed(2)}
            </div>
          </div>
        ))}
      </div>

      <div className="flex justify-between items-center mb-4 text-xl font-bold">
        <span>Total:</span>
        <span className="text-blue-600">${getTotalPrice().toFixed(2)}</span>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}

      <div className="flex gap-3">
        <button
          onClick={handleSubmit}
          disabled={loading}
          className="flex-1 bg-blue-500 text-white py-3 rounded-lg font-semibold hover:bg-blue-600 transition-colors disabled:bg-gray-300 disabled:cursor-not-allowed"
        >
          {loading ? 'Procesando...' : 'Crear Pedido'}
        </button>
        <button
          onClick={onClearCart}
          className="px-6 py-3 bg-gray-200 text-gray-700 rounded-lg font-semibold hover:bg-gray-300 transition-colors"
        >
          Limpiar
        </button>
      </div>
    </div>
  );
}
