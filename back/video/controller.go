// videohandler/videohandler.go
package videohandler

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
)

const (
	rtcpPLIInterval = time.Second * 3
)

type VideoHandler struct {
	peerConnection *webrtc.PeerConnection
	uploader       *s3manager.Uploader
	videoFile      *os.File
}

func NewVideoHandler(peerConnection *webrtc.PeerConnection) *VideoHandler {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		log.Fatal(err)
	}
	uploader := s3manager.NewUploader(sess)

	return &VideoHandler{
		peerConnection: peerConnection,
		uploader:       uploader,
	}
}

// ReceiveTrack is a method that receives a track from a peer and adds it to the peer connection
func (vh *VideoHandler) ReceiveTrack(peerConnectionMap map[string]chan *webrtc.Track, peerID string) {
	if _, ok := peerConnectionMap[peerID]; !ok {
		peerConnectionMap[peerID] = make(chan *webrtc.Track, 1)
	}
	localTrack := <-peerConnectionMap[peerID]
	_, err := vh.peerConnection.AddTrack(localTrack)
	if err != nil {
		log.Fatal(err)
	}
}

// CreateTrack is a method that creates a track from a peer and adds it to the peer connection
func (vh *VideoHandler) CreateTrack(peerConnectionMap map[string]chan *webrtc.Track, currentUserID string) {
	if _, err := vh.peerConnection.AddTransceiver(webrtc.RTPCodecTypeVideo); err != nil {
		log.Fatal(err)
	}

	vh.peerConnection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		// Código para manejar PLI
		go func() {
			ticker := time.NewTicker(rtcpPLIInterval)
			for range ticker.C {
				if rtcpSendErr := vh.peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()

		localTrack, newTrackErr := vh.peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
		if newTrackErr != nil {
			log.Fatal(newTrackErr)
		}

		localTrackChan := make(chan *webrtc.Track, 1)
		localTrackChan <- localTrack
		if existingChan, ok := peerConnectionMap[currentUserID]; ok {
			existingChan <- localTrack
		} else {
			peerConnectionMap[currentUserID] = localTrackChan
		}

		var err error
		
		vh.videoFile, err = os.Create("video.temp")
		if err != nil {
			log.Fatal(err)
		}
		defer vh.videoFile.Close()

		rtpBuf := make([]byte, 1400)
		for {
			i, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				log.Fatal(readErr)
			}

			if _, writeErr := vh.videoFile.Write(rtpBuf[:i]); writeErr != nil && writeErr != io.ErrClosedPipe {
				log.Fatal(writeErr)
			}

			if _, writeErr := localTrack.Write(rtpBuf[:i]); writeErr != nil && writeErr != io.ErrClosedPipe {
				log.Fatal(writeErr)
			}
		}
	})
}
func (vh *VideoHandler) UploadVideo() {
	// Reabrir el archivo para leer desde el principio
	file, err := os.Open("video.temp")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, uploadErr := vh.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("your-bucket-name"),
		Key:    aws.String("video.mp4"),
		Body:   file,
	})
	if uploadErr != nil {
		log.Fatal(uploadErr)
	}
	// Eliminar el archivo temporal después de la subida
	os.Remove("video.temp")
}
