# Get the base image
FROM golang:1.23

# Set the working directory
WORKDIR /src

# Copy the go modules manifest files
COPY go.mod go.sum ./

# Download go depnednecies 
RUN go mod download

# Copy the entire project into the container
COPY . .

# Install `gin` to use hot-reload
RUN go install github.com/codegangsta/gin@latest

# Add the GO bin directory to your path
ENV PATH=$PATH:/go/bin

# Expose the port the app listens on
EXPOSE 8080

CMD ["gin", "--port", "8080", "--appPort", "8081", "run", "main.go"]

