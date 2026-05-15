import { useState } from 'react';
import { ethers } from 'ethers';
import { App } from 'antd';
import { ABIS } from '../constants/contracts_abi';
import { getWithdrawSignature } from '../api/vault';

export const useWithdraw = (
    vaultAddress: string,
    tokenAddress: string,
    onSuccess?: () => void // 提现成功后的回调（比如用来刷新页面上的余额）
) => {
    const [loading, setLoading] = useState(false);
    const { message } = App.useApp();

    /**
     * handleWithdraw 执行提现流程
     * @param amountStr - 提现数量 (人类可读格式, 如 "10.5")
     * @param currency - 提现币种 (如 "USDT"), 用于后端余额扣减
     */
    const handleWithdraw = async (amountStr: string, currency: string = "USDT") => {
        if (!(window as any).ethereum) return message.error("请安装 MetaMask 钱包");

        if (!vaultAddress || vaultAddress === "" || vaultAddress === "0x") {
            return message.error("系统合约地址尚未加载，请稍后再试");
        }

        setLoading(true);
        const msgKey = 'withdraw_process';

        try {
            const provider = new ethers.BrowserProvider((window as any).ethereum);
            const signer = await provider.getSigner();

            // 1. 转换单位为 Wei (发给后端的参数)
            const amountWei = ethers.parseEther(amountStr).toString();

            message.loading({ content: '正在请求后端风控校验与签名...', key: msgKey });

            // 2. 获取后端签名数据 (后端会从链上读取正确的 nonce)
            const res: any = await getWithdrawSignature({
                amount: amountWei,
                currency: currency
            });
            // 注意: axios 拦截器已经提取了 response.data.data，所以 res 就是最终数据
            const { signature, nonce, amount: sigAmount } = res;

            message.loading({ content: '签名获取成功，等待钱包确认交易...', key: msgKey });

            // 3. 实例化合约并执行
            const vaultContract = new ethers.Contract(vaultAddress, ABIS.VAULT, signer);

            const tx = await vaultContract.withdraw(
                tokenAddress,
                sigAmount,
                nonce,
                signature
            );

            message.loading({ content: '交易已提交，等待区块链确认...', key: msgKey });

            // 4. 等待回执
            const receipt = await tx.wait();

            if (receipt.status === 1) {
                message.success({
                    content: '提现成功！资金已汇入您的钱包',
                    key: msgKey,
                    duration: 3
                });

                // 回调通知父组件（例如重新拉取钱包余额和系统余额）
                onSuccess?.();
            } else {
                throw new Error("链上交易执行失败");
            }

        } catch (error: any) {
            console.error("提现异常:", error);

            // 优化错误提示
            let errorMsg = "提现操作失败";
            if (error.code === 4001 || error.message?.includes("user rejected")) {
                errorMsg = "用户取消了交易签名";
            } else if (error.reason) {
                errorMsg = `合约报错: ${error.reason}`;
            } else if (error.response?.data?.message) {
                // 捕获后端返回的错误（例如：余额不足）
                errorMsg = error.response.data.message;
            } else if (error.message) {
                errorMsg = error.message;
            }

            message.error({ content: errorMsg, key: msgKey, duration: 4 });
        } finally {
            setLoading(false);
        }
    };

    return { handleWithdraw, loading };
};
