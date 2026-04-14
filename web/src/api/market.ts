import request from '../utils/request';

export interface KlineItem {
    t: number; o: string; h: string; l: string; c: string; v: string;
}

export const getKLines = (symbol: string, period: string, limit: number = 200) =>
    request.get('/market/kline', { params: { symbol, period, limit } });

export const getDepth = (symbol: string, limit: number = 20) =>
    request.get('/market/depth', { params: { symbol, limit } });
