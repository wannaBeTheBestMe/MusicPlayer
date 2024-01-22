package playback

import "C"
import (
	"fmt"
	"log"
	"sync"
	"unsafe"
)

//#include <main.h>
//static void print_string(const char* str) {
//    printf("%s\n", str);
//}
import "C"

var (
	PDMu     sync.Mutex
	PDecoder *C.ma_decoder
)

var PDevice *C.ma_device

func PlayAudio(fileName string) {
	file_name := C.CString(fileName)
	defer C.free(unsafe.Pointer(file_name))

	PDecoder = C.create_decoder(file_name)
	if PDecoder == nil {
		log.Fatal("PlayAudio: Could not create decoder")
	}

	PDevice = C.play_audio(file_name, PDecoder)
	if PDevice == nil {
		log.Fatal("PlayAudio: PDevice is nil")
	}

	defer CleanCP()

	PauseAudio(PDevice)
	select {}
}

func CleanCP() {
	if PDevice != nil {
		C.ma_device_uninit(PDevice)
		C.free(unsafe.Pointer(PDevice))
		PDevice = nil
	}
	if PDecoder != nil {
		C.destroy_decoder(PDecoder)
		PDecoder = nil
	}
}

func SetLoopPlayback(loopPlayback bool) {
	PDMu.Lock()
	defer PDMu.Unlock()

	C.set_loop_playback(C.bool(loopPlayback))
}

func SetSilent(silent bool) {
	C.set_silent(C.bool(silent))
}

func GetCurrentPosition(pDecoder *C.ma_decoder) float64 {
	return float64(C.get_current_position(pDecoder))
}

func GetTotalLength(pDecoder *C.ma_decoder) float64 {
	PDMu.Lock()
	defer PDMu.Unlock()

	return float64(C.get_total_length(pDecoder))
}

func PauseAudio(pDevice *C.ma_device) {
	C.pause_audio(pDevice)
}

func ResumeAudio(pDevice *C.ma_device) {
	C.resume_audio(pDevice)
}

func SeekToTime(pDecoder *C.ma_decoder, s uint64) {
	PDMu.Lock()
	defer PDMu.Unlock()

	C.seek_to_time(pDecoder, C.ulonglong(s))
}

func SetVolume(volume float32) {
	C.set_volume(C.float(volume))
}

func GetTotalLengthFromPath(fileName string) float64 {
	file_name := C.CString(fileName)
	//if fileName == "C:\\Users\\Asus\\Music\\MusicPlayer\\Vangelis - Juno to Jupiter (2021) [24 Bit Hi-Res]\\01. Vangelis- Atlas’ push.flac" {
	//	C.print_string(file_name)
	//}
	defer C.free(unsafe.Pointer(file_name))

	pDecoder := C.create_decoder(file_name)
	if pDecoder == nil {
		fmt.Println("GetTotalLengthFromPath: Could not create decoder")
		return 0
	}
	defer C.destroy_decoder(pDecoder)

	return GetTotalLength(pDecoder)
}
