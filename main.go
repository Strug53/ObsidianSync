package main

import (
	"fmt"
	"time"
)

func consoleGUI() {

	var Reset = "\033[0m"
	var Red = "\033[31m"
	var Green = "\033[32m"
	var Yellow = "\033[33m"
	var Blue = "\033[34m"

	fmt.Println("OBSIDIAN SYNC")
	for {
		time.Sleep(15 * time.Second)
		fmt.Println(Blue + "MENU:" + Reset)
		fmt.Println("\t" + Green + "- See all devices")
		fmt.Println("\t" + Red + "- Connect to device")
		fmt.Println("\t" + Yellow + "- ??" + Reset)
		/*
			for _, el := range devices {
				fmt.Println("\t -", el.IP, "|", el.Name)
			}*/
	}

}

func main() {

}
