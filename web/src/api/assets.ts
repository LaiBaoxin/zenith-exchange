import request from '../utils/request';

export const getBalance = () =>
    request.get('/assets/balance');
