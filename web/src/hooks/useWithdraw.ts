import { useState } from 'react';
import { ethers } from 'ethers';
import { message } from 'antd';
import { getWithdrawSignature } from '../api/vault';
import { ABIS } from '../constants/contracts_abi';

// 接收 vaultAddress 和 tokenAddress
export const useWithdraw = (vaultAddress: string, tokenAddress: string) => {
    const [loading, setLoading] = useState(false);

    const handleWithdraw = async (amountStr: string) => {
        // 环境检查
        if (!(window as any).ethereum) return message.error("请安装 MetaMask");

       // 获取合约地址
        if (!vaultAddress || vaultAddress === "" || vaultAddress === "0x") {
            return message.error("金库合约地址尚未就绪，请稍后");
        }

        setLoading(true);
        const msgKey = 'withdraw_process';

        try {
            const provider = new ethers.BrowserProvider((window as any).ethereum);
            const signer = await provider.getSigner();

            // 转换单位
            const amountWei = ethers.parseEther(amountStr).toString();

            message.loading({ content: '正在请求后端安全签名...', key: msgKey });

            // 获取签名数据
            const { signature, nonce, amount: sigAmount } = await getWithdrawSignature(amountWei);

            message.loading({ content: '等待钱包确认交易...', key: msgKey });

            // 实例化合约
            const vaultContract = new ethers.Contract(vaultAddress, ABIS.VAULT, signer);

            /**
             * 严格按照合约 withdraw 方法参数顺序:
             * function withdraw(address token, uint256 amount, uint256 nonce, bytes memory signature)
             */
            const tx = await vaultContract.withdraw(
                tokenAddress, // 从 config 传入的代币地址
                sigAmount,    // 后端返回确认的数量
                nonce,        // 后端生成的随机值
                signature     // 后端生成的 EIP-712 或标准签名
            );

            message.loading({ content: '交易已提交，等待区块确认...', key: msgKey });
            await tx.wait();

            message.success({ content: '提现成功！资金已汇入您的钱包', key: msgKey, duration: 3 });

        } catch (error: any) {
            console.error("Withdraw Error:", error);
            // 提取合约报错或网络报错信息
            const errorMsg = error.reason || error.data?.message || error.message || "提现交易失败";
            message.error({ content: errorMsg, key: msgKey, duration: 4 });
        } finally {
            setLoading(false);
        }
    };

    return { handleWithdraw, loading };
};
