# SmartDecode
[![License](https://img.shields.io/badge/license-GPLv3-blue.svg)](https://raw.githubusercontent.com/AnnaOpss/smartDecode/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/annaopss/smartdecode)](https://goreportcard.com/report/github.com/AnnaOpss/smartDecode)
[![Coverage](https://gocover.io/_badge/github.com/annaopss/smartdecode?nocache=smartdecode)](https://gocover.io/github.com/AnnaOpss/smartDecode)

This package provides a finate state machine to smartly decode strings.

It used to be a package of [Wapty](https://github.com/empijei/wapty), but then grow as a standalone tool.

```
Usage of smartDecode:
  -codec string
    	Sets the decoder/encoder codec. Multiple codecs can be specified and comma separated:
    		they will be applied one on the output of the previous as in a pipeline.
    		 (default "smart")
  -encode
    	Sets the decoder to an encoder instead. If this flag is present, a codec needs to be specified.
```
