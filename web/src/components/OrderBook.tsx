import { Typography, Flex } from 'antd';
const { Text } = Typography;

export const OrderBook = ({ isDark }: { isDark: boolean }) => {
    const borderColor = isDark ? '#1f1f1f' : '#f0f0f0';
    const textColor = isDark ? '#ffffff' : '#000000';

    const orders = [
        { price: '1.0542', amount: '842.1', type: 'sell' },
        { price: '1.0531', amount: '120.0', type: 'sell' },
        { price: '1.0515', amount: '430.0', type: 'buy' },
        { price: '1.0509', amount: '75.2', type: 'buy' },
        { price: '1.0498', amount: '210.5', type: 'buy' },
    ];

    return (
        <Flex vertical style={{ height: '100%' }}>
            <div style={{ padding: '8px 12px', borderBottom: `1px solid ${borderColor}` }}>
                <Text strong style={{ color: textColor, fontSize: '12px' }}>实时成交</Text>
            </div>
            <div style={{ flex: 1, overflow: 'hidden' }}>
                <Flex justify="space-between" style={{ padding: '4px 12px', opacity: 0.5 }}>
                    <Text style={{ fontSize: '10px', color: textColor }}>价格</Text>
                    <Text style={{ fontSize: '10px', color: textColor }}>数量</Text>
                </Flex>
                {orders.map((item, i) => (
                    <Flex key={i} justify="space-between" style={{ padding: '4px 12px', fontSize: '11px', fontFamily: 'monospace' }}>
                        <Text strong style={{ color: item.type === 'sell' ? '#f23645' : '#22ab94' }}>{item.price}</Text>
                        <Text style={{ color: textColor }}>{item.amount}</Text>
                    </Flex>
                ))}
            </div>
        </Flex>
    );
};
