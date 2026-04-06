// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/MockToken.sol";
import "../src/ZenithVault.sol";

contract DeployAll is Script {
    function run() external {
        // 从环境变量读取私钥
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployerAddr = vm.addr(deployerPrivateKey);

        // 开始广播交易
        vm.startBroadcast(deployerPrivateKey);

        // 部署 MockToken
        MockToken token = new MockToken("Zenith Test Token", "ZNT");
        console.log("MockToken deployed at:", address(token));

        ZenithVault vault = new ZenithVault(deployerAddr);
        console.log("ZenithVault deployed at:", address(vault));

        // 初始给部署者 Mint 一些代币方便测试
        token.mint(deployerAddr, 10000 ether);

        vm.stopBroadcast();

        console.log("---------------------------");
        console.log("Copy these to your Frontend:");
        console.log("VAULT_ADDR =", address(vault));
        console.log("TOKEN_ADDR =", address(token));
        console.log("---------------------------");
    }
}