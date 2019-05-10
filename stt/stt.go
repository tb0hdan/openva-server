package stt

import (
	"github.com/tb0hdan/openva-server/api"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func GoogleSTTToOpenVASTT(resp *speechpb.StreamingRecognizeResponse) api.StreamingRecognizeResponse {
	results := make([]*api.StreamingRecognitionResult, 0)
	for _, res := range resp.Results {
		alternatives := make([]*api.SpeechRecognitionAlternative, 0)
		for _, alt := range res.Alternatives {
			words := make([]*api.WordInfo, 0)
			for _, word := range alt.Words {
				wrd := &api.WordInfo{
					StartTime: word.StartTime,
					EndTime:   word.EndTime,
					Word:      word.Word,
				}
				words = append(words, wrd)
			}

			alternative := &api.SpeechRecognitionAlternative{
				Transcript: alt.Transcript,
				Confidence: alt.Confidence,
				Words:      words,
			}
			alternatives = append(alternatives, alternative)
		}

		result := &api.StreamingRecognitionResult{
			Alternatives: alternatives,
			IsFinal:      res.IsFinal,
			Stability:    res.Stability,
		}
		results = append(results, result)
	}

	return api.StreamingRecognizeResponse{
		Results:         results,
		SpeechEventType: api.StreamingRecognizeResponse_SpeechEventType(resp.SpeechEventType),
	}
}
