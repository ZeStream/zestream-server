package constants

type FILE_TYPE int

const (
	Audio192K FILE_TYPE = iota
	ImagePng
	Video5M
	Video3M
	Video1M
	Video800K
	Video400K
)

type FFMPEG_KWARGS int

const (
	Preset FFMPEG_KWARGS = iota
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
	VideoFrames
	ScreenShot
)

type MP4BOX_ARGS int

const (
	Dash MP4BOX_ARGS = iota
	Rap
	FragRap
	BsSwitching
	Profile
	Out
)

const DashOutputExt = ".mpd"
const MP4Box = "MP4Box"
const DEFAULT_THUMBNAIL_TIMESTAMP = "00:00:02"
const Overlay = "overlay"
const Scale = "scale"

var AudioFileTypeMap = map[FILE_TYPE]string{
	Audio192K: "_audio192k.m4a",
}

var ImageFileTypeMap = map[FILE_TYPE]string{
	ImagePng: ".png",
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

var VideoKwargs = map[FFMPEG_KWARGS]string{
	Preset:             "preset",
	Tune:               "tune",
	FpsMode:            "fps_mode",
	AudioExclusion:     "an",
	VideoCodec:         "c:v",
	ConstantRateFactor: "crf",
	MaxRate:            "maxrate",
	BufferSize:         "bufsize",
	HWAccel:            "hwaccel",
	VideoFormat:        "f",
	VideoFrames:        "frames:v",
	ScreenShot:         "ss",
}

var AudioKwargs = map[FFMPEG_KWARGS]string{
	AudioCodec:        "c:a",
	AudioBitrate:      "b:a",
	AllowSoftEncoding: "allow_sw",
	VideoExclusion:    "vn",
	HWAccel:           "hwaccel",
}

var FFmpegConfig = map[FFMPEG_KWARGS]string{
	Preset:         "fast",
	Tune:           "film",
	FpsMode:        "passthrough",
	AudioExclusion: "",
	VideoExclusion: "",
	AudioCodec:     "aac",
	// TODO: add support for different encoders
	// https://stackoverflow.com/a/50703794
	VideoCodec:         "libx264",
	ConstantRateFactor: "23",
	AllowSoftEncoding:  "1",
	VideoFormat:        "mp4",
	HWAccel:            "auto",
	VideoFrames:        "1",
}

var Mp4BoxArgs = map[MP4BOX_ARGS]string{
	Dash:        "dash",
	Rap:         "rap",
	FragRap:     "frag-rap",
	BsSwitching: "bs-switching",
	Profile:     "profile",
	Out:         "out",
}

var Mp4BoxConfig = map[MP4BOX_ARGS]string{
	Dash:        "10000",
	Rap:         "",
	FragRap:     "",
	BsSwitching: "no",
	Profile:     "dashavc264:live",
}
