FROM golang:1.21.4-alpine3.18

# Change mirror for the image
RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.18/main/ > /etc/apk/repositories
RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.18/community/ >> /etc/apk/repositories

# Install necessary tools
RUN apk add --no-cache bash make

WORKDIR /IO-IMS
# Copy all files from the current directory to the working directory
ADD ./ ./

# Set the Go language environment variables
# Enable Go Module mode
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn \
    mode=release

RUN bash -c "ls && pwd && tree ."

# Download dependencies
RUN go mod download

# Start the service and provide the terminal
ENTRYPOINT ["/bin/sh", "-c", "make run"]


