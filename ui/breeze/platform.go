package breeze

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/egoist/mygo/ui"
)

// Platform is the desktop whose controls Breeze UI draws.
type Platform uint8

const (
	// Auto follows the system the app runs on: Mac on macOS, Windows on Windows, Linux elsewhere.
	Auto Platform = iota
	// Mac draws the controls of macOS 26 and later.
	Mac
	// Windows draws the controls of Windows 11 (WinUI 3, Fluent 2).
	Windows
	// Linux draws the controls of GNOME's libadwaita.
	Linux
)

// Platforms lists the platforms Breeze UI draws, in a stable order.
var Platforms = []Platform{Mac, Windows, Linux}

func (p Platform) String() string {
	switch p {
	case Mac:
		return "macOS"
	case Windows:
		return "Windows"
	case Linux:
		return "Linux"
	default:
		return "Auto"
	}
}

var forced atomic.Uint32

// SetPlatform makes every window draw p's controls instead of the system's, as for screenshots of another platform. SetPlatform(Auto) follows the system again.
func SetPlatform(p Platform) { forced.Store(uint32(p)) }

// Current returns the platform whose controls are drawn outside a Preview: the one SetPlatform set, else the system's.
func Current() Platform {
	if p := Platform(forced.Load()); p != Auto {
		return p
	}
	return system(runtime.GOOS)
}

func system(goos string) Platform {
	switch goos {
	case "darwin", "ios":
		return Mac
	case "windows":
		return Windows
	default:
		return Linux
	}
}

// previews holds the platform each window's Preview is drawing, keyed by the window's context. A view builds on one goroutine at a time, but windows may build on different ones.
var previews sync.Map // *ui.Context -> Platform

// Preview builds body with p's controls in this window, as a gallery showing the three platforms side by side does. Widgets keep the platform they were built with, so it does not matter that they paint after Preview returns.
func Preview(c *ui.Context, p Platform, body func()) {
	prev, had := previews.Load(c)
	if p == Auto {
		previews.Delete(c)
	} else {
		previews.Store(c, p)
	}
	defer func() {
		if had {
			previews.Store(c, prev)
		} else {
			previews.Delete(c)
		}
	}()
	body()
}

// PlatformOf returns the platform whose controls widgets built now in c's window get: the Preview's, else Current.
func PlatformOf(c *ui.Context) Platform {
	if p, ok := previews.Load(c); ok {
		return p.(Platform)
	}
	return Current()
}
