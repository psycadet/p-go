# p-go

Simplistic image/file service.

Created in order to satisfy all my needs for similar sites (khm, imgur).


<p align="center"><img src ="https://i.imgur.com/juX1G8q.png" /></p>
*Gopher Artwork by Ashley McNamara*

Rewritten from [Python3](https://github.com/whiteshtef/p) to Go.

## Deployment

### 1. Direct
    
`go get && go run main.go`
    

The images are stored in the `./storage/` directory.


### 2. Docker

`sudo docker build -t p .`

`sudo docker run -v /machinedirectory:/go/src/github.com/whiteshtef/p/storage -p 80:80 p`


## Usage

   Open in browser, select an image, hit submit - bam. You are redirected to the hosted version of the image you just uploaded.


