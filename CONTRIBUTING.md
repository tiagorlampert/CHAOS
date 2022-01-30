# Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new pull request (PR)

# TODO

- [x] ~~Remote Shell~~
- [x] ~~Take Screenshot~~
- [x] ~~File Explorer~~
- [x] ~~Download File~~
- [x] ~~Upload File~~
- [x] ~~Restart~~
- [x] ~~Shutdown~~
- [ ] Persistence (install at startup)
- [ ] Lock screen
- [X] ~~Open URL~~
- [ ] Kill process

# Local Development
```bash
# Install dependencies
$ sudo apt install golang git -y

# Get this repository
$ git clone https://github.com/tiagorlampert/CHAOS

# Go into the repository
$ cd CHAOS/

# Run
$ PORT=8080 DATABASE_NAME=chaos go run cmd/chaos/main.go
```

# Build Docker Image:

```bash
docker build -t tiagorlampert/chaos:v5.0.0 --build-arg APP_VERSION=v5.0.0 .

docker run -it --rm -e PORT=8080 -p 8080:8080 tiagorlampert/chaos:v5.0.0
```

# Deploy on heroku (manual deployment)
```bash
$ git clone https://github.com/tiagorlampert/CHAOS
$ cd CHAOS/

$ heroku container:login
$ heroku create
$ heroku container:push web
$ heroku container:release web
$ heroku open

# Can be called from a url to test
# https://dashboard.heroku.com/new?button-url=https://github.com/tiagorlampert/CHAOS&template=https://github.com/tiagorlampert/CHAOS/tree/{branch_with_deploy}
```