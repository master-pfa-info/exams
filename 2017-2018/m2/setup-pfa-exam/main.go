package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	goVersion = "1.9.2"
)

var (
	setupGo  = flag.Bool("init", false, "init workspace")
	saveWork = flag.Bool("save", false, "save work for grade")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("exam-pfa: ")

	flag.Parse()

	switch {
	case *setupGo:
		doSetupGo()
	case *saveWork:
		doSaveWork()
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func doSetupGo() {
	goroot, err := installGo(goVersion)
	if err != nil {
		log.Fatalf("could not install Go-%v: %v", goVersion, err)
	}

	log.Printf("goroot=%q\n", goroot)
	gopath := getGoPath()
	log.Printf("gopath=%q\n", gopath)
	srcdir := filepath.Join(gopath, "src")
	err = os.MkdirAll(srcdir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("srcdir=%q\n", srcdir)

	usr, err := user.Current()
	if err != nil {
		log.Fatalf("could not get current user: %v", err)
	}

	usrdir := filepath.Join(srcdir, "uca.fr", usr.Username)
	log.Printf("finals=%q\n", usrdir)
	err = os.MkdirAll(usrdir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	mbrotDir := filepath.Join(usrdir, "mbrot")
	err = os.MkdirAll(mbrotDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(
		filepath.Join(mbrotDir, "main.go"),
		[]byte(mandel),
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func doSaveWork() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	gopath := getGoPath()
	srcdir := filepath.Join(gopath, "src")
	usrdir := filepath.Join(srcdir, "uca.fr", usr.Username)

	log.Printf("gopath=%q\n", gopath)
	log.Printf("srcdir=%q\n", srcdir)
	log.Printf("finals=%q\n", usrdir)

	log.Printf("saving everything under %q...", usrdir)
	tmpdir, err := ioutil.TempDir("", "pfa-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	save := filepath.Join(tmpdir, usr.Username+"-pfa.tar.gz")
	cmd := exec.Command("tar", "zcvf", save, "src")
	cmd.Dir = gopath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	final := filepath.Join(usr.HomeDir, "M_"+usr.Username, usr.Username+"-pfa.tar.gz")
	cmd = exec.Command("/bin/cp", save, final)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("please send an email to binet@cern.ch with the file %q", final)
	log.Printf("specifying your name (%s).", usr.Name)
}

func installGo(v string) (string, error) {
	log.Printf("downloading go-%v...", v)
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("could not get current user: %v", err)
	}

	var gotar io.ReadCloser
	if true {
		burl := "https://golang.org/dl/go" + v + ".linux-amd64.tar.gz"
		resp, err := http.Get(burl)
		if err != nil {
			return "", err
		}
		gotar = resp.Body
		//defer resp.Body.Close()
	} else {
		f, err := os.Open("/home/SCIENCES/" + usr.Username + "/Public/Sebinet/go" + v + ".linux-amd64.tar.gz")
		if err != nil {
			log.Fatal(err)
		}
		gotar = f
	}
	defer gotar.Close()

	goroot := filepath.Join(usr.HomeDir, "M_"+usr.Username, "go-"+v)

	err = os.MkdirAll(goroot, 0755)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("tar", "zxf", "-")
	cmd.Dir = goroot
	cmd.Stdin = gotar
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	goroot = filepath.Join(goroot, "go")
	gopath := getGoPath()

	os.Setenv("GOROOT", goroot)
	os.Setenv("PATH", filepath.Join(goroot, "bin")+":"+os.Getenv("PATH"))

	fname := filepath.Join(usr.HomeDir, ".bashrc")
	err = appendFile(
		fname,
		[]byte(fmt.Sprintf(`
### AUTOMATICALLY added by setup-pfa
export GOROOT=%q
export GOPATH=%q
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
`,
			goroot,
			gopath,
		)),
	)
	if err != nil {
		log.Fatalf("could not modify bash_profile: %v", err)
	}

	return goroot, nil
}

func appendFile(fname string, data []byte) error {
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Seek(0, 2)
	if err != nil {
		log.Fatalf("could not seek to the end: %v", err)
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func getGoRoot() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("could not get current user: %v", err)
	}
	goroot := filepath.Join(usr.HomeDir, "M_"+usr.Username, "go-"+goVersion)

	return goroot
}

func getGoPath() string {
	if true {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("could not get current user: %v", err)
		}
		gopath := filepath.Join(usr.HomeDir, "M_"+usr.Username, "go")
		os.Setenv("GOPATH", gopath)
		return gopath
	}

	p := os.Getenv("GOPATH")
	if p != "" {
		return p
	}
	raw, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSuffix(string(raw), "\n")
}

const mandel = `package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"
)

const (
	output = "out.png"
	width  = 2048
	height = 1024
)

func main() {
	f, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	start := time.Now()
	img := create(width, height)
	delta := time.Since(start)
	fmt.Printf("time=%v\n", delta)

	if err = png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}

// create fills one pixel at a time.
//
// time=??? <<< put the timing you find here.
func create(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			m.Set(i, j, pixel(i, j, width, height))
		}
	}
	return m
}

// create1 creates a Mandelbrot image.
//
// time=???
func create1(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))
	// ???
	return m
}

// create2 creates a Mandelbrot image.
//
// time=???
func create2(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))
	// ???
	return m
}

// create3 creates a Mandelbrot image.
//
// time=???
func create3(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))
	// ???
	return m
}

// pixel returns the color of a Mandelbrot fractal at the given point.
func pixel(i, j, width, height int) color.Color {
	const complexity = 1024

	xi := norm(i, width, -1.0, 2)
	yi := norm(j, height, -1, 1)

	const maxI = 1000
	x, y := 0., 0.

	for i := 0; (x*x+y*y < complexity) && i < maxI; i++ {
		x, y = x*x-y*y+xi, 2*x*y+yi
	}

	return color.Gray{uint8(x)}
}

func norm(x, total int, min, max float64) float64 {
	return (max-min)*float64(x)/float64(total) - max
}
`
