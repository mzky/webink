//go:build !release

package blink

import (
	"github.com/mzky/weblink/internal/devtools"
)

func (v *View) ShowDevTools(devtoolsCallbacks ...func(devtools *View)) {

	if !v.Mb.Resource.IsExist("__devtools__") {
		v.Mb.Resource.Bind("__devtools__", devtools.FS)
	}

	var callback WkeOnShowDevtoolsCallback = func(hwnd WkeHandle, param uintptr) uintptr {

		view := NewView(v.Mb, hwnd, WKE_WINDOW_TYPE_POPUP, v)

		v.DevTools = view

		for _, cb := range devtoolsCallbacks {
			cb(view)
		}

		view.ForceReload() // 必须刷新才会加载

		return 0
	}

	v.Mb.CallFunc("wkeShowDevtools", uintptr(v.Hwnd), StringToWCharPtr("http://__devtools__/inspector.html"), CallbackToPtr(callback), 0)
}
