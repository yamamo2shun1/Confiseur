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
	Stick  [][]string `toml:"stick"`
	Led    [][]byte   `toml:"led"`
}

var err error
var d *hid.Device
var layouts Layouts
var remapRows []byte = make([]byte, 16)

var maxRows = 5
var maxColumns = 13
var isStk = false

func GetSettingPath() string {
	path := "nopath"

	hid.Enumerate(hid.VendorIDAny, hid.ProductIDAny, func(info *hid.DeviceInfo) error {
		if strings.Contains(info.ProductStr, "C4NDY") && info.Usage == 1 {
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
		if (val == 0x13 || val == 0x1B) && i > 1 {
			continue
		}
		if (val == 0x14 || val == 0x1C) && i > 2 {
			continue
		}
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
			if (val == 0x13 || val == 0x1B) && j > 3 {
				continue
			}
			if (val == 0x14 || val == 0x1C) && j > 2 {
				continue
			}

			if val == 0x14 || val == 0x1C {
				fmt.Fprintf(w, "%02X\t", remapRows[j])
			} else {
				fmt.Fprintf(w, "%s\t", KN[remapRows[j]])
			}
		}
		fmt.Fprintln(w, "]\t")
	}
	w.Flush()
}

func writeKeymap(val byte) {
	for i := range remapRows {
		remapRows[i] = 0x00
	}
	for i := 0; i < maxRows; i++ {
		if (val == 0x03 || val == 0x0B) && i > 1 {
			continue
		}
		if (val == 0x04 || val == 0x0C) && i > 2 {
			continue
		}

		remapRows[0] = 0x00
		remapRows[1] = 0xF0 + byte(i)
		remapRows[2] = val
		for j := 0; j < maxColumns; j++ {
			if (val == 0x03 || val == 0x0B) && j > 3 {
				continue
			}
			if (val == 0x04 || val == 0x0C) && j > 2 {
				continue
			}

			switch val {
			case 0x01:
				remapRows[j+3] = SC[layouts.Layout1.Normal[i][j]]
			case 0x02:
				remapRows[j+3] = SC[layouts.Layout1.Upper[i][j]]
			case 0x03:
				remapRows[j+3] = SC[layouts.Layout1.Stick[i][j]]
			case 0x04:
				remapRows[j+3] = layouts.Layout1.Led[i][j]
			case 0x09:
				remapRows[j+3] = SC[layouts.Layout2.Normal[i][j]]
			case 0x0A:
				remapRows[j+3] = SC[layouts.Layout2.Upper[i][j]]
			case 0x0B:
				remapRows[j+3] = SC[layouts.Layout2.Stick[i][j]]
			case 0x0C:
				remapRows[j+3] = layouts.Layout2.Led[i][j]
			}
		}

		if _, err := d.Write(remapRows); err != nil {
			log.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func remap(inputfile string) {
	fmt.Println("-- Remap Layout ScanCode ---")

	_, err = toml.DecodeFile(inputfile, &layouts)
	if err != nil {
		log.Fatal(err)
	}

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
	if len(layouts.Layout1.Upper) != 0 {
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
	}
	if len(layouts.Layout1.Stick) != 0 {
		fmt.Println("  Stick ->")
		for i := 0; i < 2; i++ {
			fmt.Println(layouts.Layout1.Stick[i])
		}
		for i := 0; i < 2; i++ {
			fmt.Print("[ ")
			for j := 0; j < 4; j++ {
				if SC[layouts.Layout1.Stick[i][j]] == 0x00 {
					if err := hid.Exit(); err != nil {
						log.Fatal(err)
					}
				} else {
					fmt.Printf("0x%02X ", SC[layouts.Layout1.Stick[i][j]])
				}
			}
			fmt.Println("]")
		}
	}
	if len(layouts.Layout1.Led) != 0 {
		fmt.Println("  Led ->")
		for i := 0; i < 3; i++ {
			fmt.Print("[ ")
			for j := 0; j < 3; j++ {
				fmt.Printf("0x%02X ", layouts.Layout1.Led[i][j])
			}
			fmt.Println("]")
		}
	}
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
	if len(layouts.Layout2.Upper) != 0 {
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
	}
	if len(layouts.Layout2.Stick) != 0 {
		fmt.Println("  Stick ->")
		for i := 0; i < 2; i++ {
			fmt.Println(layouts.Layout2.Stick[i])
		}
		for i := 0; i < 2; i++ {
			fmt.Print("[ ")
			for j := 0; j < 4; j++ {
				fmt.Printf("0x%02X ", SC[layouts.Layout2.Stick[i][j]])
			}
			fmt.Println("]")
		}
	}
	if len(layouts.Layout2.Led) != 0 {
		fmt.Println("  Led ->")
		for i := 0; i < 3; i++ {
			fmt.Print("[ ")
			for j := 0; j < 3; j++ {
				fmt.Printf("0x%02X ", layouts.Layout2.Led[i][j])
			}
			fmt.Println("]")
		}
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
		isStk = true
		maxRows = 4
		maxColumns = 10
	}

	checkFlag := flag.Bool("check", false, "Show information on C4NDY KeyVLM/STK connected to PC/Mac.")
	loadFlag := flag.Bool("load", false, "Show the current key names of the keyboard.")
	remapFlag := flag.Bool("remap", false, "Write the keyboard with the keymap set in layouts.toml.")
	saveFlag := flag.Bool("save", false, "Save the keymap written by \"-remap\" to the memory area.")
	verFlag := flag.Bool("version", false, "Show the version of the tool installed.")
	inputfile := flag.String("file", "layouts.toml", "Write the keymap set in the specified .toml to the keyboard.")

	flag.Parse()

	if *checkFlag {
		checkHid()
		fmt.Println("")
	} else if *loadFlag {
		initKN()

		// check current hardware layout
		fmt.Println("--- Current Hardware Layout ScanCode ---")
		fmt.Println("::Layout1::")
		fmt.Println("  Normal ->")
		loadKeymap(0x11)
		if isStk {
			fmt.Println("  Upper ->")
			loadKeymap(0x12)
			fmt.Println("  Stick ->")
			loadKeymap(0x13)
			fmt.Println("  Led ->")
			loadKeymap(0x14)
		}
		fmt.Println("")
		fmt.Println("::Layout2::")
		fmt.Println("  Normal ->")
		loadKeymap(0x19)
		if isStk {
			fmt.Println("  Upper ->")
			loadKeymap(0x1A)
			fmt.Println("  Stick ->")
			loadKeymap(0x1B)
			fmt.Println("  Led ->")
			loadKeymap(0x1C)
		}
		fmt.Println("")
	} else if *remapFlag {
		remap(*inputfile)

		fmt.Println("remap layout1&2(Normal/Upper)")

		writeKeymap(0x01)
		if isStk {
			writeKeymap(0x02)
			writeKeymap(0x03)
			writeKeymap(0x04)
		}
		writeKeymap(0x09)
		if isStk {
			writeKeymap(0x0A)
			writeKeymap(0x0B)
			writeKeymap(0x0C)
		}
		fmt.Println("")
	} else if *saveFlag {
		saveToFlash()
		fmt.Println("")
	} else if *verFlag {
		fmt.Println("C4NDY KeyConfigurator v1.1.0!")
		fmt.Println("")
	}

	time.Sleep(100 * time.Millisecond)

	// Finalize the hid package.
	if err := hid.Exit(); err != nil {
		log.Fatal(err)
	}
}
