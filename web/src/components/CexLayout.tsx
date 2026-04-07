import { useState, useEffect, useCallback } from 'react';
import { Layout, Button, ConfigProvider, theme, Typography, App } from 'antd';
import { SunOutlined, MoonOutlined, LogoutOutlined } from '@ant-design/icons';
import { useConnect, useAccount, useSignMessage, useDisconnect } from 'wagmi';
import { injected } from 'wagmi/connectors';
import TradingTerminal from './TradingTerminal';
import WelcomeView from './WelcomeView';
import { login } from '../api/system.ts';

const { Header, Content } = Layout;
const { Text } = Typography;

const CexLayout = () => {
    const { message: msgApi } = App.useApp();
    const { connectAsync } = useConnect();
    const { disconnect } = useDisconnect();
    const { isConnected: isWagmiConnected } = useAccount();
    const { signMessageAsync } = useSignMessage();

    const [isDark, setIsDark] = useState(() => localStorage.getItem('theme') !== 'light');
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [displayAddress, setDisplayAddress] = useState<string>('');
    // 增加全局 loading 状态，防止重复点击 [修复问题二]
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        const token = localStorage.getItem('zenith_auth_token');
        const storedAddr = localStorage.getItem('user_address');

        if (isWagmiConnected && token && storedAddr) {
            setIsLoggedIn(true);
            setDisplayAddress(storedAddr);
        } else {
            setIsLoggedIn(false);
        }
    }, [isWagmiConnected]);

    const toggleTheme = useCallback(() => {
        setIsDark(prev => {
            const next = !prev;
            localStorage.setItem('theme', next ? 'dark' : 'light');
            return next;
        });
    }, []);

    const handleConnect = async () => {
        if (loading) return; // 防止并发点击

        setLoading(true);
        const hide = msgApi.loading('身份验证中...', 0);

        try {
            // 连接钱包
            const conn = await connectAsync({ connector: injected() });
            const walletAddress = conn.accounts[0];

            // 签名验证
            const messageToSign = `Welcome to Zenith Exchange!\nAddress: ${walletAddress}\nTimestamp: ${Date.now()}`;
            await signMessageAsync({ message: messageToSign });

            // 登录后端
            const res = await login(walletAddress);
            console.log("Res",res)
            if (res && res.token) {

                const finalAddr = res.walletAddress || walletAddress;

                // 统一存储键名
                window.localStorage.setItem('zenith_auth_token', res.token);
                window.localStorage.setItem('user_address', finalAddr);

                setDisplayAddress(finalAddr);
                setIsLoggedIn(true);
                msgApi.success('登录成功');
            } else {
                console.error("响应结构异常:", res);
                throw new Error('未获取到有效 Token，请检查后端返回结构');
            }
        } catch (error: any) {
            console.error("Login Error:", error);
            // 捕获签名拒绝或业务错误
            const errorMsg = error.code === 4001 ? "用户取消签名" : (error.message || "登录失败");
            msgApi.error(errorMsg);
        } finally {
            setLoading(false);
            hide();
        }
    };

    const handleLogout = () => {
        localStorage.removeItem('zenith_auth_token');
        localStorage.removeItem('user_address');
        disconnect();
        setIsLoggedIn(false);
        msgApi.info('已安全退出');
    };

    return (
        <ConfigProvider theme={{ algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm }}>
            <Layout style={{ minHeight: '100vh' }}>
                <Header style={{
                    display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                    padding: '0 24px', background: isDark ? '#141414' : '#fff',
                    borderBottom: `1px solid ${isDark ? '#303030' : '#f0f0f0'}`, zIndex: 10
                }}>
                    <div style={{ fontSize: '18px', fontWeight: 'bold', color: isDark ? '#fff' : '#1677ff', letterSpacing: '1px' }}>
                        ZENITH EXCHANGE
                    </div>

                    <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
                        <Button type="text" icon={isDark ? <SunOutlined /> : <MoonOutlined />} onClick={toggleTheme} />
                        {isLoggedIn && displayAddress ? (
                            <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                                <Text strong style={{ fontFamily: 'monospace' }}>
                                    {displayAddress.slice(0, 6)}...{displayAddress.slice(-4)}
                                </Text>
                                <Button type="text" size="small" icon={<LogoutOutlined />} onClick={handleLogout} danger />
                            </div>
                        ) : (
                            // 头部按钮同步 loading 状态
                            <Button type="primary" shape="round" onClick={handleConnect} loading={loading}>
                                连接钱包
                            </Button>
                        )}
                    </div>
                </Header>

                <Content style={{ display: 'flex', flexDirection: 'column', background: isDark ? '#000' : '#f0f2f5' }}>
                    {isLoggedIn ? (
                        <TradingTerminal isDark={isDark} />
                    ) : (
                        <WelcomeView
                            onConnect={handleConnect}
                            isDark={isDark}
                            loading={loading}
                        />
                    )}
                </Content>
            </Layout>
        </ConfigProvider>
    );
};

export default CexLayout;
