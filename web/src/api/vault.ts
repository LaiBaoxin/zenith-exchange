import request from '../utils/request';

export interface WithdrawSignResponse {
    signature: string;
    nonce: number;
    amount: string;
    token: string;
}

export const getWithdrawSignature = (amount: string): Promise<WithdrawSignResponse> => {
    return request.post('/v1/vault/withdraw-sign', { amount });
};
