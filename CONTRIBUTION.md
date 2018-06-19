
# Contributing Guide

## Prerequisite

- `docker` service is running
- `docker-compose` is installed
- `go` is installed
- `GOPATH` is set
- `PATH` variable includes `GOBIN`
- `make` tool is installed (most proabably already installed)

## Build

Go to the root directory of project

    make build
    
It will result is `github-flair-web:latest` docker image. Check with `docker images` if it exists

## Run

    make run
    
It will build the image if not already build and run them. It may take some minutes first time to fetch `mongo` image.

Now, visit `http://localhost:443/github/<your-git-username>.png`.

<b>Note: It is `http` on port 443</b>

## Code Edit

This project uses Github API to fetch stats about profile. 
Github limits the API usage if you are accessing it without authentication.
Checkout [this](https://developer.github.com/v3/#rate-limiting) link for more info. 

Generate your personal access token from [here](https://github.com/settings/tokens) and put it in `secrets/access_token.txt`
    
    mkdir -p secrets
    echo $GITHUB_ACCESS_TOKEN > secrets/access_token.txt

Make sure to remove all the permission from token. So that no one can misuse it if leaked accidently.
Now, make your changes and try to run the modified code

    make build && make run

After your code is successfully working. Run the test, before making PR

    make test
    
## SSL Support

To provide SSL support just put these file in `secrets` directory

 - crt-bundle.pem
 - ssl-private.key (without password)
 
It will automatically use <b>https</b> if the pair is correct

## Port confusion

The application runs on port 8443 no matter if it `http` or `https` but in `docker-compose.yml` the 8443 port is mapped to 443 on host.
