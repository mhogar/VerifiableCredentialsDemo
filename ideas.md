# Final Project Components
- https://www.xtseminars.co.uk/post/introduction-to-the-future-of-identity-dids-vcs

### Issuer
- MS example app good starting place
- front-end issue request generator
    - generate QR-code that contains issuer DID and required fields (VC template)
- back-end server
    - receive request from user and verify signature
    - take body and add own DID and sign with private key
    - return to user

### User
- front-end mobile app
- create ID
    - scan issuer QR-code (or click generate on issuer site on mobile device)
    - prompt user if they trust issuer (gather fields from DID on blockchain)
    - enter fields requested by issuer
    - add public key and sign with private key
    - send to issuer
    - receive and store returned signed ID
- verify ID
    - scan verifier QR-code
    - prompt user if they trust verifier (gather fields from DID on blockchain and presentation details)
    - send ID to verifier
    - receive result

### Verifier
- MS example app good starting place
- front-end presentation request generator
    - generate QR-code that contains verifier DID and other presentation details
- back-end server
    - receive request from user and verify signature
    - verify issuer signature and DID
    - do something with result (change front-end)
    - send result to user

### Use Cases
- design verifier for specific use case
- university exam
    - pass around sheet with verifier QR code instead of sign in sheet
        - improve privacy by not being able to see details of others in the exam

### Other Considerations
- what blockchain to use to store DIDs for issuer and verifier
- how to generate DID JSONs. Do we need all those fields for this prototype?
