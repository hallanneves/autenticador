package autenticador

import (
	"encoding/json"
	"log"
	"runtime"
)

//LogError loga a requisicao
func LogError(params ...interface{}) {
	log.Println("=> " + MyCaller())
	for _, param := range params {
		retorno, err := json.Marshal(param)
		if err != nil {
			log.Println("Nao consegui expandir o json")
		}
		log.Println("==> " + string(retorno))
	}
}

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

// MyCaller returns the function that called it :)
func MyCaller() string {
	// Skip GetCallerFunctionName and the function to get the caller of
	return getFrame(3).Function
}
