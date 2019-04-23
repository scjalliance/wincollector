To build with the appropriate variable set, run `GOOS=windows go build -ldflags "-X main.VERSION=$(git rev-list -1 HEAD)"` or some appropriate approximation of this.
