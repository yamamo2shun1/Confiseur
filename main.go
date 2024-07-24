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
	Normal [][]string `toml:"normal"`
	Upper  [][]string `toml:"upper"`
}

var err error
var d *hid.Device
var layouts Layouts
var remapRows []byte = make([]byte, 16)

var maxRows = 5
var maxColumns = 13

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

	if strings.Contains(s, "C4NDY STK") {
		maxRows = 4
		maxColumns = 10
	}

	// Read the Serial Number String.
	s, err = d.GetSerialNbr()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Serial Number String: %s\n", s)

	/*
		// Read Indexed String 1.
		s, err = d.GetIndexedStr(1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Indexed String 1: %s\n", s)
	*/
}

func loadKeymap(val byte) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

	for i := 0; i < maxRows; i++ {
		remapRows[0] = 0x00
		remapRows[1] = 0xF0 + byte(i)
		remapRows[2] = val
		if _, err := d.Write(remapRows); err != nil {
			log.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)

		if _, err := d.Read(remapRows); err != nil {
			log.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)

		fmt.Fprint(w, "[\t")
		for j := 0; j < maxColumns; j++ {
			fmt.Fprintf(w, "%s\t", KN[remapRows[j]])
		}
		fmt.Fprintln(w, "]\t")
	}
	w.Flush()
}

func writeKeymap(val byte) {
	for i := 0; i < maxRows; i++ {
		remapRows[0] = 0x00
		remapRows[1] = 0xF0 + byte(i)
		remapRows[2] = val
		for j := 0; j < maxColumns; j++ {
			switch val {
			case 0x01:
				remapRows[j+3] = SC[layouts.Layout1.Normal[i][j]]
			case 0x02:
				remapRows[j+3] = SC[layouts.Layout1.Upper[i][j]]
			case 0x03:
				remapRows[j+3] = SC[layouts.Layout2.Normal[i][j]]
			case 0x04:
				remapRows[j+3] = SC[layouts.Layout2.Upper[i][j]]
			}
		}

		if _, err := d.Write(remapRows); err != nil {
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

	fmt.Println(layouts)

	fmt.Println("::Layout1::")
	fmt.Println("  Normal ->")
	for i := 0; i < maxRows; i++ {
		fmt.Println(layouts.Layout1.Normal[i])
	}
	for i := 0; i < maxRows; i++ {
		fmt.Print("[ ")
		for j := 0; j < maxColumns; j++ {
			if SC[layouts.Layout1.Normal[i][j]] == 0x00 {
				if err := hid.Exit(); err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Printf("0x%02X ", SC[layouts.Layout1.Normal[i][j]])
			}
		}
		fmt.Println("]")
	}
	fmt.Println("  Upper ->")
	for i := 0; i < maxRows; i++ {
		fmt.Println(layouts.Layout1.Upper[i])
	}
	for i := 0; i < maxRows; i++ {
		fmt.Print("[ ")
		for j := 0; j < maxColumns; j++ {
			if SC[layouts.Layout1.Upper[i][j]] == 0x00 {
				if err := hid.Exit(); err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Printf("0x%02X ", SC[layouts.Layout1.Upper[i][j]])
			}
		}
		fmt.Println("]")
	}
	fmt.Println("")

	fmt.Println("")
	fmt.Println("::Layout2::")
	fmt.Println("  Normal ->")
	for i := 0; i < maxRows; i++ {
		fmt.Println(layouts.Layout2.Normal[i])
	}
	for i := 0; i < maxRows; i++ {
		fmt.Print("[ ")
		for j := 0; j < maxColumns; j++ {
			fmt.Printf("0x%02X ", SC[layouts.Layout2.Normal[i][j]])
		}
		fmt.Println("]")
	}
	fmt.Println("  Upper ->")
	for i := 0; i < maxRows; i++ {
		fmt.Println(layouts.Layout2.Upper[i])
	}
	for i := 0; i < maxRows; i++ {
		fmt.Print("[ ")
		for j := 0; j < maxColumns; j++ {
			fmt.Printf("0x%02X ", SC[layouts.Layout2.Upper[i][j]])
		}
		fmt.Println("]")
	}
	fmt.Println("")
}

func saveToFlash() {
	remapRows[0] = 0x00
	remapRows[1] = 0xF5
	if _, err := d.Write(remapRows); err != nil {
		log.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	if _, err := d.Read(remapRows); err != nil {
		log.Fatal(err)
	}
	if remapRows[1] == 0xF5 {
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

	// Read the Product String.
	s, err := d.GetProductStr()
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(s, "C4NDY STK") {
		maxRows = 4
		maxColumns = 10
	}

	switch os.Args[1] {
	case "check":
		checkHid()
		fmt.Println("")
	case "help":
		fmt.Println("META OPTIONS")
		fmt.Println("  help     Display this information.")
		fmt.Println("  ver      Show the version of the tool installed.")
		fmt.Println("  check    Show information on C4NDY KeyVLM/STK connected to PC/Mac.")
		fmt.Println("  load     Show the current key names of the keyboard.")
		fmt.Println("  remap    Write the keyboard with the keymap set in layouts.toml.")
		fmt.Println("  remap -f <file>")
		fmt.Println("           Write the keymap set in the specified .toml to the keydoad.")
		fmt.Println("  save     Save the keymap written by \"-r or --remap\" to the memory area.")
		fmt.Println("")
	case "load":
		initKN()

		// check current hardware layout
		fmt.Println("--- Current Hardware Layout ScanCode ---")
		fmt.Println("::Layout1::")
		fmt.Println("  Normal ->")
		loadKeymap(0x05)
		fmt.Println("  Upper ->")
		loadKeymap(0x06)
		fmt.Println("")
		fmt.Println("::Layout2::")
		fmt.Println("  Normal ->")
		loadKeymap(0x07)
		fmt.Println("  Upper ->")
		loadKeymap(0x08)
		fmt.Println("")
	case "remap":
		remap()

		fmt.Println("remap layout1&2(Normal/Upper)")

		writeKeymap(0x01)
		writeKeymap(0x02)
		writeKeymap(0x03)
		writeKeymap(0x04)
		fmt.Println("")
	case "save":
		saveToFlash()
		fmt.Println("")
	case "ver":
		fmt.Println("C4NDY KeyConfigurator v0.4!")
		fmt.Println("")
	default:
	}

	time.Sleep(100 * time.Millisecond)

	// Finalize the hid package.
	if err := hid.Exit(); err != nil {
		log.Fatal(err)
	}
}
