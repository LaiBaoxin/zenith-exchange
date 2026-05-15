import request from '../utils/request';

export const placeOrder = (data: { symbol: string, side: 'buy' | 'sell', price: number, amount: number }) =>
    request.post('/order/place', data);

export const cancelOrder = (orderId: number) =>
    request.post('/order/cancel', { order_id: orderId.toString() });

export const getTodayOrders = (symbol?: string) =>
    request.get('/order/today', { params: { symbol } });

export const getAllOrders = (params: { symbol?: string, page: number, page_size: number }) =>
    request.get('/order/list', { params });

export const getOrderDetail = (orderId: number) =>
    request.get(`/order/detail/${orderId}`);
