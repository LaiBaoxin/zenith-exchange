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
        injected(), // 自动识别 MetaMask
    ],
    transports: {
        [anvil.id]: http(import.meta.env.VITE_ANVIL_URL || 'http://127.0.0.1:8545'),
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
