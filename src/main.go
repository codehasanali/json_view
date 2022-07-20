package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/fatih/color"
	"github.com/nwidger/jsoncolor"
	"github.com/spf13/pflag"
)

func main() {
	display_help := pflag.BoolP("Yardım", "h", false, "yardım")

	pflag.Parse()

	if *display_help {
		displayUsage()
		os.Exit(0)
	}

	log.SetHandler(cli.New(os.Stderr))

	log.SetLevel(log.WarnLevel)

	r, err := openStdinOrFile()

	checkError("Girdi alınırken hata oluştu.", err)

	prettyPrint(r)
}



func prettyPrint(data io.Reader) {
	b, _ := ioutil.ReadAll(data)

	f := jsoncolor.NewFormatter()

	f.SpaceColor = color.New(color.FgRed, color.Bold)
	f.CommaColor = color.New(color.FgWhite, color.Bold)
	f.ColonColor = color.New(color.FgYellow)
	f.ObjectColor = color.New(color.FgBlue, color.Bold)
	f.ArrayColor = color.New(color.FgWhite)
	f.FieldColor = color.New(color.FgGreen)
	f.StringColor = color.New(color.FgMagenta, color.Bold)
	f.TrueColor = color.New(color.FgCyan, color.Bold)
	f.FalseColor = color.New(color.FgHiRed)
	f.NumberColor = color.New(color.FgHiYellow)
	f.NullColor = color.New(color.FgWhite, color.Bold)
	f.StringQuoteColor = color.New(color.FgBlue, color.Bold)

	dst := &bytes.Buffer{}
	err := f.Format(dst, b)
	checkError("Json renklendirme hatası oluştu. Json dosyasın da hata olabilir!", err)

	fmt.Println(dst.String())
}

func displayUsage() {
	fmt.Printf("Kullanımı: json_view [<flags>] [DOSYA]\n\n")
	fmt.Printf("Örnek:\n\tjson_view dosya.json\n")
	fmt.Printf("\tcat dosya.json | json_view \n\n")
	pflag.PrintDefaults()
}

func openStdinOrFile() (io.Reader, error) {
	if len(pflag.Args()) == 1 {
		r, err := os.Open(pflag.Arg(0))
		if err != nil {
			return nil, err
		}
		return r, err
	}

	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeNamedPipe == 0 {
		return nil, fmt.Errorf("Json dosyası  belirtilmedi.")
	}

	return os.Stdin, err
}

func checkError(message string, err error) {
	if err != nil {
		log.WithError(err).Fatal(message)
	}
}