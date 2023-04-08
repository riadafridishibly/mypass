# My Pass - Dead Simple Password Manager

## Features

- Generate password
- Save passwords to different namespaces
- Save ssh keys
- One master password (we may support split keys)
- Public Key Encryption

## Commands

Higher level command preview,

```text
# Generate password
$ mypass generate [--size=N --no-special --no-number --no-lower --no-upper]

# List public keys
$ mypass pubkey list --json
$ mypass pubkey remove 'publicKey' (remove from all items)
$ mypass pubkey add 'publicKey' (encrypt all items with this pubkey)

# Add new password
$ mypass add [password | ssh] --title='' --namespace='' --username='' --host='' --port='' --url='' --password='' --extra='{}'

$ mypass list --ns=''

id=SOME-ID title='This is the production server' tags=tag1,tag2 --username=''

$ mypass remove <item-id>

$ mypass update <item-id>
```

## Libraries to look into

- https://github.com/manifoldco/promptui
- https://github.com/lithammer/fuzzysearch

## TODOS

- Validate database schema
