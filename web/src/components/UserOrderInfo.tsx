import { useEffect, useState, useCallback } from 'react';
import { Tabs, Typography, App, Spin, Flex, Tag, Button, Progress, Modal, Descriptions, theme, ConfigProvider, Space } from 'antd';
import { InfoCircleOutlined, CopyOutlined, ShoppingOutlined } from '@ant-design/icons';
import { getTodayOrders, getAllOrders, cancelOrder, getOrderDetail } from '../api/order';

const { Text } = Typography;

const STATUS_MAP: Record<number, { text: string; color: string }> = {
    0: { text: '挂单中', color: 'orange' },
    1: { text: '部分成交', color: 'blue' },
    2: { text: '已成交', color: 'green' },
    3: { text: '已撤单', color: 'default' },
};

export default function UserOrderInfo({ symbol, refreshTrigger, isDark }: { symbol: string, refreshTrigger: number, isDark: boolean }) {
    const { message } = App.useApp();
    const [activeTab, setActiveTab] = useState('today');
    const [loading, setLoading] = useState(false);
    const [orders, setOrders] = useState<any[]>([]);

    const [isModalOpen, setIsModalOpen] = useState(false);
    const [detailLoading, setDetailLoading] = useState(false);
    const [selectedOrder, setSelectedOrder] = useState<any>(null);

    const bgColor = isDark ? '#000000' : '#ffffff';
    const borderColor = isDark ? '#1f1f1f' : '#f0f0f0';
    const cardBg = isDark ? '#141414' : '#ffffff';
    const subTextColor = isDark ? 'rgba(255,255,255,0.45)' : 'rgba(0,0,0,0.45)';

    const fetchOrders = useCallback(async () => {
        setLoading(true);
        try {
            let res: any;
            if (activeTab === 'today') {
                res = await getTodayOrders(symbol);
                setOrders(Array.isArray(res?.data) ? res.data : (Array.isArray(res) ? res : []));
            } else {
                res = await getAllOrders({ symbol, page: 1, page_size: 20 });
                setOrders(res?.data?.list || []);
            }
        } catch (e) {
            console.error("Fetch error", e);
        } finally {
            setLoading(false);
        }
    }, [activeTab, symbol, refreshTrigger]);

    useEffect(() => { fetchOrders(); }, [fetchOrders]);

    useEffect(() => {
        const timer = setInterval(() => fetchOrders(), 5000);
        return () => clearInterval(timer);
    }, [fetchOrders]);

    const showDetail = async (order: any) => {
        setSelectedOrder(order);
        setIsModalOpen(true);
        setDetailLoading(true);
        try {
            const res = await getOrderDetail(order.id);
            if (res?.data) setSelectedOrder(res.data);
        } catch (e) {
            message.error("获取详细数据失败");
        } finally {
            setDetailLoading(false);
        }
    };

    const handleCancel = async (e: React.MouseEvent, id: number) => {
        e.stopPropagation();
        try {
            await cancelOrder(id);
            message.success("撤单成功");
            fetchOrders();
        } catch (e: any) {
            message.error("撤单失败");
        }
    };

    return (
        <ConfigProvider theme={{ algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm }}>
            <div style={{ background: bgColor, height: '100%', display: 'flex', flexDirection: 'column', transition: '0.3s' }}>
                <Tabs
                    centered
                    activeKey={activeTab}
                    onChange={setActiveTab}
                    size="small"
                    tabBarStyle={{ marginBottom: 0, borderBottom: `1px solid ${borderColor}` }}
                    items={[{ key: 'today', label: '今日委托' }, { key: 'history', label: '历史账本' }]}
                />

                <div style={{ flex: 1, overflowY: 'auto', padding: '12px' }}>
                    {loading ? <Flex justify="center" align="center" style={{ height: '100%' }}><Spin /></Flex> : (
                        orders.length === 0 ? (
                            <Flex vertical align="center" justify="center" style={{ height: '80%', opacity: 0.3 }}>
                                <ShoppingOutlined style={{ fontSize: 40 }} />
                                <Text>暂无数据</Text>
                            </Flex>
                        ) : orders.map(order => {
                            const sideColor = order.side === 'buy' ? '#22ab94' : '#f23645';
                            const progress = Number(((order.filled_amount / order.amount) * 100).toFixed(1));

                            return (
                                <div key={order.id} onClick={() => showDetail(order)} style={{
                                    background: cardBg, border: `1px solid ${borderColor}`, borderRadius: 6,
                                    padding: '10px', marginBottom: 10, cursor: 'pointer', transition: '0.2s'
                                }} className="order-card-hover">
                                    <Flex justify="space-between" align="center" style={{ marginBottom: 8 }}>
                                        <Space>
                                            <Tag color={sideColor} style={{ fontWeight: 'bold', margin: 0 }}>{order.side === 'buy' ? '买入' : '卖出'}</Tag>
                                            <Text strong style={{ fontSize: 13 }}>{order.symbol}</Text>
                                        </Space>
                                        <Text style={{ fontSize: 11, color: subTextColor }}>
                                            {new Date(order.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                                        </Text>
                                    </Flex>

                                    <Flex justify="space-between" style={{ marginBottom: 4 }}>
                                        <Text style={{ fontSize: 12, color: subTextColor }}>价格</Text>
                                        <Text strong>{Number(order.price).toFixed(2)}</Text>
                                    </Flex>

                                    <Flex justify="space-between" style={{ marginBottom: 6 }}>
                                        <Text style={{ fontSize: 12, color: subTextColor }}>数量</Text>
                                        <Text>{Number(order.amount).toFixed(4)}</Text>
                                    </Flex>

                                    <Progress percent={progress} size="small" strokeColor={sideColor} showInfo={false} style={{ marginBottom: 4 }} />

                                    <Flex justify="space-between" align="center" style={{ marginTop: 4 }}>
                                        <Tag color={STATUS_MAP[order.status]?.color} style={{ fontSize: 10 }}>
                                            {STATUS_MAP[order.status]?.text}
                                        </Tag>
                                        {order.status <= 1 && (
                                            <Button type="link" danger size="small" onClick={(e) => handleCancel(e, order.id)} style={{ padding: 0, height: 'auto', fontSize: 12 }}>
                                                撤单
                                            </Button>
                                        )}
                                    </Flex>
                                </div>
                            );
                        })
                    )}
                </div>

                <Modal
                    title={<Space><InfoCircleOutlined />订单详情</Space>}
                    open={isModalOpen}
                    onCancel={() => setIsModalOpen(false)}
                    footer={<Button type="primary" onClick={() => setIsModalOpen(false)}>确定</Button>}
                    width={400}
                    centered
                >
                    <Spin spinning={detailLoading}>
                        {selectedOrder && (
                            <Descriptions column={1} size="small" bordered style={{ marginTop: 16 }}>
                                <Descriptions.Item label="订单ID">{selectedOrder.id}</Descriptions.Item>
                                <Descriptions.Item label="交易对">{selectedOrder.symbol}</Descriptions.Item>
                                <Descriptions.Item label="类型">{selectedOrder.type || '限价'}</Descriptions.Item>
                                <Descriptions.Item label="委托价">{Number(selectedOrder.price).toFixed(2)}</Descriptions.Item>
                                <Descriptions.Item label="委托量">{Number(selectedOrder.amount).toFixed(4)}</Descriptions.Item>
                                <Descriptions.Item label="已成交">{Number(selectedOrder.filled_amount).toFixed(4)}</Descriptions.Item>
                                <Descriptions.Item label="下单时间">{new Date(selectedOrder.created_at).toLocaleString()}</Descriptions.Item>
                                <Descriptions.Item label="消息哈希">
                                    <Text copyable={{ icon: <CopyOutlined /> }} style={{ fontSize: 11 }}>
                                        {selectedOrder.msg_hash ? `${selectedOrder.msg_hash.slice(0, 10)}...` : 'N/A'}
                                    </Text>
                                </Descriptions.Item>
                            </Descriptions>
                        )}
                    </Spin>
                </Modal>

                <style>{`
                    .order-card-hover:hover {
                        border-color: #1677ff !important;
                        box-shadow: 0 2px 8px ${isDark ? 'rgba(0,0,0,0.5)' : 'rgba(0,0,0,0.05)'};
                    }
                `}</style>
            </div>
        </ConfigProvider>
    );
}
