<!--
order: 1
-->

# Concepts

This module is responsible for the minting and burning of liquid staking receipt tokens, collectively referred to as `bnemo`. Delegated nemo can be converted to delegator-specific `bnemo`. Ie, 100 NEMO delegated to validator `nemovaloper123` can be converted to 100 `bnemo-nemovaloper123`. Similarly, 100 `bnemo-nemovaloper123` can be converted back to a delegation of 100 NEMO to  `nemovaloper123`. In this design, all validators can permissionlessly participate in liquid staking while users retain the delegator specific slashing risk and voting rights of their original validator. Note that because each `bnemo` denom is validator specific, this module does not specify a fungibility mechanism for `bnemo` denoms. 