# imageMock

This is a really simple mock image server to host mock images at a specified size for development purposes.

To build:
* go get
* go build

To run:
./imageMock

It hosts on port 8080 and will host an 800x600 image at the root URL.  To generate a custom size, pass it <width>x<height> (i.e. http://localhost:8080/1024x768)
