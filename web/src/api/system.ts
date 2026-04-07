import request from '../utils/request';

export interface SystemConfig {
    vault_address: string;
    token_address: string;
    chain_id: number;
}

// 获取后端配置
export const getSystemConfig = (): Promise<{ data: SystemConfig }> => {
    return request.get('/system/config');
};

// 登录接口
export const login = (address: string): Promise<{ token: string, walletAddress: string }> => {
    return request.post('/auth/login', { address });
};
