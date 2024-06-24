FROM golang:latest

LABEL description="This project consists in creating a web forum."
LABEL authors="@lodiouf, @jgoudiab, @mabdelwe, @mouhamadoudiouf"

WORKDIR /forum

COPY . /forum/

# Garbage Collector configuration
ENV GOGC=200

RUN go build -o forum .

# RUN apk update && apk add --no-cache bash 

EXPOSE 9090

 CMD [ "./forum" ]

#HOW TO RUN THE DOCKERFILE

#1st STEP 
#make your image  [ docker build -t <image> . ]
#execut a container from your image [ docker run -d -p 8080:9090 <image> ]
