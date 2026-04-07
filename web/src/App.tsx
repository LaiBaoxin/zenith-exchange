import { ConfigProvider, theme, App as AntdApp } from 'antd';
import CexLayout from './components/CexLayout';

const App = () => {
    return (
        <ConfigProvider
            theme={{
                algorithm: theme.darkAlgorithm, // 强制暗黑模式
                token: {
                    colorPrimary: '#1677ff', // Zenith 蓝
                    colorBgBase: '#000000',  // 纯黑背景
                    colorBgContainer: '#141414', // 容器背景
                    borderRadius: 4, // 偏方正的专业交易风格
                },
            }}
        >
            <AntdApp>
                <CexLayout />
            </AntdApp>
        </ConfigProvider>
    );
};

export default App;
