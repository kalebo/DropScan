package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/gographics/imagick.v2/imagick"
)

//go:generate go-bindata index.html

func process_image(blob []byte, basename string) {
	const qr float64 = float64(65535) // quantum range max

	mw := imagick.NewMagickWand()
	//mw.SetSize(640, 480)
	mw.ReadImageBlob(blob)
	mw.SetFormat("png")

	c := mw.GetImageHeight()
	r := mw.GetImageWidth()
	nr, nc := (r*30)/100, (c*30)/100
	//mw.WriteImage("before.png")

	// TODO: only resize for images larger than a threshold
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
	mw.TrimImage(60)

	mw.WriteImage(basename + ".png")
	fmt.Println("Finished!")
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	indexTemplate.Execute(w, templateValues{IP.String(), PORT})
}

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("Got request!")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Malformated payload")
	}
	t := time.Now()
	// TODO: reduce potenial for name collsions by allowing users to provide basename
	process_image(body, "img_" + t.Format("20060102150405"))
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

type templateValues struct {
	IP string
	Port string
}

func init() {
    indexhtml, _ := Asset("index.html")
	indexTemplate, _ = template.New("index").Parse(string(indexhtml))
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
//	process_image(blob, "after")

	fmt.Printf("Server starting at http://%v%v\n\n", IP, PORT)
	log.Fatal(http.ListenAndServe(PORT, router))

}
