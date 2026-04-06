import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';

import { WagmiProvider, createConfig, http } from 'wagmi';
import { anvil } from 'wagmi/chains';
import { injected } from 'wagmi/connectors';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient();

// 完整的 Wagmi 配置
const config = createConfig({
    chains: [anvil],
    connectors: [
        injected(), // 允许连接 MetaMask 或浏览器钱包
    ],
    transports: {
        [anvil.id]: http(import.meta.env.VITE_ANVIL_URL),
    },
});

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <WagmiProvider config={config}>
            <QueryClientProvider client={queryClient}>
                <App />
            </QueryClientProvider>
        </WagmiProvider>
    </React.StrictMode>
);
