package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/shopspring/decimal"
	"github.com/sstallion/go-hid"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type Layouts struct {
	Layout1 Layout
	Layout2 Layout
}

type Layout struct {
	Normal    [][][]string `toml:"normal"`
	Upper     [][][]string `toml:"upper"`
	Stick     [][][]string `toml:"stick"`
	Led       [][]byte     `toml:"led"`
	Intensity []float64    `toml:"intensity"`
}

var err error

var hidDevices []*hid.Device
var connectedDeviceNum = 0

var layouts Layouts
var remapRows []byte = make([]byte, 32)

var maxRows = 5
var maxColumns = 13
var isStk = false

func getConnectedDeviceNum() int {
	return connectedDeviceNum
}

func getSettingPaths() []string {
	var path []string

	err := hid.Enumerate(hid.VendorIDAny, hid.ProductIDAny, func(info *hid.DeviceInfo) error {
		if strings.Contains(info.ProductStr, "C4NDY") && info.Usage == 1 {
			path = append(path, info.Path)
		}
		//fmt.Printf("ProductStr: %s/Usage: %d\n", info.ProductStr, info.Usage)
		return nil
	})
	if err != nil {
		return nil
	}

	return path
}

func getConnectedDeviceList() []string {
	var deviceList []string

	for _, device := range hidDevices {
		deviceName, _ := device.GetProductStr()
		deviceList = append(deviceList, deviceName)
	}

	return deviceList
}

func checkHid() {
	for i := 0; i < connectedDeviceNum; i++ {
		fmt.Printf("::%d::\n", i)
		// Read the Manufacturer String.
		s, err := hidDevices[i].GetMfrStr()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Manufacturer String: %s\n", s)

		// Read the Product String.
		s, err = hidDevices[i].GetProductStr()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Product String: %s\n", s)

		if strings.Contains(s, "C4NDY STK") {
			maxRows = 4
			maxColumns = 10
		}

		// Read the Serial Number String.
		s, err = hidDevices[i].GetSerialNbr()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Serial Number String: %s\n", s)
		fmt.Println("")
	}
}

func openConnectedHIDDevices() {
	settingPaths := getSettingPaths()
	connectedDeviceNum = len(settingPaths)
	hidDevices = make([]*hid.Device, connectedDeviceNum)

	for i, settingPath := range settingPaths {
		hidDevices[i], _ = hid.OpenPath(settingPath)
	}
}

func checkKeyboardType(index int) {
	// Read the Product String.
	s, err := hidDevices[index].GetProductStr()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	if strings.Contains(s, "C4NDY KeyVLM") {
		isStk = false
		maxRows = 5
		maxColumns = 13
	}
	if strings.Contains(s, "C4NDY STK") {
		isStk = true
		maxRows = 4
		maxColumns = 10
	}
}

func loadKeymap(index int, val byte) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

	for i := 0; i < maxRows; i++ {
		if (val == 0x13 || val == 0x1B) && i > 1 {
			continue
		}
		if (val == 0x14 || val == 0x1C) && i > 2 {
			continue
		}
		if (val == 0x15 || val == 0x1D) && i > 0 {
			continue
		}
		remapRows[0] = 0x00
		remapRows[1] = 0xF0 + byte(i)
		remapRows[2] = val
		if _, err := hidDevices[index].Write(remapRows); err != nil {
			log.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)

		if _, err := hidDevices[index].Read(remapRows); err != nil {
			log.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)

		_, err := fmt.Fprint(w, "[\t")
		if err != nil {
			return
		}
		for j := 0; j < maxColumns; j++ {
			if (val == 0x13 || val == 0x1B) && j > 8 {
				continue
			}
			if (val == 0x14 || val == 0x1C) && j > 2 {
				continue
			}
			if (val == 0x15 || val == 0x1D) && j > 0 {
				continue
			}

			if val == 0x14 || val == 0x1C {
				_, err := fmt.Fprintf(w, "%02X\t", remapRows[j])
				if err != nil {
					return
				}
			} else if val == 0x15 || val == 0x1D {
				rate, _ := decimal.NewFromFloat(float64(remapRows[j]) / 255.0).Round(2).Float64()
				_, err := fmt.Fprintf(w, "%f\t", rate)
				if err != nil {
					return
				}
			} else {
				_, err := fmt.Fprintf(w, "{%s, %02X}\t", KEYNAME[remapRows[2*j]], remapRows[2*j+1])
				if err != nil {
					return
				}
			}
		}
		_, err = fmt.Fprintln(w, "]\t")
		if err != nil {
			return
		}
	}
	err := w.Flush()
	if err != nil {
		return
	}
}

