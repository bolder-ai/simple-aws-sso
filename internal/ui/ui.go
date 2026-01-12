package ui

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	successIcon = color.GreenString("✓")
	errorIcon   = color.RedString("✗")
	infoIcon    = color.CyanString("→")
	waitIcon    = color.YellowString("◌")
)

func Success(format string, args ...any) {
	fmt.Printf("%s %s\n", successIcon, fmt.Sprintf(format, args...))
}

func Error(format string, args ...any) {
	fmt.Printf("%s %s\n", errorIcon, color.RedString(format, args...))
}

func Info(format string, args ...any) {
	fmt.Printf("%s %s\n", infoIcon, fmt.Sprintf(format, args...))
}

func Wait(format string, args ...any) {
	fmt.Printf("%s %s\n", waitIcon, color.YellowString(format, args...))
}

func Bold(format string, args ...any) string {
	return color.New(color.Bold).Sprintf(format, args...)
}

func Cyan(format string, args ...any) string {
	return color.CyanString(format, args...)
}

func Green(format string, args ...any) string {
	return color.GreenString(format, args...)
}

func Yellow(format string, args ...any) string {
	return color.YellowString(format, args...)
}

func Dim(format string, args ...any) string {
	return color.New(color.Faint).Sprintf(format, args...)
}

func Print(format string, args ...any) {
	fmt.Printf(format, args...)
}

func Println(args ...any) {
	fmt.Println(args...)
}
