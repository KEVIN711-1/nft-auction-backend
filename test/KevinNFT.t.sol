// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/KevinNFT.sol";

contract KevinNFTTest is Test {
    KevinNFT nft;

    address owner = address(this);
    address user1 = address(0x1);
    address user2 = address(0x2);

    receive() external payable {}

    function setUp() public {
        nft = new KevinNFT(
            "KevinNFT",
            "KNFT",
            0.01 ether,
            10
        );
    }

    function testMintNFT() public {
        vm.deal(user1, 1 ether);

        vm.prank(user1);
        uint256 tokenId = nft.mintNFT{value: 0.01 ether}("ipfs://token1");

        assertEq(tokenId, 1);
        assertEq(nft.ownerOf(1), user1);
        assertEq(nft.totalSupply(), 1);
    }

    function testMintFailInsufficientETH() public {
        vm.deal(user1, 1 ether);

        vm.prank(user1);
        vm.expectRevert("Insufficient payment");
        nft.mintNFT{value: 0.001 ether}("ipfs://token1");
    }

    function testOwnerMint() public {
        uint256 tokenId = nft.ownerMint("ipfs://ownerMint", user2);

        assertEq(tokenId, 1);
        assertEq(nft.ownerOf(1), user2);
    }

    function testWithdraw() public {
        vm.deal(user1, 1 ether);

        vm.prank(user1);
        nft.mintNFT{value: 0.01 ether}("ipfs://token1");

        uint256 balanceBefore = owner.balance;
        nft.withdraw();
        uint256 balanceAfter = owner.balance;

        assertEq(balanceAfter - balanceBefore, 0.01 ether);
    }
}
