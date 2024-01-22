package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/sstallion/go-hid"
)

var keyList [91]string = [91]string{
	"A",           // a
	"B",           // b
	"C",           // c
	"D",           // d
	"E",           // e
	"F",           // f
	"G",           // g
	"H",           // h
	"I",           // i
	"J",           // j
	"K",           // k
	"L",           // l
	"M",           // m
	"N",           // n
	"O",           // o
	"P",           // p
	"Q",           // q
	"R",           // r
	"S",           // s
	"T",           // t
	"U",           // u
	"V",           // v
	"W",           // w
	"X",           // x
	"Y",           // y
	"Z",           // z
	"1",           // 1!
	"2",           // 2@
	"3",           // 3#
	"4",           // 4$
	"5",           // 5%
	"6",           // 6^
	"7",           // 7&
	"8",           // 8*
	"9",           // 9(
	"0",           // 0)
	"Enter",       // Enter
	"Esc",         // Escape
	"BS",          // Backspace
	"Tab",         // Tab
	"Space",       // Space
	"Minus",       // -_
	"Equal",       // =+
	"O_SBracket",  // [{
	"C_SBracket",  // ]}
	"Backslash",   // \|
	"Colon",       // ;:
	"Apostrophe",  // '"
	"Backquote",   // `~
	"Comma",       // ,<
	"Period",      // .>
	"Slash",       // /?
	"CapsLock",    // Caps Lock
	"F1",          // F1
	"F2",          // F2
	"F3",          // F3
	"F4",          // F4
	"F5",          // F5
	"F6",          // F6
	"F7",          // F7
	"F8",          // F8
	"F9",          // F9
	"F10",         // F10
	"F11",         // F11
	"F12",         // F12
	"PrintScreen", // Print Screen
	"ScrollLock",  // Scroll Lock
	"Pause",       // Pause
	"Ins",         // Insert
	"Home",        // Home
	"PageUp",      // Page Up
	"Del",         // Delete
	"End",         // End
	"PageDown",    // Page Down
	"Right",       // Right
	"Left",        // Left
	"Down",        // Down
	"Up",          // Up
	"NumLock",     // Num Lock
	"Katakana",    // カタカナ ひらがな
	"Yen",         // ￥|
	"Henkan",      // 変換
	"Muhenkan",    // 無変換
	"L_Control",   // Left Control
	"L_Shift",     // Left Shift
	"L_Alt",       // Left Alt
	"L_Gui",       // Left GUI
	"R_Control",   // Right Control
	"R_Shift",     // Right Shift
	"R_Alt",       // Right Alt
	"R_Gui",       // Right GUI
}

var SC = map[string]byte{
	"A":           0x04, // a
	"B":           0x05, // b
	"C":           0x06, // c
	"D":           0x07, // d
	"E":           0x08, // e
	"F":           0x09, // f
	"G":           0x0A, // g
	"H":           0x0B, // h
	"I":           0x0C, // i
	"J":           0x0D, // j
	"K":           0x0E, // k
	"L":           0x0F, // l
	"M":           0x10, // m
	"N":           0x11, // n
	"O":           0x12, // o
	"P":           0x13, // p
	"Q":           0x14, // q
	"R":           0x15, // r
	"S":           0x16, // s
	"T":           0x17, // t
	"U":           0x18, // u
	"V":           0x19, // v
	"W":           0x1A, // w
	"X":           0x1B, // x
	"Y":           0x1C, // y
	"Z":           0x1D, // z
	"1":           0x1E, // 1!
	"2":           0x1F, // 2@
	"3":           0x20, // 3#
	"4":           0x21, // 4$
	"5":           0x22, // 5%
	"6":           0x23, // 6^
	"7":           0x24, // 7&
	"8":           0x25, // 8*
	"9":           0x26, // 9(
	"0":           0x27, // 0)
	"Enter":       0x28, // Enter
	"Esc":         0x29, // Escape
	"BS":          0x2A, // Backspace
	"Tab":         0x2B, // Tab
	"Space":       0x2C, // Space
	"Minus":       0x2D, // -_
	"Equal":       0x2E, // =+
	"O_SBracket":  0x2F, // [{
	"C_SBracket":  0x30, // ]}
	"Backslash":   0x31, // \|
	"Colon":       0x33, // ;:
	"Apostrophe":  0x34, // '"
	"Backquote":   0x35, // `~
	"Comma":       0x36, // ,<
	"Period":      0x37, // .>
	"Slash":       0x38, // /?
	"CapsLock":    0x39, // Caps Lock
	"F1":          0x3A, // F1
	"F2":          0x3B, // F2
	"F3":          0x3C, // F3
	"F4":          0x3D, // F4
	"F5":          0x3E, // F5
	"F6":          0x3F, // F6
	"F7":          0x40, // F7
	"F8":          0x41, // F8
	"F9":          0x42, // F9
	"F10":         0x43, // F10
	"F11":         0x44, // F11
	"F12":         0x45, // F12
	"PrintScreen": 0x46, // Print Screen
	"ScrollLock":  0x47, // Scroll Lock
	"Pause":       0x48, // Pause
	"Ins":         0x49, // Insert
	"Home":        0x4A, // Home
	"PageUp":      0x4B, // Page Up
	"Del":         0x4C, // Delete
	"End":         0x4D, // End
	"PageDown":    0x4E, // Page Down
	"Right":       0x4F, // Right
	"Left":        0x50, // Left
	"Down":        0x51, // Down
	"Up":          0x52, // Up
	"NumLock":     0x53, // Num Lock
	"Katakana":    0x88, // カタカナ ひらがな
	"Yen":         0x89, // ￥|
	"Henkan":      0x8A, // 変換
	"Muhenkan":    0x8B, // 無変換
	"L_Control":   0xE0, // Left Control
	"L_Shift":     0xE1, // Left Shift
	"L_Alt":       0xE2, // Left Alt
	"L_Gui":       0xE3, // Left GUI
	"R_Control":   0xE4, // Right Control
	"R_Shift":     0xE5, // Right Shift
	"R_Alt":       0xE6, // Right Alt
	"R_Gui":       0xE7, // Right GUI
	"LNPH":        0xFE, // Line Phono Switch
	"LAYOUT":      0xFF, // Layout Switch
}

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

		fmt.Print("[ ")
		for j := 0; j < 13; j++ {
			fmt.Printf("0x%02X ", remapRow[j])
		}
		fmt.Println("]")
	}
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
		fmt.Println("C4NDY KeyVLM Configurator v0.2!")
		fmt.Println("")
	default:
	}

	time.Sleep(100 * time.Millisecond)

	// Finalize the hid package.
	if err := hid.Exit(); err != nil {
		log.Fatal(err)
	}
}
