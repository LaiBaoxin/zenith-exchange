// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/MockToken.sol";
import "../src/ZenithVault.sol";

contract DeployAll is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployerAddr = vm.addr(deployerPrivateKey);

        vm.deal(deployerAddr, 1000000 ether);

        // 开始广播交易
        vm.startBroadcast(deployerPrivateKey);

        // 部署 MockToken 
        MockToken token = new MockToken("Zenith Test Token", "ZNT");
        
        // 部署 ZenithVault
        ZenithVault vault = new ZenithVault(deployerAddr);

        vm.stopBroadcast();

        console.log("---------------------------");
        console.log("Deployer Address:", deployerAddr);
        console.log("ETH Balance:", deployerAddr.balance / 1e18, "ETH");
        console.log("ZNT Balance:", token.balanceOf(deployerAddr) / 1e18, "ZNT");
        console.log("---------------------------");
        console.log("Copy these to your Frontend:");
        console.log("VAULT_ADDR =", address(vault));
        console.log("TOKEN_ADDR =", address(token));
        console.log("---------------------------");
    }
}