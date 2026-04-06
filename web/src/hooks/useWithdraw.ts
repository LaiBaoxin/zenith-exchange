import { useState } from 'react';
import { ethers } from 'ethers';
import { message } from 'antd'; // 使用 antd 的全局提示
import { getWithdrawSignature } from '../api/vault';
import { ABIS as VAULT_ABI} from '../constants/contracts_abi.ts';

export const useWithdraw = (vaultAddress: string, tokenAddress: string) => {
    const [loading, setLoading] = useState(false);

    const handleWithdraw = async (amountStr: string, nonce: number) => {
        if (!(window as any).ethereum) return message.error("未检测到钱包环境");

        setLoading(true);
        try {
            const provider = new ethers.BrowserProvider((window as any).ethereum);
            const signer = await provider.getSigner();
            const network = await provider.getNetwork();
            const amountWei = ethers.parseEther(amountStr).toString();

            // 获取签名
            const { signature } = await getWithdrawSignature({
                user: await signer.getAddress(),
                token: tokenAddress,
                amount: amountWei,
                nonce,
                vault_addr: vaultAddress,
                chain_id: Number(network.chainId)
            });

            // 合约交互
            const vaultContract = new ethers.Contract(vaultAddress, VAULT_ABI.VAULT, signer);
            const tx = await vaultContract.withdraw(tokenAddress, amountWei, nonce, signature);

            message.loading({ content: '交易已提交，等待上链...', key: 'tx_wait' });
            await tx.wait();
            message.success({ content: '提现成功！', key: 'tx_wait', duration: 3 });

        } catch (error: any) {
            console.error(error);
            message.error(error.response?.data?.error || error.reason || "交易执行失败");
        } finally {
            setLoading(false);
        }
    };

    return { handleWithdraw, loading };
};
