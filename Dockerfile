############################
# STEP 1 build executable binary
############################
FROM golang:1.22
# Create appuser.
ENV USER=appuser
ENV UID=10001 
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"
WORKDIR $GOPATH/src/
COPY . .
# Fetch dependencies.
RUN go get -d -v
# RUN go mod download
RUN apt-get update
RUN apt-get install sqlite3 libsqlite3-dev
RUN go build -o /go/bin/innovite
RUN mkdir -p /etc/temp
CMD ["/go/bin/innovite"]
#ENTRYPOINT ["/go/bin/innovite"]

# TODO: Ideally for production we should use scratch image which is light weight.
# We want the docker image to only contain the static go binary, but the sqlite needs cgo, hence TODO

############################
# STEP 2 build a small image
############################
#FROM scratch
#RUN apk add build-base
#RUN apk --no-cache add make git gcc libtool musl-dev ca-certificates dumb-init
#RUN go install github.com/mattn/go-sqlite3
# Import the user and group files from the builder.
#COPY --from=builder /etc/passwd /etc/passwd
#COPY --from=builder /etc/group /etc/group
# Copy our static executable.
#COPY --from=builder /go/bin/innovite /go/bin/innovite
#COPY --from=builder /etc/config.yml .
# Use an unprivileged user.
#USER appuser:appuser
# Run the hello binary.
#ENTRYPOINT ["/go/bin/innovite"]
