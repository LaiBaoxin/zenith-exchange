import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { Layout, Flex, Typography, Spin } from 'antd';
import { TradePanel } from "./TradePanel.tsx";
import { OrderBook } from "./OrderBook.tsx";
import { UserOrder } from "./UserOrders.tsx";
import { KLineChart, type Period } from "./KLineChart.tsx";
// import { type SystemConfig } from '../api/system.ts';
import { getKLines } from '../api/market.ts';

const { Text } = Typography;

interface TradingTerminalProps {
    isDark: boolean;
    config: any;
}

const PERIOD_TEXT_MAP: Record<Period, string> = {
    '1m': '1分钟线',
    '5m': '5分钟线',
    '15m': '15分钟线',
    '1h': '1小时线',
    '1d': '日K线'
};

export default function TradingTerminal({ isDark, config: rawConfig }: TradingTerminalProps) {
    const [refreshTrigger, setRefreshTrigger] = useState(0);

    // 【关键修复】：自动识别 config 的层级
    // 如果 rawConfig 里有 data 属性，就用 rawConfig.data，否则用 rawConfig 本身
    const safeConfig = useMemo(() => {
        if (!rawConfig) return null;
        return rawConfig.data ? rawConfig.data : rawConfig;
    }, [rawConfig]);

    const [klineData, setKlineData] = useState<any[]>([]);
    const [currentPeriod, setCurrentPeriod] = useState<Period>('1m');
    const [klineLoading, setKlineLoading] = useState(true);

    useEffect(() => {
        if (safeConfig) {
            console.log("交易终端已接收有效配置:", safeConfig.vault_address);
        }
    }, [safeConfig]);

    const handleTradeSuccess = () => {
        setRefreshTrigger(prev => prev + 1);
    };

    const fetchKlines = useCallback(async (period: Period, showLoading = false) => {
        if (showLoading) setKlineLoading(true);
        try {
            const res = await getKLines('BTC_USDT', period, 500) as any;
            const rawData = res.data || res || [];
            setKlineData(Array.isArray(rawData) ? rawData : []);
        } catch (e) {
            console.error("K线数据加载失败", e);
        } finally {
            setKlineLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchKlines(currentPeriod, true);
        const timer = setInterval(() => fetchKlines(currentPeriod, false), 3000);
        return () => clearInterval(timer);
    }, [currentPeriod, fetchKlines]);

    const bgColor = isDark ? '#000000' : '#f5f5f5';
    const panelBg = isDark ? '#141414' : '#ffffff';
    const borderColor = isDark ? '#1f1f1f' : '#f0f0f0';
    const headerBg = isDark ? '#1a1a1a' : '#fafafa';
    const textColor = isDark ? '#ffffff' : '#000000';

    const panelStyle: React.CSSProperties = {
        background: panelBg,
        border: `1px solid ${borderColor}`,
        borderRadius: '4px',
        display: 'flex',
        flexDirection: 'column',
        overflow: 'hidden',
        transition: 'all 0.3s'
    };

    return (
        <Layout style={{ height: 'calc(100vh - 64px)', background: bgColor, overflow: 'hidden', padding: '8px' }}>
            <Flex style={{ height: '100%', gap: '8px' }}>
                <Flex vertical style={{ flex: 1, gap: '8px', minWidth: 0 }}>
                    <div style={{ ...panelStyle, flex: 1, position: 'relative', minHeight: '450px' }}>
                        <div style={{ padding: '8px 12px', borderBottom: `1px solid ${borderColor}`, background: headerBg }}>
                            <Text strong style={{ fontSize: '12px', color: textColor }}>
                                BTC/USDT - {PERIOD_TEXT_MAP[currentPeriod]}
                            </Text>
                        </div>

                        <div style={{ flex: 1, position: 'relative' }}>
                            <KLineChart
                                isDark={isDark}
                                data={klineData}
                                currentPeriod={currentPeriod}
                                onPeriodChange={(p) => setCurrentPeriod(p)}
                            />

                            {klineLoading && klineData.length === 0 && (
                                <div style={{
                                    position: 'absolute', top: 45, left: 0, right: 0, bottom: 0,
                                    display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center',
                                    background: panelBg, zIndex: 10
                                }}>
                                    <Spin size="large" />
                                    <Text style={{ marginTop: 12, color: isDark ? '#666' : '#ccc' }}>
                                        正在同步 {PERIOD_TEXT_MAP[currentPeriod]} 数据...
                                    </Text>
                                </div>
                            )}
                        </div>
                    </div>

                    <div style={{ ...panelStyle, height: 'auto', flexShrink: 0 }}>
                        <TradePanel
                            isDark={isDark}
                            // 【关键修复】：传递处理后的 safeConfig
                            config={safeConfig}
                            symbol="BTC_USDT"
                            onTradeSuccess={handleTradeSuccess}
                        />
                    </div>
                </Flex>

                <Flex vertical style={{ width: '320px', flexShrink: 0, gap: '8px' }}>
                    <div style={{ ...panelStyle, flex: 6 }}>
                        <OrderBook isDark={isDark} symbol="BTC_USDT" />
                    </div>
                    <div style={{ ...panelStyle, flex: 4 }}>
                        <div style={{ padding: '8px 12px', borderBottom: `1px solid ${borderColor}`, background: headerBg }}>
                            <Text strong style={{ fontSize: '12px', color: textColor }}>当前委托</Text>
                        </div>
                        <div style={{ flex: 1, overflowY: 'auto' }}>
                            <UserOrder symbol="BTC_USDT" refreshTrigger={refreshTrigger} isDark={isDark} />
                        </div>
                    </div>
                </Flex>
            </Flex>
        </Layout>
    );
}
