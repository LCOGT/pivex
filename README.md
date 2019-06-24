# pivex
Utility that exports Pivotal stories to a Google Slides deck in the software
team shared Google Drive.

## Requirements
- **Pivotal API key**: Can be found under a user's
[profile settings](https://www.pivotaltracker.com/profile).
- **Google Slides API OAuth 2.0 client ID**: Download the
[Client ID](http://developers.lco.gtn/pivex-gslides-oauth2.0-client-id.json).

## Build
Build a binary to use:
```
make
```

## Installation
Either build the binary manually or download the pre-built
[asset](https://github.com/LCOGT/pivex/releases).

## Usage
When the program is run for the first time, you'll need to specify the file
names for both the Pivotal API token and the Google Drive OAuth 2.0 Client ID:
```
pivex --google-client-id-file client-id.json --pivotal-api-token-file pivotal-api
```

You will be prompted with an OAuth page in your browser to allow access to your
Google Drive account. Once you have allowed access, paste the generated
authorization code into the application prompt.

These credentials are now stored on your machine and will not have to be
retrieved again.

Once the application has finished running, the slides will be generated and
available on the software team drive.

To generate the slides again in the future, just run the application:
```
pivex
```

## License
[GPLv3](LICENSE)
