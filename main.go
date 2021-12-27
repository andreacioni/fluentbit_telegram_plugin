package main

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
)

import "C"

const PluginName = "telegram"
const PlugingDesc = "Telegram"

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	// Gets called only once when the plugin.so is loaded
	log.Printf("[%s] registering plugin", PluginName)
	return output.FLBPluginRegister(def, PluginName, PlugingDesc)
}

//export FLBPluginInit
func FLBPluginInit(plugin unsafe.Pointer) int {
	// Gets called only once for each instance you have configured.
	apiKey := output.FLBPluginConfigKey(plugin, "api_key")
	chatId := output.FLBPluginConfigKey(plugin, "chat_id")

	log.Printf("[%s] api_key = %q, chat_id = %q", PluginName, apiKey, chatId)

	// Set the context to point to any Go variable
	output.FLBPluginSetContext(plugin, apiKey)
	output.FLBPluginSetContext(plugin, chatId)

	return output.FLB_OK
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, tag *C.char) int {
	// Gets called with a batch of records to be written to an instance.

	// Type assert context back into the original type for the Go variable
	apiKey := output.FLBPluginConfigKey(ctx, "api_key")
	chatId := output.FLBPluginConfigKey(ctx, "chat_id")

	log.Printf("[%s] flush called %s %s", PluginName, apiKey, chatId)

	dec := output.NewDecoder(data, int(length))

	for {
		ret, ts, record := output.GetRecord(dec)
		if ret != 0 {
			break
		}

		// Print record keys and values
		timestamp := ts.(output.FLBTime)
		str := fmt.Sprintf("%s [%s, {", C.GoString(tag), timestamp.String())

		for k, v := range record {
			str += fmt.Sprintf("\"%s\": %v, ", k, v)
		}
		str += "}"

		if err := SendTelegramMessage(apiKey, chatId, str); err != nil {
			log.Printf("[%s] telegram notification failed: %+v", PluginName, err)
			return output.FLB_ERROR
		}
	}

	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	return output.FLB_OK
}

func main() {
}
