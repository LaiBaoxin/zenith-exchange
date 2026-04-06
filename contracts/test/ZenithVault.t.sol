// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {ZenithVault} from "../src/ZenithVault.sol";
import {MockToken} from "./MockToken.sol"; // 已经从外部导入了
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {MessageHashUtils} from "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";

contract ZenithVaultTest is Test {
    using MessageHashUtils for bytes32;

    ZenithVault public vault;
    MockToken public token;

    uint256 internal signerPrivKey = 0xABC123;
    address internal backendSigner;
    address internal user = address(0x1337);

    function setUp() public {
        backendSigner = vm.addr(signerPrivKey);
        vault = new ZenithVault(backendSigner);
        
        // 实例化外部导入的 MockToken
        token = new MockToken("Zenith Token", "ZNT");
        
        // 给用户转账
        token.transfer(user, 1000 ether);
    }

    function test_FullFlow() public {
        // 1. 充值
        vm.startPrank(user);
        token.approve(address(vault), 100 ether);
        vault.deposit(address(token), 100 ether);
        assertEq(vault.balances(user, address(token)), 100 ether);

        // 2. 准备提现数据
        uint256 withdrawAmount = 50 ether;
        uint256 nonce = 0;

        // 3. 模拟后端签名逻辑
        bytes32 msgHash = keccak256(
            abi.encodePacked(user, address(token), withdrawAmount, nonce, address(vault), block.chainid)
        );
        bytes32 ethSignedHash = msgHash.toEthSignedMessageHash();
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPrivKey, ethSignedHash);
        bytes memory signature = abi.encodePacked(r, s, v);

        // 4. 执行提现
        vault.withdraw(address(token), withdrawAmount, nonce, signature);
        vm.stopPrank();

        // 5. 校验结果
        assertEq(vault.balances(user, address(token)), 50 ether);
        assertEq(token.balanceOf(user), 950 ether);
    }

    function test_RevertOn_InvalidSignature() public {
        vm.startPrank(user);
        token.approve(address(vault), 100 ether);
        vault.deposit(address(token), 100 ether);

        // 使用错误的私钥签名来模拟黑客攻击
        uint256 wrongKey = 0xBAD;
        bytes32 msgHash = keccak256(
        abi.encodePacked(user, address(token), uint256(50 ether), uint256(0), address(vault), block.chainid));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(wrongKey, msgHash.toEthSignedMessageHash());
        bytes memory signature = abi.encodePacked(r, s, v);

        vm.expectRevert("Invalid signature");
        vault.withdraw(address(token), 50 ether, 0, signature);
        vm.stopPrank();
    }
}