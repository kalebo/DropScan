package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/gographics/imagick.v2/imagick"
)

func process_image(blob []byte) {
	const qr float64 = float64(65535) // quantum range max

	mw := imagick.NewMagickWand()
	//mw.SetSize(640, 480)
	mw.ReadImageBlob(blob)
	mw.SetFormat("png")

	c := mw.GetImageHeight()
	r := mw.GetImageWidth()
	nr, nc := (r*15)/100, (c*15)/100
	//mw.WriteImage("before.png")

	fmt.Printf("Resizing from %vx%v to %vx%v ...\n", r, c, nr, nc)
	mw.ResizeImage(nr, nc, imagick.FILTER_LANCZOS2, 1)

	fmt.Println("Processing filters...")
	kernel := imagick.NewKernelInfoBuiltIn(imagick.KERNEL_DOG, "15,100,0")
	kernel.Scale(1.0, imagick.KERNEL_NORMALIZE_VALUE)

	// filters
	mw.MorphologyImage(imagick.MORPHOLOGY_CONVOLVE, 1, kernel)
	mw.NegateImage(false)
	mw.NormalizeImage()
	mw.BlurImage(0, 0.5)
	mw.LevelImage(0.6*qr, 0.1, 0.91*qr)

	mw.WriteImage("after.png")
	fmt.Println("Finished!")
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	indexTemplate.Execute(w, IP.String())
}

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("Got request!")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Malformated payload")
	}
	process_image(body)
}

// Borrowed from from SO user Mr. Wang from Next Door
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

var indexTemplate *template.Template
var PORT string
var IP net.IP

func init() {
	indexTemplate, _ = template.ParseFiles("index.html")
	PORT = ":8080"
	IP = GetOutboundIP()
}

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/upload", Upload)

	//	blob, _ := ioutil.ReadFile("before.png")
	//	process_image(blob)

	fmt.Printf("Server starting at %v%v\n\n", IP, PORT)
	log.Fatal(http.ListenAndServe(PORT, router))

}
