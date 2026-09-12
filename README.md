<a href="https://pkg.go.dev/github.com/kandayin/go-espeak-ng"><img src="https://pkg.go.dev/badge/github.com/kandayin/go-espeak-ng" alt="PkgGoDev"></a>
<a href="https://ci.djangulo.com/teams/djangulo/pipelines/go-espeak"><img src="https://ci.djangulo.com/api/v1/teams/djangulo/pipelines/go-espeak/jobs/test-unit/badge" alt="CI status"></a>

# go-espeak-ng

Golang C bindings for the `espeak` voice synthesizer.

There is a live demo of its usage at <a rel="noopener noreferrer" target="_blank" href="https://go-espeak-demo.djangulo.com">https://go-espeak-demo.djangulo.com</a>, source code in [examples/demo](https://github.com/kandayin/go-espeak-ng/tree/main/examples/demo).

Sub-package `native` contains a mostly C implementation, minimizing the amount of Go used. This implementation is slightly faster than the go implementation, with the inconvenience of being a black box from the input to the `.wav`.

## Requirements

- Go >= 1.15 with `cgo` support
- <a target="_blank" rel="noopener noreferrer" href="https://github.com/espeak-ng/espeak-ng">`espeak-ng`</a>.

## Install

### Install requirements

Arch

```bash
~# pacman -S espeak-ng
```

Ubuntu

```bash
~# apt-get install espeak-ng
```

### Install go-espeak-ng

```bash
~$ go get -u github.com/kandayin/go-espeak-ng
```

## Usage

[examples/basic-usage](https://github.com/kandayin/go-espeak-ng/tree/main/examples/basic-usage).

```golang
package main

import (
	"github.com/kandayin/go-espeak-ng"
)

func main() {

	// need to call terminate so espeak can clean itself out
	defer espeak.Terminate()
	params := espeak.NewParameters().WithDir(".")
	espeak.TextToSpeech(
		"Hello World!", // Text to speak
		nil,            // voice to use, nil == DefaultVoice (en-us male)
		"hello.wav",    // if "" or "play", it plays to default audio out
		params,         // Parameters for voice modulation, nil == DefaultParameters
	)

	// get a random spanish voice
	v, _ := espeak.VoiceFromSpec(&espeak.Voice{Languages: "es"})
	espeak.TextToSpeech("¡Hola mundo!", v, "hola.wav", params)
}

```
