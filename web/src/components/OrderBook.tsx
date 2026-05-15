import { useEffect, useState } from 'react';
import { Typography, Flex, Spin, ConfigProvider, theme } from 'antd';
import { getDepth } from '../api/market';

const { Text } = Typography;

interface PriceLevel {
    price: string;
    amount: string;
}

interface DepthResponse {
    bids: PriceLevel[];
    asks: PriceLevel[];
}

export const OrderBook = ({ isDark, symbol = "BTC_USDT" }: { isDark: boolean, symbol?: string }) => {
    const borderColor = isDark ? '#1f1f1f' : '#f0f0f0';
    const textColor = isDark ? '#ffffff' : '#000000';

    // 状态管理
    const [depth, setDepth] = useState<DepthResponse>({ bids: [], asks: [] });
    const [loading, setLoading] = useState(true);
    // 【新增】最新成交价状态，初始化为一个默认值
    const [latestPrice, setLatestPrice] = useState("0.00");

    useEffect(() => {
        const fetchDepth = async () => {
            try {
                const res = await getDepth(symbol, 15) as any;
                const data = res;

                // 更新深度列表
                setDepth({
                    bids: data?.bids || [],
                    asks: data?.asks || []
                });

                // 【核心修复】计算并更新中间显示的最新价格
                // 逻辑：如果盘口有数据，取卖一和买一的平均值作为中间参考价
                if (data?.bids?.length > 0 && data?.asks?.length > 0) {
                    const buyOne = parseFloat(data.bids[0].price);
                    const sellOne = parseFloat(data.asks[0].price);
                    const midPrice = (buyOne + sellOne) / 2;

                    // 格式化为两位小数
                    setLatestPrice(midPrice.toLocaleString(undefined, {
                        minimumFractionDigits: 2,
                        maximumFractionDigits: 2
                    }));
                }
            } catch (e) {
                console.error("获取深度失败", e);
            } finally {
                setLoading(false);
            }
        };

        fetchDepth();
        const timer = setInterval(fetchDepth, 2000);
        return () => clearInterval(timer);
    }, [symbol]);

    const maxAmount = Math.max(
        ...depth.asks.map(item => parseFloat(item.amount)),
        ...depth.bids.map(item => parseFloat(item.amount)),
        1
    );

    const renderLevel = (item: PriceLevel, type: 'bid' | 'ask') => {
        const amountNum = parseFloat(item.amount);
        const percent = Math.min((amountNum / maxAmount) * 100, 100);

        return (
            <div key={`${type}-${item.price}`} style={{ position: 'relative', height: '22px', cursor: 'pointer' }} className="depth-row">
                <div style={{
                    position: 'absolute',
                    right: 0,
                    top: 1,
                    bottom: 1,
                    width: `${percent}%`,
                    background: type === 'ask' ? 'rgba(242, 54, 69, 0.1)' : 'rgba(34, 171, 148, 0.1)',
                    transition: 'width 0.3s ease',
                    zIndex: 0
                }} />

                <Flex justify="space-between" align="center" style={{ padding: '0 12px', height: '100%', position: 'relative', zIndex: 1 }}>
                    <Text style={{
                        color: type === 'ask' ? '#f23645' : '#22ab94',
                        fontSize: '12px',
                        fontFamily: 'monospace',
                        fontWeight: 500
                    }}>
                        {parseFloat(item.price).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </Text>
                    <Text style={{ color: textColor, fontSize: '12px', fontFamily: 'monospace', opacity: 0.9 }}>
                        {parseFloat(item.amount).toFixed(4)}
                    </Text>
                </Flex>
            </div>
        );
    };

    return (
        <ConfigProvider theme={{ algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm }}>
            <Flex vertical style={{ width: '100%', height: '100%', overflow: 'hidden', background: isDark ? '#141414' : '#fff' }}>
                <div style={{ padding: '8px 12px', borderBottom: `1px solid ${borderColor}` }}>
                    <Text strong style={{ color: textColor, fontSize: '13px' }}>盘口深度</Text>
                </div>

                {loading && depth.asks.length === 0 ? (
                    <Flex justify="center" align="center" style={{ flex: 1 }}><Spin size="small" /></Flex>
                ) : (
                    <div style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
                        <Flex justify="space-between" style={{ padding: '4px 12px', opacity: 0.5, borderBottom: `1px solid ${borderColor}` }}>
                            <Text style={{ fontSize: '11px', color: textColor }}>价格(USDT)</Text>
                            <Text style={{ fontSize: '11px', color: textColor }}>数量({symbol.split('_')[0]})</Text>
                        </Flex>

                        {/* 卖盘区域 */}
                        <div style={{
                            flex: 1,
                            display: 'flex',
                            flexDirection: 'column-reverse',
                            justifyContent: 'flex-start',
                            overflow: 'hidden',
                            paddingBottom: '2px'
                        }}>
                            {depth.asks.map(item => renderLevel(item, 'ask'))}
                        </div>

                        {/* 【重点修改】中间价格显示区 - 现在绑定了 latestPrice 变量 */}
                        <div style={{
                            padding: '10px 12px',
                            background: isDark ? 'rgba(255,255,255,0.02)' : '#fafafa',
                            borderTop: `1px solid ${borderColor}`,
                            borderBottom: `1px solid ${borderColor}`,
                            margin: '2px 0',
                            zIndex: 2
                        }}>
                            <Flex align="baseline" gap={4}>
                                <Text strong style={{ fontSize: '18px', color: '#22ab94' }}>
                                    {latestPrice}
                                </Text>
                                <Text style={{ fontSize: '12px', color: textColor, opacity: 0.6 }}>
                                    ≈ ${latestPrice}
                                </Text>
                            </Flex>
                        </div>

                        {/* 买盘区域 */}
                        <div style={{
                            flex: 1,
                            display: 'flex',
                            flexDirection: 'column',
                            justifyContent: 'flex-start',
                            overflow: 'hidden',
                            paddingTop: '2px'
                        }}>
                            {depth.bids.map(item => renderLevel(item, 'bid'))}
                        </div>
                    </div>
                )}
                <style>{`
                    .depth-row:hover { background: ${isDark ? 'rgba(255,255,255,0.05)' : 'rgba(0,0,0,0.03)'}; }
                `}</style>
            </Flex>
        </ConfigProvider>
    );
};
