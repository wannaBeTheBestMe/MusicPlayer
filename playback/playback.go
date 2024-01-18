package playback

//#include <main.h>
import "C"
import (
	"log"
	"unsafe"
)

var PDecoder *C.ma_decoder
var PDevice *C.ma_device

func PlayAudio(fileName string) {
	file_name := C.CString(fileName)
	defer C.free(unsafe.Pointer(file_name))

	PDecoder = C.create_decoder(file_name)
	if PDecoder == nil {
		log.Fatal("Could not create decoder")
	}
	defer C.destroy_decoder(PDecoder)

	PDevice = C.play_audio(file_name, PDecoder)
	if PDevice == nil {
		log.Fatal("PDevice is nil")
	}
	defer C.free(unsafe.Pointer(PDevice))
	defer C.ma_device_uninit(PDevice)

	PauseAudio(PDevice)
	select {}
}

func SetSilent(silent bool) {
	C.set_silent(C.bool(silent))
}

func GetCurrentPosition(PDecoder *C.ma_decoder) float64 {
	return float64(C.get_current_position(PDecoder))
}

func GetTotalLength(PDecoder *C.ma_decoder) float64 {
	return float64(C.get_total_length(PDecoder))
}

func PauseAudio(PDevice *C.ma_device) {
	C.pause_audio(PDevice)
}

func ResumeAudio(PDevice *C.ma_device) {
	C.resume_audio(PDevice)
}

func SeekToTime(PDecoder *C.ma_decoder, s uint64) {
	C.seek_to_time(PDecoder, C.ulonglong(s))
}

func SetVolume(volume float32) {
	C.set_volume(C.float(volume))
}
