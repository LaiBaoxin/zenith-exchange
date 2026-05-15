import request from '../utils/request';

export const getWithdrawSignature = (params: { amount: string, currency: string }) => {
    return request.post('/vault/withdraw-sign', params);
};

