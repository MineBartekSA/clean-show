# Showcase repository

This project aims to showcase my ability to design and write a simple REST backend server utilizing dependency injection and clean architecture

Please note, this is not a production-ready server.
This project is a very simple and basic implementation of an online store backend

The following things could or **should** be implemented for this project to become production ready:

- Sessions are only valid for 30 minutes and there is no mechanism to ensure the security of the tokens through browser fingerprinting.
- This entire project was not written in mind with any localization.
  A production environment for a global store page should include such features.
- The current models lack flexibility in the following ways:
  - The Product description is a single string value.
    Products should also include key-value properties to present a product's attributes better
  - Images for a Product shouldn't be in a single string value.
    Image metadata should be stored in another table that points to a proper cdn
  - Orders shouldn't contain addresses in a single string value.
    Addresses should be stored in another table and be assigned to specific accounts
  - Orders should allow users to attach additional notes for the order items.
  - Ordered products shouldn't be stored in a single string value.
    Products ordered should be stored in another table with additional data
  - Order shipping information should be expanded.
  - The audit log should contain more information about events
- Account E-Mails should be verified through a proper verification e-mail
- Accounts should implement two-factor authentication (2FA)
- Account password recovery should be implemented

### Implementations

- [Golang using Gin and Fx](../golang/gin)
