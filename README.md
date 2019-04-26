# DropScan

This was a quick project to scratch one of my own itches. Its purpose is to streamline
the process of taking images of whiteboards, handwritten notes and line drawings
with my phone's camera and then cleaning them up so I can insert those images into documents.

It is structured as a simple web app that does the image processing on the
server side.

## Building 

To build this project you will need the following 

### 1. Install the ImageMagick MagickWand library
This can be installed in debian derivatives as follows:

``` sh
sudo apt install libmagickwand-dev
```

### 2. Go`get` the following go libraries and binaries:

```sh
go get "github.com/julienschmidt/httprouter"
go get "gopkg.in/gographics/imagick.v2/imagick"
```

And to install the go-bindata utilitity:
```sh
go get -u github.com/go-bindata/go-bindata/...
```
### 3. Generate assets and build

```sh
go generate
go build
```


