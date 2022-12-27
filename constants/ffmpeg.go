package constants

type FILE_TYPE int

const (
	Audio192K FILE_TYPE = iota
	Video5M
	Video3M
	Video1M
	Video800K
	Video400K
)

type KWARGS int

const (
	Preset KWARGS = iota
	Tune
	FpsMode
	AudioExclusion
	VideoExclusion
	AudioCodec
	VideoCodec
	AudioBitrate
	MaxRate
	ConstantRateFactor
	BufferSize
	AllowSoftEncoding
	VideoFormat
	HWAccel
)

var AudioFileTypeMap = map[FILE_TYPE]string{
	Audio192K: "_audio192k.m4a",
}

var VideoFileTypeMap = map[FILE_TYPE]string{
	Video5M:   "_5000.mp4",
	Video3M:   "_3000.mp4",
	Video1M:   "_1500.mp4",
	Video800K: "_800.mp4",
	Video400K: "_400.mp4",
}

var AudioBitrateMap = map[FILE_TYPE]string{
	Audio192K: "192k",
}

var VideoBitrateMap = map[FILE_TYPE]string{
	Video5M:   "5000k",
	Video3M:   "3000k",
	Video1M:   "1500k",
	Video800K: "800k",
	Video400K: "400k",
}

var VideoBufferSizeMap = map[FILE_TYPE]string{
	Video5M:   "12000k",
	Video3M:   "6000k",
	Video1M:   "3000k",
	Video800K: "2000k",
	Video400K: "1000k",
}

var VideoKwargs = map[KWARGS]string{
	Preset:             "preset",
	Tune:               "tune",
	FpsMode:            "fps_mode",
	AudioExclusion:     "an",
	VideoCodec:         "c:v",
	ConstantRateFactor: "crf",
	MaxRate:            "maxrate",
	BufferSize:         "bufsize",
	VideoFormat:        "f",
	HWAccel:            "hwaccel",
}

var AudioKwargs = map[KWARGS]string{
	AudioCodec:        "c:a",
	AudioBitrate:      "b:a",
	AllowSoftEncoding: "allow_sw",
	VideoExclusion:    "vn",
	HWAccel:           "hwaccel",
}

var FFmpegConfig = map[KWARGS]string{
	Preset:         "fast",
	Tune:           "film",
	FpsMode:        "passthrough",
	AudioExclusion: "",
	VideoExclusion: "",
	AudioCodec:     "aac",
	// TODO: add support for different encoders
	// https://stackoverflow.com/a/50703794
	VideoCodec:         "h264_videotoolbox",
	ConstantRateFactor: "23",
	AllowSoftEncoding:  "1",
	VideoFormat:        "mp4",
	HWAccel:            "auto",
}
