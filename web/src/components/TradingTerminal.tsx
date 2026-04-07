import { Layout, Flex, Typography } from 'antd';
import { TradePanel } from "./TradePanel.tsx";
import { OrderBook } from "./OrderBook.tsx";

const { Text } = Typography;

export default function TradingTerminal({ isDark }: { isDark: boolean }) {
    const bgColor = isDark ? '#000000' : '#ffffff';
    const borderColor = isDark ? '#1f1f1f' : '#f0f0f0';

    const panelStyle: React.CSSProperties = {
        background: bgColor,
        border: `1px solid ${borderColor}`,
        borderRadius: '4px',
        display: 'flex',
        flexDirection: 'column',
        overflow: 'hidden'
    };

    return (
        <Layout style={{ height: '100vh', background: bgColor, overflow: 'hidden', padding: '8px' }}>
            <Flex style={{ height: '100%', gap: '8px' }}>

                {/* 左侧主区域 */}
                <Flex vertical style={{ flex: 1, gap: '8px', minWidth: 0 }}>
                    {/* K线区域：flex: 1 会自动占据剩余的所有高度 */}
                    <div style={{ ...panelStyle, flex: 1, alignItems: 'center', justifyContent: 'center' }}>
                        <Text style={{ color: isDark ? '#1a1a1a' : '#f5f5f5', fontSize: '24px', letterSpacing: '8px' }}>K-LINE CHART</Text>
                    </div>

                    {/* 交易面板：高度设为 auto，防止撑破屏幕产生滚动 */}
                    <div style={{ ...panelStyle, height: 'auto', flexShrink: 0 }}>
                        <TradePanel isDark={isDark} />
                    </div>
                </Flex>

                {/* 右侧盘口：宽度稍微收窄，更紧凑 */}
                <div style={{ width: '260px', ...panelStyle, flexShrink: 0 }}>
                    <OrderBook isDark={isDark} />
                </div>
            </Flex>
        </Layout>
    );
}
