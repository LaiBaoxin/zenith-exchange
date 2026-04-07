import { useState } from 'react';
import { Tabs, Typography, InputNumber, Button, Flex, Divider, Row, Col, Space, ConfigProvider, theme } from 'antd';
import { WalletOutlined, PlusOutlined, MinusOutlined, CopyOutlined } from '@ant-design/icons';
import { colors } from '../assets/css/TradePanel.styles.ts';

const { Text } = Typography;

export const TradePanel = ({ isDark }: { isDark: boolean }) => {
    const [activeTab, setActiveTab] = useState('1');
    const [amount, setAmount] = useState<number | null>(1.0);
    const [percent, setPercent] = useState<number | null>(null);

    const themeMode = isDark ? 'dark' : 'light';
    const c = colors[themeMode];
    const balance = 1.2540;

    const renderForm = () => (
        <Flex vertical style={{ padding: '8px 12px' }} gap={8}>
            {/* 余额与输入 */}
            <Flex vertical gap={6}>
                <Flex justify="space-between" align="center">
                    <Space size={4}><WalletOutlined style={{ color: c.secondary, fontSize: '11px' }} /><Text style={{ color: c.secondary, fontSize: '11px' }}>余额</Text></Space>
                    <Text strong style={{ color: c.text, fontSize: '12px' }}>{balance} BTC</Text>
                </Flex>

                <div style={{ border: `1px solid ${c.border}`, borderRadius: '4px', background: c.innerBg, padding: '2px 8px' }}>
                    <Flex align="center">
                        <InputNumber value={amount} onChange={(v) => {setAmount(v); setPercent(null);}} variant="borderless" controls={false} style={{ flex: 1, fontSize: '16px', fontWeight: '600', color: c.text }} />
                        <Text strong style={{ color: c.text, fontSize: '12px', marginRight: '6px' }}>BTC</Text>
                        <Divider type="vertical" style={{ height: '16px', margin: 0 }} />
                        <Space size={8} style={{ marginLeft: '8px' }}>
                            <PlusOutlined style={{ cursor: 'pointer', color: '#22ab94', fontSize: '12px' }} onClick={() => setAmount((amount || 0) + 0.1)} />
                            <MinusOutlined style={{ cursor: 'pointer', color: c.secondary, fontSize: '12px' }} onClick={() => setAmount(Math.max(0, (amount || 0) - 0.1))} />
                        </Space>
                    </Flex>
                </div>

                <Row gutter={4}>
                    {[25, 50, 100].map(p => (
                        <Col span={8} key={p}>
                            <Button block size="small" style={{ fontSize: '10px', height: '22px', background: percent === p ? '#22ab94' : 'transparent', color: percent === p ? '#fff' : c.secondary }} onClick={() => {setPercent(p); setAmount(Number((balance * p / 100).toFixed(4)));}}>
                                {p}%
                            </Button>
                        </Col>
                    ))}
                </Row>
            </Flex>

            {/* 合约与按钮 */}
            <Flex vertical gap={8}>
                <div style={{ border: `1px solid ${isDark ? '#22ab9422' : '#f0f0f0'}`, padding: '4px 8px', borderRadius: '4px', background: isDark ? 'rgba(34, 171, 148, 0.03)' : '#fafafa' }}>
                    <Flex align="center" justify="space-between">
                        <Text style={{ fontSize: '10px', color: c.secondary }}>Vault: <Text strong style={{ color: '#22ab94', fontFamily: 'monospace' }}>0x5fbD...0aa3</Text></Text>
                        <CopyOutlined style={{ color: '#22ab94', cursor: 'pointer', fontSize: '10px' }} />
                    </Flex>
                </div>
                <Flex gap={6}>
                    <Button danger type="primary" style={{ flex: 1, height: '30px', fontSize: '12px' }}>授权</Button>
                    <Button type="primary" style={{ flex: 2, height: '30px', fontSize: '12px', fontWeight: 'bold', background: '#22ab94', borderColor: '#22ab94' }}>确认存入</Button>
                </Flex>
            </Flex>
        </Flex>
    );

    return (
        <ConfigProvider theme={{ algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm }}>
            <Tabs centered activeKey={activeTab} onChange={setActiveTab} tabBarStyle={{ marginBottom: 0, height: '32px' }}
                  items={[
                      { key: '1', label: <span style={{ fontSize: '12px' }}>存入</span>, children: renderForm() },
                      { key: '2', label: <span style={{ fontSize: '12px' }}>提现</span>, children: <div style={{ padding: 10, fontSize: '11px', textAlign: 'center', color: c.secondary }}>建设中</div> },
                  ]}
            />
        </ConfigProvider>
    );
};
