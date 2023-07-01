import { ethers } from "hardhat";

async function main() {
  const tokenName = "Nemo-wrapped ATOM";
  const tokenSymbol = "kATOM";
  const tokenDecimals = 6;

  const ERC20NemoWrappedCosmosCoin = await ethers.getContractFactory(
    "ERC20NemoWrappedCosmosCoin"
  );
  const token = await ERC20NemoWrappedCosmosCoin.deploy(
    tokenName,
    tokenSymbol,
    tokenDecimals
  );

  await token.deployed();

  console.log(
    `Token "${tokenName}" (${tokenSymbol}) with ${tokenDecimals} decimals is deployed to ${token.address}!`
  );
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
