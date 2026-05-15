import { useEffect, useState, useCallback, useMemo } from 'react';
import {
    Tabs, Typography, InputNumber, Button, Flex,
    Space, ConfigProvider, theme, App, Modal
} from 'antd';
import { WalletOutlined, PlusOutlined, ArrowDownOutlined } from '@ant-design/icons';
import { colors } from '../assets/css/TradePanel.styles.ts';
import { getBalance } from '../api/assets';
import { placeOrder } from '../api/order';
import { useWithdraw } from '../hooks/useWithdraw.ts';
import { useDeposit } from '../hooks/useDeposit.ts';
import type { SystemConfig } from "../api/system.ts";

const { Text } = Typography;

interface TradePanelProps {
    isDark: boolean;
    symbol?: string;
    config: SystemConfig | null;
    onTradeSuccess?: () => void;
}

export const TradePanel = ({ isDark, symbol = "BTC_USDT", config, onTradeSuccess }: TradePanelProps) => {
    const { message } = App.useApp();
    const [activeTab, setActiveTab] = useState('1');
    const [amount, setAmount] = useState<number | null>(0);
    const [allBalances, setAllBalances] = useState<any[]>([]);
    const [price, setPrice] = useState<number | null>(0);
    const [orderLoading, setOrderLoading] = useState(false);
    const [isDepositModalOpen, setIsDepositModalOpen] = useState(false);
    const [depositAmount, setDepositAmount] = useState<number | null>(0);

    const themeMode = isDark ? 'dark' : 'light';
    const c = colors[themeMode];

    const [baseCurrency, quoteCurrency] = useMemo(() => {
        const parts = symbol.split('_');
        return [parts[0] || 'BTC', parts[1] || 'USDT'];
    }, [symbol]);

    // 提现与充值 Hook
    const { handleWithdraw, loading: withdrawLoading } = useWithdraw(
        config?.vault_address || "",
        config?.token_address || "",
        () => updateBalance()
    );

    const { handleDeposit, loading: depositLoading } = useDeposit(
        config?.vault_address || "",
        config?.token_address || "",
        () => {
            updateBalance();
            setIsDepositModalOpen(false);
            setDepositAmount(0);
        }
    );

    const getBalanceBySymbol = useCallback((curName: string) => {
        const asset = allBalances.find(a => a.currency?.toUpperCase() === curName?.toUpperCase());
        return parseFloat(asset?.available || "0");
    }, [allBalances]);

    const currentAvailableValue = useMemo(() => {
        if (activeTab === '1') return getBalanceBySymbol(quoteCurrency);
        if (activeTab === '3') return getBalanceBySymbol(quoteCurrency); // 提现: 显示 USDT (链上代币映射)
        return getBalanceBySymbol(baseCurrency);
    }, [activeTab, quoteCurrency, baseCurrency, getBalanceBySymbol]);

    const updateBalance = useCallback(async () => {
        try {
            const res = await getBalance() as any;
            const data = res.data || res;
            if (Array.isArray(data)) setAllBalances(data);
        } catch (e) { console.error("刷新余额失败"); }
    }, []);

    useEffect(() => {
        updateBalance();
        const timer = setInterval(updateBalance, 5000);
        return () => clearInterval(timer);
    }, [updateBalance]);

    const handlePlaceOrder = async (side: 'buy' | 'sell') => {
        if (!amount || amount <= 0) return message.warning("请输入数量");
        setOrderLoading(true);
        try {
            await placeOrder({ symbol, side, price: price || 0, amount });
            message.success("订单提交成功");
            setAmount(0);
            updateBalance();
            if (onTradeSuccess) onTradeSuccess();
        } catch (e: any) {
            message.error(e.response?.data?.msg || "交易失败");
        } finally { setOrderLoading(false); }
    };

    const renderTradeForm = (side: 'buy' | 'sell') => {
        const balanceCur = side === 'buy' ? quoteCurrency : baseCurrency;
        const isBalanceEmpty = currentAvailableValue <= 0;

        return (
            <Flex vertical style={{ padding: '12px' }} gap={12}>
                <Flex vertical gap={6}>
                    <Flex justify="space-between" align="center">
                        <Space size={4}>
                            <WalletOutlined style={{ color: c.secondary, fontSize: '11px' }} />
                            <Text style={{ color: c.secondary, fontSize: '11px' }}>可用余额</Text>
                        </Space>
                        <Space>
                            <Text strong style={{ color: c.text, fontSize: '12px' }}>
                                {currentAvailableValue.toFixed(4)} {balanceCur}
                            </Text>
                            <Button size="small" type="link" style={{fontSize: '10px'}} icon={<PlusOutlined />} onClick={() => setIsDepositModalOpen(true)}>充值</Button>
                        </Space>
                    </Flex>

                    <div style={{ border: `1px solid ${c.border}`, borderRadius: '4px', background: c.innerBg, padding: '4px 8px' }}>
                        <Flex align="center">
                            <Text type="secondary" style={{ fontSize: '12px', width: '40px' }}>价格</Text>
                            <InputNumber value={price} onChange={v => setPrice(v)} variant="borderless" controls={false} style={{ flex: 1, color: c.text, fontWeight: 'bold' }} />
                            <Text style={{ fontSize: '12px', color: c.secondary }}>{quoteCurrency}</Text>
                        </Flex>
                    </div>

                    <div style={{ border: `1px solid ${c.border}`, borderRadius: '4px', background: c.innerBg, padding: '4px 8px' }}>
                        <Flex align="center">
                            <Text type="secondary" style={{ fontSize: '12px', width: '40px' }}>数量</Text>
                            <InputNumber value={amount} onChange={v => setAmount(v)} variant="borderless" controls={false} style={{ flex: 1, color: c.text, fontSize: '16px', fontWeight: 'bold' }} />
                            <Text strong style={{ color: c.text, fontSize: '12px' }}>{baseCurrency}</Text>
                        </Flex>
                    </div>

                    <Button
                        type="primary"
                        block
                        loading={orderLoading}
                        disabled={isBalanceEmpty || (amount || 0) > currentAvailableValue}
                        onClick={() => handlePlaceOrder(side)}
                        style={{ height: '40px', fontWeight: 'bold', background: side === 'buy' ? '#22ab94' : '#f23645', borderColor: 'transparent' }}
                    >
                        {isBalanceEmpty ? "余额不足，请先充值" : (side === 'buy' ? `买入 ${baseCurrency}` : `卖出 ${baseCurrency}`)}
                    </Button>
                </Flex>
            </Flex>
        );
    };

    return (
        <ConfigProvider
            theme={{
                algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
                components: {
                    Modal: {
                        // 强制 Modal 内容背景色跟随主题，解决“白边”问题
                        contentBg: isDark ? '#141414' : '#fff',
                        headerBg: isDark ? '#141414' : '#fff',
                    }
                }
            }}
        >
            <div style={{ background: isDark ? '#000' : '#fff', height: '100%', borderTop: `1px solid ${c.border}` }}>
                <Tabs centered activeKey={activeTab} onChange={key => { setActiveTab(key); setAmount(0); }}
                      items={[
                          { key: '1', label: '买入', children: renderTradeForm('buy') },
                          { key: '2', label: '卖出', children: renderTradeForm('sell') },
                          {
                              key: '3',
                              label: '提现',
                              children: (
                                  <Flex vertical style={{padding: '20px'}} gap={16}>
                                      <Flex justify="space-between" align="center">
                                          <Text style={{ color: c.secondary, fontSize: '11px' }}>可提现余额 (链上代币)</Text>
                                          <Text strong style={{ color: c.text, fontSize: '12px' }}>
                                              {currentAvailableValue.toFixed(4)} {quoteCurrency}
                                          </Text>
                                      </Flex>

                                      <div style={{ border: `1px solid ${c.border}`, borderRadius: '4px', background: c.innerBg, padding: '8px' }}>
                                          <InputNumber
                                              value={amount}
                                              onChange={v => setAmount(v)}
                                              variant="borderless"
                                              placeholder="输入提现数量"
                                              style={{ width: '100%', fontSize: '16px', color: c.text, fontWeight: 'bold' }}
                                          />
                                          <Flex justify="flex-end">
                                              <Button size="small" type="link" onClick={() => setAmount(currentAvailableValue)}>全部提现</Button>
                                          </Flex>
                                      </div>

                                      <Button
                                          danger
                                          type="primary"
                                          block
                                          size="large"
                                          disabled={!config?.vault_address || currentAvailableValue <= 0 || (amount || 0) > currentAvailableValue}
                                          loading={withdrawLoading}
                                          onClick={() => {
                                              if (!amount || amount <= 0) return message.warning("请输入提现数量");
                                              handleWithdraw(amount.toString(), quoteCurrency);
                                          }}
                                          style={{ height: '40px', fontWeight: 'bold' }}
                                      >
                                          {currentAvailableValue <= 0 ? "暂无可用余额提现" : `确认提现 ${quoteCurrency}`}
                                      </Button>
                                  </Flex>
                              )
                          },
                      ]}
                />

                <Modal
                    title={<Text style={{ color: c.text }}>充值资产</Text>}
                    open={isDepositModalOpen}
                    onCancel={() => setIsDepositModalOpen(false)}
                    footer={null}
                    destroyOnClose
                    centered
                    // v6 版本的 styles 接口，修正了背景边框显示逻辑
                    styles={{
                        body: {
                            padding: '10px 0',
                            backgroundColor: 'transparent',
                        },
                        mask: {
                            backdropFilter: 'blur(4px)',
                        }
                    }}
                >
                    <Flex vertical gap={20}>
                        <Flex vertical gap={4}>
                            <Text type="secondary" style={{ fontSize: '12px' }}>代币将从小狐狸钱包存入交易所金库以增加可用余额</Text>
                            <Text strong style={{ color: c.text }}>币种: {quoteCurrency} (链上代币)</Text>
                        </Flex>

                        <div style={{
                            border: `1px solid ${c.border}`,
                            borderRadius: '8px',
                            background: isDark ? '#1f1f1f' : '#f5f5f5',
                            padding: '12px 16px'
                        }}>
                            <Text type="secondary" style={{ fontSize: '11px', display: 'block', marginBottom: '4px' }}>充值数量</Text>
                            <Flex align="center">
                                <InputNumber
                                    variant="borderless"
                                    style={{
                                        flex: 1,
                                        fontSize: '24px',
                                        fontWeight: 'bold',
                                        color: c.text,
                                        padding: 0
                                    }}
                                    placeholder="0.00"
                                    min={0}
                                    autoFocus
                                    value={depositAmount}
                                    onChange={v => setDepositAmount(v)}
                                    controls={false}
                                />
                                <Text strong style={{ color: c.text, fontSize: '16px' }}>
                                    {quoteCurrency}
                                </Text>
                            </Flex>
                        </div>

                        <Button
                            type="primary"
                            block
                            size="large"
                            loading={depositLoading}
                            disabled={!depositAmount || depositAmount <= 0}
                            onClick={() => handleDeposit(depositAmount!.toString())}
                            style={{
                                height: '50px',
                                background: '#22ab94',
                                borderColor: 'transparent',
                                fontWeight: 'bold',
                                borderRadius: '8px'
                            }}
                        >
                            {depositLoading ? '正在链上处理...' : '立即充值'}
                        </Button>
                        <Text type="secondary" style={{ fontSize: '11px', textAlign: 'center' }}>
                            <ArrowDownOutlined /> 交易将通过 MetaMask 授权上链
                        </Text>
                    </Flex>
                </Modal>
            </div>
        </ConfigProvider>
    );
};
