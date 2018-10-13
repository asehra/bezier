Curve Coding Test
=================

The Problem
-----------

In our simplified model of a prepaid card the card holds a balance in GBP and users can make
transactions in GBP. 

### Epic 1: User can have a card with some balance

To be able to use the card the user must first load some money onto the
card. This increases the available balance on the card. 

#### Stories:
1. A card can be created on the system
1. A card can be topped up

### Epic 2: User can make purchases on card at a Merchant
When a user goes and makes a
transaction using the card, the merchant (e.g your local independent coffee shop) sends the
card an Authorization request to check if the user loaded enough money on the account to
pay for their coffee.
If the user has loaded enough, the Authorization request is approved and
the amount the merchant requested is earmarked (or blocked). 

#### Stories:
1. Merchant can be registered on system
1. Merchants can generate authorizations that block funds on card
1. Authorizations can be declined on insufficient funds

The merchant at a later point in the future can decide to Capture ​
the transaction at which point we send the merchant the money. 

### Epic 4: User can check account status
At any point the user should be able to see how much he loaded onto the card, the
current available balance, and the current amount that is blocked on there (waiting to be
captured).

#### Stories
1. Contributions to a card can be retreived
1. Available balance on a card can be retreived
1. Blocked balance on a card can be retreived

### Epic 5: Merchant can capture (confirm) or reverse (decline) inbound transaction
The merchant can decide to only capture part of the amount or capture the amount multiple
times. 

In this model, the merchant can’t capture more than we initially authorized him to. 

The merchant can decide to reverse the whole or part of the initial Authorization at which point they
can no longer capture the full amount (only the amount that is still authorized). 

#### Stories:
1. Merchent can retrieve authorized transactions
1. Merchant can fully capture transaction
1. Merchant can part capture transaction
1. Merchant can fully reverse transaction
1. Merchant can part reverse transaction whereby, capturable amount is also reduced.

### Epic 6: Merchant can refund captured funds.
The merchant can Refund ​the user after they capture the funds. They can’t refund the user more than they
captured. 

#### Stories
1. Merchant can refund captured funds. 

The user can then use the refunded amount to buy more coffee.


*** FOR EXTRA POINTS ***
Please deploy and allow us to interact with the card as a web service. We would like to be able
to create new cards and load money, make transactions with it as described above, and see the
available and blocked balance.

*** MORE EXTRA POINTS ***
Create a transaction statement. We would like to know where we’re spending our money. This is
your chance to show your creative side.

