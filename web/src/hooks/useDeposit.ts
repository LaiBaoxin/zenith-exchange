import { useState } from 'react';
import { ethers } from 'ethers';
import { App } from 'antd';
import { ABIS } from '../constants/contracts_abi';

export const useDeposit = (
    vaultAddress: string,
    tokenAddress: string,
    onSuccess?: () => void
) => {
    const [loading, setLoading] = useState(false);
    const [loadingText, setLoadingText] = useState('');
    const { message, modal } = App.useApp();

    const handleDeposit = async (amountStr: string) => {
        // 1. 基础环境校验
        if (!(window as any).ethereum) {
            return modal.error({ title: '未检测到钱包', content: '请安装并登录 MetaMask' });
        }

        if (!vaultAddress || vaultAddress === "0x" || !tokenAddress) {
            return message.error("系统合约地址尚未加载，请确保部署脚本已运行");
        }

        setLoading(true);
        const msgKey = 'deposit_process';

        try {
            // 2. 初始化连接 - 使用 "any" 模式防止 Anvil 重启导致的 ChainId 报错
            const provider = new ethers.BrowserProvider((window as any).ethereum, "any");
            const signer = await provider.getSigner();
            const userAddress = await signer.getAddress();
            const amountWei = ethers.parseUnits(amountStr, 18); // 默认 18 位精度

            // 3. 实例化代币合约进行余额预检
            const tokenContract = new ethers.Contract(tokenAddress, ABIS.TOKEN, signer);

            setLoadingText('检查链上余额...');
            const realBalance = await tokenContract.balanceOf(userAddress);

            // 如果余额不足，直接中断并弹框提醒（解决重启后没打钱的问题）
            if (realBalance < amountWei) {
                modal.warning({
                    title: '账户余额不足',
                    content: `链上真实余额为 ${ethers.formatEther(realBalance)} ZNT。重启 Anvil 后请重新运行发币脚本。`,
                    okText: '知道了',
                });
                setLoading(false);
                return;
            }

            // --- 阶段 1: 授权 (Approve) ---
            setLoadingText('正在请求授权...');
            message.loading({ content: '请在钱包中确认代币授权金额...', key: msgKey });

            const approveTx = await tokenContract.approve(vaultAddress, amountWei);
            message.loading({ content: '授权交易已提交，等待链上确认...', key: msgKey });
            await approveTx.wait();

            // --- 阶段 2: 充值 (Deposit) ---
            setLoadingText('正在确认充值...');
            message.loading({ content: '授权成功！请在钱包确认充值交易...', key: msgKey });

            const vaultContract = new ethers.Contract(vaultAddress, ABIS.VAULT, signer);
            const depositTx = await vaultContract.deposit(tokenAddress, amountWei);

            message.loading({ content: '充值交易已提交，等待区块打包...', key: msgKey });
            const receipt = await depositTx.wait();

            // 4. 成功逻辑
            if (receipt && receipt.status === 1) {
                message.success({ content: '充值成功！资产已存入金库', key: msgKey, duration: 3 });
                // 仅在完全成功时调用，触发父组件 setIsOpen(false)
                onSuccess?.();
            } else {
                throw new Error("区块链执行失败 (Transaction Reverted)");
            }

        } catch (error: any) {
            console.error("充值异常详情:", error);

            let errorContent = "发生未知错误，请检查控制台";

            // 精细化错误分类
            if (error.code === 4001 || error.message?.includes("user rejected")) {
                errorContent = "您已取消了交易签名";
                message.warning({ content: errorContent, key: msgKey });
            } else if (error.message?.includes("Failed to fetch") || error.code === -32603) {
                errorContent = "无法连接到本地节点。请检查 Anvil 是否运行，并尝试在 MetaMask 设置中【重置账户】。";
                modal.error({ title: '连接异常', content: errorContent });
            } else if (error.message?.includes("execution reverted")) {
                errorContent = "合约逻辑报错（Reverted）。通常是由于余额不足、权限错误或地址不正确。";
                modal.error({ title: '充值失败', content: errorContent });
            } else {
                errorContent = error.reason || error.message || errorContent;
                modal.error({ title: '充值失败', content: errorContent });
            }

            message.destroy(msgKey);
        } finally {
            setLoading(false);
            setLoadingText('');
        }
    };

    return { handleDeposit, loading, loadingText };
};
