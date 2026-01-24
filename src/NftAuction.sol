// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";  // 添加这一行

contract NftAuction is IERC721Receiver {
    struct Auction {
        address seller;
        uint256 duration;      // 秒
        uint256 startPrice;    // wei单位
        uint256 startTime;
        bool ended;
        address highestBidder;
        uint256 highestBid;    // wei单位
        address nftContract;
        uint256 tokenId;
        bool useERC20;         // 是否使用ERC20
        address erc20Token;    // ERC20代币地址（如果使用）
    }
    
    mapping(uint256 => Auction) public auctions;
    uint256 public nextAuctionId;
    address public admin;
    
    // ERC20白名单（安全考虑）
    mapping(address => bool) public allowedERC20Tokens;
    
    event AuctionCreated(uint256 auctionId, address seller, uint256 tokenId, uint256 startPrice);
    event NewBid(uint256 auctionId, address bidder, uint256 amount);
    event AuctionEnded(uint256 auctionId, address winner, uint256 finalPrice);
    
    constructor() {
        admin = msg.sender;
    }
    
    // 允许特定的ERC20代币
    function allowERC20Token(address token) external {
        require(msg.sender == admin, "Only admin");
        allowedERC20Tokens[token] = true;
    }
    
    // 创建拍卖（ETH拍卖）
    function createAuctionETH(
        uint256 _duration,
        uint256 _startPrice,   // wei单位，如 1 ETH = 1000000000000000000
        address _nftAddress,
        uint256 _tokenId
    ) external onlyAdmin {
        _createAuction(_duration, _startPrice, _nftAddress, _tokenId, false, address(0));
    }
    
    // 创建拍卖（ERC20拍卖）
    function createAuctionERC20(
        uint256 _duration,
        uint256 _startPrice,   // ERC20最小单位
        address _nftAddress,
        uint256 _tokenId,
        address _erc20Token
    ) external onlyAdmin {
        require(allowedERC20Tokens[_erc20Token], "Token not allowed");
        _createAuction(_duration, _startPrice, _nftAddress, _tokenId, true, _erc20Token);
    }
    
    function _createAuction(
        uint256 _duration,
        uint256 _startPrice,
        address _nftAddress,
        uint256 _tokenId,
        bool _useERC20,
        address _erc20Token
    ) internal {
        require(_duration >= 60, "Min 60 seconds");
        require(_startPrice > 0, "Start price > 0");
        
        IERC721(_nftAddress).safeTransferFrom(msg.sender, address(this), _tokenId);
        
        auctions[nextAuctionId] = Auction({
            seller: msg.sender,
            duration: _duration,
            startPrice: _startPrice,
            startTime: block.timestamp,
            ended: false,
            highestBidder: address(0),
            highestBid: 0,
            nftContract: _nftAddress,
            tokenId: _tokenId,
            useERC20: _useERC20,
            erc20Token: _erc20Token
        });
        
        emit AuctionCreated(nextAuctionId, msg.sender, _tokenId, _startPrice);
        nextAuctionId++;
    }
    
    // 出价（ETH拍卖）
    function placeBidETH(uint256 _auctionId) external payable auctionActive(_auctionId) {
        Auction storage auction = auctions[_auctionId];
        require(!auction.useERC20, "This auction accepts ERC20, not ETH");
        require(msg.value > auction.highestBid, "Bid too low");
        require(msg.value >= auction.startPrice, "Bid below start price");
        
        _placeBid(_auctionId, msg.value);
    }
    
    // 出价（ERC20拍卖）
    function placeBidERC20(uint256 _auctionId, uint256 _amount) external auctionActive(_auctionId) {
        Auction storage auction = auctions[_auctionId];
        require(auction.useERC20, "This auction accepts ETH, not ERC20");
        require(_amount > auction.highestBid, "Bid too low");
        require(_amount >= auction.startPrice, "Bid below start price");
        
        // 转移ERC20到合约
        IERC20(auction.erc20Token).transferFrom(msg.sender, address(this), _amount);
        
        _placeBid(_auctionId, _amount);
    }
    
    function _placeBid(uint256 _auctionId, uint256 _amount) internal {
        Auction storage auction = auctions[_auctionId];
        
        // 退还前一个出价
        if (auction.highestBid > 0) {
            _refundPreviousBidder(auction);
        }
        
        // 更新拍卖状态
        auction.highestBid = _amount;
        auction.highestBidder = msg.sender;
        
        emit NewBid(_auctionId, msg.sender, _amount);
    }
    
    function _refundPreviousBidder(Auction storage auction) internal {
        if (auction.useERC20) {
            IERC20(auction.erc20Token).transfer(auction.highestBidder, auction.highestBid);
        } else {
            // payable(auction.highestBidder).transfer(auction.highestBid);
            (bool success, ) = payable(auction.highestBidder).call{value: auction.highestBid}("");
            require(success, "ETH transfer failed");
        }
    }
    
    // 结束拍卖
    function endAuction(uint256 _auctionId) external auctionExists(_auctionId) {
        Auction storage auction = auctions[_auctionId];
        require(!auction.ended, "Already ended");
        require(block.timestamp >= auction.startTime + auction.duration, "Not yet ended");
        
        auction.ended = true;
        
        if (auction.highestBidder != address(0)) {
            // 转移NFT给赢家
            IERC721(auction.nftContract).safeTransferFrom(
                address(this),
                auction.highestBidder,
                auction.tokenId
            );
            
            // 转移资金给卖家
            if (auction.useERC20) {
                IERC20(auction.erc20Token).transfer(auction.seller, auction.highestBid);
            } else {
                // payable(auction.seller).transfer(auction.highestBid);
                (bool success, ) = payable(auction.seller).call{value: auction.highestBid}("");
                require(success, "ETH transfer failed");
            }
            
            emit AuctionEnded(_auctionId, auction.highestBidder, auction.highestBid);
        } else {
            // 无人出价，退回NFT
            IERC721(auction.nftContract).safeTransferFrom(
                address(this),
                auction.seller,
                auction.tokenId
            );
        }
    }
    
    // 辅助函数：将ETH转换为wei（前端使用）
    function ethToWei(uint256 ethAmount) public pure returns (uint256) {
        return ethAmount * 10**18;
    }
    
    // 辅助函数：将wei转换为ETH（前端使用）
    function weiToEth(uint256 weiAmount) public pure returns (uint256) {
        return weiAmount / 10**18;
    }
    
    // 修改器和辅助函数
    modifier onlyAdmin() { require(msg.sender == admin, "Only admin"); _; }
    modifier auctionExists(uint256 id) { require(id < nextAuctionId, "Not exist"); _; }
    modifier auctionActive(uint256 id) { 
        require(id < nextAuctionId, "Not exist");
        Auction storage a = auctions[id];
        require(!a.ended, "Ended");
        require(block.timestamp < a.startTime + a.duration, "Expired");
        _;
    }
    
    function onERC721Received(address, address, uint256, bytes calldata) 
        external pure returns (bytes4) 
    {
        return this.onERC721Received.selector;
    }
    
    // 提取意外发送的ETH（仅管理员）
    function withdrawETH(uint256 amount) external onlyAdmin {
        // payable(admin).transfer(amount);

        (bool success, ) = payable(admin).call{value: amount}("");
        require(success, "ETH transfer failed");
    }
    
    // 提取意外发送的ERC20（仅管理员）
    function withdrawERC20(address token, uint256 amount) external onlyAdmin {
        IERC20(token).transfer(admin, amount);
    }
}