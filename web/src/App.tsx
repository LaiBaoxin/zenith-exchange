import React, { useState } from 'react';
import { Card, Input, Button, Typography, Space, message, Divider, Tabs, Statistic } from 'antd';
import { DownloadOutlined, UploadOutlined, SafetyCertificateOutlined } from '@ant-design/icons';
import { CONTRACT_ADDRESSES } from './constants/addresses';
import { useWithdraw } from './hooks/useWithdraw';
import { useAccount, useWriteContract, useReadContract } from 'wagmi';
import { parseEther, formatEther } from 'viem';
import { ABIS } from './constants/contracts_abi';

const { Text } = Typography;

const App: React.FC = () => {
    const { address, isConnected } = useAccount();
    const [amount, setAmount] = useState<string>("10.0");
    const { writeContractAsync, isPending: txLoading } = useWriteContract();

    // 1. 提现 Hook (需要你内部实现 axios 调用 Go 后端)
    const { handleWithdraw, loading: withdrawing } = useWithdraw(
        CONTRACT_ADDRESSES.VAULT,
        CONTRACT_ADDRESSES.ZNT_TOKEN
    );

    // 2. 读取合约内的存款余额
    const { data: vaultBalance } = useReadContract({
        address: CONTRACT_ADDRESSES.VAULT,
        abi: ABIS.VAULT,
        functionName: 'balances',
        args: address ? [address, CONTRACT_ADDRESSES.ZNT_TOKEN] : undefined,
    });

    // 3. 读取 Nonce (用于提现校验)
    const { data: nonce } = useReadContract({
        address: CONTRACT_ADDRESSES.VAULT,
        abi: ABIS.VAULT,
        functionName: 'nonces',
        args: address ? [address] : undefined,
    });

    // 4. 授权逻辑
    const onApprove = async () => {
        try {
            await writeContractAsync({
                address: CONTRACT_ADDRESSES.ZNT_TOKEN,
                abi: ABIS.TOKEN,
                functionName: 'approve',
                args: [CONTRACT_ADDRESSES.VAULT, parseEther(amount)],
            });
            message.success("授权成功！");
        } catch (e: any) { message.error("授权错误: " + e.shortMessage); }
    };

    // 5. 存入逻辑
    const onDeposit = async () => {
        try {
            await writeContractAsync({
                address: CONTRACT_ADDRESSES.VAULT,
                abi: ABIS.VAULT,
                functionName: 'deposit',
                args: [CONTRACT_ADDRESSES.ZNT_TOKEN, parseEther(amount)],
            });
            message.success("存入成功！");
        } catch (e: any) { message.error("存入失败: " + e.shortMessage); }
    };

    const onWithdraw = async () => {
        if (!amount || parseFloat(amount) <= 0) return message.warning("请输入有效金额");
        await handleWithdraw(amount, Number(nonce || 0));
    };

    return (
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh', backgroundColor: '#f0f2f5' }}>
            <Card
                style={{ width: 500, borderRadius: 16, boxShadow: '0 10px 25px rgba(0,0,0,0.05)' }}
                title={<Space><SafetyCertificateOutlined style={{ color: '#1890ff' }} /> <Text strong>Zenith Exchange 资产中心</Text></Space>}
            >
                <Space direction="vertical" size="large" style={{ width: '100%' }}>
                    <div style={{ padding: '12px', background: '#fafafa', borderRadius: 8 }}>
                        <Text type="secondary" style={{ fontSize: 12 }}>当前连接钱包</Text>
                        <div style={{ fontWeight: 'bold', fontSize: '13px' }}>{isConnected ? address : '未连接钱包'}</div>
                    </div>

                    <Tabs defaultActiveKey="1" centered items={[
                        {
                            key: '1',
                            label: '存入资产',
                            children: (
                                <div style={{ paddingTop: 16 }}>
                                    <Text type="secondary">输入 ZNT 数量</Text>
                                    <Input size="large" prefix="ZNT" value={amount} onChange={e => setAmount(e.target.value)} style={{ margin: '8px 0 20px 0' }} />
                                    <Button block icon={<UploadOutlined />} onClick={onApprove} loading={txLoading} style={{ marginBottom: 12 }}>
                                        1. 授权金库 (Approve)
                                    </Button>
                                    <Button type="primary" block size="large" onClick={onDeposit} loading={txLoading}>
                                        2. 确认存入 (Deposit)
                                    </Button>
                                </div>
                            )
                        },
                        {
                            key: '2',
                            label: '安全提现',
                            children: (
                                <div style={{ paddingTop: 16 }}>
                                    <Statistic title="合约中可用余额" value={vaultBalance ? formatEther(vaultBalance as bigint) : '0.00'} suffix="ZNT" />
                                    <Divider />
                                    <Button type="primary" danger block size="large" icon={<DownloadOutlined />} loading={withdrawing} onClick={onWithdraw}>
                                        安全提现 (后端签名验证)
                                    </Button>
                                </div>
                            )
                        }
                    ]} />
                </Space>
            </Card>
        </div>
    );
};

export default App;
