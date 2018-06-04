# pivex

### Prerequisites
- Go	
	If you don't have Go environment created already:
	```
	mkdir -p ~/go/{src,bin,pkg}
	export GOPATH="~/go"
	```
- Pivotal API key
- Google API credentials	
	Download the `pivex` JSON credentials from [here](https://console.developers.google.com/apis/credentials?project=lco-internal&folder&organizationId=78492096084)
	and rename it to `api-creds.json` then:
	```
	mkdir ~/.pivex
	mv pivex-creds.json ~/.pivex
	```

## Installation
```
go get
go install
```

## Usage
```
pivex
```

## Examples
