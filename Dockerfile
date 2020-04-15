FROM golang:alpine
#FROM golang:onbuild

# Set necessary environmet variables needed for our image
#ENV databaseName=golang-db/
#CONNECTIONSTRING=mongodb://localhost:27017

# Move to working directory /build

WORKDIR /src/golang-assignment

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
#WORKDIR "/go/src/golang-assignment"

# Copy binary from build to main folder
#RUN cp /build/main .

# Export necessary port
EXPOSE 8000

# Command to run when starting the container
CMD ["./main"]
