pragma solidity ^0.4.24;

contract FundManager {
	uint devFund = 0;
	uint conFund = 0;
	address huntContrAddress;
	address projManager;
	Admin[] admins;
	uint numadmins = 0;
	mapping(address => bool) hashmins;

	struct Admin {
		address addr;
		string name;
	}

	event donationReceived(address donor, uint amount, uint8 perc);
	event costsPayed(address amdin, uint amount, string category, string purpose);
	event bonusDisbursed(address admin, uint amount, string purpose);
	event newlyRecruited(address admin, string name);
	event honorablyDischarged(address admin, string name);

	// Accept donations with a cerntain % allocated to community
	// development ('perc'), and the rest to conservation efforts 
	function donate(uint8 perc) public payable {
		require(msg.value % 100 == 0, "Donations needs to be multiple of 100!");
		require(perc >= 0 && perc <= 100, "Invalid percentage!");
		devFund += (msg.value*perc)/100;
		conFund += (msg.value*(100-perc))/100;
		emit donationReceived(msg.sender, msg.value, perc);
	}

	// exchange coins received as payment for hunter services for ether
	// only the contract representing the hunter payment can invoke this
	function redeemEther(address recipient, uint amount) public {
		require(msg.sender == huntContrAddress, "Access restricted!");
		require(amount <= conFund, "Out of funds!");
		conFund -= amount;
		recipient.transfer(amount);
	}

	// allows wwf employees to withdraw funds to cover operation costs
	// of conservation efforts, has to provide a category to track fund allocation:
	// e.g. one of the impact goals, salary, etc..., also provide detailed purpose
	function payOpCosts(address recipient, uint amount, string category, string purpose) public {
		require(hashmins[msg.sender], "Access restricted!");
		require(amount <= conFund, "Out of funds!");
		conFund -= amount;
		emit costsPayed(msg.sender, amount, category, purpose);
		recipient.transfer(amount);
	}

	// allows wwf employees to disburs bonus community development fund
	// ****add interaction with voting, community, etc...****
	function payCommProj(address recipient, uint amount, string purpose) public {
		require(hashmins[msg.sender], "Access restricted!");
		require(amount <= conFund, "Out of funds!");
		emit bonusDisbursed(msg.sender, amount, purpose);
		devFund -= amount;
		recipient.transfer(amount);
	}

	// modify wwf employee list, only allowed by project manager
	function addEmployee(address account, string name) public returns (uint){
		require(msg.sender == projManager);
		emit newlyRecruited(account, name);
		hashmins[account] = true;
		uint id = admins.push(Admin(account, name));
		numadmins++;
		return id;
	}

	// modify wwf employee list, only allowed by project manager
	function delEmployee(uint id) public {
		require(msg.sender == projManager);
		require(id >= 0 && id < numadmins, "Invalid admin id!");
		emit honorablyDischarged(admins[id].addr, admins[id].name);
		hashmins[admins[id].addr] = false;
		admins[id] = admins[numadmins-1];
		delete admins[numadmins-1];
		numadmins--;
	}

	// change project manager, careful (only reversibel by new address)
	function switchManager(address account) public {
		require(msg.sender == projManager);
		projManager = account;
	}

	// provide viewing access to funds
	function getConFund() public view returns(uint){
		return conFund;
	}

	// provide viewing access to funds
	function getDevFund() public view returns(uint){
		return devFund;
	}

	function getNumEmp() public view returns(uint){
		return numadmins;
	}

    // only use this to iterate through array, employee ids are not static
	function getEmployee(uint id) public view returns (address, string) {
		return (admins[id].addr, admins[id].name);
	}
	
	function setHuntContrAddr(address account) public {
	    require(huntContrAddress == address(0), "Can only be set once!");
	    huntContrAddress = account;
	}
	
	function setProjManager(address account) public {
	    require(projManager == address(0), "Can only be set once!");
	    projManager = account;
	}
}