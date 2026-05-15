import { useState, useEffect, useCallback, useRef } from 'react';
import { Layout, Button, ConfigProvider, theme, Typography, App, Spin, Flex } from 'antd';
import { SunOutlined, MoonOutlined, LogoutOutlined } from '@ant-design/icons';
import { useConnect, useAccount, useSignMessage, useDisconnect } from 'wagmi';
import { injected } from 'wagmi/connectors';
import TradingTerminal from './TradingTerminal';
import WelcomeView from './WelcomeView';
import { login, getSystemConfig, type SystemConfig } from '../api/system.ts';

const { Header, Content } = Layout;
const { Text } = Typography;

const CexLayout = () => {
    const { message: msgApi } = App.useApp();
    const { connectAsync } = useConnect();
    const { disconnect } = useDisconnect();
    const { address, isConnected: isWagmiConnected } = useAccount();
    const { signMessageAsync } = useSignMessage();

    const [isDark, setIsDark] = useState(() => localStorage.getItem('theme') !== 'light');
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [displayAddress, setDisplayAddress] = useState<string>('');
    const [loading, setLoading] = useState(false);

    const [config, setConfig] = useState<SystemConfig | null>(null);
    const [configLoading, setConfigLoading] = useState(true);
    const [showLogin, setShowLogin] = useState(false);

    const isProcessingLogin = useRef(false);

    useEffect(() => {
        document.body.style.backgroundColor = isDark ? '#000000' : '#f0f2f5';
    }, [isDark]);

    useEffect(() => {
        const onExpired = () => {
            setIsLoggedIn(false);
            setShowLogin(false);
        };
        window.addEventListener('auth:expired', onExpired);
        return () => window.removeEventListener('auth:expired', onExpired);
    }, []);

    useEffect(() => {
        const fetchConfig = async () => {
            try {
                const res = await getSystemConfig();
                if (res) setConfig(res as any);
            } catch (e) {
                msgApi.error("基础配置加载失败");
            } finally {
                setConfigLoading(false);
            }
        };
        fetchConfig();
    }, [msgApi]);

    useEffect(() => {
        if (isProcessingLogin.current) return;
        const token = localStorage.getItem('zenith_auth_token');
        const storedAddr = localStorage.getItem('user_address');

        if (isWagmiConnected && address && token && storedAddr?.toLowerCase() === address.toLowerCase()) {
            setIsLoggedIn(true);
            setDisplayAddress(address);
        } else {
            setIsLoggedIn(false);
        }
    }, [isWagmiConnected, address]);

    const toggleTheme = useCallback(() => {
        setIsDark(prev => {
            const next = !prev;
            localStorage.setItem('theme', next ? 'dark' : 'light');
            return next;
        });
    }, []);

    const handleLoginProcess = async () => {
        if (loading) return;
        isProcessingLogin.current = true;
        setLoading(true);
        const hide = msgApi.loading('正在进行安全验证...', 0);
        try {
            let currentAddr = address;
            if (!isWagmiConnected) {
                const conn = await connectAsync({ connector: injected() });
                currentAddr = conn.accounts[0];
            }
            if (!currentAddr) throw new Error("未能获取钱包地址");
            const messageToSign = `Welcome to Zenith Exchange!\nAddress: ${currentAddr}\nTimestamp: ${Date.now()}`;
            await signMessageAsync({ message: messageToSign });
            const res: any = await login(currentAddr);
            if (res && res.token) {
                localStorage.setItem('zenith_auth_token', res.token);
                localStorage.setItem('user_address', currentAddr);
                setDisplayAddress(currentAddr);
                setIsLoggedIn(true);
                setShowLogin(false);
                msgApi.success('身份验证成功');
            }
        } catch (error: any) {
            msgApi.error(error.message || "登录失败");
            disconnect();
        } finally {
            setLoading(false);
            hide();
            setTimeout(() => { isProcessingLogin.current = false; }, 500);
        }
    };

    const handleLogout = () => {
        localStorage.removeItem('zenith_auth_token');
        localStorage.removeItem('user_address');
        disconnect();
        setIsLoggedIn(false);
        setShowLogin(false);
        msgApi.info('已安全退出');
    };

    return (
        <ConfigProvider theme={{
            algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
            token: {
                colorBgLayout: isDark ? '#000000' : '#f0f2f5',
                colorBgBase: isDark ? '#000000' : '#ffffff',
            },
            components: {
                Modal: {
                    contentBg: isDark ? '#141414' : '#ffffff',
                    headerBg: isDark ? '#141414' : '#ffffff',
                },
                Layout: {
                    headerBg: isDark ? '#141414' : '#ffffff',
                }
            }
        }}>
            <App style={{ minHeight: '100vh', backgroundColor: isDark ? '#000000' : '#f0f2f5' }}>
                <Layout style={{ minHeight: '100vh', background: 'transparent' }}>
                    <Header style={{
                        display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                        padding: '0 24px',
                        borderBottom: `1px solid ${isDark ? '#303030' : '#f0f0f0'}`,
                        zIndex: 10,
                        transition: 'all 0.3s'
                    }}>
                        <div style={{ fontSize: '18px', fontWeight: 'bold', color: isDark ? '#fff' : '#1677ff' }}>
                            ZENITH EXCHANGE
                        </div>

                        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
                            <Button
                                type="text"
                                icon={isDark ? <SunOutlined /> : <MoonOutlined />}
                                onClick={toggleTheme}
                            />
                            {isLoggedIn ? (
                                <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                                    <Text strong style={{ fontFamily: 'monospace' }}>
                                        {displayAddress.slice(0, 6)}...{displayAddress.slice(-4)}
                                    </Text>
                                    <Button type="text" size="small" icon={<LogoutOutlined />} onClick={handleLogout} danger />
                                </div>
                            ) : (
                                <Button type="primary" shape="round" onClick={() => setShowLogin(true)}>
                                    连接钱包
                                </Button>
                            )}
                        </div>
                    </Header>

                    <Content style={{
                        display: 'flex',
                        flexDirection: 'column',
                        flex: 1,
                        transition: 'all 0.3s'
                    }}>
                        {showLogin && !isLoggedIn ? (
                            <WelcomeView onConnect={handleLoginProcess} isDark={isDark} loading={loading} />
                        ) : configLoading ? (
                            <Flex justify="center" align="center" style={{ flex: 1 }}>
                                <Spin size="large" tip="正在初始化交易所配置..." />
                            </Flex>
                        ) : (
                            <TradingTerminal isDark={isDark} config={config} />
                        )}
                    </Content>
                </Layout>
            </App>
        </ConfigProvider>
    );
};

export default CexLayout;
