import request from '../utils/request';

// 定义接口请求参数类型
export interface WithdrawSignParams {
    user: string;
    token: string;
    amount: string;
    nonce: number;
    vault_addr: string;
    chain_id: number;
}

// 定义返回类型
export interface WithdrawSignResponse {
    signature: string;
}

/**
 * 获取提现授权签名
 */
export const getWithdrawSignature = async (params: WithdrawSignParams): Promise<WithdrawSignResponse> => {
    return await request.post('/v1/vault/withdraw-sign', params) as unknown as Promise<WithdrawSignResponse>;
};
