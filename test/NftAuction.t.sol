// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/NftAuction.sol";
import "../src/KevinNFT.sol";

contract NftAuctionTest is Test {
    NftAuction auction;
    KevinNFT nft;

    // 使用 EOA，而不是测试合约本身
    address admin   = address(0x1);
    address bidder1 = address(0x2);
    address bidder2 = address(0x3);

    receive() external payable {}

    function setUp() public {
        // 给账户打钱
        vm.deal(admin, 10 ether);
        vm.deal(bidder1, 10 ether);
        vm.deal(bidder2, 10 ether);

        // 用 admin 部署拍卖合约
        vm.prank(admin);
        auction = new NftAuction();

        // 用 admin 部署 NFT 合约
        vm.prank(admin);
        nft = new KevinNFT(
            "KevinNFT",
            "KNFT",
            0.01 ether,
            100
        );

        // admin mint NFT（tokenId = 1）
        vm.prank(admin);
        nft.mintNFT{value: 0.01 ether}("ipfs://auctionNFT");

        // admin 授权拍卖合约
        vm.prank(admin);
        nft.approve(address(auction), 1);
    }

    function testCreateAuctionETH() public {
        vm.prank(admin);
        auction.createAuctionETH(
            300,
            0.1 ether,
            address(nft),
            1
        );

        (
            address seller,
            uint256 duration,
            uint256 startPrice,
            ,
            bool ended,
            ,
            ,
            ,
            uint256 tokenId,
            bool useERC20,

        ) = auction.auctions(0);

        assertEq(seller, admin);
        assertEq(duration, 300);
        assertEq(startPrice, 0.1 ether);
        assertEq(ended, false);
        assertEq(tokenId, 1);
        assertEq(useERC20, false);

        // NFT 已托管到拍卖合约
        assertEq(nft.ownerOf(1), address(auction));
    }

    function testPlaceBidETH() public {
        vm.prank(admin);
        auction.createAuctionETH(
            300,
            0.1 ether,
            address(nft),
            1
        );

        vm.prank(bidder1);
        auction.placeBidETH{value: 0.2 ether}(0);

        (, , , , , address highestBidder, uint256 highestBid, , , , ) =
            auction.auctions(0);

        assertEq(highestBidder, bidder1);
        assertEq(highestBid, 0.2 ether);
    }

    function testOutbidAndRefund() public {
        vm.prank(admin);
        auction.createAuctionETH(
            300,
            0.1 ether,
            address(nft),
            1
        );

        vm.prank(bidder1);
        auction.placeBidETH{value: 0.2 ether}(0);

        uint256 beforeRefund = bidder1.balance;

        vm.prank(bidder2);
        auction.placeBidETH{value: 0.3 ether}(0);

        uint256 afterRefund = bidder1.balance;

        assertEq(afterRefund - beforeRefund, 0.2 ether);
    }

    function testEndAuction() public {
        vm.prank(admin);
        auction.createAuctionETH(
            300,
            0.1 ether,
            address(nft),
            1
        );

        vm.prank(bidder1);
        auction.placeBidETH{value: 0.2 ether}(0);

        uint256 adminBalanceBefore = admin.balance;

        vm.warp(block.timestamp + 400);

        vm.prank(admin);
        auction.endAuction(0);

        assertEq(nft.ownerOf(1), bidder1);
        assertEq(admin.balance - adminBalanceBefore, 0.2 ether);
    }
}
