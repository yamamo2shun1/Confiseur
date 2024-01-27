package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/sstallion/go-hid"
)

type Layouts struct {
	Layout1 Layout
	Layout2 Layout
}

type Layout struct {
	Rows [][]string `toml:"rows"`
}

var err error
var d *hid.Device
var layouts Layouts
var remapRow []byte = make([]byte, 16)

func GetSettingPath() string {
	path := "nopath"

	hid.Enumerate(hid.VendorIDAny, hid.ProductIDAny, func(info *hid.DeviceInfo) error {
		if strings.Contains(info.ProductStr, "Setting") {
			path = info.Path
		}
		return nil
	})

	return path
}

func checkHid() {
	// Read the Manufacturer String.
	s, err := d.GetMfrStr()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Manufacturer String: %s\n", s)

	// Read the Product String.
	s, err = d.GetProductStr()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Product String: %s\n", s)

	// Read the Serial Number String.
	s, err = d.GetSerialNbr()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Serial Number String: %s\n", s)

	// Read Indexed String 1.
	s, err = d.GetIndexedStr(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Indexed String 1: %s\n", s)
}

func loadKeymap(val byte) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

	for i := 0; i < 5; i++ {
		remapRow[0] = 0x00
		remapRow[1] = 0xF0 + byte(i)
		remapRow[2] = val
		if _, err := d.Write(remapRow); err != nil {
			log.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)

		if _, err := d.Read(remapRow); err != nil {
			log.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)

		fmt.Fprint(w, "[\t")
		for j := 0; j < 13; j++ {
			fmt.Fprintf(w, "%s\t", KN[remapRow[j]])
		}
		fmt.Fprintln(w, "]\t")
	}
	w.Flush()
}

func writeKeymap(val byte) {
	for i := 0; i < 5; i++ {
		remapRow[0] = 0x00
		remapRow[1] = 0xF0 + byte(i)
		remapRow[2] = val
		for j := 0; j < 13; j++ {
			if val == 0x01 {
				remapRow[j+3] = SC[layouts.Layout1.Rows[i][j]]
			} else if val == 0x02 {
				remapRow[j+3] = SC[layouts.Layout2.Rows[i][j]]
			}
		}

		if _, err := d.Write(remapRow); err != nil {
			log.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func remap() {
	fmt.Println("-- Remap Layout ScanCode ---")

	inputfile := flag.String("f", "layouts.toml", "flag for input .toml file.")
	flag.Parse()

	_, err = toml.DecodeFile(*inputfile, &layouts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("::Layout1::")
	for i := 0; i < 5; i++ {
		fmt.Println(layouts.Layout1.Rows[i])
	}
	for i := 0; i < 5; i++ {
		fmt.Print("[ ")
		for j := 0; j < 13; j++ {
			if SC[layouts.Layout1.Rows[i][j]] == 0x00 {
				if err := hid.Exit(); err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Printf("0x%02X ", SC[layouts.Layout1.Rows[i][j]])
			}
		}
		fmt.Println("]")
	}
	fmt.Println("")

	fmt.Println("")
	fmt.Println("::Layout2::")
	for i := 0; i < 5; i++ {
		fmt.Println(layouts.Layout2.Rows[i])
	}

	for i := 0; i < 5; i++ {
		fmt.Print("[ ")
		for j := 0; j < 13; j++ {
			fmt.Printf("0x%02X ", SC[layouts.Layout2.Rows[i][j]])
		}
		fmt.Println("]")
	}
	fmt.Println("")
}

func saveToFlash() {
	remapRow[0] = 0x00
	remapRow[1] = 0xF5
	if _, err := d.Write(remapRow); err != nil {
		log.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	if _, err := d.Read(remapRow); err != nil {
		log.Fatal(err)
	}
	if remapRow[1] == 0xF5 {
		fmt.Println("Finish.")
	}
}

func main() {
	// Initialize the hid package.
	if err := hid.Init(); err != nil {
		log.Fatal(err)
	}

	// Open the device using the VID and PID.
	d, err = hid.OpenPath(GetSettingPath())
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "check":
		checkHid()
		fmt.Println("")
	case "load":
		initKN()

		// check current hardware layout
		fmt.Println("--- Current Hardware Layout ScanCode ---")
		fmt.Println("::Layout1::")
		loadKeymap(0x03)
		fmt.Println("")
		fmt.Println("::Layout2::")
		loadKeymap(0x04)
		fmt.Println("")
	case "remap":
		remap()

		fmt.Println("remap layout1&2")

		writeKeymap(0x01)
		writeKeymap(0x02)
		fmt.Println("")
	case "save":
		saveToFlash()
		fmt.Println("")
	case "ver":
		fmt.Println("C4NDY KeyVLM Configurator v0.3!")
		fmt.Println("")
	default:
	}

	time.Sleep(100 * time.Millisecond)

	// Finalize the hid package.
	if err := hid.Exit(); err != nil {
		log.Fatal(err)
	}
}
