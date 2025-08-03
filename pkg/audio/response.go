package audio

//Response Body: {
//    "text":"Make me a story about Strocki the ostrich who met in friendly Puma.",
//    "language":"english",
//    "task":"transcribe",
//    "duration":11.449999809265137,
//    "words":null,
//    "segments":[
//        {"id":0,
//         "avg_logprob":-0.6582053899765015,
//         "compression_ratio":0.9436619877815247,
//         "end":11.5,
//         "no_speech_prob":0.15403856337070465,
//         "seek":0,
//         "start":0.0,
//         "temperature":0.0,
//         "text":" Make me a story about Strocki the ostrich who met in friendly Puma.",
//         "tokens":[50364,4387,385,257,1657,466,745,17799,72,264,44024,480,567,1131,294,9208,430,5544,13,50939]
//        }]}

type AudioResponse struct {
	Text             string    `json:"text"`
	Language         string    `json:"language"`
	Task             string    `json:"task"`
	Duration         float32   `json:"duration"`
	Words            []Word    `json:"words"`
	Segments         []Segment `json:"segments"`
	CompressionRatio float32   `json:"compression_ratio"`
	NoSpeechProb     float32   `json:"no_speech_prob"`
	AvgLogProb       float32   `json:"avg_logprob"`
	Seek             int       `json:"seek"`
	Start            float32   `json:"start"`
	End              float32   `json:"end"`
	Temperature      float32   `json:"temperature"`
	TextSegment      string    `json:"text_segment"`
	Tokens           []int     `json:"tokens"`
}

type Word struct {
	Word  string  `json:"word"`
	Start float32 `json:"start"`
	End   float32 `json:"end"`
}

type Segment struct {
	ID               int     `json:"id"`
	AvgLogProb       float32 `json:"avg_logprob"`
	CompressionRatio float32 `json:"compression_ratio"`
	NoSpeechProb     float32 `json:"no_speech_prob"`
	Seek             int     `json:"seek"`
	Start            float32 `json:"start"`
	End              float32 `json:"end"`
	Temperature      float32 `json:"temperature"`
	Text             string  `json:"text"`
	Tokens           []int   `json:"tokens"`
}
