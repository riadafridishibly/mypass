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

```json
{
  "namespace-one": {
    "meta": {
      "created_at": "time",
      "updated_at": "time"
    },
    "items": [
        {
            "title": "Hello World",
            "id": "P-123",
            "type": "ssh (will be used as key)",
            "ssh": ""
        },
    ]
  },
  "namespace-two": {}
}
```
