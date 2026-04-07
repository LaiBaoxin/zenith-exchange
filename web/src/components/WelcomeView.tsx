import { Typography, Space, Button } from 'antd';

const { Title, Text } = Typography;

interface WelcomeProps {
    onConnect: () => void;
    isDark: boolean;
    loading: boolean;
}

const WelcomeView = ({ onConnect, isDark, loading }: WelcomeProps) => {
    return (
        <div style={{
            flex: 1, display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center',
            textAlign: 'center', background: isDark ? '#000' : '#f5f5f5', height: '100vh'
        }}>
            <Space direction="vertical" size={24} style={{ zIndex: 1, maxWidth: '600px' }}>
                <Title level={1} style={{ margin: 0, fontSize: '48px', color: isDark ? '#fff' : '#000' }}>
                    ZENITH EXCHANGE
                </Title>

                <Text style={{ fontSize: '18px', color: isDark ? 'rgba(255,255,255,0.65)' : 'rgba(0,0,0,0.65)' }}>
                    下一代去中心化高频交易平台
                </Text>

                {/* 这里是重点：
                  1. onClick 直接绑定父组件传来的 handleConnect
                  2. loading 状态由父组件统一控制，授权时按钮会自动转圈并禁用
                */}
                <Button
                    type="primary"
                    size="large"
                    loading={loading}
                    onClick={onConnect}
                    style={{
                        width: '240px', height: '56px', fontSize: '18px', fontWeight: 'bold', borderRadius: '8px'
                    }}
                >
                    {loading ? '正在验证身份...' : '立即连接钱包登录'}
                </Button>

                <div style={{ marginTop: '24px' }}>
                    <Text type="secondary" style={{ fontSize: '12px' }}>
                        支持 MetaMask / WalletConnect · 资产由智能合约托管
                    </Text>
                </div>
            </Space>
        </div>
    );
};

export default WelcomeView;
