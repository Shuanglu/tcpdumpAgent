package tcpdump

import (
	"fmt"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	log "github.com/sirupsen/logrus"
)

var (
	snaplenU    uint32        = 262144
	snaplen     int32         = 262144
	promiscuous bool          = false
	timeout     time.Duration = 30 * time.Second
	handle      *pcap.Handle
)

/*
func FindAllDevs() ([]pcap.Interface, error) {
	devs, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}
	return devs, nil
}
*/

func Capture(filePath string, durationTime int) error {
	pcapHandle, err := pcap.OpenLive("any", snaplen, promiscuous, timeout)
	if err != nil {
		log.Warn("Failed to open the pcap handle due to : %q", err)
		return err
	}
	fileHandle, err := os.Create(filePath)
	if err != nil {
		log.Warn("Failed to create the file due to: %q", err)
		return err
	}
	defer fileHandle.Close()
	writer := pcapgo.NewWriter(fileHandle)
	writer.WriteFileHeader(snaplenU, pcapHandle.LinkType())
	packetSource := gopacket.NewPacketSource(pcapHandle, pcapHandle.LinkType())
	//durationTimeInt, err := strconv.Atoi(durationTime)
	end := time.After(time.Duration(durationTime) * time.Duration(time.Second))
	log.Info(fmt.Sprintf("Start capturing of %s", filePath))
	/*
		for {
			select {
			case <-end:
				if err != nil {
					log.Warn("There was error while writing the packet: %q", err)
					return err
				}
				log.Info("break")
				goto end
			default:
				log.Info("working on the capture_1")
				for packet := range packetSource.Packets() {
					log.Info("working on the capture_3")

					if packet != nil {
						err = writer.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
						if err != nil {
							log.Warn("Failed to write the packet due to: %q", err)
						}
					} else {
						continue
					}
				}
				log.Info("working on the capture_2")
			}
		}
	*/

	for {
		select {
		case <-packetSource.Packets():
			for packet := range packetSource.Packets() {
				select {
				case <-end:
					goto final
				default:
					err = writer.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
					if err != nil {
						log.Warn("Failed to write the packet due to: %q", err)
					}
				}
			}
			//goto final
		case <-end:
			if err != nil {
				log.Warn("There was error while writing the packet: %q", err)
				return err
			}
			goto final
		default:
			continue
		}
	}
final:
	log.Info(fmt.Sprintf("Complete the capturing of %s", filePath))
	return nil
}
