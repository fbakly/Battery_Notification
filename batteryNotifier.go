package main

import (
	"os"
	"time"

	"github.com/distatus/battery"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/mqu/go-notify"
)

var beepPath string = "/home/fbakly/.scripts/resources/beep.mp3"
var inboxPath string = "/home/fbakly/.scripts/resources/filling-your-inbox.mp3"
var notificationTime time.Duration = 2

func playAudio(pathToAudio string) {
	f, err := os.Open(pathToAudio)
	if err == nil {
		streamer, format, err := mp3.Decode(f)
		if err == nil {
			defer streamer.Close()
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

			done := make(chan bool)
			speaker.Play(beep.Seq(streamer, beep.Callback(func() {
				done <- true
			})))
			<-done
			close(done)
			f.Close()
		}
	}
}

func sendNotification(title string, content string, icon string, urgency notify.NotifyUrgency, finished chan bool) {
	notification := notify.NotificationNew(title, content, icon)
	notify.NotificationSetUrgency(notification, urgency)
	notification.Show()
	time.Sleep(notificationTime * time.Second)
	notify.NotificationClose(notification)
	finished <- true
}

func notifyBatteryLevel(title string, content string, icon string, urgency notify.NotifyUrgency, notified *bool) {
	if !*(notified) {
		finished := make(chan bool)
		*notified = true
		go sendNotification(title, content, icon, urgency, finished)
		playAudio(beepPath)
		<-finished
	}
}

func notifyCharge(isCharging bool, chargeNotified *bool) {
	if isCharging && !(*chargeNotified) {
		finished := make(chan bool)
		*chargeNotified = true
		go sendNotification("Charging Battery", "", "dialog-information", notify.NOTIFY_URGENCY_NORMAL, finished)
		playAudio(inboxPath)
		<-finished
	} else if !isCharging && *chargeNotified {
		finished := make(chan bool)
		*chargeNotified = false
		go sendNotification("Battery Discharging", "", "dialog-information", notify.NOTIFY_URGENCY_NORMAL, finished)
		playAudio(inboxPath)
		<-finished
	}
}

func getBatteryState(batteries []*battery.Battery, isCharging *bool) {
	length := len(batteries)

	*isCharging = false
	for index := 0; index < length; index++ {
		if batteries[index].State.String() == "Charging" {
			*isCharging = true
		}
	}
}

func main() {
	chargeNotified := false
	isCharging := false
	notified := false
	notify.Init("Battery Percentage Notifier")
	for {
		percentage := 0
		batteries, err := battery.GetAll()

		if err != nil {
			return
		}

		getBatteryState(batteries, &isCharging)

		for _, battery := range batteries {
			percentage += int((battery.Current / battery.Full) * 100)
		}

		percentage /= len(batteries)

		if percentage == 100 {
			notifyBatteryLevel("Battery Full", "", "dialog-information", notify.NOTIFY_URGENCY_NORMAL, &notified)
		} else if percentage <= 20 {
			notifyBatteryLevel("Battery Low", "Charge Now !!!", "dialog-information", notify.NOTIFY_URGENCY_CRITICAL, &notified)
		} else {
			notified = false
		}

		notifyCharge(isCharging, &chargeNotified)
	}
}
