## Time synchronization
Building: The time is based on the local machine
Running: The time is based on the ntp time.google.com

## Encryption algorithm
### Encrypt
1. Generate a random key ENCRYPTION_KEY
2. ENCRYPTION_KEY_X = b64Encode(ENCRYPTION_KEY XOR "RDPCloud"))
3. ENCRYPTED_{SECRET_VARIABLE} = b64Encode({SECRET_VARIABLE} XOR ENCRYPTION_KEY_X)
4. SIGNATURE = b64Encode((concatenate all ENCRYPTED_{SECRET_VARIABLE} + "SIGNATURE") XOR ENCRYPTION_KEY_X)

### Decrypt
1. Use the key ENCRYPTION_KEY
2. ENCRYPTION_KEY_X = b64Encode(ENCRYPTION_KEY XOR "RDPCloud"))
3. Verify signature: SIGNATURE == b64Encode((concatenate all ENCRYPTED_{SECRET_VARIABLE} + "SIGNATURE") XOR ENCRYPTION_KEY_X)
4. {SECRET_VARIABLE} = b64Decode(ENCRYPTED_{SECRET_VARIABLE}) XOR ENCRYPTION_KEY_X

## Obfuscation
garble is used for obfuscation

## Notes
- Different clients/servers accross different languages does not have to be exactly the same. Example: client/go offers a UI (within a webserver) but client/php does not. Features are implemented based on the need and each language can have different flavors with different functionalities
