package main

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
)

type AWSPolly struct {
	PollyClient *polly.Client
}

func (service AWSPolly) Text2Speech(fileName string) error {
	readFile, err := os.ReadFile(fileName)
	if err != nil {
		log.Print("Error while reading the file", err)
	}
	fileString := string(readFile[:])

	speechInput := &polly.SynthesizeSpeechInput{
		OutputFormat: types.OutputFormatMp3,
		Text:         aws.String(fileString),
		VoiceId:      types.VoiceIdAditi,
	}
	synthesizeOutput, err := service.PollyClient.SynthesizeSpeech(context.TODO(), speechInput)
	if err != nil {
		log.Print("Error while synthesizing the output")
	}
	audioName := strings.Split(fileName, ".")
	name := audioName[0] + ".mp3"
	audioFile, err := os.Create(name)
	if err != nil {
		log.Println("Error while creating the audio file")
	}
	_, err = io.Copy(audioFile, synthesizeOutput.AudioStream)
	if err != nil {
		log.Print("Error while copy the text file")
	}
	defer audioFile.Close()
	return err
}

func main() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Println("Error whlie loading the config file")
	}

	service := AWSPolly{
		PollyClient: polly.NewFromConfig(sdkConfig),
	}
	err = service.Text2Speech("hello.txt")
	if err != nil {
		log.Print("Error in giving output")
	} else {
		log.Print("Sucess")
	}
}
