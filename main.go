package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/empijei/cli/lg"
)

var flagEncode bool      // -encode
var flagCodeclist string // -codec

func main() {
	flag.BoolVar(&flagEncode, "encode", false, "Sets the decoder to an encoder instead")
	flag.StringVar(&flagCodeclist, "codec", "smart",
		`Sets the decoder/encoder codec. Multiple codecs can be specified and comma separated:
	they will be applied one on the output of the previous as in a pipeline.
	`)

	flag.Parse()

	args := flag.Args()

	// TODO: better handle help output
	if len(args) != 1 {
		lg.Error("Please provide a single string")
		os.Exit(2)
	}

	buf := args[0]
	sequence := strings.Split(flagCodeclist, ",")
	for _, codec := range sequence {
		//This is to avoid printing twice the final result
		//if i < len(sequence)-2 {
		//fmt.Fprintln(os.Stderr, buf)
		//}
		if out, codecUsed, err := DecodeEncode(buf, flagEncode, codec); err == nil {
			lg.Debugf("Codec %s\n", codecUsed)
			lg.Infof("%s\n", out)
			buf = out
		} else {
			lg.Error(err.Error())
			os.Exit(2)
		}
	}
}

// DecodeEncode takes an input string `buf` and decodes/encodes it (depending on the
// `encode` parameter) with the given `codec`. It returns the encoded/decoded string
// or an error if the process failed.
func DecodeEncode(buf string, encode bool, codec string) (out string, codecUsed string, err error) {
	// Build list of available codecs
	var codecNames []string
	for _, cc := range codecs {
		codecNames = append(codecNames, cc.name)
	}
	codecNamesStr := strings.Join(codecNames, ", ")

	var c CodecC
	if codec == "smart" {
		if encode {
			err = fmt.Errorf("Cannot 'smart' encode, please specify a codec")
			return
		}
		c = SmartDecode(buf)
	} else {
		for _, cc := range codecs {
			if cc.name == codec {
				c = cc.codecCons(buf)
			}
		}
		if c == nil {
			err = fmt.Errorf("Codec not found: '%s'. Supported codecs are: %s\n", codec, codecNamesStr)
			return
		}
	}
	codecUsed = c.Name()
	if encode {
		out = c.Encode()
	} else {
		out = c.Decode()
	}
	return
}
