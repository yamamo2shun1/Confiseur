package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sstallion/go-hid"
)

const VERSION = "v0.14.1"

func main() {
	checkFlag := flag.Bool("check", false, "Show information on C4NDY KeyVLM/STK connected to PC/Mac.")
	listFlag := flag.Bool("list", false, "Show connected device list.")
	id := flag.Int("id", 0, "Select connected device ID.")
	loadFlag := flag.Bool("load", false, "Show the current key names of the keyboard.")
	remapFile := flag.String("remap", "", "Write the keyboard with the keymap set in toml.")
	saveFlag := flag.Bool("save", false, "Save the keymap written by \"-remap\" to the memory area.")
	restartFlag := flag.Bool("restart", false, "Restart the keyboard immediately.")
	ledColor := flag.Int("led", -1, "Set LED RGB value for checking color.")
	ledIntensity := flag.Float64("intensity", -1.0, "Set LED intensity.")
	factoryresetFlag := flag.Bool("factoryreset", false, "Reset all settings to factory defaults.")
	verFlag := flag.Bool("version", false, "Show the version of the tool installed.")

	flag.Parse()

	// Initialize the hid package.
	if err := hid.Init(); err != nil {
		log.Fatal(err)
	}

	// Open the device using the VID and PID.
	openConnectedHIDDevices()

	if getConnectedDeviceNum() > 0 {
		checkKeyboardType(*id)

		if *checkFlag {
			checkHid()
			fmt.Println("")
		} else if *listFlag {
			deviceList := getConnectedDeviceList()
			for i, deviceName := range deviceList {
				fmt.Printf("%d: %s\n", i, deviceName)
			}
		} else if *loadFlag {
			swapKeyCodeAndName()

			// check current hardware layout
			fmt.Println("--- Current Hardware Layout ScanCode ---")
			fmt.Println("::Layout1::")
			fmt.Println("  Normal ->")
			loadKeymap(*id, 0x11)
			if isStk {
				fmt.Println("  Upper ->")
				loadKeymap(*id, 0x12)
				fmt.Println("  Stick ->")
				loadKeymap(*id, 0x13)
				fmt.Println("  Led ->")
				loadKeymap(*id, 0x14)
				fmt.Println("  Intensity ->")
				loadKeymap(*id, 0x15)
			}
			fmt.Println("")
			fmt.Println("::Layout2::")
			fmt.Println("  Normal ->")
			loadKeymap(*id, 0x19)
			if isStk {
				fmt.Println("  Upper ->")
				loadKeymap(*id, 0x1A)
				fmt.Println("  Stick ->")
				loadKeymap(*id, 0x1B)
				fmt.Println("  Led ->")
				loadKeymap(*id, 0x1C)
				fmt.Println("  Intensity ->")
				loadKeymap(*id, 0x1D)
			}
			fmt.Println("")
		} else if *saveFlag {
			saveToFlash(*id)
			fmt.Println("")
		} else if *restartFlag {
			restart(*id)
			fmt.Println("")
		} else if *factoryresetFlag {
			factoryReset(*id)
			fmt.Println("")
		} else if *ledColor >= 0 && *ledColor <= 0xFFFFFF {
			checkLEDColor(*id, *ledColor)
		} else if *ledIntensity >= 0.0 && *ledIntensity <= 1.0 {
			changeLEDIntensity(*id, *ledIntensity)
		} else if *remapFile != "" {
			if _, err := os.Stat(*remapFile); os.IsNotExist(err) {
				fmt.Printf("::ERROR:: \"%s\" is not existed.\n", *remapFile)
				os.Exit(0)
			}

			remap(*remapFile)

			fmt.Println("remap layout1&2(Normal/Upper)")

			writeKeymap(*id, 0x01)
			if isStk {
				writeKeymap(*id, 0x02)
				writeKeymap(*id, 0x03)
				writeKeymap(*id, 0x04)
				writeKeymap(*id, 0x05)
			}
			writeKeymap(*id, 0x09)
			if isStk {
				writeKeymap(*id, 0x0A)
				writeKeymap(*id, 0x0B)
				writeKeymap(*id, 0x0C)
				writeKeymap(*id, 0x0D)
			}
			fmt.Println("")
		}
		time.Sleep(
			100 * time.Millisecond)
	} else {
		fmt.Println("::WARNING:: C4NDY KeyVLM/STK is not found.")
	}

	if *verFlag {
		fmt.Printf("C4NDY Confiseur %s!\n", VERSION)
		fmt.Println("")
	}

	// Finalize the hid package.
	if err := hid.Exit(); err != nil {
		log.Fatal(err)
	}
}
