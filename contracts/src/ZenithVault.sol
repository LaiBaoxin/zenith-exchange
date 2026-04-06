// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {MessageHashUtils} from "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title ZenithVault
 * @dev 极点交易所资产金库 - 采用后端签名授权提现机制
 */
contract ZenithVault is Ownable {
    using ECDSA for bytes32;
    using MessageHashUtils for bytes32;

    address public backendSigner;
    mapping(address => mapping(address => uint256)) public balances;
    mapping(address => uint256) public nonces;

    event Deposit(address indexed user, address indexed token, uint256 amount);
    event Withdraw(address indexed user, address indexed token, uint256 amount, uint256 nonce);

    // 注意：OpenZeppelin v5 的 Ownable 构造函数需要传入初始 Owner 地址
    constructor(address _initialSigner) Ownable(msg.sender) {
        backendSigner = _initialSigner;
    }

    function setSigner(address _newSigner) external onlyOwner {
        backendSigner = _newSigner;
    }

    function deposit(address token, uint256 amount) external {
        require(amount > 0, "Amount must be greater than zero");
        IERC20(token).transferFrom(msg.sender, address(this), amount);
        balances[msg.sender][token] += amount;
        emit Deposit(msg.sender, token, amount);
    }

    function withdraw(
        address token,
        uint256 amount,
        uint256 nonce,
        bytes calldata signature
    ) external {
        require(balances[msg.sender][token] >= amount, "Insufficient balance");
        require(nonce == nonces[msg.sender], "Invalid nonce");

        // 构造结构化数据哈希
        bytes32 messageHash = keccak256(
            abi.encodePacked(msg.sender, token, amount, nonce, address(this), block.chainid)
        );

        // 使用 MessageHashUtils 转化为以太坊签名消息哈希，再用 ECDSA 恢复地址
        address recoveredSigner = messageHash.toEthSignedMessageHash().recover(signature);
        require(recoveredSigner == backendSigner, "Invalid signature");

        nonces[msg.sender]++;
        balances[msg.sender][token] -= amount;
        IERC20(token).transfer(msg.sender, amount);

        emit Withdraw(msg.sender, token, amount, nonce);
    }
}