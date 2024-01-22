#define MINIAUDIO_IMPLEMENTATION

#include "miniaudio-0.11.21/miniaudio.h"

#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <math.h>
//#include <unistd.h>  // For sleep() function

float g_linear_volume = 1.0;
static float g_exponential_volume = 1.0;
static bool g_silent = false;
static bool g_loop_playback = false;

void set_loop_playback(bool loop_playback) {
    g_loop_playback = loop_playback;
}

void set_silent(bool silent) {
    g_silent = silent;
}

double get_current_position(ma_decoder *pDecoder) {
    if (pDecoder == NULL) {
        return -1.0;
    }

    ma_uint64 currentFrame;
    if (ma_decoder_get_cursor_in_pcm_frames(pDecoder, &currentFrame) != MA_SUCCESS) {
        return -1.0;
    }

    return (double) currentFrame / (double) pDecoder->outputSampleRate;
}

double get_total_length(ma_decoder *pDecoder) {
    if (pDecoder == NULL) {
        return -1.0;
    }

    ma_uint64 totalFrames;
    if (ma_decoder_get_length_in_pcm_frames(pDecoder, &totalFrames) != MA_SUCCESS) {
        return -1.0;
    }
    return (double) totalFrames / (double) pDecoder->outputSampleRate;
}

void pause_audio(ma_device *pDevice) {
    if (pDevice != NULL) {
        ma_device_stop(pDevice);
    }
}

void resume_audio(ma_device *pDevice) {
    if (pDevice != NULL) {
        ma_device_start(pDevice);
    }
}

void set_volume(float volume) {
    g_linear_volume = volume;
    g_exponential_volume = exp(g_linear_volume - 1.0);
}

void data_callback(ma_device *pDevice, void *pOutput, const void *pInput, ma_uint32 frameCount) {
    ma_decoder *pDecoder = (ma_decoder *) pDevice->pUserData;
    if (pDecoder == NULL) {
        memset(pOutput, 0, frameCount * pDecoder->outputChannels * sizeof(float)); // Silence
        return;
    }

    while (frameCount > 0) {
        ma_uint64 framesRead = 0;
        ma_result result = ma_decoder_read_pcm_frames(pDecoder, pOutput, frameCount, &framesRead);
        if (result != MA_SUCCESS) {
            // Handle the error appropriately. For example, you might want to fill the rest of the buffer with silence.
            memset(pOutput, 0, frameCount * pDecoder->outputChannels * sizeof(float));
            return;
        }

        if (framesRead < frameCount) {
            if (g_loop_playback) {
                ma_decoder_seek_to_pcm_frame(pDecoder, 0);
                continue;
            } else {
                memset((void *) ((float *) pOutput + framesRead * pDecoder->outputChannels), 0,
                       (frameCount - framesRead) * pDecoder->outputChannels * sizeof(float));
                break;
            }
        }

        if (g_silent) {
            memset(pOutput, 0, framesRead * pDecoder->outputChannels * sizeof(float));
        } else {
            float *pOut = (float *) pOutput;
            for (ma_uint32 i = 0; i < framesRead * pDecoder->outputChannels; ++i) {
                pOut[i] *= g_exponential_volume;
            }
        }

        frameCount -= framesRead;
        pOutput = (void *) ((float *) pOutput + framesRead * pDecoder->outputChannels);
    }

    (void) pInput;
}

void seek_to_time(ma_decoder *pDecoder, ma_uint64 targetTimeInSeconds) {
    ma_result result;
    ma_uint64 frameIndex;

    frameIndex = targetTimeInSeconds * pDecoder->outputSampleRate * pDecoder->outputChannels;
    result = ma_decoder_seek_to_pcm_frame(pDecoder, frameIndex);
    if (result != MA_SUCCESS) {
        printf("Failed to seek to PCM frame.\n");
    }
}

ma_device *play_audio(const char *filename, ma_decoder *pDecoder) {
    ma_result result;
    ma_device_config deviceConfig;
    ma_device *pDevice = malloc(sizeof(ma_device));  // Allocate memory for the device

    if (pDevice == NULL) {
        printf("Failed to allocate memory for playback device.\n");
        return NULL;
    }

    result = ma_decoder_init_file(filename, NULL, pDecoder);
    if (result != MA_SUCCESS) {
        printf("Could not initialize decoder.\n");
        free(pDevice);
        return NULL;
    }

    deviceConfig = ma_device_config_init(ma_device_type_playback);
    deviceConfig.playback.format = pDecoder->outputFormat;
    deviceConfig.playback.channels = pDecoder->outputChannels;
    deviceConfig.sampleRate = pDecoder->outputSampleRate;
    deviceConfig.dataCallback = data_callback;
    deviceConfig.pUserData = pDecoder;

    result = ma_device_init(NULL, &deviceConfig, pDevice);
    if (result != MA_SUCCESS) {
        printf("Could not initialize playback device.\n");
        ma_decoder_uninit(pDecoder);
        free(pDevice);
        return NULL;
    }

    result = ma_device_start(pDevice);
    if (result != MA_SUCCESS) {
        printf("Failed to start playback device.\n");
        ma_device_uninit(pDevice);
        ma_decoder_uninit(pDecoder);
        free(pDevice);
        return NULL;
    }

    return pDevice;
}

ma_decoder *create_decoder(const char *filename) {
    ma_decoder *pDecoder = (ma_decoder *) malloc(sizeof(ma_decoder));
    if (pDecoder == NULL) {
        return NULL;
    }
    if (ma_decoder_init_file(filename, NULL, pDecoder) != MA_SUCCESS) {
        free(pDecoder);
        return NULL;
    }

    return pDecoder;
}

void destroy_decoder(ma_decoder *pDecoder) {
    if (pDecoder != NULL) {
        ma_decoder_uninit(pDecoder);
        free(pDecoder);
    }
}

//int main() {
//    const char* filename = "C:\\Users\\Asus\\Music\\MusicPlayer\\Gladiator Soundtrack\\17. Now We Are Free.flac";  // Replace with your audio file path
//    ma_decoder decoder;
//    ma_device* pDevice;
//
//    pDevice = play_audio(filename, &decoder);
//    if (pDevice == NULL) {
//        return -1;  // Error occurred in play_audio
//    }
//
//    // Play for 5 seconds
//    sleep(5);
//
//    // Seek to 10 seconds into the audio
//    seek_to_time(&decoder, 10);
//
//    // Play for another 5 seconds
//    sleep(5);
//
//    // Clean up
//    ma_device_uninit(pDevice);
//    ma_decoder_uninit(&decoder);
//    free(pDevice);
//
//    return 0;
//}
