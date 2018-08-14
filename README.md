# pivex

### Prerequisites
- **Go**	
	Make sure Go is installed. If you don't have Go environment created already:
	```
	mkdir -p ~/go/{src,bin,pkg}
	export GOPATH="~/go"
	```
- **Pivex credentials directory**  
	Credentials used in `pivex` are stored under `~/.pivex`, create this
	directory:
	```
	mkdir ~/.pivex
	```
  - **Pivotal API key**  
	Can be found under a user's [profile settings](https://www.pivotaltracker.com/profile).
	Copy the API key into a file and put it in the `pivex` credentials
	directory:
	```
	echo 'my-api-key' > ~/.pivex/pivotal-token
	```
  - **Google API credentials**  
	Select `pivex` from the list of
	[Google APIs](https://console.developers.google.com/apis/credentials?project=lco-internal&folder&organizationId=78492096084)
	and click `DOWNLOAD JSON`. Then rename the JSON credentials to
	`api-creds.json` and move them into the `pivex` credentials directory:
	```
	mv pivex-creds.json ~/.pivex
	```

## Installation
In the project directory:
```
go get
go install
```

## Usage
```
pivex
```
The first time you run this program, you will be prompted with an OAuth consent
screen.
