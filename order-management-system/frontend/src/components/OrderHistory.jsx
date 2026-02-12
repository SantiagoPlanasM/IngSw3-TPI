import { useState, useEffect } from 'react';
import { orderService } from '../services/api';

const statusColors = {
  PENDING: 'bg-yellow-100 text-yellow-800',
  CONFIRMED: 'bg-blue-100 text-blue-800',
  SHIPPED: 'bg-green-100 text-green-800',
  CANCELLED: 'bg-red-100 text-red-800',
};

const statusLabels = {
  PENDING: 'Pendiente',
  CONFIRMED: 'Confirmado',
  SHIPPED: 'Enviado',
  CANCELLED: 'Cancelado',
};

export default function OrderHistory({ refreshTrigger }) {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [expandedOrder, setExpandedOrder] = useState(null);

  useEffect(() => {
    loadOrders();
  }, [refreshTrigger]);

  const loadOrders = async () => {
    try {
      setLoading(true);
      const response = await orderService.getAll();
      setOrders(response.data.sort((a, b) => b.id - a.id));
    } catch (err) {
      setError('Error al cargar pedidos');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleConfirm = async (orderId) => {
    if (!confirm('¿Confirmar este pedido? Se reducirá el stock.')) return;
    
    try {
      await orderService.confirm(orderId);
      loadOrders();
    } catch (err) {
      alert(err.response?.data?.error || 'Error al confirmar pedido');
    }
  };

  const handleShip = async (orderId) => {
    if (!confirm('¿Marcar este pedido como enviado?')) return;
    
    try {
      await orderService.ship(orderId);
      loadOrders();
    } catch (err) {
      alert(err.response?.data?.error || 'Error al enviar pedido');
    }
  };

  const handleCancel = async (orderId) => {
    if (!confirm('¿Cancelar este pedido? Se devolverá el stock si estaba confirmado.')) return;
    
    try {
      await orderService.cancel(orderId);
      loadOrders();
    } catch (err) {
      alert(err.response?.data?.error || 'Error al cancelar pedido');
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
        {error}
      </div>
    );
  }

  if (orders.length === 0) {
    return (
      <div className="bg-gray-50 rounded-lg p-8 text-center">
        <p className="text-gray-500 text-lg">📋 No hay pedidos registrados</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {orders.map((order) => (
        <div key={order.id} className="bg-white rounded-lg shadow-md overflow-hidden">
          <div 
            className="p-4 cursor-pointer hover:bg-gray-50 transition-colors"
            onClick={() => setExpandedOrder(expandedOrder === order.id ? null : order.id)}
          >
            <div className="flex justify-between items-start">
              <div className="flex-1">
                <div className="flex items-center gap-3 mb-2">
                  <h3 className="font-semibold text-lg">Pedido #{order.id}</h3>
                  <span className={`px-3 py-1 rounded-full text-xs font-semibold ${statusColors[order.status]}`}>
                    {statusLabels[order.status]}
                  </span>
                </div>
                <p className="text-sm text-gray-600">
                  Cliente: <span className="font-medium">{order.user?.name}</span>
                </p>
                <p className="text-sm text-gray-500">
                  {new Date(order.created_at).toLocaleString('es-AR')}
                </p>
              </div>
              <div className="text-right">
                <p className="text-2xl font-bold text-blue-600">
                  ${order.total.toFixed(2)}
                </p>
                <p className="text-xs text-gray-500">
                  {order.items?.length || 0} producto(s)
                </p>
              </div>
            </div>
          </div>

          {expandedOrder === order.id && (
            <div className="border-t border-gray-200 p-4 bg-gray-50">
              <h4 className="font-semibold mb-3">Productos:</h4>
              <div className="space-y-2 mb-4">
                {order.items?.map((item, idx) => (
                  <div key={idx} className="flex justify-between items-center bg-white p-3 rounded">
                    <div>
                      <p className="font-medium">{item.product?.name}</p>
                      <p className="text-sm text-gray-600">
                        Cantidad: {item.quantity} × ${item.price.toFixed(2)}
                      </p>
                    </div>
                    <p className="font-semibold text-blue-600">
                      ${(item.quantity * item.price).toFixed(2)}
                    </p>
                  </div>
                ))}
              </div>

              <div className="flex gap-2 pt-3 border-t border-gray-200">
                {order.status === 'PENDING' && (
                  <>
                    <button
                      onClick={() => handleConfirm(order.id)}
                      className="flex-1 bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 transition-colors"
                    >
                      ✓ Confirmar
                    </button>
                    <button
                      onClick={() => handleCancel(order.id)}
                      className="flex-1 bg-red-500 text-white py-2 px-4 rounded hover:bg-red-600 transition-colors"
                    >
                      ✕ Cancelar
                    </button>
                  </>
                )}
                {order.status === 'CONFIRMED' && (
                  <>
                    <button
                      onClick={() => handleShip(order.id)}
                      className="flex-1 bg-green-500 text-white py-2 px-4 rounded hover:bg-green-600 transition-colors"
                    >
                      🚚 Enviar
                    </button>
                    <button
                      onClick={() => handleCancel(order.id)}
                      className="flex-1 bg-red-500 text-white py-2 px-4 rounded hover:bg-red-600 transition-colors"
                    >
                      ✕ Cancelar
                    </button>
                  </>
                )}
              </div>
            </div>
          )}
        </div>
      ))}
    </div>
  );
}
