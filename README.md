# sm
Simple secret management tool for server configuration

[ ![Codeship Status for agilebits/sm](https://app.codeship.com/projects/33899e80-fae5-0134-b168-721cf569a862/status?branch=master)](https://app.codeship.com/projects/211385)

## How to build

```
go get -u -v github.com/agilebits/sm
cd ~/go/src/github.com/agilebits/sm
go install
```
## Encrypt/decrypt data on development machines

```
cat app-config.yml | sm encrypt > app-config.sm
cat app-config.sm | sm decrypt
```

On the first run, the utility will generate a new master key and store it in `~/.sm/masterkey` file. The `masterkey` must be saved and copied across all developer machines.


## Encrypt/decrypt data with Amazon Web Service KMS

First, you have to create a master key using AWS IAM and give yourself permissions to use this key for encryption and decryption.

```
export AWS_REGION='us-east-1'
export KMS_KEY_ID='arn:aws:kms:us-east-1:123123123123:key/d845cfa3-0719-4631-1d00-10ab63e40ddf'

# encrypt the file and pipe stdout to a file
cat app-config.yml | sm encrypt \
	--env aws \
	--region $AWS_REGION \
	--master $KMS_KEY_ID \
	> app-config.sm

# encrypt the file and write the output to a file
cat app-config.yml | sm encrypt \
  --env aws \
  --region $AWS_REGION \
  --master $KMS_KEY_ID \
  --out app-config.sm

# encrypt the file using settings from a configuration file
cat app-config.yml | sm encrypt \
  --config config.yml
  --out app-config.sm

# decrypt the file specified via stdin
cat app-config.sm | sm decrypt

# decrypt the file speified via flag
sm decrypt --input app-config.sm

# decrypt the file and write the output to a file
cat app-config.sm | sm decrypt --out app-config.yml
```

## Use jq to validate JSON files

For example:
```
export AWS_REGION=us-east-1
export KMS_KEY_ID=alias/YOUR-KEY-ALIAS

jq --compact-output . < config.json | sm encrypt \
        --env aws \
        --region $AWS_REGION \
        --master $KMS_KEY_ID \
        > config.sm

sm decrypt < config.sm | jq

```

## Using with GIT

You can integrate SM with `git diff` by adding the following in your repository's `.gitattributes` file:

```
*.sm diff=sm
```

The above change assumes that all encrypted files will have the ending `.sm`. It can be safely committed if that is the case.

Developers will also need to add the following stanza to their `.git/config` file - either globally or on a per-repository basis:

```
[diff "sm"]
    textconv = sm decrypt --input
```