func writeKeymap(index int, val byte) {
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
		if (val == 0x05 || val == 0x0D) && i > 0 {
			continue
		}

		remapRows[0] = 0x00
		remapRows[1] = 0xF0 + byte(i)
		remapRows[2] = val
		for j := 0; j < maxColumns; j++ {
			if (val == 0x03 || val == 0x0B) && j > 8 {
				continue
			}
			if (val == 0x04 || val == 0x0C) && j > 2 {
				continue
			}
			if (val == 0x05 || val == 0x0D) && j > 0 {
				continue
			}

			switch val {
			case 0x01:
				remapRows[(2*j)+3] = KEYCODE[layouts.Layout1.Normal[i][j][0]]
				if modifiers, err := strconv.ParseUint(layouts.Layout1.Normal[i][j][1], 2, 8); err == nil {
					remapRows[(2*j+1)+3] = byte(modifiers)
				}
			case 0x02:
				remapRows[(2*j)+3] = KEYCODE[layouts.Layout1.Upper[i][j][0]]
				if modifiers, err := strconv.ParseUint(layouts.Layout1.Upper[i][j][1], 2, 8); err == nil {
					remapRows[(2*j+1)+3] = byte(modifiers)
				}
			case 0x03:
				remapRows[(2*j)+3] = KEYCODE[layouts.Layout1.Stick[i][j][0]]
				if modifiers, err := strconv.ParseUint(layouts.Layout1.Stick[i][j][1], 2, 8); err == nil {
					remapRows[(2*j+1)+3] = byte(modifiers)
				}
			case 0x04:
				remapRows[j+3] = layouts.Layout1.Led[i][j]
			case 0x05:
				remapRows[j+3] = byte(layouts.Layout1.Intensity[j] * 255)
			case 0x09:
				remapRows[(2*j)+3] = KEYCODE[layouts.Layout2.Normal[i][j][0]]
				if modifiers, err := strconv.ParseUint(layouts.Layout2.Normal[i][j][1], 2, 8); err == nil {
					remapRows[(2*j+1)+3] = byte(modifiers)
				}
			case 0x0A:
				remapRows[(2*j)+3] = KEYCODE[layouts.Layout2.Upper[i][j][0]]
				if modifiers, err := strconv.ParseUint(layouts.Layout2.Upper[i][j][1], 2, 8); err == nil {
					remapRows[(2*j+1)+3] = byte(modifiers)
				}
			case 0x0B:
				remapRows[(2*j)+3] = KEYCODE[layouts.Layout2.Stick[i][j][0]]
				if modifiers, err := strconv.ParseUint(layouts.Layout2.Stick[i][j][1], 2, 8); err == nil {
					remapRows[(2*j+1)+3] = byte(modifiers)
				}
			case 0x0C:
				remapRows[j+3] = layouts.Layout2.Led[i][j]
			case 0x0D:
				remapRows[j+3] = byte(layouts.Layout2.Intensity[j] * 255)
			}
		}

		if _, err := hidDevices[index].Write(remapRows); err != nil {
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
			if KEYCODE[layouts.Layout1.Normal[i][j][0]] == 0x00 {
				if err := hid.Exit(); err != nil {
					log.Fatal(err)
				}
			} else {
				if modifiers, err := strconv.ParseUint(layouts.Layout1.Normal[i][j][1], 2, 8); err == nil {
					fmt.Printf("{0x%02X, 0x%02X} ", KEYCODE[layouts.Layout1.Normal[i][j][0]], modifiers)
				}
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
				if KEYCODE[layouts.Layout1.Upper[i][j][0]] == 0x00 {
					if err := hid.Exit(); err != nil {
						log.Fatal(err)
					}
				} else {
					if modifiers, err := strconv.ParseUint(layouts.Layout1.Upper[i][j][1], 2, 8); err == nil {
						fmt.Printf("{0x%02X, 0x%02X} ", KEYCODE[layouts.Layout1.Upper[i][j][0]], modifiers)
					}
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
			for j := 0; j < 9; j++ {
				if KEYCODE[layouts.Layout1.Stick[i][j][0]] == 0x00 {
					if err := hid.Exit(); err != nil {
						log.Fatal(err)
					}
				} else {
					if modifiers, err := strconv.ParseUint(layouts.Layout1.Stick[i][j][1], 2, 8); err == nil {
						fmt.Printf("{0x%02X, 0x%02X} ", KEYCODE[layouts.Layout1.Stick[i][j][0]], modifiers)
					}
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
	if len(layouts.Layout1.Intensity) != 0 {
		fmt.Println("  Intensity ->")
		fmt.Printf("%f", layouts.Layout1.Intensity)
		fmt.Println("]")
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
			if modifiers, err := strconv.ParseUint(layouts.Layout2.Normal[i][j][1], 2, 8); err == nil {
				fmt.Printf("{0x%02X, 0x%02X} ", KEYCODE[layouts.Layout2.Normal[i][j][0]], modifiers)
			}
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
				if modifiers, err := strconv.ParseUint(layouts.Layout2.Upper[i][j][1], 2, 8); err == nil {
					fmt.Printf("{0x%02X, 0x%02X} ", KEYCODE[layouts.Layout2.Upper[i][j][0]], modifiers)
				}
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
			for j := 0; j < 9; j++ {
				if modifiers, err := strconv.ParseUint(layouts.Layout2.Stick[i][j][1], 2, 8); err == nil {
					fmt.Printf("{0x%02X, 0x%02X} ", KEYCODE[layouts.Layout2.Stick[i][j][0]], modifiers)
				}
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
	if len(layouts.Layout2.Intensity) != 0 {
		fmt.Println("  Intensity ->")
		fmt.Printf("%f", layouts.Layout2.Intensity)
		fmt.Println("]")
	}
	fmt.Println("")
}

func saveToFlash(index int) {
	remapRows[0] = 0x00
	remapRows[1] = 0xF5
	if _, err := hidDevices[index].Write(remapRows); err != nil {
		log.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	if _, err := hidDevices[index].Read(remapRows); err != nil {
		log.Fatal(err)
	}
	if remapRows[1] == 0xF5 {
		fmt.Println("Finish.")
	}
}

func restart(index int) {
	remapRows[0] = 0x00
	remapRows[1] = 0xF6
	if _, err := hidDevices[index].Write(remapRows); err != nil {
		log.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	if _, err := hidDevices[index].Read(remapRows); err != nil {
		log.Fatal(err)
	}
	if remapRows[1] == 0xF6 {
		fmt.Println("Finish.")
	}
}

func checkLEDColor(index int, value int) {
	remapRows[0] = 0x00
	remapRows[1] = 0xF7
	remapRows[2] = byte((value >> 16) & 0x000000FF)
	remapRows[3] = byte((value >> 8) & 0x000000FF)
	remapRows[4] = byte(value & 0x000000FF)

	if _, err := hidDevices[index].Write(remapRows); err != nil {
		log.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
}

func changeLEDIntensity(index int, value float64) {
	remapRows[0] = 0x00
	remapRows[1] = 0xF8
	remapRows[2] = byte(value * 255)

	if _, err := hidDevices[index].Write(remapRows); err != nil {
		log.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
}

func factoryReset(index int) {
	remapRows[0] = 0x00
	remapRows[1] = 0xF9
	if _, err := hidDevices[index].Write(remapRows); err != nil {
		log.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	if _, err := hidDevices[index].Read(remapRows); err != nil {
		log.Fatal(err)
	}
	if remapRows[1] == 0xF8 {
		fmt.Println("Finish.")
	}
}
