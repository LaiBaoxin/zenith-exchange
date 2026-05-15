import { useEffect, useState, useCallback } from 'react';
import { Tag, Button, Tabs, App, Typography, Flex, Spin, Progress, Modal, Descriptions, Empty, Pagination, ConfigProvider, theme, Space } from 'antd';
import { CopyOutlined } from '@ant-design/icons';
import { getTodayOrders, getAllOrders, cancelOrder, getOrderDetail } from '../api/order';

const { Text } = Typography;

const STATUS_MAP: Record<number, { text: string; color: string }> = {
    0: { text: '挂单中', color: 'orange' },
    1: { text: '部分成交', color: 'blue' },
    2: { text: '已成交', color: 'green' },
    3: { text: '已撤单', color: 'default' },
};

export const UserOrder = ({ symbol, refreshTrigger, isDark }: { symbol: string, refreshTrigger: number, isDark: boolean }) => {
    const { message } = App.useApp();
    const [activeTab, setActiveTab] = useState('today');
    const [loading, setLoading] = useState(false);
    const [orders, setOrders] = useState<any[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);

    const [isModalOpen, setIsModalOpen] = useState(false);
    const [detailLoading, setDetailLoading] = useState(false);
    const [selectedOrder, setSelectedOrder] = useState<any | null>(null);

    const bgColor = isDark ? '#141414' : '#ffffff';
    const borderColor = isDark ? '#303030' : '#f0f0f0';
    const subTextColor = isDark ? 'rgba(255, 255, 255, 0.45)' : 'rgba(0, 0, 0, 0.45)';

    const formatOrderDate = (dateStr: string) => {
        const d = new Date(dateStr);
        if (activeTab === 'today') {
            return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false });
        }
        return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
    };

    const fetchOrders = useCallback(async () => {
        setLoading(true);
        try {
            let res: any;
            if (activeTab === 'today') {
                res = await getTodayOrders(symbol);
                const list = Array.isArray(res) ? res : [];
                setOrders(list);
                setTotal(list.length);
            } else {
                res = await getAllOrders({ symbol, page, page_size: 10 });
                setOrders(res?.list || []);
                setTotal(res?.total || 0);
            }
        } catch (e) {
            console.error("数据加载失败:", e);
            setOrders([]);
        } finally {
            setLoading(false);
        }
    }, [activeTab, symbol, page, refreshTrigger]);

    useEffect(() => {
        fetchOrders();
    }, [fetchOrders]);

    useEffect(() => {
        const timer = setInterval(() => fetchOrders(), 5000);
        return () => clearInterval(timer);
    }, [fetchOrders]);

    const handleShowDetail = async (order: any) => {
        setSelectedOrder(order);
        setIsModalOpen(true);
        setDetailLoading(true);
        try {
            const res = await getOrderDetail(order.id);
            if (res) {
                setSelectedOrder(res);
            }
        } catch (e) {
            console.error("获取详情错误:", e);
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
        <ConfigProvider theme={{
            algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
            token: { borderRadius: 4, colorBgContainer: isDark ? '#1d1d1d' : '#ffffff' }
        }}>
            <div style={{ width: '100%', height: '100%', overflow: 'hidden', background: bgColor, display: 'flex', flexDirection: 'column', color: isDark ? '#fff' : '#000' }}>
                <Tabs
                    activeKey={activeTab}
                    onChange={(key) => { setActiveTab(key); setPage(1); setOrders([]); }}
                    size="small"
                    tabBarStyle={{ marginBottom: 0, padding: '0 12px' }}
                    items={[
                        { key: 'today', label: '今日委托' },
                        { key: 'history', label: '历史账本' },
                    ]}
                />

                <div style={{ flex: 1, overflowY: 'auto', padding: '8px' }}>
                    {loading && orders.length === 0 ? (
                        <Flex justify="center" style={{ paddingTop: 40 }}><Spin size="small" /></Flex>
                    ) : (orders.length === 0) ? (
                        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无记录" />
                    ) : (
                        orders.map(order => {
                            const progress = ((order.filled_amount / order.amount) * 100) || 0;
                            const sideColor = order.side === 'buy' ? '#22ab94' : '#f23645';
                            return (
                                <div key={order.id} onClick={() => handleShowDetail(order)} className="order-item-card"
                                     style={{
                                         background: isDark ? '#1d1d1d' : '#ffffff',
                                         border: `1px solid ${borderColor}`,
                                         padding: '10px', marginBottom: '8px', borderRadius: '4px', cursor: 'pointer'
                                     }}>
                                    <Flex justify="space-between" align="baseline" style={{ marginBottom: 6 }}>
                                        <Space size={4}>
                                            <Tag color={order.side === 'buy' ? 'cyan' : 'volcano'} style={{ border: 'none', fontSize: '10px', margin: 0 }}>
                                                {order.side === 'buy' ? '买入' : '卖出'}
                                            </Tag>
                                            <Text strong style={{ fontSize: '13px' }}>{order.symbol}</Text>
                                        </Space>
                                        <Text style={{ fontSize: '11px', color: subTextColor }}>
                                            {formatOrderDate(order.created_at)}
                                        </Text>
                                    </Flex>

                                    <Flex justify="space-between" style={{ fontSize: '12px', marginBottom: 2 }}>
                                        <span style={{ color: subTextColor }}>价格</span>
                                        <span style={{ fontWeight: 500 }}>{Number(order.price).toLocaleString()}</span>
                                    </Flex>
                                    <Flex justify="space-between" style={{ fontSize: '12px', marginBottom: 6 }}>
                                        <span style={{ color: subTextColor }}>数量</span>
                                        <span>{Number(order.amount).toFixed(4)}</span>
                                    </Flex>

                                    <Progress percent={progress} size={[100, 2]} showInfo={false} strokeColor={sideColor} railColor={isDark ? '#333' : '#f0f0f0'} style={{ marginBottom: 8, marginTop: 4 }} />

                                    <Flex justify="space-between" align="center">
                                        <Tag color={STATUS_MAP[order.status]?.color} style={{ fontSize: '10px', borderRadius: '2px' }}>
                                            {STATUS_MAP[order.status]?.text}
                                        </Tag>
                                        {(order.status === 0 || order.status === 1) && (
                                            <Button type="link" danger size="small" onClick={(e) => handleCancel(e, order.id)} style={{ height: 20, padding: 0, fontSize: '12px' }}>
                                                撤单
                                            </Button>
                                        )}
                                    </Flex>
                                </div>
                            );
                        })
                    )}
                </div>

                {activeTab === 'history' && total > 10 && (
                    <div style={{ padding: '8px', textAlign: 'center', borderTop: `1px solid ${borderColor}` }}>
                        <Pagination size="small" current={page} total={total} pageSize={10} onChange={setPage} simple />
                    </div>
                )}

                <ConfigProvider theme={{
                    algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
                    token: { colorBgElevated: isDark ? '#1f1f1f' : '#ffffff' },
                    components: {
                        Descriptions: {
                            colorFillAlter: isDark ? 'rgba(255,255,255,0.04)' : '#fafafa',
                            colorSplit: isDark ? '#303030' : '#f0f0f0',
                        }
                    }
                }}>
                    <Modal
                        title="委托详情"
                        open={isModalOpen}
                        onCancel={() => setIsModalOpen(false)}
                        footer={null}
                        centered
                        width={400}
                        styles={{ body: { padding: '12px 20px 24px' } }}
                    >
                        <Spin spinning={detailLoading}>
                            {selectedOrder && (
                                <Descriptions column={1} size="small" bordered labelStyle={{ width: '100px' }}>
                                    <Descriptions.Item label="订单编号">{selectedOrder.id}</Descriptions.Item>
                                    <Descriptions.Item label="委托类型">{selectedOrder.type?.toUpperCase()}</Descriptions.Item>
                                    <Descriptions.Item label="委托价格">{Number(selectedOrder.price).toLocaleString()}</Descriptions.Item>
                                    <Descriptions.Item label="委托数量">{selectedOrder.amount}</Descriptions.Item>
                                    <Descriptions.Item label="已成交量">{selectedOrder.filled_amount}</Descriptions.Item>
                                    <Descriptions.Item label="创建时间">{new Date(selectedOrder.created_at).toLocaleString()}</Descriptions.Item>
                                    <Descriptions.Item label="交易Hash">
                                        <Text ellipsis copyable={{ icon: <CopyOutlined /> }} style={{ width: 180, fontSize: '11px' }}>
                                            {selectedOrder.msg_hash || '---'}
                                        </Text>
                                    </Descriptions.Item>
                                </Descriptions>
                            )}
                        </Spin>
                    </Modal>
                </ConfigProvider>

                <style>{`
                    .order-item-card:hover {
                        border-color: #1677ff !important;
                        background: ${isDark ? '#262626' : '#f8faff'} !important;
                    }
                    .ant-tabs-nav::before { border-bottom: 1px solid ${borderColor} !important; }
                `}</style>
            </div>
        </ConfigProvider>
    );
};
